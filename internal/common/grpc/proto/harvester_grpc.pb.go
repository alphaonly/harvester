// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: proto/harvester.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Service_AddMetric_FullMethodName      = "/harvester.Service/AddMetric"
	Service_AddMetricMulti_FullMethodName = "/harvester.Service/AddMetricMulti"
	Service_GetMetric_FullMethodName      = "/harvester.Service/GetMetric"
	Service_GetMetricMulti_FullMethodName = "/harvester.Service/GetMetricMulti"
)

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	AddMetric(ctx context.Context, in *AddMetricRequest, opts ...grpc.CallOption) (*AddMetricResponse, error)
	AddMetricMulti(ctx context.Context, opts ...grpc.CallOption) (Service_AddMetricMultiClient, error)
	GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
	GetMetricMulti(ctx context.Context, opts ...grpc.CallOption) (Service_GetMetricMultiClient, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) AddMetric(ctx context.Context, in *AddMetricRequest, opts ...grpc.CallOption) (*AddMetricResponse, error) {
	out := new(AddMetricResponse)
	err := c.cc.Invoke(ctx, Service_AddMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) AddMetricMulti(ctx context.Context, opts ...grpc.CallOption) (Service_AddMetricMultiClient, error) {
	stream, err := c.cc.NewStream(ctx, &Service_ServiceDesc.Streams[0], Service_AddMetricMulti_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceAddMetricMultiClient{stream}
	return x, nil
}

type Service_AddMetricMultiClient interface {
	Send(*AddMetricRequest) error
	Recv() (*AddMetricResponse, error)
	grpc.ClientStream
}

type serviceAddMetricMultiClient struct {
	grpc.ClientStream
}

func (x *serviceAddMetricMultiClient) Send(m *AddMetricRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serviceAddMetricMultiClient) Recv() (*AddMetricResponse, error) {
	m := new(AddMetricResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *serviceClient) GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, Service_GetMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetMetricMulti(ctx context.Context, opts ...grpc.CallOption) (Service_GetMetricMultiClient, error) {
	stream, err := c.cc.NewStream(ctx, &Service_ServiceDesc.Streams[1], Service_GetMetricMulti_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceGetMetricMultiClient{stream}
	return x, nil
}

type Service_GetMetricMultiClient interface {
	Send(*GetMetricRequest) error
	Recv() (*GetMetricResponse, error)
	grpc.ClientStream
}

type serviceGetMetricMultiClient struct {
	grpc.ClientStream
}

func (x *serviceGetMetricMultiClient) Send(m *GetMetricRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serviceGetMetricMultiClient) Recv() (*GetMetricResponse, error) {
	m := new(GetMetricResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	AddMetric(context.Context, *AddMetricRequest) (*AddMetricResponse, error)
	AddMetricMulti(Service_AddMetricMultiServer) error
	GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
	GetMetricMulti(Service_GetMetricMultiServer) error
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) AddMetric(context.Context, *AddMetricRequest) (*AddMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMetric not implemented")
}
func (UnimplementedServiceServer) AddMetricMulti(Service_AddMetricMultiServer) error {
	return status.Errorf(codes.Unimplemented, "method AddMetricMulti not implemented")
}
func (UnimplementedServiceServer) GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedServiceServer) GetMetricMulti(Service_GetMetricMultiServer) error {
	return status.Errorf(codes.Unimplemented, "method GetMetricMulti not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_AddMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).AddMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_AddMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).AddMetric(ctx, req.(*AddMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_AddMetricMulti_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServiceServer).AddMetricMulti(&serviceAddMetricMultiServer{stream})
}

type Service_AddMetricMultiServer interface {
	Send(*AddMetricResponse) error
	Recv() (*AddMetricRequest, error)
	grpc.ServerStream
}

type serviceAddMetricMultiServer struct {
	grpc.ServerStream
}

func (x *serviceAddMetricMultiServer) Send(m *AddMetricResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serviceAddMetricMultiServer) Recv() (*AddMetricRequest, error) {
	m := new(AddMetricRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Service_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetMetric(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetMetricMulti_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServiceServer).GetMetricMulti(&serviceGetMetricMultiServer{stream})
}

type Service_GetMetricMultiServer interface {
	Send(*GetMetricResponse) error
	Recv() (*GetMetricRequest, error)
	grpc.ServerStream
}

type serviceGetMetricMultiServer struct {
	grpc.ServerStream
}

func (x *serviceGetMetricMultiServer) Send(m *GetMetricResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serviceGetMetricMultiServer) Recv() (*GetMetricRequest, error) {
	m := new(GetMetricRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "harvester.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddMetric",
			Handler:    _Service_AddMetric_Handler,
		},
		{
			MethodName: "GetMetric",
			Handler:    _Service_GetMetric_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AddMetricMulti",
			Handler:       _Service_AddMetricMulti_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "GetMetricMulti",
			Handler:       _Service_GetMetricMulti_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/harvester.proto",
}