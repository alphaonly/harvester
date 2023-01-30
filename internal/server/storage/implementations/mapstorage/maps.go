package mapstorage

import (
	"context"
	"errors"
	"sync"

	metricsjson "github.com/alphaonly/harvester/internal/server/metricsJSON"
	M "github.com/alphaonly/harvester/internal/server/metricvalue"
	S "github.com/alphaonly/harvester/internal/server/storage/interfaces"
)

type MapStorage struct {
	mutex      *sync.Mutex
	metricsMap *metricsjson.MetricsMapType
}

func New() (sr *S.Storage) {
	map_ := make(metricsjson.MetricsMapType)
	var mapStorage S.Storage = MapStorage{
		mutex:      &sync.Mutex{},
		metricsMap: &map_,
	}
	return &mapStorage
}

//Имплементация интерфейса storage

// type Storage interface {
// 	GetMetric(ctx context.Context, name string) (mv *interfaces.MetricValue, err error)
// 	SaveMetric(ctx context.Context, name string, mv *interfaces.MetricValue) (err error)
// 	GetAllMetrics(ctx context.Context) (mvList *map[string]interfaces.MetricValue, err error)
// }

func (m MapStorage) GetMetric(ctx context.Context, name string) (mv *M.MetricValue, err error) {

	if m.metricsMap == nil || len(*m.metricsMap) == 0 {
		return nil, errors.New("404 - not found")
	}

	if value := (*m.metricsMap)[name]; value == nil {
		return nil, errors.New("404 - not found")
	} else {
		return &value, nil
	}

}
func (m MapStorage) SaveMetric(ctx context.Context, name string, mv *M.MetricValue) (r error) {

	(*m.metricsMap)[name] = *mv

	return nil
}

func (m MapStorage) GetAllMetrics(ctx context.Context) (mvList *metricsjson.MetricsMapType, err error) {

	if m.metricsMap == nil {
		return nil, errors.New("map was not initialized")
	}

	return m.metricsMap, nil

}

func (m MapStorage) SaveAllMetrics(ctx context.Context, mvList *metricsjson.MetricsMapType) (err error) {
	(*m.metricsMap) = *mvList

	return nil
}
