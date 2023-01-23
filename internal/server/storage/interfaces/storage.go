package storage

import (
	"context"

	M "github.com/alphaonly/harvester/internal/server/interfaces"
)

type Gauge float64
type Counter int64
type Storage interface {
	GetMetric(ctx context.Context, name string) (mv *M.MetricValue, err error)
	SaveMetric(ctx context.Context, name string, mv *M.MetricValue) (err error)
	GetAllMetrics(ctx context.Context) (mvList *map[string]M.MetricValue, err error)
}
