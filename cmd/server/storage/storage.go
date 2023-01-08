package storage

import (
	"context"
	"errors"
	"github.com/golang-collections/collections/stack"
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
}

type GaugeValue struct {
	value Gauge
}

func (v *GaugeValue) GetValue() interface{} {
	return v.value
}
func (v *GaugeValue) SetValue(value interface{}) {
	v.value = value.(Gauge)
}
func (v *GaugeValue) GetString() string {
	return strconv.FormatFloat(float64(v.value), 'E', -1, 64)
}

type CounterValue struct {
	value Counter
}

func (v *CounterValue) GetValue() interface{} {
	return v.value
}
func (v *CounterValue) SetValue(value interface{}) {
	v.value = value.(Counter)
}
func (v *CounterValue) GetString() string {
	return strconv.FormatUint(uint64(v.value), 10)
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
	metricsStack    *stack.Stack
	mapMetricsStack *stack.Stack
	metricsMap      *map[string]MetricValue
}

type metricMemRepository interface {
	GetMetric(ctx context.Context, PollCount Counter) (*MetricsStorage, error)
	GetCurrentMetric(ctx context.Context) (*MetricsStorage, error)
	AddCurrentToMap(ctx context.Context, name string, value MetricValue) (r error)
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
	m.metricsStorage.metricsStack = &stack.Stack{}

	maps := make(map[string]MetricValue)
	m.metricsStorage.metricsMap = &maps

	m.metricsStorage.mapMetricsStack = &stack.Stack{}
	return m
}

func (m DataServer) GetMetric(ctx context.Context, PollCount Counter) (ms *Metrics, r error) {

	return nil, errors.New("no data")
}
func (m DataServer) AddCurrentToMap(ctx context.Context, name string, value MetricValue) (r error) {

	mp := *m.metricsStorage.metricsMap
	mp[name] = value

	return nil
}

func (m DataServer) SaveMetric(ctx context.Context, metrics Metrics) (r error) {
	ms := m.metricsStorage

	//if ms.ID.IsZero() && metrics.PollCount == 1 {
	//	ms.ID = time.Now()
	//}

	lenBefore := ms.mapMetricsStack.Len()

	ms.metricsStack.Push(metrics)
	ms.mapMetricsStack.Push(ms.metricsMap)

	mp := make(map[string]MetricValue)
	ms.metricsMap = &mp

	//msmap := *ms.metricsMap
	//msmap[metrics.PollCount] = metrics
	//
	//fmt.Println(ms.metricsStack.Len())

	lenAfter := ms.mapMetricsStack.Len()
	if lenAfter == lenBefore {
		return errors.New("Stack adding error")
	}

	return nil
}

func (m DataServer) DeleteMetric(ctx context.Context, PollCount Counter) (*MetricsStorage, error) {

	return nil, nil
}
func (m DataServer) GetAllMetrics(ctx context.Context, PollCount Counter) (*stack.Stack, error) {
	return nil, nil
}
func (m DataServer) GetCurrentMetric(ctx context.Context) (Metrics, error) {
	var (
		err error
		ms  = m.metricsStorage
	)
	stack := ms.metricsStack
	if &stack == nil {
		return Metrics{}, nil
	}
	if stack.Len() == 0 {
		return Metrics{}, nil
	}
	currentMetrics := stack.Peek().(Metrics)
	if &currentMetrics == nil {
		err = errors.New("unexpectedly no data in stack")
	}

	return currentMetrics, err
}
