package countervalue

import (
	"strconv"

	interfaces "github.com/alphaonly/harvester/internal/server/interfaces/MetricValue"
)

type CounterValue struct {
	value    interfaces.Counter
	valueInt int64
}

func (v *CounterValue) New(c interfaces.Counter) *CounterValue {
	v.value = c
	v.valueInt = int64(c)
	return v
}
func NewInt(c int64) *CounterValue {
	return &CounterValue{
		value:    interfaces.Counter(c),
		valueInt: c,
	}
}

func (v *CounterValue) SetValue(counter interfaces.MetricValue) {
	v.value = interfaces.Counter(counter.(*CounterValue).value)
	v.valueInt = int64(v.value)
}
func (v *CounterValue) GetValue() interfaces.MetricValue {
	return v
}
func (v *CounterValue) GetInternalValue() interface{} {
	return v.valueInt
}
func (v *CounterValue) GetString() string {
	return strconv.FormatUint(uint64(v.value), 10)
}
func (v *CounterValue) AddValue(v1 interfaces.MetricValue) interfaces.MetricValue {
	sumVal := int64(v1.(*CounterValue).value) + v.valueInt
	return &CounterValue{value: interfaces.Counter(sumVal), valueInt: sumVal}
}

// type MetricValue interface {
// 	SetValue(MetricValue)
// 	GetValue() MetricValue         //Gauge or Counter
// 	GetInternalValue() interface{} //float64 or int64
// 	GetString() string
// 	AddValue(MetricValue) MetricValue // increment  current and return incremented( for counter only)
// }

// check
var m interfaces.MetricValue = &CounterValue{}
