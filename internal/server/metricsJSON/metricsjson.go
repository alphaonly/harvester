package metricsjson

import (
	"encoding/json"
	"strconv"

	"github.com/alphaonly/harvester/internal/server/metricvalue"
	countervalue "github.com/alphaonly/harvester/internal/server/metricvalue/MetricValue/implementations/CounterValue"
	gaugevalue "github.com/alphaonly/harvester/internal/server/metricvalue/MetricValue/implementations/Gaugevalue"
)

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, пID    string   `json:"id"`              // имя метрикиринимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsMapType map[string]metricvalue.MetricValue

func (m MetricsMapType) UnmarshalJSON(b []byte) error {
	data := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		var dst metricvalue.MetricValue
		// populate dst with an instance of the actual type you want to unmarshal into
		if _, err := strconv.Atoi(string(v)); err != nil {
			dst = &gaugevalue.GaugeValue{} // notice the dereference
		} else {
			dst = &countervalue.CounterValue{}
		}

		if err := json.Unmarshal(v, dst); err != nil {
			return err
		}
		m[k] = dst
	}
	return nil
}
