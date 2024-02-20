// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: proto/oidc/user/user.proto

package user

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
	UserService_CreateUser_FullMethodName  = "/proto.oidc.user.UserService/CreateUser"
	UserService_GetUser_FullMethodName     = "/proto.oidc.user.UserService/GetUser"
	UserService_DeleteUser_FullMethodName  = "/proto.oidc.user.UserService/DeleteUser"
	UserService_UpdateUser_FullMethodName  = "/proto.oidc.user.UserService/UpdateUser"
	UserService_ListUser_FullMethodName    = "/proto.oidc.user.UserService/ListUser"
	UserService_LinkUsers_FullMethodName   = "/proto.oidc.user.UserService/LinkUsers"
	UserService_UnlinkUsers_FullMethodName = "/proto.oidc.user.UserService/UnlinkUsers"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	// Create user
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	// Get user
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	// Delete user
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
	// Update user
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error)
	// List users
	ListUser(ctx context.Context, in *ListUserRequest, opts ...grpc.CallOption) (*ListUserResponse, error)
	LinkUsers(ctx context.Context, in *LinkUsersRequest, opts ...grpc.CallOption) (*LinkUsersResponse, error)
	UnlinkUsers(ctx context.Context, in *UnlinkUsersRequest, opts ...grpc.CallOption) (*UnlinkUsersResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, UserService_CreateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, UserService_GetUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	out := new(DeleteUserResponse)
	err := c.cc.Invoke(ctx, UserService_DeleteUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error) {
	out := new(UpdateUserResponse)
	err := c.cc.Invoke(ctx, UserService_UpdateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ListUser(ctx context.Context, in *ListUserRequest, opts ...grpc.CallOption) (*ListUserResponse, error) {
	out := new(ListUserResponse)
	err := c.cc.Invoke(ctx, UserService_ListUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) LinkUsers(ctx context.Context, in *LinkUsersRequest, opts ...grpc.CallOption) (*LinkUsersResponse, error) {
	out := new(LinkUsersResponse)
	err := c.cc.Invoke(ctx, UserService_LinkUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UnlinkUsers(ctx context.Context, in *UnlinkUsersRequest, opts ...grpc.CallOption) (*UnlinkUsersResponse, error) {
	out := new(UnlinkUsersResponse)
	err := c.cc.Invoke(ctx, UserService_UnlinkUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	// Create user
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	// Get user
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	// Delete user
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	// Update user
	UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserResponse, error)
	// List users
	ListUser(context.Context, *ListUserRequest) (*ListUserResponse, error)
	LinkUsers(context.Context, *LinkUsersRequest) (*LinkUsersResponse, error)
	UnlinkUsers(context.Context, *UnlinkUsersRequest) (*UnlinkUsersResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserServiceServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserServiceServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUserServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserServiceServer) ListUser(context.Context, *ListUserRequest) (*ListUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUser not implemented")
}
func (UnimplementedUserServiceServer) LinkUsers(context.Context, *LinkUsersRequest) (*LinkUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LinkUsers not implemented")
}
func (UnimplementedUserServiceServer) UnlinkUsers(context.Context, *UnlinkUsersRequest) (*UnlinkUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnlinkUsers not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_DeleteUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_UpdateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ListUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ListUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_ListUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ListUser(ctx, req.(*ListUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_LinkUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LinkUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).LinkUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_LinkUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).LinkUsers(ctx, req.(*LinkUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UnlinkUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnlinkUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UnlinkUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_UnlinkUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UnlinkUsers(ctx, req.(*UnlinkUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.oidc.user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _UserService_GetUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UserService_DeleteUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserService_UpdateUser_Handler,
		},
		{
			MethodName: "ListUser",
			Handler:    _UserService_ListUser_Handler,
		},
		{
			MethodName: "LinkUsers",
			Handler:    _UserService_LinkUsers_Handler,
		},
		{
			MethodName: "UnlinkUsers",
			Handler:    _UserService_UnlinkUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/oidc/user/user.proto",
}
