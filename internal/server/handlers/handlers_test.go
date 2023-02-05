package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/alphaonly/harvester/cmd/httplib/increment4/fork"
	"github.com/alphaonly/harvester/internal/schema"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"time"

	"net/url"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"
)

func TestHandleMetric(t *testing.T) {

	type want struct {
		code        int
		response    string
		contentType string
	}
	type requestParams struct {
		method string
		url    string
		want   want
	}

	data := url.Values{}

	metricsRequestsParam := make(map[string]requestParams)

	contentType := "text/plain"
	urlPrefix := ""
	//Check Url Ok
	urlStr := urlPrefix + "/update/gauge/Alloc/2.36912E+05"
	r1 := requestParams{method: http.MethodPost, url: urlStr,
		want: want{code: http.StatusOK, response: `{"status":"ok"}`, contentType: contentType}}

	//Check Url bad unknown metric
	urlStr = urlPrefix + "/update/gauge/Alerrorloc/2.36912E+05"
	r3 := requestParams{method: http.MethodPost, url: urlStr,
		want: want{code: http.StatusOK, response: ``, contentType: contentType}}
	//Check Url bad method
	urlStr = urlPrefix + "/update/gauge/Alloc/2.36912E+05"
	r5 := requestParams{method: http.MethodGet, url: urlStr,
		want: want{code: http.StatusMethodNotAllowed, response: `{"status":"ok"}`, contentType: contentType}}
	//Check Url empty metric
	urlStr = urlPrefix + "/update/gauge//2.36912E+05"
	r6 := requestParams{method: http.MethodPost, url: urlStr,
		want: want{code: http.StatusNotFound, response: `{"status":"ok"}`, contentType: contentType}}
	//Check Url empty metric value
	urlStr = urlPrefix + "/update/gauge/counter/"
	r7 := requestParams{method: http.MethodPost, url: urlStr,
		want: want{code: http.StatusNotFound, response: `{"status":"ok"}`, contentType: contentType}}

	//var r4 requestParams

	metricsRequestsParam["r1"] = r1

	metricsRequestsParam["r3"] = r3

	metricsRequestsParam["r5"] = r5
	metricsRequestsParam["r6"] = r6
	metricsRequestsParam["r7"] = r7

	tests := []struct {
		name string
		ID   string
		want want
	}{
		{
			name: "test#1 positive",
			ID:   "r1",
			want: metricsRequestsParam["r1"].want,
		},

		{
			name: "test#3 negative",
			ID:   "r3",
			want: metricsRequestsParam["r3"].want,
		},

		{
			name: "test#5 negative",
			ID:   "r5",
			want: metricsRequestsParam["r5"].want,
		},
		{
			name: "test#6 negative",
			ID:   "r6",
			want: metricsRequestsParam["r6"].want,
		},
		{
			name: "test#7 negative",
			ID:   "r7",
			want: metricsRequestsParam["r7"].want,
		},
	}
	fmt.Println("start!")

	s := mapstorage.New()
	h := New(s)

	r := h.NewRouter()

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("url from test:" + metricsRequestsParam[tt.ID].url)

			request := httptest.NewRequest(metricsRequestsParam[tt.ID].method, metricsRequestsParam[tt.ID].url, bytes.NewBufferString(data.Encode()))

			w := httptest.NewRecorder()

			//h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//
			//	handlers := Handlers{}
			//	handlers.HandlePostMetric(w, r)
			//})

			r.ServeHTTP(w, request)

			response := w.Result()
			if response.StatusCode != tt.want.code {
				t.Errorf("error code %v want %v", response.StatusCode, tt.want.code)
				fmt.Println(response)
				fmt.Println(w.Body.String())

			}

			if (response.StatusCode == http.StatusOK) &&
				(response.Header.Get("Content-type") != tt.want.contentType) {
				t.Errorf("error contentType %v want %v", response.Header.Get("Content-type"), tt.want.contentType)
			}
			err := response.Body.Close()
			if err != nil {
				t.Errorf("response body close error: %v response", response.Body)
			}
		})

	}

}

func TestHandlePostMetricJSON(t *testing.T) {
	type Iteration4Suite struct {
		suite.Suite

		serverAddress string
		serverProcess *fork.BackgroundProcess
		agentProcess  *fork.BackgroundProcess

		knownEncodingLibs []string

		rnd *rand.Rand
	}
	var suite *Iteration4Suite

	type want struct {
		code         int
		responseBody []byte
		contentType  string
	}
	type requestParams struct {
		method string
		URL    url.URL
		want   want
	}

	errRedirectBlocked := errors.New("HTTP redirect blocked")
	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})
	httpc := resty.New().
		SetHostURL("http://localhost:8080/update").
		SetRedirectPolicy(redirPolicy)

	tests := []struct {
		name   string
		method string
		value  float64
		delta  int64
		update int
		ok     bool
		static bool
	}{
		{method: "counter", name: "PollCount"},
		{method: "gauge", name: "RandomValue"},
		{method: "gauge", name: "Alloc"},
		{method: "gauge", name: "BuckHashSys", static: true},
		{method: "gauge", name: "Frees"},
		{method: "gauge", name: "GCCPUFraction", static: true},
		{method: "gauge", name: "GCSys", static: true},
		{method: "gauge", name: "HeapAlloc"},
		{method: "gauge", name: "HeapIdle"},
		{method: "gauge", name: "HeapInuse"},
		{method: "gauge", name: "HeapObjects"},
		{method: "gauge", name: "HeapReleased", static: true},
		{method: "gauge", name: "HeapSys", static: true},
		{method: "gauge", name: "LastGC", static: true},
		{method: "gauge", name: "Lookups", static: true},
		{method: "gauge", name: "MCacheInuse", static: true},
		{method: "gauge", name: "MCacheSys", static: true},
		{method: "gauge", name: "MSpanInuse", static: true},
		{method: "gauge", name: "MSpanSys", static: true},
		{method: "gauge", name: "Mallocs"},
		{method: "gauge", name: "NextGC", static: true},
		{method: "gauge", name: "NumForcedGC", static: true},
		{method: "gauge", name: "NumGC", static: true},
		{method: "gauge", name: "OtherSys", static: true},
		{method: "gauge", name: "PauseTotalNs", static: true},
		{method: "gauge", name: "StackInuse", static: true},
		{method: "gauge", name: "StackSys", static: true},
		{method: "gauge", name: "Sys", static: true},
		{method: "gauge", name: "TotalAlloc"},
	}

	req := httpc.R().
		SetHeader("Content-Type", "application/json")

	timer := time.NewTimer(time.Minute)

cont:
	for ok := 0; ok != len(tests); {
		// suite.T().Log("tick", len(tests)-ok)
		select {
		case <-timer.C:
			break cont
		default:
		}
		for i, tt := range tests {
			if tt.ok {
				continue
			}
			var (
				resp *resty.Response
				err  error
			)
			time.Sleep(100 * time.Millisecond)

			var result schema.MetricsJSON
			resp, err = req.
				SetBody(&schema.MetricsJSON{
					ID:    tt.name,
					MType: tt.method,
				}).
				SetResult(&result).
				Post("/value/")

			dumpErr := suite.Assert().NoErrorf(err, "Ошибка при попытке сделать запрос с получением значения %s", tt.name)
			if resp.StatusCode() == http.StatusNotFound {
				continue
			}
			dumpErr = dumpErr && suite.Assert().Containsf(resp.Header().Get("Content-Type"), "application/json",
				"Заголовок ответа Content-Type содержит несоответствующее значение")
			dumpErr = dumpErr && suite.Assert().True(((result.MType == "gauge" && result.Value != nil) || (result.MType == "counter" && result.Delta != nil)),
				"Получен не однозначный результат (тип метода не соответствует возвращаемому значению) '%q %s'", req.Method, req.URL)
			dumpErr = dumpErr && suite.Assert().True(result.MType != "gauge" || result.Value != nil,
				"Получен не однозначный результат (возвращаемое значение value=nil не соответствет типу gauge) '%q %s'", req.Method, req.URL)
			dumpErr = dumpErr && suite.Assert().True(result.MType != "counter" || result.Delta != nil,
				"Получен не однозначный результат (возвращаемое значение delta=nil не соответствет типу counter) '%q %s'", req.Method, req.URL)
			dumpErr = dumpErr && suite.Assert().False(result.Delta == nil && result.Value == nil,
				"Получен результат без данных (Dalta == nil && Value == nil) '%q %s'", req.Method, req.URL)
			dumpErr = dumpErr && suite.Assert().False(result.Delta != nil && result.Value != nil,
				"Получен не однозначный результат (Dalta != nil && Value != nil) '%q %s'", req.Method, req.URL)
			dumpErr = dumpErr && suite.Assert().Equalf(http.StatusOK, resp.StatusCode(),
				"Несоответствие статус кода ответа ожидаемому в хендлере %q: %q ", req.Method, req.URL)
			dumpErr = dumpErr && suite.Assert().True(result.MType == "gauge" || result.MType == "counter",
				"Получен ответ с неизвестным значением типа: %q, '%q %s'", result.MType, req.Method, req.URL)

			if !dumpErr {
				//dump := dumpRequest(req.RawRequest, true)
				//suite.T().Logf("Оригинальный запрос:\n\n%s", dump)
				//dump = dumpResponse(resp.RawResponse, true)
				//suite.T().Logf("Оригинальный ответ:\n\n%s", dump)
				return
			}

			switch tt.method {
			case "gauge":
				if (tt.update != 0 && *result.Value != tt.value) || tt.static {
					tests[i].ok = true
					ok++
					suite.T().Logf("get %s: %q, value: %f", tt.method, tt.name, *result.Value)
				}
				tests[i].value = *result.Value
			case "counter":
				if (tt.update != 0 && *result.Delta != tt.delta) || tt.static {
					tests[i].ok = true
					ok++
					suite.T().Logf("get %s: %q, value: %d", tt.method, tt.name, *result.Delta)
				}
				tests[i].delta = *result.Delta
			}

			tests[i].update++
		}
	}
	for _, tt := range tests {
		suite.Run(tt.method+"/"+tt.name, func() {
			suite.Assert().Truef(tt.ok, "Отсутствует изменение метрики: %s, тип: %s", tt.name, tt.method)
		})
	}

}
