package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	i "github.com/alphaonly/harvester/internal/server/interfaces"
	c "github.com/alphaonly/harvester/internal/server/interfaces/MetricValue/implementations/CounterValue"
	g "github.com/alphaonly/harvester/internal/server/interfaces/MetricValue/implementations/GaugeValue"
	"github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Storage *interfaces.Storage
}

func (h *Handlers) SetDataServer(storage *interfaces.Storage) {
	h.Storage = storage

}
func (h *Handlers) HandleGetMetricFieldList(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ms, err := (*h.Storage).GetAllMetrics(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("<h1><ul>"))
	for key := range *ms {
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

	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")

	if metricName == "" {
		http.Error(w, "is not parsed, empty metric name ", http.StatusNotFound)
	}
	if metricType == "" {
		http.Error(w, metricType+"is not recognized type", http.StatusNotImplemented)
		return

	}

	if h.Storage == nil {
		http.Error(w, "dataServer is not initialized", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	metricsValue, err := (*h.Storage).GetMetric(ctx, metricName)
	if err != nil {
		http.Error(w, "404 - not found", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricValue := (*metricsValue).GetString()

	_, err = w.Write([]byte(metricValue))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "plain/text; charset=utf-8")

}

func (h *Handlers) HandleGetMetricValueJSON(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	requestByteData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unrecognized json request ", http.StatusBadRequest)
		return
	}
	var requestMetricsJson interfaces.MetricsJSON

	err = json.Unmarshal(requestByteData, &requestMetricsJson)
	if err != nil {
		http.Error(w, "Error json-marshal request data", http.StatusBadRequest)
		return
	}
	if !(requestMetricsJson.Delta == nil && requestMetricsJson.Value == nil) {
		http.Error(w, "not empty values in json-marshal request data", http.StatusBadRequest)
		return
	}

	metricType := requestMetricsJson.MType
	metricName := requestMetricsJson.ID

	if metricName == "" {
		http.Error(w, "is not parsed, empty metric name ", http.StatusNotFound)
		return
	}
	if metricType == "" {
		http.Error(w, metricType+"is not recognized type", http.StatusNotImplemented)
		return
	}

	if h.Storage == nil {
		http.Error(w, "dataServer is not initialized", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	metricsValue, err := (*h.Storage).GetMetric(ctx, metricName)
	if err != nil {
		http.Error(w, "404 - not found", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var responseMetricsJSon = requestMetricsJson

	switch requestMetricsJson.MType {
	case "agent.gauge":
		{
			v := (*metricsValue).GetInternalValue().(float64)
			responseMetricsJSon.Value = &v
		}
	case "agent.counter":
		{
			v := (*metricsValue).GetInternalValue().(int64)
			responseMetricsJSon.Delta = &v
		}
	default:
		{
			http.Error(w, "unknown metric type", http.StatusInternalServerError)
			return
		}
	}
	responseByteData, err := json.Marshal(responseMetricsJSon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	_, err = w.Write(responseByteData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

}

func (h *Handlers) HandlePostMetric(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")
	metricValue := chi.URLParam(r, "VALUE")

	log.Println("metricType :" + metricType +
		" metricName :" + metricName +
		" metricValue :" + metricValue)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if h.Storage == nil {
		http.Error(w, "data storage not initiated", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodPost:
		{

			if metricName == "" {
				http.Error(w, "not parsed, empty metric name!"+metricName, http.StatusNotFound)
				return
			}

			if metricValue == "" {
				http.Error(w, "not parsed, empty metric value", http.StatusBadRequest)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			switch metricType {
			case "gauge":
				{
					float64Value, err := strconv.ParseFloat(metricValue, 64)
					if err != nil {
						http.Error(w, "value:"+metricValue+" not parsed, value cast error", http.StatusBadRequest)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					var m i.MetricValue = *(g.GaugeValue{}.NewFloat(float64Value))

					err2 := (*h.Storage).SaveMetric(r.Context(), metricName, m)
					if err2 != nil {
						http.Error(w, "internal value add error", http.StatusInternalServerError)
						return
					}

				}
			case "counter":
				{
					intValue, err := strconv.ParseInt(metricValue, 10, 64)
					if err != nil {
						http.Error(w, "value: "+metricValue+" not parsed", http.StatusBadRequest)
						return
					}

					prevMetricValue, err := (*h.Storage).GetMetric(r.Context(), metricName)

					if err != nil || prevMetricValue == nil {

						*prevMetricValue = c.CounterValue{}

					}
					sum := c.CounterValue{}.NewInt(intValue).AddValue(*prevMetricValue)

					(*h.Storage).SaveMetric(r.Context(), metricName, &sum)

				}
			default:
				http.Error(w, metricType+" not recognized type", http.StatusNotImplemented)
				return
			}

		}
	default:
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

}
func (h *Handlers) HandlePostMetricJSON(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	dataServer := h.Storage

	if dataServer == nil {
		http.Error(w, "storage not initiated", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodPost:
		{

			byteData, err := io.ReadAll(r.Body)
			if err != nil {

				http.Error(w, "unrecognized request body", http.StatusBadRequest)
				return
			}
			var metricsJson interfaces.MetricsJSON
			err = json.Unmarshal(byteData, &metricsJson)

			if err != nil {
				http.Error(w, "Unrecognized json", http.StatusBadRequest)
				return
			}

			if metricsJson.ID == "" {
				http.Error(w, "not parsed, empty metric name!"+metricsJson.ID, http.StatusNotFound)
				return
			}

			if metricsJson.Delta == nil && metricsJson.Value == nil {
				http.Error(w, "not parsed, empty metric value", http.StatusBadRequest)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			//logic
			switch metricsJson.MType {
			case "gauge":
				{
					var gValue i.MetricValue = i.GaugeValue{}.NewFloat(*metricsJson.Value)
					err := (*h.Storage).SaveMetric(r.Context(), metricsJson.ID, &gValue)
					if err != nil {
						http.Error(w, "internal value add error", http.StatusInternalServerError)
						return
					}

				}
			case "counter":
				{

					prevMetricValue, err := (*h.Storage).GetMetric(r.Context(), metricsJson.ID)

					if err != nil || prevMetricValue == nil {

						*prevMetricValue = i.CounterValue{}

					}
					sum := i.CounterValue{}.NewInt(*metricsJson.Delta).AddValue(*prevMetricValue)

					(*h.Storage).SaveMetric(r.Context(), metricsJson.ID, &sum)

				}
			default:
				http.Error(w, metricsJson.MType+" not recognized type", http.StatusNotImplemented)
				return
			}

			byteData, err2 := json.Marshal(metricsJson)
			if err2 != nil || byteData == nil {
				http.Error(w, " json response forming error", http.StatusInternalServerError)
				return
			}
			w.Write(byteData)
		}
	default:
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

}

func (h *Handlers) HandlePostErrorPattern(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "Unknown request", http.StatusNotFound)

}
func (h *Handlers) HandlePostErrorPatternNoName(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "Unknown request", http.StatusNotFound)

}

func (h *Handlers) NewRouter() chi.Router {

	r := chi.NewRouter()
	//
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.HandleGetMetricFieldList)
		r.Get("/value", h.HandleGetMetricValueJSON)
		r.Get("/value/{TYPE}/{NAME}", h.HandleGetMetricValue)
		r.Post("/update", h.HandlePostMetricJSON)
		r.Post("/updateViaUrl/{TYPE}/{NAME}/{VALUE}", h.HandlePostMetric)
		r.Post("/update/{TYPE}/{NAME}/", h.HandlePostErrorPattern)
		r.Post("/update/{TYPE}/", h.HandlePostErrorPatternNoName)

	})

	return r
}
