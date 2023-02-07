package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/alphaonly/harvester/internal/schema"
	mVal "github.com/alphaonly/harvester/internal/server/metricvalueInt"
	stor "github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	MemKeeper stor.Storage
}

func New(storage stor.Storage) *Handlers {
	return &Handlers{MemKeeper: storage}
}

func (h *Handlers) HandleGetMetricFieldList(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	ms, err := h.MemKeeper.GetAllMetrics(r.Context())

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
	metricsValue, err := h.MemKeeper.GetMetric(ctx, metricName)
	if err != nil {
		http.Error(w, "404 - not found", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricValue := metricsValue.GetString()

	_, err = w.Write([]byte(metricValue))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "plain/text")

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

	var requestMetricsJSON schema.MetricsJSON

	err = json.Unmarshal(requestByteData, &requestMetricsJSON)
	if err != nil {
		http.Error(w, "Error json-marshal request data", http.StatusBadRequest)
		return
	}
	if !(requestMetricsJSON.Delta == nil && requestMetricsJSON.Value == nil) {
		http.Error(w, "not empty values in json-marshal request data", http.StatusBadRequest)
		return
	}

	if requestMetricsJSON.ID == "" {
		http.Error(w, "is not parsed, empty metric name ", http.StatusNotFound)
		return
	}
	if requestMetricsJSON.MType == "" {
		http.Error(w, requestMetricsJSON.MType+"is not recognized type", http.StatusNotImplemented)
		return
	}

	if h.MemKeeper == nil {
		http.Error(w, "dataServer is not initialized", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	metricsValue, err := h.MemKeeper.GetMetric(ctx, requestMetricsJSON.ID)
	if err != nil {
		http.Error(w, "404 - not found", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var responseMetricsJSON = requestMetricsJSON

	switch requestMetricsJSON.MType {
	case "gauge":
		{
			v := metricsValue.GetInternalValue().(float64)
			responseMetricsJSON.Value = &v
		}
	case "counter":
		{
			v := metricsValue.GetInternalValue().(int64)
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

	w.Header().Set("Content-Type", "application/json")

}

func (h *Handlers) HandlePostMetric(w http.ResponseWriter, r *http.Request) {
	log.Println("HandlePostMetric invoked")

	metricType := chi.URLParam(r, "TYPE")
	metricName := chi.URLParam(r, "NAME")
	metricValue := chi.URLParam(r, "VALUE")

	log.Println("server:received data via URL: type :" + metricType +
		" name :" + metricName +
		" value :" + metricValue)

	w.Header().Set("Content-Type", "text/plain")

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

					var m mVal.MetricValue = mVal.NewFloat(float64Value)

					err = h.MemKeeper.SaveMetric(r.Context(), metricName, &m)
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
					prevMetricValue, err := h.MemKeeper.GetMetric(r.Context(), metricName)
					if err != nil || prevMetricValue == nil {
						prevMetricValue = mVal.NewCounterValue()
					}
					sum := mVal.NewInt(intValue).AddValue(prevMetricValue)
					err = h.MemKeeper.SaveMetric(r.Context(), metricName, &sum)
					if err != nil {
						http.Error(w, "value: "+metricValue+" not saved in memStorage", http.StatusInternalServerError)
						return

					}

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
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostMetricJSON invoked")
		//validation
		if h.MemKeeper == nil {
			http.Error(w, "storage not initiated", http.StatusInternalServerError)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
			return
		}

		byteData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unrecognized request body:"+err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Server:json received:" + string(byteData))
		var mj schema.MetricsJSON
		err = json.Unmarshal(byteData, &mj)

		if err != nil {
			http.Error(w, "unmarshal error:", http.StatusBadRequest)
			log.Println("unmarshal error:" + err.Error())
			return
		}

		if mj.ID == "" {
			http.Error(w, "not parsed, empty metric name!"+mj.ID, http.StatusNotFound)
			log.Println("Error not parsed, empty metric name: 404")
			return
		}

		//запрос пост в базу от агента
		switch mj.MType {
		case "gauge":
			{
				if mj.Value != nil {
					mjVal := *mj.Value
					//пишем если есть значение
					mv := mVal.MetricValue(mVal.NewFloat(mjVal))
					err := h.MemKeeper.SaveMetric(r.Context(), mj.ID, &mv)
					if err != nil {
						http.Error(w, "internal value add error", http.StatusInternalServerError)
						return
					}
				}
				//читаем  для ответа
				var f float64 = 0
				gv, err := h.MemKeeper.GetMetric(r.Context(), mj.ID)
				if err != nil {
					log.Println("value not found")
				} else {
					f = gv.GetInternalValue().(float64)
				}
				mj.Value = &f
			}
		case "counter":
			{
				if mj.Delta != nil {
					mjVal := *mj.Delta
					//пишем если есть значение
					prevMetricValue, err := h.MemKeeper.GetMetric(r.Context(), mj.ID)
					if err != nil {
						prevMetricValue = mVal.NewCounterValue()
					}
					sum := mVal.NewInt(mjVal).AddValue(prevMetricValue)
					err = h.MemKeeper.SaveMetric(r.Context(), mj.ID, &sum)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
				//читаем для ответа
				var i int64 = 0
				cv, err := h.MemKeeper.GetMetric(r.Context(), mj.ID)
				if err != nil {
					log.Println("value not found")
				} else {
					i = cv.GetInternalValue().(int64)
				}
				mj.Delta = &i

			}
		default:
			http.Error(w, mj.MType+" not recognized type", http.StatusNotImplemented)
			return
		}

		//перевод в json ответа
		byteData, err = json.Marshal(mj)
		if err != nil || byteData == nil {
			http.Error(w, " json response forming error", http.StatusInternalServerError)
			return
		}
		//response
		if next != nil {
			r.Body = io.NopCloser(bytes.NewReader(byteData))
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(byteData)
		if err != nil {
			log.Println("response writing error")
			http.Error(w, "response writing error", http.StatusInternalServerError)
			return
		}

	}
}

func (h *Handlers) HandlePostErrorPattern(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unknown request,HandlePostErrorPattern invoked", http.StatusNotFound)
	log.Println("Chi rounting error, unknown route to get handler")
	log.Println("HandlePostErrorPattern invoked")
	r.Body.Close()
}
func (h *Handlers) HandlePostErrorPatternNoName(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "Unknown request,HandlePostErrorPattern invoked", http.StatusNotFound)
	log.Println("Chi rounting error, unknown route to get handler")
	log.Println("HandlePostErrorPatternNoName invoked")
	r.Body.Close()
}

func (h *Handlers) NewRouter() chi.Router {

	// d := compression.Deflator{
	// 	Level: flate.BestSpeed,
	// }

	// var postJsonCompressedScenario = d.DeCompressionHandler(h.HandlePostMetricJSON(d.CompressionHandler(d.WriteResponseBodyHandler(nil))))
	var postJsonAndGetDataIncrement4Scenario = h.HandlePostMetricJSON(nil)

	r := chi.NewRouter()
	//
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.HandleGetMetricFieldList)
		r.Get("/value/{TYPE}/{NAME}", h.HandleGetMetricValue)
		r.Post("/value", postJsonAndGetDataIncrement4Scenario)
		r.Post("/value/", postJsonAndGetDataIncrement4Scenario)
		r.Post("/update", postJsonAndGetDataIncrement4Scenario)
		r.Post("/update/", postJsonAndGetDataIncrement4Scenario)
		r.Post("/update/{TYPE}/{NAME}/{VALUE}", h.HandlePostMetric)

		r.Post("/update/{TYPE}/{NAME}/", h.HandlePostErrorPattern)
		r.Post("/update/{TYPE}/", h.HandlePostErrorPatternNoName)

	})

	return r
}
