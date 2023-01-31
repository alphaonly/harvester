package agentjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/alphaonly/harvester/internal/agent"
)

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, пID    string   `json:"id"`              // имя метрикиринимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
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
			var val = float64(value.(agent.Gauge))
			j.Value = &val
		case "agent.counter", "counter":
			var val = int64(value.(agent.Counter))
			j.Delta = &val
		default:
			panic(errors.New("unknown type"))
		}
	}
	return j
}
func (j *MetricsJSON) GetMetricJSON(baseURL *url.URL, name string, MType string) (mj MetricsJSON) {

	metricsJSONRequest := NewMetricJSON(name, MType, nil)
	data, err := json.Marshal(metricsJSONRequest)
	if err != nil {
		log.Fatal(err)
	}

	URL := (*baseURL).JoinPath("value")
	request, err := http.NewRequest(http.MethodGet, URL.String(), bytes.NewBuffer(data))
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
