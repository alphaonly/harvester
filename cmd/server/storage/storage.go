package storage

import (
	"context"
	"errors"
	"github.com/golang-collections/collections/stack"
	"reflect"
)

type Gauge float64
type Counter int64

//type MetricsMap map[string]interface{}
//type MetricsGaugeKeys map[string]

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

func (m Metrics) StringValue(field string) (value string, err error) {

	v, err := m.GetValue(field)
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (m Metrics) GetValue(field string) (v interface{}, err error) {

	r := reflect.ValueOf(m)

	value := reflect.Indirect(r).FieldByName(field)
	if value.IsZero() {
		return nil, errors.New("no value (isZero==true)")
	}

	switch field {
	case "PollCount":
		v = Counter(value.Uint())
	default:
		v = Gauge(value.Float())
	}
	return v, nil

}

type MetricsStorage struct {
	metricsStack *stack.Stack
	metricsMap   *map[Counter]Metrics
}

type metricMemRepository interface {
	GetMetric(ctx context.Context, PollCount Counter) (*MetricsStorage, error)
	GetCurrentMetric(ctx context.Context) (*MetricsStorage, error)
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
	//msmap := make(map[Counter]Metrics)
	//m.metricsStorage.metricsMap = &msmap
	return m
}

func (m DataServer) GetMetric(ctx context.Context, PollCount Counter) (ms *Metrics, r error) {

	return nil, errors.New("no data")
}
func (m DataServer) SaveMetric(ctx context.Context, metrics Metrics) (r error) {
	ms := m.metricsStorage

	//if ms.ID.IsZero() && metrics.PollCount == 1 {
	//	ms.ID = time.Now()
	//}

	lenBefore := ms.metricsStack.Len()

	ms.metricsStack.Push(metrics)

	//msmap := *ms.metricsMap
	//msmap[metrics.PollCount] = metrics
	//
	//fmt.Println(ms.metricsStack.Len())

	lenAfter := ms.metricsStack.Len()
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
	currentMetrics := stack.Peek().(Metrics)
	if &currentMetrics == nil {
		err = errors.New("Unexpectedly no data in stack")
	}

	return currentMetrics, err
}
