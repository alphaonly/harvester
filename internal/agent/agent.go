package agent

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"log"
	"net/http"
	"runtime"

	"bytes"
	"strconv"
	"time"

	"math/rand"
	"net/url"
)

type Configuration struct {
	PollInterval   int
	ReportInterval int
	ServerHost     string
	ServerPort     string
	UseJSON        bool
}

type gauge float64
type counter int64

type Metrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	RandomValue   gauge

	PollCount counter
}

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, пID    string   `json:"id"`              // имя метрикиринимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (j MetricsJSON) newMetricJson(name string, MType string, value interface{}) (ret MetricsJSON) {

	j = MetricsJSON{}
	if name == "" || MType == "" {
		panic(errors.New("name or type is empty"))
	}

	j.ID = name
	j.MType = MType

	if value != nil {
		switch MType {
		case "agent.gauge", "gauge":
			var val = float64(value.(gauge))
			j.Value = &val
		case "agent.counter", "counter":
			var val = int64(value.(counter))
			j.Delta = &val
		default:
			panic(errors.New("unknown type"))
		}
	}
	return j
}

type Agent struct {
	Configuration *Configuration
	baseURL       url.URL
}

func NewAgent(c *Configuration) Agent {
	return Agent{
		Configuration: c,
		baseURL: url.URL{
			Scheme: "http",
			Host:   c.ServerHost + ":" + c.ServerPort,
		},
	}
}

func (a Agent) GetMetricJson(name string, MType string) (mj MetricsJSON) {

	metricsJsonRequest := MetricsJSON{}.newMetricJson(name, MType, nil)
	data, err := json.Marshal(metricsJsonRequest)
	if err != nil {
		log.Fatal(err)
	}

	URL := a.baseURL.JoinPath("value")
	request, err := http.NewRequest(http.MethodGet, URL.String(), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Accept", "application/json; charset=utf-8")

	client := http.Client{}
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var metricsJsonResponse MetricsJSON
	err = json.Unmarshal(responseData, &metricsJsonResponse)
	if err != nil {
		log.Fatal(err)
	}
	return metricsJsonResponse

}

func AddCounterData(urlPref *url.URL, val counter, name string, data *map[sendData]bool) {
	URL := urlPref.
		JoinPath("counter").
		JoinPath(name).
		JoinPath(strconv.FormatUint(uint64(val), 10)) //value float
	sd := sendData{
		url:  *URL,
		body: bytes.NewBufferString(url.Values{}.Encode()), //need to transer something
	}
	(*data)[sd] = true

}
func AddGaugeData(urlPref *url.URL, val gauge, name string, data *map[sendData]bool) {

	URL := urlPref.
		JoinPath("gauge").
		JoinPath(name).
		JoinPath(strconv.FormatFloat(float64(val), 'E', -1, 64)) //value float

	sd := sendData{
		url:  *URL,
		body: bytes.NewBufferString(url.Values{}.Encode()), //need to transer something
	}
	(*data)[sd] = true

}

type sendData struct {
	url  url.URL
	body io.Reader
}

func (data sendData) sendDataURL(client *http.Client) error {

	request, err := http.NewRequest(http.MethodPost, data.url.String(), data.body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("url from agent):%s", data.url.String())
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("response from server:" + response.Status)
	return err
}

func (a Agent) Update(ctx context.Context, metrics *Metrics) {

	var m runtime.MemStats

	ticker := time.NewTicker(time.Duration(a.Configuration.PollInterval) * time.Second)

	defer ticker.Stop()

repeatAgain:
	select {
	case <-ticker.C:
		{
			runtime.ReadMemStats(&m)

			metrics.Alloc = gauge(m.Alloc)
			metrics.BuckHashSys = gauge(m.BuckHashSys)
			metrics.Frees = gauge(m.Frees)
			metrics.GCCPUFraction = gauge(m.GCCPUFraction)
			metrics.GCSys = gauge(m.GCSys)
			metrics.HeapAlloc = gauge(m.HeapAlloc)
			metrics.HeapIdle = gauge(m.HeapIdle)
			metrics.HeapInuse = gauge(m.HeapInuse)
			metrics.HeapObjects = gauge(m.HeapObjects)
			metrics.HeapReleased = gauge(m.HeapReleased)
			metrics.HeapSys = gauge(m.HeapSys)
			metrics.LastGC = gauge(m.LastGC)
			metrics.Lookups = gauge(m.Lookups)
			metrics.MCacheInuse = gauge(m.MCacheInuse)
			metrics.MCacheSys = gauge(m.MCacheSys)
			metrics.MSpanInuse = gauge(m.MSpanInuse)
			metrics.MSpanSys = gauge(m.MSpanSys)
			metrics.Mallocs = gauge(m.Mallocs)
			metrics.NextGC = gauge(m.NextGC)
			metrics.NumForcedGC = gauge(m.NumForcedGC)
			metrics.NumGC = gauge(m.NumGC)
			metrics.OtherSys = gauge(m.OtherSys)
			metrics.PauseTotalNs = gauge(m.PauseTotalNs)
			metrics.StackInuse = gauge(m.StackInuse)
			metrics.StackSys = gauge(m.StackSys)
			metrics.Sys = gauge(m.Sys)
			metrics.TotalAlloc = gauge(m.TotalAlloc)

			goto repeatAgain
		}

	case <-ctx.Done():
		{
			log.Println("Metrics reading cancelled by context")
			return
		}
	}

}
func (a Agent) prepareData(metrics *Metrics) map[sendData]bool {
	m := make(map[sendData]bool)

	switch a.Configuration.UseJSON {
	case true:
		{
			//Mocked up
			//url := a.baseURL.
			//	JoinPath("update")
		}
	default:
		{
			URL := a.baseURL.
				JoinPath("update")

			AddGaugeData(URL, metrics.Alloc, "Alloc", &m)
			AddGaugeData(URL, metrics.Frees, "Frees", &m)
			AddGaugeData(URL, metrics.GCCPUFraction, "GCCPUFraction", &m)
			AddGaugeData(URL, metrics.GCSys, "GCSys", &m)
			AddGaugeData(URL, metrics.HeapAlloc, "HeapAlloc", &m)
			AddGaugeData(URL, metrics.HeapIdle, "HeapIdle", &m)
			AddGaugeData(URL, metrics.HeapInuse, "HeapInuse", &m)
			AddGaugeData(URL, metrics.HeapObjects, "HeapObjects", &m)
			AddGaugeData(URL, metrics.HeapReleased, "HeapReleased", &m)
			AddGaugeData(URL, metrics.HeapSys, "HeapSys", &m)
			AddGaugeData(URL, metrics.LastGC, "LastGC", &m)
			AddGaugeData(URL, metrics.Lookups, "Lookups", &m)
			AddGaugeData(URL, metrics.MCacheSys, "MCacheSys", &m)
			AddGaugeData(URL, metrics.MSpanInuse, "MSpanInuse", &m)
			AddGaugeData(URL, metrics.MSpanSys, "MSpanSys", &m)
			AddGaugeData(URL, metrics.Mallocs, "Mallocs", &m)
			AddGaugeData(URL, metrics.NextGC, "NextGC", &m)
			AddGaugeData(URL, metrics.NumForcedGC, "NumForcedGC", &m)
			AddGaugeData(URL, metrics.NumGC, "NumGC", &m)
			AddGaugeData(URL, metrics.OtherSys, "OtherSys", &m)
			AddGaugeData(URL, metrics.PauseTotalNs, "PauseTotalNs", &m)
			AddGaugeData(URL, metrics.StackInuse, "StackInuse", &m)
			AddGaugeData(URL, metrics.StackSys, "StackSys", &m)
			AddGaugeData(URL, metrics.Sys, "Sys", &m)
			AddGaugeData(URL, metrics.TotalAlloc, "TotalAlloc", &m)
			AddGaugeData(URL, metrics.RandomValue, "RandomValue", &m)
			AddCounterData(URL, metrics.PollCount, "", &m)

		}

	}

	return m
}
func (a Agent) Send(ctx context.Context, client *http.Client, metrics *Metrics) {

	ticker := time.NewTicker(time.Duration(a.Configuration.ReportInterval) * time.Second)
	defer ticker.Stop()

repeatAgain:
	select {
	case <-ticker.C:
		{
			metrics.PollCount++
			metrics.RandomValue = gauge(rand.Int63())
			dataPackage := a.prepareData(metrics)

			for key := range dataPackage {
				err := key.sendDataURL(client)
				if err != nil {
					log.Println(err)
					return
				}
			}
			goto repeatAgain
		}

	case <-ctx.Done():
		{
			break
		}
	}

}
func (a Agent) Run(ctx context.Context, client *http.Client) {

	metrics := Metrics{}

	go a.Update(ctx, &metrics)
	go a.Send(ctx, client, &metrics)

}
