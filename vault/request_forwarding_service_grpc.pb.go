// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package vault

import (
	context "context"
	forwarding "github.com/morevault/vaultum/helper/forwarding"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RequestForwardingClient is the client API for RequestForwarding service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RequestForwardingClient interface {
	ForwardRequest(ctx context.Context, in *forwarding.Request, opts ...grpc.CallOption) (*forwarding.Response, error)
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error)
}

type requestForwardingClient struct {
	cc grpc.ClientConnInterface
}

func NewRequestForwardingClient(cc grpc.ClientConnInterface) RequestForwardingClient {
	return &requestForwardingClient{cc}
}

func (c *requestForwardingClient) ForwardRequest(ctx context.Context, in *forwarding.Request, opts ...grpc.CallOption) (*forwarding.Response, error) {
	out := new(forwarding.Response)
	err := c.cc.Invoke(ctx, "/vault.RequestForwarding/ForwardRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestForwardingClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error) {
	out := new(EchoReply)
	err := c.cc.Invoke(ctx, "/vault.RequestForwarding/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RequestForwardingServer is the server API for RequestForwarding service.
// All implementations must embed UnimplementedRequestForwardingServer
// for forward compatibility
type RequestForwardingServer interface {
	ForwardRequest(context.Context, *forwarding.Request) (*forwarding.Response, error)
	Echo(context.Context, *EchoRequest) (*EchoReply, error)
	mustEmbedUnimplementedRequestForwardingServer()
}

// UnimplementedRequestForwardingServer must be embedded to have forward compatible implementations.
type UnimplementedRequestForwardingServer struct {
}

func (UnimplementedRequestForwardingServer) ForwardRequest(context.Context, *forwarding.Request) (*forwarding.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForwardRequest not implemented")
}
func (UnimplementedRequestForwardingServer) Echo(context.Context, *EchoRequest) (*EchoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedRequestForwardingServer) mustEmbedUnimplementedRequestForwardingServer() {}

// UnsafeRequestForwardingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RequestForwardingServer will
// result in compilation errors.
type UnsafeRequestForwardingServer interface {
	mustEmbedUnimplementedRequestForwardingServer()
}

func RegisterRequestForwardingServer(s grpc.ServiceRegistrar, srv RequestForwardingServer) {
	s.RegisterService(&RequestForwarding_ServiceDesc, srv)
}

func _RequestForwarding_ForwardRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(forwarding.Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestForwardingServer).ForwardRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vault.RequestForwarding/ForwardRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestForwardingServer).ForwardRequest(ctx, req.(*forwarding.Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestForwarding_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestForwardingServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vault.RequestForwarding/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestForwardingServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RequestForwarding_ServiceDesc is the grpc.ServiceDesc for RequestForwarding service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RequestForwarding_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vault.RequestForwarding",
	HandlerType: (*RequestForwardingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ForwardRequest",
			Handler:    _RequestForwarding_ForwardRequest_Handler,
		},
		{
			MethodName: "Echo",
			Handler:    _RequestForwarding_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vault/request_forwarding_service.proto",
}
