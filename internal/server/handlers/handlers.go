package handlers

import (
	"compress/flate"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/alphaonly/harvester/internal/server/compression"
	J "github.com/alphaonly/harvester/internal/server/metricsJSON"
	M "github.com/alphaonly/harvester/internal/server/metricvalue"
	C "github.com/alphaonly/harvester/internal/server/metricvalue/MetricValue/implementations/CounterValue"
	G "github.com/alphaonly/harvester/internal/server/metricvalue/MetricValue/implementations/Gaugevalue"
	S "github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	MemKeeper *S.Storage
}

func New(storage *S.Storage) *Handlers {
	return &Handlers{MemKeeper: storage}
}

func (h *Handlers) HandleGetMetricFieldList(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ms, err := (*h.MemKeeper).GetAllMetrics(r.Context())

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
		return
	}
	if metricType == "" {
		http.Error(w, metricType+"is not recognized type", http.StatusNotImplemented)
		return

	}

	if h.MemKeeper == nil {
		http.Error(w, "dataServer is not initialized", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	metricsValue, err := (*h.MemKeeper).GetMetric(ctx, metricName)
	if err != nil {
		http.Error(w, "404 - not found", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricValue := (*metricsValue).GetString()

	_, err = w.Write([]byte(metricValue))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	var requestMetricsJSON J.MetricsJSON

	err = json.Unmarshal(requestByteData, &requestMetricsJSON)
	if err != nil {
		http.Error(w, "Error json-marshal request data", http.StatusBadRequest)
		return
	}
	if !(requestMetricsJSON.Delta == nil && requestMetricsJSON.Value == nil) {
		http.Error(w, "not empty values in json-marshal request data", http.StatusBadRequest)
		return
	}

	metricType := requestMetricsJSON.MType
	metricName := requestMetricsJSON.ID

	if metricName == "" {
		http.Error(w, "is not parsed, empty metric name ", http.StatusNotFound)
		return
	}
	if metricType == "" {
		http.Error(w, metricType+"is not recognized type", http.StatusNotImplemented)
		return
	}

	if h.MemKeeper == nil {
		http.Error(w, "dataServer is not initialized", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	metricsValue, err := (*h.MemKeeper).GetMetric(ctx, metricName)
	if err != nil {
		http.Error(w, "404 - not found", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var responseMetricsJSON = requestMetricsJSON

	switch requestMetricsJSON.MType {
	case "gauge":
		{
			v := (*metricsValue).GetInternalValue().(float64)
			responseMetricsJSON.Value = &v
		}
	case "counter":
		{
			v := (*metricsValue).GetInternalValue().(int64)
			responseMetricsJSON.Delta = &v
		}
	default:
		{
			http.Error(w, "unknown metric type", http.StatusInternalServerError)
			return
		}
	}
	responseByteData, err := json.Marshal(responseMetricsJSON)
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

	if h.MemKeeper == nil {
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
					var m M.MetricValue = G.NewFloat(float64Value)

					err = (*h.MemKeeper).SaveMetric(r.Context(), metricName, &m)
					if err != nil {
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

					prevMetricValue, err := (*h.MemKeeper).GetMetric(r.Context(), metricName)

					if err != nil || prevMetricValue == nil {

						prevMetricValue = C.NewCounterValue()

					}
					sum := C.NewInt(intValue).AddValue(*prevMetricValue)

					(*h.MemKeeper).SaveMetric(r.Context(), metricName, &sum)

				}
			default:
				http.Error(w, metricType+" not recognized type", http.StatusNotImplemented)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	default:
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

}

func (h *Handlers) HandlePostMetricJSON(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if h.MemKeeper == nil {
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
				// d := CM.Deflator{
				// 	Level:           flate.BestSpeed,
				// 	ContentEncoding: r.Header.Get("Content-Encoding"),
				// }

				// byteData, err = d.Decompress(byteData)
				// if err != nil {
				// 	http.Error(w, err.Error(), http.StatusInternalServerError)
				// 	return
				// }

				var metricsJSON J.MetricsJSON
				err = json.Unmarshal(byteData, &metricsJSON)

				if err != nil {
					http.Error(w, "Unrecognized json", http.StatusBadRequest)
					return
				}

				if metricsJSON.ID == "" {
					http.Error(w, "not parsed, empty metric name!"+metricsJSON.ID, http.StatusNotFound)
					log.Println("Error not parsed, empty metric name: 404")
					return
				}

				if metricsJSON.Delta == nil && metricsJSON.Value == nil {
					http.Error(w, "not parsed, empty metric value", http.StatusBadRequest)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				//logic
				switch metricsJSON.MType {
				case "gauge":
					{
						var m M.MetricValue = G.NewFloat(*metricsJSON.Value)
						err := (*h.MemKeeper).SaveMetric(r.Context(), metricsJSON.ID, &m)
						if err != nil {
							http.Error(w, "internal value add error", http.StatusInternalServerError)
							return
						}
					}
				case "counter":
					{
						var prevMetricValue M.MetricValue = &C.CounterValue{}
						prevMetricValuePtr, err := (*h.MemKeeper).GetMetric(r.Context(), metricsJSON.ID)

						if err != nil || prevMetricValuePtr == nil {
							prevMetricValuePtr = &prevMetricValue
						}
						sum := C.NewInt(*metricsJSON.Delta).AddValue(*prevMetricValuePtr)
						(*h.MemKeeper).SaveMetric(r.Context(), metricsJSON.ID, &sum)
					}
				default:
					http.Error(w, metricsJSON.MType+" not recognized type", http.StatusNotImplemented)
					return
				}

				byteData, err2 := json.Marshal(metricsJSON)
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

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *Handlers) HandlePostErrorPattern(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unknown request", http.StatusNotFound)
	r.Body.Close()
}
func (h *Handlers) HandlePostErrorPatternNoName(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "Unknown request", http.StatusNotFound)
	r.Body.Close()
}

func (h *Handlers) NewRouter() chi.Router {

	d := compression.Deflator{
		Level: flate.BestSpeed,
		// ContentEncoding:r.Header.Get("Content-Encoding"),
		ContentEncoding: "deflate",
	}

	var postJsonCompressedScenario http.HandlerFunc = 
	d.DeCompressionHandler(h.HandlePostMetricJSON(nil))

	r := chi.NewRouter()
	//
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.HandleGetMetricFieldList)
		r.Get("/value", h.HandleGetMetricValueJSON)
		r.Get("/value/{TYPE}/{NAME}", h.HandleGetMetricValue)
		r.Post("/update", postJsonCompressedScenario)
		r.Post("/update/{TYPE}/{NAME}/{VALUE}", h.HandlePostMetric)
		r.Post("/update/{TYPE}/{NAME}/", h.HandlePostErrorPattern)
		r.Post("/update/{TYPE}/", h.HandlePostErrorPatternNoName)

	})

	return r
}
