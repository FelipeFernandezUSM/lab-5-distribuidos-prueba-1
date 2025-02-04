// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: fulcrum.proto

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

// FulcrumClient is the client API for Fulcrum service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FulcrumClient interface {
	ApplyCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error)
	ProcessCommandMessage(ctx context.Context, in *Mensaje, opts ...grpc.CallOption) (*Notificacion, error)
	ApplyPropagation(ctx context.Context, in *Propagation, opts ...grpc.CallOption) (*PropagationResponse, error)
}

type fulcrumClient struct {
	cc grpc.ClientConnInterface
}

func NewFulcrumClient(cc grpc.ClientConnInterface) FulcrumClient {
	return &fulcrumClient{cc}
}

func (c *fulcrumClient) ApplyCommand(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error) {
	out := new(CommandResponse)
	err := c.cc.Invoke(ctx, "/fulcrum.Fulcrum/ApplyCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fulcrumClient) ProcessCommandMessage(ctx context.Context, in *Mensaje, opts ...grpc.CallOption) (*Notificacion, error) {
	out := new(Notificacion)
	err := c.cc.Invoke(ctx, "/fulcrum.Fulcrum/ProcessCommandMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fulcrumClient) ApplyPropagation(ctx context.Context, in *Propagation, opts ...grpc.CallOption) (*PropagationResponse, error) {
	out := new(PropagationResponse)
	err := c.cc.Invoke(ctx, "/fulcrum.Fulcrum/ApplyPropagation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FulcrumServer is the server API for Fulcrum service.
// All implementations must embed UnimplementedFulcrumServer
// for forward compatibility
type FulcrumServer interface {
	ApplyCommand(context.Context, *CommandRequest) (*CommandResponse, error)
	ProcessCommandMessage(context.Context, *Mensaje) (*Notificacion, error)
	ApplyPropagation(context.Context, *Propagation) (*PropagationResponse, error)
	mustEmbedUnimplementedFulcrumServer()
}

// UnimplementedFulcrumServer must be embedded to have forward compatible implementations.
type UnimplementedFulcrumServer struct {
}

func (UnimplementedFulcrumServer) ApplyCommand(context.Context, *CommandRequest) (*CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApplyCommand not implemented")
}
func (UnimplementedFulcrumServer) ProcessCommandMessage(context.Context, *Mensaje) (*Notificacion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessCommandMessage not implemented")
}
func (UnimplementedFulcrumServer) ApplyPropagation(context.Context, *Propagation) (*PropagationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApplyPropagation not implemented")
}
func (UnimplementedFulcrumServer) mustEmbedUnimplementedFulcrumServer() {}

// UnsafeFulcrumServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FulcrumServer will
// result in compilation errors.
type UnsafeFulcrumServer interface {
	mustEmbedUnimplementedFulcrumServer()
}

func RegisterFulcrumServer(s grpc.ServiceRegistrar, srv FulcrumServer) {
	s.RegisterService(&Fulcrum_ServiceDesc, srv)
}

func _Fulcrum_ApplyCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FulcrumServer).ApplyCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fulcrum.Fulcrum/ApplyCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FulcrumServer).ApplyCommand(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fulcrum_ProcessCommandMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Mensaje)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FulcrumServer).ProcessCommandMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fulcrum.Fulcrum/ProcessCommandMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FulcrumServer).ProcessCommandMessage(ctx, req.(*Mensaje))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fulcrum_ApplyPropagation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Propagation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FulcrumServer).ApplyPropagation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fulcrum.Fulcrum/ApplyPropagation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FulcrumServer).ApplyPropagation(ctx, req.(*Propagation))
	}
	return interceptor(ctx, in, info, handler)
}

// Fulcrum_ServiceDesc is the grpc.ServiceDesc for Fulcrum service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Fulcrum_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fulcrum.Fulcrum",
	HandlerType: (*FulcrumServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ApplyCommand",
			Handler:    _Fulcrum_ApplyCommand_Handler,
		},
		{
			MethodName: "ProcessCommandMessage",
			Handler:    _Fulcrum_ProcessCommandMessage_Handler,
		},
		{
			MethodName: "ApplyPropagation",
			Handler:    _Fulcrum_ApplyPropagation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fulcrum.proto",
}
