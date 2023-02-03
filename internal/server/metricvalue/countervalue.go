package countervalue

import (
	"encoding/json"
	"errors"
	"strconv"

	interfaces "github.com/alphaonly/harvester/internal/server/metricvalue"
)

type CounterValueJSON struct {
	Value string `json:"value"`
}

type CounterValue struct {
	value    interfaces.Counter
	valueInt int64
}

func NewCounterValue() *interfaces.MetricValue {
	m := interfaces.MetricValue(&CounterValue{})
	return &m
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

func (v *CounterValue) MarshalJSON() ([]byte, error) {
	cj := CounterValueJSON{Value: v.GetString()}

	return json.Marshal(cj)
}

func (v *CounterValue) UnmarshalJSON(data []byte) error {
	cj := &CounterValueJSON{}

	err := json.Unmarshal(data, cj)
	if err != nil {
		return errors.New("CounterValue unmarshal error")
	}
	v.valueInt, err = strconv.ParseInt(cj.Value, 10, 64)
	if err != nil {
		return errors.New("unmarshalled CounterValue parse error")
	}
	v.value = interfaces.Counter(v.valueInt)

	return nil

}

// type MetricValue interface {
// 	SetValue(MetricValue)
// 	GetValue() MetricValue         //Gauge or Counter
// 	GetInternalValue() interface{} //float64 or int64
// 	GetString() string
// 	AddValue(MetricValue) MetricValue // increment  current and return incremented( for counter only)
// }

// check
//var m interfaces.MetricValue = &CounterValue{}
//