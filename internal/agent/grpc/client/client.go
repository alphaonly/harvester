package client

import (
	pb "github.com/alphaonly/harvester/internal/common/grpc/proto"
	"github.com/alphaonly/harvester/internal/common/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Client pb.ServiceClient
	conn   *grpc.ClientConn
}

func (g GRPCClient) Close() {
	g.conn.Close()
}

func NewGRPCClient(address string) *GRPCClient {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	logging.LogFatal(err)

	return &GRPCClient{
		Client: pb.NewServiceClient(conn),
		conn:   conn}
}
