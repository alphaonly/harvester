package gaugevalue

import (
	"encoding/json"
	"errors"
	"strconv"

	interfaces "github.com/alphaonly/harvester/internal/server/metricvalue"
)

type GaugeValueJSON struct {
	Value float64 `json:"value"`
}

type GaugeValue struct {
	value      interfaces.Gauge
	valueFloat float64
}

func (v *GaugeValue) New(g interfaces.Gauge) *GaugeValue {
	v.value = g
	v.valueFloat = float64(g)
	return v
}
func NewFloat(g float64) *GaugeValue {
	return &GaugeValue{
		value:      interfaces.Gauge(g),
		valueFloat: g,
	}
}
func (v *GaugeValue) SetValue(gauge interfaces.MetricValue) {
	v.value = interfaces.Gauge(gauge.(*GaugeValue).value)
	v.valueFloat = float64(v.value)
}
func (v *GaugeValue) GetValue() interfaces.MetricValue {
	return v
}
func (v *GaugeValue) GetInternalValue() interface{} {
	return v.valueFloat
}
func (v *GaugeValue) GetString() string {
	return strconv.FormatFloat(float64(v.value), 'f', -1, 64)
}

func (v *GaugeValue) AddValue(v1 interfaces.MetricValue) interfaces.MetricValue {
	return v //Mocked, as it's needed for counter only
}
func (v *GaugeValue) MarshalJSON() ([]byte, error) {
	cj := GaugeValueJSON{Value: v.valueFloat}

	return json.Marshal(cj)
}

func (v *GaugeValue) UnmarshalJSON(data []byte) error {
	var cj = &GaugeValueJSON{}

	err := json.Unmarshal(data, cj)
	if err != nil {
		return errors.New("CounterValue unmarshal error")
	}
	v.valueFloat = cj.Value

	v.value = interfaces.Gauge(v.valueFloat)

	return nil

}

// check
//var m interfaces.MetricValue = &GaugeValue{}
