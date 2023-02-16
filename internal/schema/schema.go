package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	mVal "github.com/alphaonly/harvester/internal/server/metricvalueInt"

	"io"
	"log"
	"net/http"
	"net/url"
)

type PreviousBytes []byte
type ContextKey int

const PKey1 ContextKey = 123455

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрикИ
	MType string   `json:"type"`            // параметр, пID    string   `json:"id"`              // имя метрикиринимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func NewMetricJSON(name string, MType string, value interface{}) (ret MetricsJSON) {

	j := MetricsJSON{}
	if name == "" || MType == "" {
		panic(errors.New("name or type is empty"))
	}

	j.ID = name
	j.MType = MType

	if value != nil {
		switch MType {
		case "agent.gauge", "gauge":
			var val = float64(value.(mVal.Gauge))
			j.Value = &val
		case "agent.counter", "counter":
			var val = int64(value.(mVal.Counter))
			j.Delta = &val
		default:
			panic(errors.New("unknown type"))
		}
	}
	return j
}
func GetMetricJSONWithPOST(baseURL *url.URL, name string, MType string) (mj MetricsJSON) {

	metricsJSONRequest := NewMetricJSON(name, MType, nil)
	data, err := json.Marshal(metricsJSONRequest)
	if err != nil {
		log.Fatal(err)
	}

	URL := (*baseURL).JoinPath("value")
	request, err := http.NewRequest(http.MethodPost, URL.String(), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Accept", "application/json; charset=utf-8")

	client := http.Client{}
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var metricsJSONResponse MetricsJSON
	err = json.Unmarshal(responseData, &metricsJSONResponse)
	if err != nil {
		log.Fatal(err)
	}
	response.Body.Close()
	return metricsJSONResponse

}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}
