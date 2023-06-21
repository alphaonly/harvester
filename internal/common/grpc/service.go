package service

import (
	"context"
	"fmt"
	"io"
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
	storage storage.Storage      		//a storage to receive data
}

//NewGRPCService - a factory to Metric gRPC server service, receives used storage implementation
func NewGRPCService(storage storage.Storage) pb.ServiceServer {
	return &GRPCService{storage: storage}
}
//AddMetric - adds inbound metric data to storage
func (s *GRPCService) AddMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var response pb.AddMetricResponse

	metric := schema.Metrics{
		ID:    in.Metric.Name,
		MType: common.ConvertGrpcType(in.Metric.Type),
		Delta: &in.Metric.Counter,
		Value: &in.Metric.Gauge,
	}
	//overwrite metric value
	s.metrics.Store(metric.ID, metric)

	return &response, nil
}
//AddMetricMulti - adds metric data from a stream to storage 
func (s *GRPCService) AddMetricMulti(in pb.Service_AddMetricMultiServer) error {
	var (
		request  = new(pb.AddMetricRequest)
		response = new(pb.AddMetricResponse)
		err      error
	)

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

		metric := schema.Metrics{
			ID:    request.Metric.Name,
			MType: common.ConvertGrpcType(request.Metric.Type),
			Delta: &request.Metric.Counter,
			Value: &request.Metric.Gauge,
		}
		//save metric
		s.metrics.Store(metric.ID, metric)
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
		val, ok := s.metrics.Load(request.Name)
		if !ok {
			err = fmt.Errorf("%w:%v", common.ErrNoMetricInStorage, request.Name)
			response.Error = err.Error()
			logging.LogPrintln(err)
			return nil
		}
		//recognize a gotten metric
		metric, ok := val.(schema.Metrics)
		if !ok {
			logging.LogFatal(common.ErrInappropriateType)
		}
		//send metric data in response
		response.Metric.Name = metric.ID
		response.Metric.Type = common.ConvertMetricType(metric.MType)
		if metric.Delta != nil {
			response.Metric.Counter = *metric.Delta
		}
		if metric.Value != nil {
			response.Metric.Gauge = *metric.Value
		}
		err = stream.Send(response)
		logging.LogFatal(err)
	}
	return err
}
