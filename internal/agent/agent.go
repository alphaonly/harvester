package agent

import (
	"context"
	"os"

	"log"
	"net/http"
	"reflect"
	"runtime"

	"bytes"
	"strconv"
	"time"

	"math/rand"
	"net/url"
)

const (
	pollInterval   = 2
	reportInterval = 3 //10

	serverHost = "127.0.0.1"
	serverPort = ":8080"
)

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

func UpdateMemStatsMetrics(ctx context.Context, metrics *Metrics) {

	var m runtime.MemStats

	ticker := time.NewTicker(pollInterval * time.Second)

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

			break repeatAgain
		}

	case <-ctx.Done():
		{
			go log.Println("Metrics reading cancelled by context")
		}
	}

}

func Run(ctx context.Context) {

	ctxMetrics, cancel := context.WithCancel(ctx)
	defer cancel()

	metrics := Metrics{}
	elements := reflect.ValueOf(&metrics).Elem()

	var elementValue string

	baseURL := url.URL{

		Scheme: "http",
		Host:   serverHost + serverPort,
	}

	go UpdateMemStatsMetrics(ctxMetrics, &metrics)

	ticker := time.NewTicker(reportInterval * time.Second)
	defer ticker.Stop()
	defer cancel()

	data := url.Values{}
	client := &http.Client{}

	for {

		metrics.PollCount++
		metrics.RandomValue = gauge(rand.Int63())

		<-ticker.C

		for i := 0; i < elements.NumField(); i++ {

			elementName := elements.Type().Field(i).Name
			elementType := elements.Type().Field(i).Type.String()

			switch elementName {
			case "PollCount":
				elementValue = strconv.FormatUint(uint64(elements.Field(i).Interface().(counter)), 10)
			default:
				elementValue = strconv.FormatFloat(float64(elements.Field(i).Interface().(gauge)), 'E', -1, 64)
			}

			url := baseURL.
				JoinPath("update").
				JoinPath(elementType).
				JoinPath(elementName).
				JoinPath(elementValue)

			go log.Println("url from agent:" + url.String())
			request, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBufferString(data.Encode()))
			if err != nil {
				log.Fatal(err)

			}

			request.Header.Set("Content-Type", "text/plain; charset=utf-8")
			request.Header.Add("Accept", "text/plain; charset=utf-8")

			response, err := client.Do(request)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			go log.Println(response.Status)
			response.Body.Close()

		}

	}
}
