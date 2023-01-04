package handlers

import (
	"context"
	"fmt"
	"github.com/alphaonly/harvester/cmd/server/storage"
	"net/http"
	"strconv"
	"strings"
)

//type gauge float64
//type counter int64

var (
	//metrics      storage.Metrics = storage.Metrics{}
	//metricsMap                   = make(storage.MetricsMap)
	metrics = storage.Metrics{}
)

type Handlers struct {
	dataServer *storage.DataServer
}

func (h *Handlers) SetDataServer(dataServer *storage.DataServer) {
	h.dataServer = dataServer

}

func (h *Handlers) HandleMetric(w http.ResponseWriter, r *http.Request) {

	var (
		gaugeValue   storage.Gauge
		counterValue storage.Counter
	)

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	dataServer := h.dataServer

	if &dataServer == nil {
		http.Error(w, "dataserver not detected", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	switch r.Method {
	case http.MethodPost:
		{
			parts := strings.SplitN(r.URL.String(), "/", 5)
			//fmt.Println(parts)

			//fmt.Println(parts[1])
			if parts[1] != "update" {
				http.Error(w, "not parsed, "+parts[1]+" bad namespace ", http.StatusNotImplemented)
				return
			}

			if parts[3] == "" {
				http.Error(w, "not parsed, empty metric name ", http.StatusNotFound)
				return
			} else {

				if parts[4] == "" {
					http.Error(w, "not parsed, empty metric value", http.StatusBadRequest)
					return
				}

				switch parts[2] {
				case "main.gauge", "gauge":
					{
						float64Value, err := strconv.ParseFloat(parts[4], 64)
						if err != nil {
							http.Error(w, "value:"+parts[4]+" not parsed, value cast error", http.StatusBadRequest)
							return
						} else {

							//metricsMap[parts[2]] = gauge(float64Value)
							//reflect.ValueOf().Field(i).SetFloat( float64Value )
							gaugeValue = storage.Gauge(float64Value)

							switch parts[3] {
							case "Alloc":
								{
									fmt.Println("Начало последовательности данных")
									//Начало последовательности данных
									//Очищаю
									metrics = storage.Metrics{}

									metrics.Alloc = gaugeValue
								}
							case "BuckHashSys":
								metrics.BuckHashSys = gaugeValue
							case "Frees":
								metrics.Frees = gaugeValue
							case "GCCPUFraction":
								metrics.GCCPUFraction = gaugeValue
							case "GCSys":
								metrics.GCCPUFraction = gaugeValue
							case "HeapAlloc":
								metrics.GCCPUFraction = gaugeValue
							case "HeapIdle":
								metrics.HeapIdle = gaugeValue
							case "HeapInuse":
								metrics.HeapInuse = gaugeValue
							case "HeapObjects":
								metrics.HeapObjects = gaugeValue
							case "HeapReleased":
								metrics.HeapReleased = gaugeValue
							case "HeapSys":
								metrics.HeapSys = gaugeValue
							case "LastGC":
								metrics.LastGC = gaugeValue
							case "Lookups":
								metrics.Lookups = gaugeValue
							case "MCacheInuse":
								metrics.MCacheInuse = gaugeValue
							case "MCacheSys":
								metrics.MCacheSys = gaugeValue
							case "MSpanInuse":
								metrics.MSpanInuse = gaugeValue
							case "MSpanSys":
								metrics.MSpanSys = gaugeValue
							case "Mallocs":
								metrics.Mallocs = gaugeValue
							case "NextGC":
								metrics.NextGC = gaugeValue
							case "NumForcedGC":
								metrics.NumForcedGC = gaugeValue
							case "NumGC":
								metrics.NumGC = gaugeValue
							case "OtherSys":
								metrics.OtherSys = gaugeValue
							case "PauseTotalNs":
								metrics.PauseTotalNs = gaugeValue
							case "StackInuse":
								metrics.StackInuse = gaugeValue
							case "StackSys":
								metrics.StackSys = gaugeValue
							case "Sys":
								metrics.Sys = gaugeValue
							case "TotalAlloc":
								metrics.TotalAlloc = gaugeValue
							case "RandomValue":
								metrics.RandomValue = gaugeValue
							default:
								{
									http.Error(w, "gauge:"+parts[3]+"unknown metric ", http.StatusOK)
									return
								}
							}
						}
					}
				case "main.counter", "counter":
					{
						intValue, err := strconv.ParseInt(parts[4], 10, 64)
						if err != nil {
							http.Error(w, "value:"+parts[4]+" not parsed", http.StatusBadRequest)
							return
						} else {

							counterValue = storage.Counter(intValue)
							//metricsMap[parts[2]] = counter(intValue)

							switch parts[3] {
							case "PollCount":
								{
									metrics.PollCount = counterValue
									//Конец последовательности, сохраняю
									fmt.Println("Конец последовательности, сохраняю")
									err := dataServer.SaveMetric(ctx, metrics)
									if err != nil {
										http.Error(w, "safe operation error", http.StatusInternalServerError)
										return
									}

								}
							default:
								http.Error(w, "counter:"+parts[3]+"unknown metric ", http.StatusOK)
								return
							}
						}
					}
				default:
					http.Error(w, "not recognized type ", http.StatusBadRequest)
					return
				}
			}
		}
	default:
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

}
