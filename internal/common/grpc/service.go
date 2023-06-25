package service

import (
	"context"
	"fmt"
	mv "github.com/alphaonly/harvester/internal/server/metricvaluei"
	"io"
	"log"
	"sync"

	common "github.com/alphaonly/harvester/internal/common/grpc/common"
	pb "github.com/alphaonly/harvester/internal/common/grpc/proto"
	"github.com/alphaonly/harvester/internal/common/logging"
	"github.com/alphaonly/harvester/internal/schema"
	storage "github.com/alphaonly/harvester/internal/server/storage/interfaces"
)

type GRPCService struct {
	pb.UnimplementedServiceServer

	metrics sync.Map
	storage storage.Storage //a storage to receive data
}

// NewGRPCService - a factory to Metric gRPC server service, receives used storage implementation
func NewGRPCService(storage storage.Storage) pb.ServiceServer {
	return &GRPCService{storage: storage}
}

// AddMetric - adds inbound metric data to storage
func (s *GRPCService) AddMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var response pb.AddMetricResponse

	metric := schema.Metrics{
		ID:    in.Metric.Name,
		MType: common.ConvertGrpcType(in.Metric.Type),
		Delta: &in.Metric.Counter,
		Value: &in.Metric.Gauge,
	}
	//overwrite metric value
	//s.metrics.Store(metric.ID, metric)

	var v mv.MetricValue

	switch metric.MType {
	case schema.GAUGE_TYPE:
		v = mv.NewGaugeMetric(*metric.Value)
	case schema.COUNTER_TYPE:
		v = mv.NewCounterMetric(*metric.Delta)
	}
	err := s.storage.SaveMetric(ctx, metric.ID, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("metric %v saved through gRPC", metric)
	return &response, err
}

// AddMetricMulti - adds metric data from a stream to storage
func (s *GRPCService) AddMetricMulti(in pb.Service_AddMetricMultiServer) error {
	var (
		request  = new(pb.AddMetricRequest)
		response = new(pb.AddMetricResponse)
		err      error
	)

	//Loop at inbound data
	for {
		//Receive request data
		request, err = in.Recv()
		if err == io.EOF {
			break
		}
		logging.LogFatal(err)

		if request.Metric.Name == "" {
			err = fmt.Errorf("%w:%v", common.ErrNoMetricName, request.Metric.Name)
			logging.LogPrintln(err)
			response.Error = err.Error()
			return err
		}
		//add metric to temp map
		switch request.Metric.Type {
		case pb.Metric_GAUGE:
			err = s.storage.SaveMetric(in.Context(), request.Metric.Name, mv.NewGaugeMetric(request.Metric.GetGauge()))
		case pb.Metric_COUNTER:
			err = s.storage.SaveMetric(in.Context(), request.Metric.Name, mv.NewCounterMetric(request.Metric.GetCounter()))
		}
		//send response
		switch {
		case err == nil:
			log.Printf(" metric %v saved through rRPC", request.Metric.Name)
		case err != nil:
			log.Printf("Unable to save  metric through rRPC due to storage problem: %v", err)
			response.Error = err.Error()
		}
		err = in.Send(response)
		logging.LogFatal(err)
	}
	return err
}
func (s *GRPCService) GetMetric(context.Context, *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	//I do not need it, left for a while
	return nil, nil
}
func (s *GRPCService) GetMetricMulti(stream pb.Service_GetMetricMultiServer) error {

	var (
		request  = new(pb.GetMetricRequest)
		response = new(pb.GetMetricResponse)
		err      error
	)

	for {
		//Receive request data
		request, err = stream.Recv()
		if err == io.EOF {
			break
		}
		logging.LogFatal(err)
		if request.Name == "" {
			err = fmt.Errorf("%w:%v", common.ErrNoMetricName, request.Name)
			logging.LogPrintln(err)
			response.Error = err.Error()
			return err
		}
		//get a value by the metric name
		mv, err := s.storage.GetMetric(stream.Context(), request.Name, common.ConvertGrpcType(pb.Metric_Type(request.GetType())))
		if err != nil {
			response.Error = err.Error()
			logging.LogPrintln(err)
			return err
		}

		//send metric data in response
		response.Metric.Name = request.Name
		response.Metric.Type = pb.Metric_Type(request.Type)
		switch request.Type {
		case pb.GetMetricRequest_COUNTER:
			response.Metric.Counter = mv.GetInternalValue().(int64)
		case pb.GetMetricRequest_GAUGE:
			response.Metric.Gauge = mv.GetInternalValue().(float64)
		}

		err = stream.Send(response)
		logging.LogFatal(err)
	}
	return err
}
