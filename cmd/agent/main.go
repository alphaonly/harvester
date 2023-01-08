package main

import (
	"context"
	"fmt"
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

func updateMemStatsMetrics(ctx context.Context, metrics *Metrics) {

	var m runtime.MemStats

	isCancelled := false
	ticker := time.NewTicker(pollInterval * time.Second)

	defer ticker.Stop()

	for !isCancelled {
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
			}

		case <-ctx.Done():
			{
				fmt.Println("Metrics reading cancelled by context")
				isCancelled = true

			}
		}

	}
}

func main() {

	// client := http.Client{}
	ctx := context.Background()
	ctxMetrics, cancel := context.WithCancel(ctx)

	metrics := Metrics{}

	elements := reflect.ValueOf(&metrics).Elem()

	var urlPrefix = "http://" + serverHost + serverPort + "/update/"
	var urlStr string
	var elementValue string
	// var packageCounter int = 0
	go updateMemStatsMetrics(ctxMetrics, &metrics)

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

				//fmt.Println(
				//	elements.Field(i).Interface(), " ",
				//	elements.Field(i).Interface().(gauge), " ",
				//	float64(elements.Field(i).Interface().(gauge)), " ",
				//	strconv.FormatFloat(float64(elements.Field(i).Interface().(gauge)), 'f', -1, 64),
				//)
				elementValue = strconv.FormatFloat(float64(elements.Field(i).Interface().(gauge)), 'E', -1, 64)

			}

			//elementType = "main.gauge"
			//elementName = "counter"
			//elementValue = ""
			urlStr = urlPrefix +
				elementType + "/" +
				elementName + "/" +
				elementValue

			urlStr = "http://127.0.0.1:8080:/update/main.gauge/counter/"
			fmt.Println("url from agent:" + urlStr)
			request, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBufferString(data.Encode()))
			if err != nil {
				log.Fatal(err)

			}
			//fmt.Println(urlStr)

			request.Header.Set("Content-Type", "text/plain; charset=utf-8")
			request.Header.Add("Accept", "text/plain; charset=utf-8")

			response, err := client.Do(request)
			if err != nil {

				os.Exit(1)
			}

			fmt.Println(response.Status)

		}

	}

}
