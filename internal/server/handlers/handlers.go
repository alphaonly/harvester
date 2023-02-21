package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/schema"
	"github.com/alphaonly/harvester/internal/server/compression"
	mVal "github.com/alphaonly/harvester/internal/server/metricvalueInt"
	"github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"
	"github.com/alphaonly/harvester/internal/signchecker"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Handlers struct {
	MemKeeper *mapstorage.MapStorage
	Signer    signchecker.Signer
	Conf      configuration.ServerConfiguration
}

func (h *Handlers) HandleGetMetricFieldListSimple(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
			return
		}
		log.Println("HandleGetMetricFieldListXXX invoked")

		ms, err := h.MemKeeper.GetAllMetrics(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//new byte buffer
		bw := bytes.NewBuffer(*new([]byte))
		_, err = bw.Write([]byte("<h1><ul>"))
		logFatal(err)
		//insert all metrics from memKeeper
		for key := range *ms {
			_, err = bw.Write([]byte(" <li>" + key + "</li>"))
			logFatal(err)
		}
		_, err = bw.Write([]byte("</ul></h1>"))
		logFatal(err)

		//compress
		var bytesData = bw.Bytes()
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			//Compression logic
			compressedByteData, err := compression.GzipCompress(bytesData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}
			bytesData = *compressedByteData
		}

		//Add header keys
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)

		//response to further handler
		if next == nil {
			//write handled body for further handle
			_, err = w.Write(bytesData)
			logFatal(err)
			return
		}
		log.Fatal(" HandleGetMetricFieldList requires next handler nil")
	}
}
func (h *Handlers) HandleGetMetricFieldList(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
			return
		}
		log.Println("HandleGetMetricFieldList invoked")

		ms, err := h.MemKeeper.GetAllMetrics(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//new byte buffer
		bw := bytes.NewBuffer(*new([]byte))
		_, err = bw.Write([]byte("<h1><ul>"))
		logFatal(err)
		//insert all metrics from memKeeper
		for key := range *ms {
			_, err = bw.Write([]byte(" <li>" + key + "</li>"))
			logFatal(err)
		}
		_, err = bw.Write([]byte("</ul></h1>"))
		logFatal(err)

		//Add header keys
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")

		log.Printf("Check response Content-Encoding in final header, value:%v", w.Header().Get("Content-Encoding"))
		log.Printf("Check response Content-Type in final header, value:%v", w.Header().Get("Content-Type"))

		//response to further handler
		if next != nil {
			//write handled body for further handle

			ctx := context.WithValue(r.Context(), schema.PKey1, schema.PreviousBytes(bw.Bytes()))
			//call further handler with context parameters
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		log.Fatal(" HandleGetMetricFieldList requires next handler not nil")
	}
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

func (h *Handlers) handlePostMetricJSONValidate(w http.ResponseWriter, r *http.Request) (ok bool) {
	if h.MemKeeper == nil {
		http.Error(w, "storage not initiated", http.StatusInternalServerError)
		return false
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func (h *Handlers) getBody(w http.ResponseWriter, r *http.Request) (b []byte, ok bool) {

	var bytesData []byte
	var err error
	var prev schema.PreviousBytes

	if p := r.Context().Value(schema.PKey1); p != nil {
		prev = p.(schema.PreviousBytes)
	}
	if prev != nil {
		//body from previous handler
		bytesData = prev
		log.Printf("got body from previous handler:%v", string(bytesData))
	} else {
		//body from request if there is no previous handler
		bytesData, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unrecognized request body:"+err.Error(), http.StatusBadRequest)
			return nil, false
		}
		log.Printf("got body from request:%v", string(bytesData))
	}
	log.Printf("Server:json body received:" + string(bytesData))
	return bytesData, true
}
func httpError(w http.ResponseWriter, err string, status int) {
	http.Error(w, err, status)
	log.Println("server:" + err)
}

func (h *Handlers) HandlePostMetricJSON(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostMetricJSON invoked")
		//validation
		if !h.handlePostMetricJSONValidate(w, r) {
			return
		}
		//Handle
		//1. get body
		bytesData, ok := h.getBody(w, r)
		if !ok {
			return
		}
		//2. JSON
		var mj schema.MetricsJSON
		err := json.Unmarshal(bytesData, &mj)
		if err != nil {
			httpError(w, "unmarshal error:", http.StatusBadRequest)
			return
		}
		//3. Валидация полученных данных
		if mj.ID == "" {
			httpError(w, "not parsed, empty metric name!"+mj.ID, http.StatusNotFound)
			return
		}
		//4.Проверяем подпись по ключу, нормально если ключ пуст в случае /update
		if mj.Delta != nil || mj.Value != nil {
			if !h.Signer.IsValidSign(mj) {
				httpError(w, "sign is not confirmed error", http.StatusBadRequest)
				log.Printf("server:sign is not confirmed error:%v", string(bytesData))
				return
			}
		}
		//Сохраняем в базу от агента и ответ обратно
		err = h.writeToStorageAndRespond(&mj, w, r)
		logFatal(err)

		//Подписываем ответ если есть значение
		if !(mj.Delta == nil && mj.Value == nil) {
			err = h.Signer.Sign(&mj)
			logFatal(err)
			//перевод в json ответа
		}
		bytesData, err = json.Marshal(mj)
		if err != nil || bytesData == nil {
			httpError(w, " json response forming error", http.StatusInternalServerError)
			return
		}
		//Set Header keys
		w.Header().Set("Content-Type", "application/json")
		//response
		if next != nil {
			//write handled body for further handle
			ctx := context.WithValue(r.Context(), schema.PKey1, schema.PreviousBytes(bytesData))
			//call further handler with context parameters
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		log.Fatal("HandlePostMetricJSON handler requires next handler not nil")
	}
}
func (h *Handlers) HandlePostErrorPattern(w http.ResponseWriter, r *http.Request) {
	log.Println("HandlePostErrorPattern invoked")

	httpError(w, "Unknown request,HandlePostErrorPattern invoked", http.StatusNotFound)
	err := r.Body.Close()
	if err != nil {
		return
	}
}
func (h *Handlers) HandlePostErrorPatternNoName(w http.ResponseWriter, r *http.Request) {
	log.Println("HandlePostErrorPatternNoName invoked")
	httpError(w, "Unknown request,HandlePostErrorPattern invoked", http.StatusNotFound)
	err := r.Body.Close()
	if err != nil {
		return
	}
}
func (h *Handlers) WriteResponseBodyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("WriteResponseBodyHandler invoked")

		//read body
		var bytesData []byte
		var err error
		var prev schema.PreviousBytes

		if p := r.Context().Value(schema.PKey1); p != nil {
			prev = p.(schema.PreviousBytes)
		}
		if prev != nil {
			//body from previous handler
			bytesData = prev
			log.Printf("got body from previous handler:%v", string(bytesData))
		} else {
			//body from request if there is no previous handler
			bytesData, err = io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}
			log.Printf("got body from request:%v", string(bytesData))
		}
		//Set flag in case compressed data
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
		}
		//Set Response Header
		w.WriteHeader(http.StatusOK)
		//write Response Body
		_, err = w.Write(bytesData)
		if err != nil {
			log.Println("byteData writing error")
			http.Error(w, "byteData writing error", http.StatusInternalServerError)
			return
		}
	}

}

func (h *Handlers) HandlePing(w http.ResponseWriter, r *http.Request) {
	log.Println("HandlePing invoked")
	log.Println("server:HandlePing:database string:" + h.Conf.DatabaseDsn)
	conn, err := pgx.Connect(r.Context(), h.Conf.DatabaseDsn)
	if err != nil {
		httpError(w, "server: ping handler: Unable to connect to database:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())
	log.Println("server: ping handler: connection established, 200 OK ")
	w.Write([]byte("200 OK"))
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) NewRouter() chi.Router {

	var (
		writePost = h.WriteResponseBodyHandler
		//writeList = h.WriteResponseBodyHandler

		compressPost = compression.GZipCompressionHandler
		//compressList = compression.GZipCompressionHandler

		handlePost = h.HandlePostMetricJSON
		//handleList = h.HandleGetMetricFieldList

		//The sequence for post JSON and respond compressed JSON if no value
		postJsonAndGetCompressed = handlePost(compressPost(writePost()))

		//The sequence for get compressed metrics html list
		//getListCompressed = handleList(compressList(writeList()))
		getListCompressed = h.HandleGetMetricFieldListSimple(nil)
	)
	r := chi.NewRouter()
	//

	// var p PingHandler
	r.Route("/", func(r chi.Router) {
		r.Get("/", getListCompressed)
		r.Get("/ping", h.HandlePing)
		r.Get("/ping/", h.HandlePing)
		r.Get("/value/{TYPE}/{NAME}", h.HandleGetMetricValue)
		r.Post("/value", postJsonAndGetCompressed)
		r.Post("/value/", postJsonAndGetCompressed)
		r.Post("/update", postJsonAndGetCompressed)
		r.Post("/update/", postJsonAndGetCompressed)
		r.Post("/update/{TYPE}/{NAME}/{VALUE}", h.HandlePostMetric)

		r.Post("/update/{TYPE}/{NAME}/", h.HandlePostErrorPattern)
		r.Post("/update/{TYPE}/", h.HandlePostErrorPatternNoName)

	})

	return r
}

func (h *Handlers) writeToStorageAndRespond(mj *schema.MetricsJSON, w http.ResponseWriter, r *http.Request) (err error) {
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
					return err
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
					return err
				}
			}
			//читаем для ответа
			var i int64 = 0
			cv, err := h.MemKeeper.GetMetric(r.Context(), mj.ID)
			if err != nil {
				log.Println("server:value not found:" + mj.ID)
			} else {
				i = cv.GetInternalValue().(int64)
			}
			mj.Delta = &i
		}
	default:
		mess := " not recognized type"
		http.Error(w, mj.MType+mess, http.StatusNotImplemented)
		return errors.New(mj.MType + mess)
	}
	return nil
}
func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
