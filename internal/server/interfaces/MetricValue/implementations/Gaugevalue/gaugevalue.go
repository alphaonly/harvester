package gaugevalue

import (
	"strconv"

	interfaces "github.com/alphaonly/harvester/internal/server/interfaces"
)

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

// check
var m interfaces.MetricValue = &GaugeValue{}
