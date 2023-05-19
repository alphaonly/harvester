package metricvalue

type Gauge float64
type Counter int64

type MetricValue interface {
	SetValue(MetricValue)
	GetValue() MetricValue         //Gauge or Counter
	GetInternalValue() interface{} //float64 or int64
	GetString() string
	AddValue(MetricValue) MetricValue // increment  current and return incremented( for counter only)

	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}
