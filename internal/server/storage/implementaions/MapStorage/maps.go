package MapStorage

import (
	"context"
	"errors"
	"github.com/alphaonly/harvester/internal/server/interfaces/MetricValue"
	"sync"
)

type MapStorage struct {
	mutex      *sync.Mutex
	metricsMap *map[string]interfaces.MetricValue
}

func New() (sr MapStorage) {
	map_ := make(map[string]interfaces.MetricValue)
	return MapStorage{
		mutex:      &sync.Mutex{},
		metricsMap: &map_,
	}
}

//Имплементация интерфейса storage

// type Storage interface {
// 	GetMetric(ctx context.Context, name string) (mv interfaces.MetricValue, err error)
// 	SaveMetric(ctx context.Context, name string, mv interfaces.MetricValue) (err error)
// 	GetAllMetrics(ctx context.Context) (mvList map[string]interfaces.MetricValue, err error)
// }

func (m MapStorage) GetMetric(ctx context.Context, name string) (mv *interfaces.MetricValue, err error) {

	if m.metricsMap == nil || len(*m.metricsMap) == 0 {
		return nil, errors.New("404 - not found")
	}

	if value := (*m.metricsMap)[name]; value == nil {
		return nil, errors.New("404 - not found")
	} else {
		return &value, nil
	}

}
func (m MapStorage) SaveMetric(ctx context.Context, name string, mv *interfaces.MetricValue) (r error) {

	(*m.metricsMap)[name] = *mv

	return nil
}

func (m MapStorage) GetAllMetrics(ctx context.Context) (mvList *map[string]interfaces.MetricValue, err error) {

	if m.metricsMap == nil {
		return nil, errors.New("no list of metrics")
	}

	return m.metricsMap, nil

}
