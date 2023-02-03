package metricsjson

import (
	"encoding/json"
	"errors"

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

func (m MetricsMapType) MarshalJSON() ([]byte, error) {
	mjArray := make([]MetricsJSON, len(m))
	i := 0
	for k, v := range m {
		switch value := v.(type) {
		case *gaugevalue.GaugeValue:
			{
				v := value.GetInternalValue().(float64)
				mjArray[i] = MetricsJSON{
					ID:    k,
					MType: "gauge",
					Value: &v,
				}
			}
		case *countervalue.CounterValue:
			{
				v := value.GetInternalValue().(int64)
				mjArray[i] = MetricsJSON{
					ID:    k,
					MType: "counter",
					Delta: &v,
				}
			}
		default:
			return nil, errors.New("undefined type in type switch metricValue")
		}
		i++
	}
	mBytes, err := json.Marshal(&mjArray)
	if err != nil {
		return nil, err
	}

	return mBytes, nil
}

func (m MetricsMapType) UnmarshalJSON(b []byte) error {
	var mjArray []MetricsJSON
	if err := json.Unmarshal(b, &mjArray); err != nil {
		return err
	}
	for _, v := range mjArray {
		switch v.MType {
		case "gauge":
			m[v.ID] = gaugevalue.NewFloat(*v.Value)
		case "counter":
			m[v.ID] = countervalue.NewInt(*v.Delta)
		default:
			return errors.New("unknown type in decoding metricsJSONArray")
		}
	}

	return nil
}
