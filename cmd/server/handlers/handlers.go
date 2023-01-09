package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alphaonly/harvester/cmd/server/storage"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	dataServer *storage.DataServer
}

func (h *Handlers) SetDataServer(dataServer *storage.DataServer) {
	h.dataServer = dataServer

}
func (h *Handlers) HandleGetMetricFieldList(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// val := reflect.ValueOf(&storage.Metrics{}).Elem()

	ms, err := h.dataServer.GetAllMetricsNames(context.Background())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("<h1><ul>"))
	for key := range ms {
		w.Write([]byte(" <li>" + key + "</li>"))
	}

	w.Write([]byte("</ul></h1>"))

	w.WriteHeader(http.StatusOK)
}
func (h *Handlers) HandleGetMetricValue(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	//metricValue := chi.URLParam(r, "value")
	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")

	if metricName == "" {
		http.Error(w, "is not parsed, empty metric name ", http.StatusNotFound)
	}
	if metricType == "" {
		http.Error(w, metricType+"is not recognized type", http.StatusNotImplemented)
		return

	}

	dataServer := h.dataServer

	if &dataServer == nil {
		http.Error(w, "dataServer is not initialized", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	metricsValue, err := dataServer.GetCurrentMetricMap(ctx, metricName)
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricValue := metricsValue.GetString()

	_, err = w.Write([]byte(metricValue))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "plain/text; charset=utf-8")

}
func (h *Handlers) HandlePostMetric(w http.ResponseWriter, r *http.Request) {

	var (
		gaugeValue   storage.Gauge
		counterValue storage.Counter
		parts        [5]string
	)

	fmt.Println(r)

	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")
	metricValue := chi.URLParam(r, "VALUE")

	fmt.Println("metricType :" + metricType)
	fmt.Println("metricName :" + metricName)
	fmt.Println("metricValue :" + metricValue)

	//w.Write([]byte("type:" + metricType))
	//w.Write([]byte("name:" + metricName))
	//w.Write([]byte("value:" + metricValue))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	dataServer := h.dataServer

	if &dataServer == nil {
		http.Error(w, "dataserver not detected", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	switch r.Method {
	case http.MethodPost:
		{

			parts[0] = ""
			parts[1] = "update"
			parts[2] = metricType
			parts[3] = metricName
			parts[4] = metricValue

			//parts = strings.SplitN(r.URL.String(), "/", 5)

			//fmt.Println(parts[1])
			if parts[1] != "update" {
				http.Error(w, "not parsed, "+parts[1]+" bad namespace ", http.StatusNotImplemented)
				return
			}

			if parts[3] == "" {
				http.Error(w, "not parsed, empty metric name!"+parts[4], http.StatusNotFound)
				fmt.Println(parts)
				return
			}

			if parts[4] == "" {
				http.Error(w, "not parsed, empty metric value", http.StatusBadRequest)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			gValue := storage.GaugeValue{}
			cValue := storage.CounterValue{}

			switch parts[2] {
			case "main.gauge", "gauge":
				{
					float64Value, err := strconv.ParseFloat(parts[4], 64)
					if err != nil {
						http.Error(w, "value:"+parts[4]+" not parsed, value cast error", http.StatusBadRequest)
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					gaugeValue = storage.Gauge(float64Value)

					gValue.SetValue(gaugeValue)
					err2 := dataServer.SaveMetricToMap(ctx, metricName, &gValue)
					if err2 != nil {
						http.Error(w, "internal value add error", http.StatusInternalServerError)
						return
					}

				}
			case "main.counter", "counter":
				{
					intValue, err := strconv.ParseInt(parts[4], 10, 64)
					if err != nil {
						http.Error(w, "value: "+parts[4]+" not parsed", http.StatusBadRequest)
						return
					}

					prevMetricValue, err := dataServer.GetCurrentMetricMap(ctx, metricName)

					if err != nil || prevMetricValue == nil {

						prevMetricValue = &storage.CounterValue{}
						// prevMetricValue.SetValue(CounterValue(0))
					}

					counterValue = storage.Counter(intValue)

					cValue.SetValue(counterValue)

					dataServer.SaveMetricToMap(ctx, metricName, cValue.AddValue(prevMetricValue))

				}
			default:
				http.Error(w, parts[2]+" not recognized type", http.StatusNotImplemented)
				return
			}

		}
	default:
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

}
func (h *Handlers) HandlePostErrorPattern(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "Unknown request", http.StatusNotFound)
	return

}
func (h *Handlers) HandlePostErrorPatternNoName(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "Unknown request", http.StatusNotFound)
	return

}

func NewRouter(ds *storage.DataServer) chi.Router {

	r := chi.NewRouter()
	h := Handlers{}
	h.SetDataServer(ds)
	//
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.HandleGetMetricFieldList)
		r.Get("/value/{TYPE}/{NAME}", h.HandleGetMetricValue)
		r.Post("/update/{TYPE}/{NAME}/{VALUE}", h.HandlePostMetric)
		r.Post("/update/{TYPE}/{NAME}/", h.HandlePostErrorPattern)
		r.Post("/update/{TYPE}/", h.HandlePostErrorPatternNoName)
	})

	return r
}
