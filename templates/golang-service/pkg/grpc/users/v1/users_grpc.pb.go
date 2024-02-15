// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package users

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

// UsersServiceClient is the client API for UsersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersServiceClient interface {
	CreateOrUpdateUser(ctx context.Context, in *CreateOrUpdateUserRequest, opts ...grpc.CallOption) (*CreateOrUpdateUserResponse, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
	GetAllUsers(ctx context.Context, in *GetAllUsersRequest, opts ...grpc.CallOption) (*GetAllUsersResponse, error)
	GetUserByName(ctx context.Context, in *GetUserByNameRequest, opts ...grpc.CallOption) (*GetUserByNameResponse, error)
}

type usersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersServiceClient(cc grpc.ClientConnInterface) UsersServiceClient {
	return &usersServiceClient{cc}
}

func (c *usersServiceClient) CreateOrUpdateUser(ctx context.Context, in *CreateOrUpdateUserRequest, opts ...grpc.CallOption) (*CreateOrUpdateUserResponse, error) {
	out := new(CreateOrUpdateUserResponse)
	err := c.cc.Invoke(ctx, "/users.v1.UsersService/CreateOrUpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	out := new(DeleteUserResponse)
	err := c.cc.Invoke(ctx, "/users.v1.UsersService/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) GetAllUsers(ctx context.Context, in *GetAllUsersRequest, opts ...grpc.CallOption) (*GetAllUsersResponse, error) {
	out := new(GetAllUsersResponse)
	err := c.cc.Invoke(ctx, "/users.v1.UsersService/GetAllUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) GetUserByName(ctx context.Context, in *GetUserByNameRequest, opts ...grpc.CallOption) (*GetUserByNameResponse, error) {
	out := new(GetUserByNameResponse)
	err := c.cc.Invoke(ctx, "/users.v1.UsersService/GetUserByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServiceServer is the server API for UsersService service.
// All implementations must embed UnimplementedUsersServiceServer
// for forward compatibility
type UsersServiceServer interface {
	CreateOrUpdateUser(context.Context, *CreateOrUpdateUserRequest) (*CreateOrUpdateUserResponse, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	GetAllUsers(context.Context, *GetAllUsersRequest) (*GetAllUsersResponse, error)
	GetUserByName(context.Context, *GetUserByNameRequest) (*GetUserByNameResponse, error)
	mustEmbedUnimplementedUsersServiceServer()
}

// UnimplementedUsersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUsersServiceServer struct {
}

func (UnimplementedUsersServiceServer) CreateOrUpdateUser(context.Context, *CreateOrUpdateUserRequest) (*CreateOrUpdateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateUser not implemented")
}
func (UnimplementedUsersServiceServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUsersServiceServer) GetAllUsers(context.Context, *GetAllUsersRequest) (*GetAllUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllUsers not implemented")
}
func (UnimplementedUsersServiceServer) GetUserByName(context.Context, *GetUserByNameRequest) (*GetUserByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByName not implemented")
}
func (UnimplementedUsersServiceServer) mustEmbedUnimplementedUsersServiceServer() {}

// UnsafeUsersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServiceServer will
// result in compilation errors.
type UnsafeUsersServiceServer interface {
	mustEmbedUnimplementedUsersServiceServer()
}

func RegisterUsersServiceServer(s grpc.ServiceRegistrar, srv UsersServiceServer) {
	s.RegisterService(&UsersService_ServiceDesc, srv)
}

func _UsersService_CreateOrUpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrUpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).CreateOrUpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.v1.UsersService/CreateOrUpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).CreateOrUpdateUser(ctx, req.(*CreateOrUpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.v1.UsersService/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_GetAllUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetAllUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.v1.UsersService/GetAllUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetAllUsers(ctx, req.(*GetAllUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_GetUserByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetUserByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.v1.UsersService/GetUserByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetUserByName(ctx, req.(*GetUserByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UsersService_ServiceDesc is the grpc.ServiceDesc for UsersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UsersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "users.v1.UsersService",
	HandlerType: (*UsersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrUpdateUser",
			Handler:    _UsersService_CreateOrUpdateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UsersService_DeleteUser_Handler,
		},
		{
			MethodName: "GetAllUsers",
			Handler:    _UsersService_GetAllUsers_Handler,
		},
		{
			MethodName: "GetUserByName",
			Handler:    _UsersService_GetUserByName_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "users/v1/users.proto",
}
