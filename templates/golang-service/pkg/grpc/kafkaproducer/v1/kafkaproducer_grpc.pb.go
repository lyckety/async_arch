// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package kafkaproducer

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

// KafkaProducerServiceClient is the client API for KafkaProducerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KafkaProducerServiceClient interface {
	Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error)
}

type kafkaProducerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKafkaProducerServiceClient(cc grpc.ClientConnInterface) KafkaProducerServiceClient {
	return &kafkaProducerServiceClient{cc}
}

func (c *kafkaProducerServiceClient) Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error) {
	out := new(SendResponse)
	err := c.cc.Invoke(ctx, "/kafkaproducer.v1.KafkaProducerService/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KafkaProducerServiceServer is the server API for KafkaProducerService service.
// All implementations must embed UnimplementedKafkaProducerServiceServer
// for forward compatibility
type KafkaProducerServiceServer interface {
	Send(context.Context, *SendRequest) (*SendResponse, error)
	mustEmbedUnimplementedKafkaProducerServiceServer()
}

// UnimplementedKafkaProducerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKafkaProducerServiceServer struct {
}

func (UnimplementedKafkaProducerServiceServer) Send(context.Context, *SendRequest) (*SendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedKafkaProducerServiceServer) mustEmbedUnimplementedKafkaProducerServiceServer() {}

// UnsafeKafkaProducerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KafkaProducerServiceServer will
// result in compilation errors.
type UnsafeKafkaProducerServiceServer interface {
	mustEmbedUnimplementedKafkaProducerServiceServer()
}

func RegisterKafkaProducerServiceServer(s grpc.ServiceRegistrar, srv KafkaProducerServiceServer) {
	s.RegisterService(&KafkaProducerService_ServiceDesc, srv)
}

func _KafkaProducerService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KafkaProducerServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kafkaproducer.v1.KafkaProducerService/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KafkaProducerServiceServer).Send(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KafkaProducerService_ServiceDesc is the grpc.ServiceDesc for KafkaProducerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KafkaProducerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kafkaproducer.v1.KafkaProducerService",
	HandlerType: (*KafkaProducerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _KafkaProducerService_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kafkaproducer/v1/kafkaproducer.proto",
}
