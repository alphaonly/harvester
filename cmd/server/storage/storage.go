package storage

import (
	"context"
	"errors"
	"reflect"
	"strconv"
)

type Gauge float64
type Counter int64

func (v Gauge) toString() (s string) {
	return strconv.FormatFloat(float64(v), 'E', -1, 64)
}
func (v Counter) toString() (s string) {
	return strconv.FormatUint(uint64(v), 10)
}

//type MetricsMap map[toString]interface{}
//type MetricsGaugeKeys map[toString]

type Metrics struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	RandomValue   Gauge

	PollCount Counter
}

type MetricValue interface {
	SetValue(interface{})
	GetValue() interface{}
	GetString() string
	AddValue(MetricValue) MetricValue
}

type GaugeValue struct {
	value      Gauge
	valueFloat float64
}

func (v *GaugeValue) GetValue() interface{} {
	return v.value
}
func (v *GaugeValue) SetValue(value interface{}) {
	v.value = value.(Gauge)
	v.valueFloat = float64(v.value)
}
func (v *GaugeValue) GetString() string {
	return strconv.FormatFloat(float64(v.value), 'f', -1, 64)
}
func (v *GaugeValue) AddValue(v1 MetricValue) MetricValue {
	ret := GaugeValue{}
	gValue := v1.(*GaugeValue)
	ret.SetValue(v.valueFloat + float64(gValue.valueFloat))

	return MetricValue(&ret)
}

type CounterValue struct {
	value    Counter
	valueInt int64
}

func (v *CounterValue) GetValue() interface{} {
	return v.value
}

func (v *CounterValue) SetValue(value interface{}) {

	v.value = value.(Counter)
	v.valueInt = int64(v.value)
}

func (v *CounterValue) GetString() string {
	return strconv.FormatUint(uint64(v.value), 10)
}
func (v *CounterValue) AddValue(v1 MetricValue) MetricValue {

	ret := CounterValue{}
	cValue := v1.(*CounterValue)
	ret.SetValue(Counter(v.valueInt + cValue.valueInt))

	return MetricValue(&ret)

}

func (m *Metrics) GetValue(field string) (mv *MetricValue, err error) {

	r := reflect.ValueOf(m)

	value := reflect.Indirect(r).FieldByName(field)
	if value.Kind() == reflect.String {
		var metricValue MetricValue
		switch field {
		case "PollCount":
			metricValue = &CounterValue{}
			metricValue.SetValue(Counter(value.Uint()))
		default:
			metricValue = &GaugeValue{}
			vf := value.Float()
			fvf := Gauge(vf)
			metricValue.SetValue(fvf)
			//v.SetValue(Gauge(value.Float()))
		}
		return &metricValue, nil
	}
	return nil, errors.New(" reflect error")
}

func (m *Metrics) StringValue(field string) (value string, err error) {

	v, err := m.GetValue(field)
	if err != nil {
		return "", err
	}

	return (*v).GetString(), nil

}

type MetricsStorage struct {
	// metricsStack       *stack.Stack
	// mapMetricsStack    *stack.Stack
	metricsMap *map[string]MetricValue
	// metricsDistinctSet *map[string]bool
}

type metricMemRepository interface {
	GetMetric(ctx context.Context, PollCount Counter) (*MetricsStorage, error)
	GetCurrentMetric(ctx context.Context) (*MetricsStorage, error)
	SaveMetricToMap(ctx context.Context, name string, value MetricValue) (r error)
	SaveMetric(ctx context.Context, metrics Metrics) error
	DeleteMetric(ctx context.Context, PollCount Counter) (*MetricsStorage, error)
	GetAllMetric(ctx context.Context, PollCount Counter) (*MetricsStorage, error)
}

type DataServer struct {
	metricsMemRepository metricMemRepository
	metricsStorage       MetricsStorage
}

func (DataServer) New() (m *DataServer) {

	m = &DataServer{}
	// m.metricsStorage.metricsStack = &stack.Stack{}

	maps := make(map[string]MetricValue)
	m.metricsStorage.metricsMap = &maps

	// m.metricsStorage.mapMetricsStack = &stack.Stack{}

	// m.metricsStorage.metricsDistinctSet = &mDistSett

	return m
}

func (m DataServer) GetMetric(ctx context.Context, PollCount Counter) (ms *Metrics, r error) {

	return nil, errors.New("no data")
}
func (m DataServer) SaveMetricToMap(ctx context.Context, name string, value MetricValue) (r error) {

	mp := *m.metricsStorage.metricsMap
	// ms := *m.metricsStorage.metricsDistinctSet

	mp[name] = value
	// ms[name] = true
	return nil
}

func (m DataServer) SaveMetric(ctx context.Context, metrics Metrics) (r error) {
	ms := m.metricsStorage

	// ms.metricsStack.Push(metrics)
	// ms.mapMetricsStack.Push(ms.metricsMap)

	mp := make(map[string]MetricValue)
	ms.metricsMap = &mp

	return nil
}

func (m DataServer) DeleteMetric(ctx context.Context, PollCount Counter) (*MetricsStorage, error) {

	return nil, nil
}
func (m DataServer) GetAllMetricsNames(ctx context.Context) (map[string]MetricValue, error) {

	if m.metricsStorage.metricsMap == nil {
		return nil, errors.New("no list of metrics names")
	}

	return *m.metricsStorage.metricsMap, nil

}
func (m DataServer) GetCurrentMetricMap(ctx context.Context, name string) (MetricValue, error) {

	ms := m.metricsStorage

	if &ms.metricsMap == nil {
		return nil, errors.New("no map initialized")
	}

	if len(*ms.metricsMap) == 0 {
		return nil, nil
	}

	mp := *ms.metricsMap
	value := mp[name]
	if value == nil {
		return nil, errors.New("404")
	}

	return value, nil

}
