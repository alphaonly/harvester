package filestorage

import (
	"context"
	"errors"

	"github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server/files"
	metricsjson "github.com/alphaonly/harvester/internal/server/metricsJSON"
	M "github.com/alphaonly/harvester/internal/server/metricvalue"
	S "github.com/alphaonly/harvester/internal/server/storage/interfaces"
)

// type Storage interface {
// 	GetMetric(ctx context.Context, name string) (mv *M.MetricValue, err error)
// 	SaveMetric(ctx context.Context, name string, mv *M.MetricValue) (err error)
// 	GetAllMetrics(ctx context.Context) (mvList *metricsjson.MetricsMapType, err error)
// 	SaveAllMetrics(ctx context.Context, mvList *metricsjson.MetricsMapType) (err error)
// }

type FileArchive struct {
	configuration *configuration.Configuration
}

func New(c *configuration.Configuration) *S.Storage {
	var s S.Storage = FileArchive{
		configuration: c,
	}
	return &s
}

func (fa FileArchive) GetMetric(ctx context.Context, name string) (mv *M.MetricValue, err error) {
	//Not supported by the implementation
	return nil, errors.New("not supported")
}
func (fa FileArchive) SaveMetric(ctx context.Context, name string, mv *M.MetricValue) (err error) {
	//Not supported by the implementation
	return errors.New("not supported")
}

// Restore data from temp dir
func (fa FileArchive) GetAllMetrics(ctx context.Context) (mvList *metricsjson.MetricsMapType, err error) {
	consumer, err := files.NewConsumer((*fa.configuration).Get("STORE_FILE"))
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	mvList, err = consumer.Read()
	if err != nil {
		emptyMap := make(metricsjson.MetricsMapType)

		return &emptyMap, err
	}
	return mvList, nil
}

// Park data to temp dir
func (fa FileArchive) SaveAllMetrics(ctx context.Context, mvList *metricsjson.MetricsMapType) (err error) {
	producer, err := files.NewProducer((*fa.configuration).Get("STORE_FILE"))
	if err != nil {
		return err
	}
	defer producer.Close()

	producer.Write(mvList)

	return nil
}

// var s S.Storage = FileArchive{}