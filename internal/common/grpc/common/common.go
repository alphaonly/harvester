package common

import (
	"errors"

	pb "github.com/alphaonly/harvester/internal/common/grpc/proto"
	"github.com/alphaonly/harvester/internal/common/logging"
	"github.com/alphaonly/harvester/internal/schema"
)

var (
	ErrNoMetricName      = errors.New("empty Name field")
	ErrNoMetricInStorage = errors.New("no metric in storage")
	ErrInappropriateType = errors.New("wrong type assertion")
	ErrBadMetricType     = errors.New("bad metric type")
)

// ConvertMetricType - converts from internal metric type value to grpc generated type value
func ConvertMetricType(typ string) pb.Metric_Type {
	switch typ {
	case schema.COUNTER_TYPE:
		return *pb.Metric_COUNTER.Enum()

	case schema.GAUGE_TYPE:
		return *pb.Metric_GAUGE.Enum()
	default:
		logging.LogFatal(ErrBadMetricType)
		return 0
	}
}

// ConvertGrpcType - converts from grpc generated type value to internal metric type value
func ConvertGrpcType(typ pb.Metric_Type) string {
	switch typ {
	case *pb.Metric_COUNTER.Enum():
		return schema.COUNTER_TYPE

	case *pb.Metric_GAUGE.Enum():
		return schema.GAUGE_TYPE
	default:
		logging.LogFatal(ErrBadMetricType)
		return ""
	}
}
