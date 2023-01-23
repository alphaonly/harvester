package interfaces

import (
	"context"

	interfaces "github.com/alphaonly/harvester/internal/server/interfaces/MetricValue"
)

type Gauge float64
type Counter int64

type MetricsJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, пID    string   `json:"id"`              // имя метрикиринимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Storage interface {
	GetMetric(ctx context.Context, name string) (mv *interfaces.MetricValue, err error)
	SaveMetric(ctx context.Context, name string, mv *interfaces.MetricValue) (err error)
	GetAllMetrics(ctx context.Context) (mvList *map[string]interfaces.MetricValue, err error)
}
