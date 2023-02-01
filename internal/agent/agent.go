package agent

import (
	"compress/flate"
	"context"
	"encoding/json"

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

func AddCounterData(common sendData, val Counter, name string, data *map[*sendData]bool) {
	URL := common.url.
		JoinPath("counter").
		JoinPath(name).
		JoinPath(strconv.FormatUint(uint64(val), 10)) //value float
	sd := sendData{
		url:  URL,
		keys: common.keys,
		body: bytes.NewBufferString(url.Values{}.Encode()), //need to transfer something
	}
	(*data)[&sd] = true

}
func AddGaugeData(common sendData, val Gauge, name string, data *map[*sendData]bool) {

	URL := common.url.
		JoinPath("gauge").
		JoinPath(name).
		JoinPath(strconv.FormatFloat(float64(val), 'E', -1, 64)) //value float

	sd := sendData{
		url:  URL,
		keys: common.keys,
		body: bytes.NewBufferString(url.Values{}.Encode()), //need to transer something
	}
	(*data)[&sd] = true

}

func AddGaugeDataJSON(common sendData, val Gauge, name string, data *map[*sendData]bool) {
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
		url:  common.url,
		keys: common.keys,
		body: bytes.NewBuffer(metricsBytes),
	}
	(*data)[&sd] = true

}
func AddCounterDataJSON(common sendData, val Counter, name string, data *map[*sendData]bool) {
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
		url:  common.url,
		keys: common.keys,
		body: bytes.NewBuffer(metricsBytes),
	}
	(*data)[&sd] = true

}

type HeaderKeys map[string]string
type sendData struct {
	url  *url.URL
	keys *HeaderKeys
	body *bytes.Buffer
}

func (data sendData) SendData(client *http.Client) error {
	retry := 10
	request, err := http.NewRequest(http.MethodPost, data.url.String(), data.body)
	if err != nil {
		log.Fatal(err)
	}
	if data.keys != nil {
		for k, v := range *data.keys {
			request.Header.Set(k, v)
		}
	}
	log.Printf("url from agent):%s", data.url.String())
again:
	response, err := client.Do(request)
	if err != nil {
		if retry != 0 {
			time.Sleep(time.Microsecond * 100)
			retry--
			log.Println("request error, retry")

			goto again
		}
		log.Fatal(err)
	}
	defer response.Body.Close()
	log.Println("response from server:" + response.Status)

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

func (a Agent) CompressData(data map[*sendData]bool) map[*sendData]bool {

	switch (*a.Configuration).Get("COMPRESS_TYPE") {
	case "deflate":
		{
			for k := range data {

				var b bytes.Buffer

				w, err := flate.NewWriter(&b, flate.BestCompression)
				if err != nil {
					log.Fatalf("failed init compress writer: %v", err)
				}
				_, err = w.Write(k.body.Bytes())
				if err != nil {
					log.Fatalf("failed write data to compress temporary buffer: %v", err)
				}

				err = w.Close()
				if err != nil {
					log.Fatalf("failed compress data: %v", err)
				}
				k.body = &b
			}
		}
	}

	return data
}
func (a Agent) prepareData(metrics *Metrics) map[*sendData]bool {
	m := make(map[*sendData]bool)
	keys := make(HeaderKeys)

	switch (*a.Configuration).Get("COMPRESS_TYPE") {
	case "deflate":
		keys["Accept-Encoding"] = "deflate"
		keys["Content-Encoding"] = "deflate"
		// default:
		// 	keys["Accept-Encoding"] = "gzip"
	}

	switch (*a.Configuration).GetBool("USE_JSON") {
	case true:
		{

			keys["Content-Type"] = "application/json; charset=utf-8"
			keys["Accept"] = "application/json; charset=utf-8"

			data := sendData{
				url:  a.baseURL.JoinPath("update"),
				keys: &keys,
			}

			AddGaugeDataJSON(data, metrics.Alloc, "Alloc", &m)
			AddGaugeDataJSON(data, metrics.Frees, "Frees", &m)
			AddGaugeDataJSON(data, metrics.GCCPUFraction, "GCCPUFraction", &m)
			AddGaugeDataJSON(data, metrics.GCSys, "GCSys", &m)
			AddGaugeDataJSON(data, metrics.HeapAlloc, "HeapAlloc", &m)
			AddGaugeDataJSON(data, metrics.HeapIdle, "HeapIdle", &m)
			AddGaugeDataJSON(data, metrics.HeapInuse, "HeapInuse", &m)
			AddGaugeDataJSON(data, metrics.HeapObjects, "HeapObjects", &m)
			AddGaugeDataJSON(data, metrics.HeapReleased, "HeapReleased", &m)
			AddGaugeDataJSON(data, metrics.HeapSys, "HeapSys", &m)
			AddGaugeDataJSON(data, metrics.LastGC, "LastGC", &m)
			AddGaugeDataJSON(data, metrics.Lookups, "Lookups", &m)
			AddGaugeDataJSON(data, metrics.MCacheSys, "MCacheSys", &m)
			AddGaugeDataJSON(data, metrics.MSpanInuse, "MSpanInuse", &m)
			AddGaugeDataJSON(data, metrics.MSpanSys, "MSpanSys", &m)
			AddGaugeDataJSON(data, metrics.Mallocs, "Mallocs", &m)
			AddGaugeDataJSON(data, metrics.NextGC, "NextGC", &m)
			AddGaugeDataJSON(data, metrics.NumForcedGC, "NumForcedGC", &m)
			AddGaugeDataJSON(data, metrics.NumGC, "NumGC", &m)
			AddGaugeDataJSON(data, metrics.OtherSys, "OtherSys", &m)
			AddGaugeDataJSON(data, metrics.PauseTotalNs, "PauseTotalNs", &m)
			AddGaugeDataJSON(data, metrics.StackInuse, "StackInuse", &m)
			AddGaugeDataJSON(data, metrics.StackSys, "StackSys", &m)
			AddGaugeDataJSON(data, metrics.Sys, "Sys", &m)
			AddGaugeDataJSON(data, metrics.TotalAlloc, "TotalAlloc", &m)
			AddGaugeDataJSON(data, metrics.RandomValue, "RandomValue", &m)
			AddCounterDataJSON(data, metrics.PollCount, "PollCount", &m)
		}
	default:
		{

			keys["Content-Type"] = "application/json; charset=utf-8"
			keys["Accept"] = "text/html; charset=utf-8"

			data := sendData{
				url:  a.baseURL.JoinPath("update"),
				keys: &keys,
			}

			AddGaugeData(data, metrics.Alloc, "Alloc", &m)
			AddGaugeData(data, metrics.Frees, "Frees", &m)
			AddGaugeData(data, metrics.GCCPUFraction, "GCCPUFraction", &m)
			AddGaugeData(data, metrics.GCSys, "GCSys", &m)
			AddGaugeData(data, metrics.HeapAlloc, "HeapAlloc", &m)
			AddGaugeData(data, metrics.HeapIdle, "HeapIdle", &m)
			AddGaugeData(data, metrics.HeapInuse, "HeapInuse", &m)
			AddGaugeData(data, metrics.HeapObjects, "HeapObjects", &m)
			AddGaugeData(data, metrics.HeapReleased, "HeapReleased", &m)
			AddGaugeData(data, metrics.HeapSys, "HeapSys", &m)
			AddGaugeData(data, metrics.LastGC, "LastGC", &m)
			AddGaugeData(data, metrics.Lookups, "Lookups", &m)
			AddGaugeData(data, metrics.MCacheSys, "MCacheSys", &m)
			AddGaugeData(data, metrics.MSpanInuse, "MSpanInuse", &m)
			AddGaugeData(data, metrics.MSpanSys, "MSpanSys", &m)
			AddGaugeData(data, metrics.Mallocs, "Mallocs", &m)
			AddGaugeData(data, metrics.NextGC, "NextGC", &m)
			AddGaugeData(data, metrics.NumForcedGC, "NumForcedGC", &m)
			AddGaugeData(data, metrics.NumGC, "NumGC", &m)
			AddGaugeData(data, metrics.OtherSys, "OtherSys", &m)
			AddGaugeData(data, metrics.PauseTotalNs, "PauseTotalNs", &m)
			AddGaugeData(data, metrics.StackInuse, "StackInuse", &m)
			AddGaugeData(data, metrics.StackSys, "StackSys", &m)
			AddGaugeData(data, metrics.Sys, "Sys", &m)
			AddGaugeData(data, metrics.TotalAlloc, "TotalAlloc", &m)
			AddGaugeData(data, metrics.RandomValue, "RandomValue", &m)
			AddCounterData(data, metrics.PollCount, "PollCount", &m)
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
			dataPackage := a.CompressData(a.prepareData(metrics))

			for key := range dataPackage {
				err := key.SendData(client)
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
