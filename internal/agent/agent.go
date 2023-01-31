package agent

import (
	"context"
	"encoding/json"
	"io"

	metricsjson "github.com/alphaonly/harvester/internal/server/metricsJSON"

	"log"
	"net/http"
	"runtime"

	"bytes"
	"strconv"
	"time"

	"math/rand"
	"net/url"

	C "github.com/alphaonly/harvester/internal/configuration"
)

type Gauge float64
type Counter int64

type Metrics struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	RandomValue   Gauge

	PollCount Counter
}

type Agent struct {
	Configuration *C.Configuration
	baseURL       url.URL
}

func NewAgent(c *C.Configuration) Agent {

	return Agent{
		Configuration: c,
		baseURL: url.URL{
			Scheme: (*c).Get("SCHEME"),
			Host:   (*c).Get("HOST") + ":" + (*c).Get("PORT"),
		},
	}
}

func AddCounterData(urlPref *url.URL, val Counter, name string, data *map[sendData]bool) {
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
func AddGaugeData(urlPref *url.URL, val Gauge, name string, data *map[sendData]bool) {

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

func AddGaugeDataJSON(urlPref *url.URL, val Gauge, name string, data *map[sendData]bool) {
	v := float64(val)
	mj := metricsjson.MetricsJSON{
		ID:    name,
		MType: "gauge",
		Value: &v,
	}

	log.Println(mj)
	metricsBytes, err := json.Marshal(mj)
	if err != nil {
		log.Fatal(err)
	}

	sd := sendData{
		url:  *urlPref,
		body: bytes.NewBuffer(metricsBytes),
	}
	(*data)[sd] = true

}
func AddCounterDataJSON(urlPref *url.URL, val Counter, name string, data *map[sendData]bool) {
	v := int64(val)
	mj := metricsjson.MetricsJSON{
		ID:    name,
		MType: "counter",
		Delta: &v,
	}
	metricsBytes, err := json.Marshal(mj)
	if err != nil {
		log.Fatal(err)
	}

	sd := sendData{
		url:  *urlPref,
		body: bytes.NewBuffer(metricsBytes),
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
	err = response.Body.Close()
	return err
}

func (a Agent) Update(ctx context.Context, metrics *Metrics) {
	var m runtime.MemStats
	ticker := time.NewTicker(time.Duration((*a.Configuration).GetInt("POLL_INTERVAL")) * time.Second)
	defer ticker.Stop()
repeatAgain:
	select {
	case <-ticker.C:
		{
			runtime.ReadMemStats(&m)

			metrics.Alloc = Gauge(m.Alloc)
			metrics.BuckHashSys = Gauge(m.BuckHashSys)
			metrics.Frees = Gauge(m.Frees)
			metrics.GCCPUFraction = Gauge(m.GCCPUFraction)
			metrics.GCSys = Gauge(m.GCSys)
			metrics.HeapAlloc = Gauge(m.HeapAlloc)
			metrics.HeapIdle = Gauge(m.HeapIdle)
			metrics.HeapInuse = Gauge(m.HeapInuse)
			metrics.HeapObjects = Gauge(m.HeapObjects)
			metrics.HeapReleased = Gauge(m.HeapReleased)
			metrics.HeapSys = Gauge(m.HeapSys)
			metrics.LastGC = Gauge(m.LastGC)
			metrics.Lookups = Gauge(m.Lookups)
			metrics.MCacheInuse = Gauge(m.MCacheInuse)
			metrics.MCacheSys = Gauge(m.MCacheSys)
			metrics.MSpanInuse = Gauge(m.MSpanInuse)
			metrics.MSpanSys = Gauge(m.MSpanSys)
			metrics.Mallocs = Gauge(m.Mallocs)
			metrics.NextGC = Gauge(m.NextGC)
			metrics.NumForcedGC = Gauge(m.NumForcedGC)
			metrics.NumGC = Gauge(m.NumGC)
			metrics.OtherSys = Gauge(m.OtherSys)
			metrics.PauseTotalNs = Gauge(m.PauseTotalNs)
			metrics.StackInuse = Gauge(m.StackInuse)
			metrics.StackSys = Gauge(m.StackSys)
			metrics.Sys = Gauge(m.Sys)
			metrics.TotalAlloc = Gauge(m.TotalAlloc)

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

	switch (*a.Configuration).GetBool("USE_JSON") {
	case true:
		{
			//Mocked up
			URL := a.baseURL.
				JoinPath("update")
			AddGaugeDataJSON(URL, metrics.Alloc, "Alloc", &m)
			AddGaugeDataJSON(URL, metrics.Frees, "Frees", &m)
			AddGaugeDataJSON(URL, metrics.GCCPUFraction, "GCCPUFraction", &m)
			AddGaugeDataJSON(URL, metrics.GCSys, "GCSys", &m)
			AddGaugeDataJSON(URL, metrics.HeapAlloc, "HeapAlloc", &m)
			AddGaugeDataJSON(URL, metrics.HeapIdle, "HeapIdle", &m)
			AddGaugeDataJSON(URL, metrics.HeapInuse, "HeapInuse", &m)
			AddGaugeDataJSON(URL, metrics.HeapObjects, "HeapObjects", &m)
			AddGaugeDataJSON(URL, metrics.HeapReleased, "HeapReleased", &m)
			AddGaugeDataJSON(URL, metrics.HeapSys, "HeapSys", &m)
			AddGaugeDataJSON(URL, metrics.LastGC, "LastGC", &m)
			AddGaugeDataJSON(URL, metrics.Lookups, "Lookups", &m)
			AddGaugeDataJSON(URL, metrics.MCacheSys, "MCacheSys", &m)
			AddGaugeDataJSON(URL, metrics.MSpanInuse, "MSpanInuse", &m)
			AddGaugeDataJSON(URL, metrics.MSpanSys, "MSpanSys", &m)
			AddGaugeDataJSON(URL, metrics.Mallocs, "Mallocs", &m)
			AddGaugeDataJSON(URL, metrics.NextGC, "NextGC", &m)
			AddGaugeDataJSON(URL, metrics.NumForcedGC, "NumForcedGC", &m)
			AddGaugeDataJSON(URL, metrics.NumGC, "NumGC", &m)
			AddGaugeDataJSON(URL, metrics.OtherSys, "OtherSys", &m)
			AddGaugeDataJSON(URL, metrics.PauseTotalNs, "PauseTotalNs", &m)
			AddGaugeDataJSON(URL, metrics.StackInuse, "StackInuse", &m)
			AddGaugeDataJSON(URL, metrics.StackSys, "StackSys", &m)
			AddGaugeDataJSON(URL, metrics.Sys, "Sys", &m)
			AddGaugeDataJSON(URL, metrics.TotalAlloc, "TotalAlloc", &m)
			AddGaugeDataJSON(URL, metrics.RandomValue, "RandomValue", &m)
			AddCounterDataJSON(URL, metrics.PollCount, "PollCount", &m)
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
			AddCounterData(URL, metrics.PollCount, "PollCount", &m)

		}

	}

	return m
}
func (a Agent) Send(ctx context.Context, client *http.Client, metrics *Metrics) {

	ticker := time.NewTicker(time.Duration((*a.Configuration).GetInt("REPORT_INTERVAL")) * time.Second)
	defer ticker.Stop()

repeatAgain:
	select {
	case <-ticker.C:
		{
			metrics.PollCount++
			metrics.RandomValue = Gauge(rand.Int63())
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
