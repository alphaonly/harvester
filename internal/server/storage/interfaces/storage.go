package storage

import (
	"context"

	metricsjson "github.com/alphaonly/harvester/internal/server/metricsJSON"
	"github.com/alphaonly/harvester/internal/server/metricvalue"
)

type Gauge float64
type Counter int64
type Storage interface {
	GetMetric(ctx context.Context, name string) (mv *metricvalue.MetricValue, err error)
	SaveMetric(ctx context.Context, name string, mv *metricvalue.MetricValue) (err error)
	GetAllMetrics(ctx context.Context) (mvList *metricsjson.MetricsMapType, err error)
	SaveAllMetrics(ctx context.Context, mvList *metricsjson.MetricsMapType) (err error)
}
