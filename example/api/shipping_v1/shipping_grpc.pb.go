// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: proto/shipping/v1/shipping.proto

package shipping_v1

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

// ShippingClient is the client API for Shipping service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShippingClient interface {
	// 商品の発送処理を行います。
	Create(ctx context.Context, in *CreateShippingRequest, opts ...grpc.CallOption) (*CreateShippingResponse, error)
	// 商品の発送ステータスを返します。
	Status(ctx context.Context, in *ShippingStatusRequest, opts ...grpc.CallOption) (*ShippingStatusResponse, error)
}

type shippingClient struct {
	cc grpc.ClientConnInterface
}

func NewShippingClient(cc grpc.ClientConnInterface) ShippingClient {
	return &shippingClient{cc}
}

func (c *shippingClient) Create(ctx context.Context, in *CreateShippingRequest, opts ...grpc.CallOption) (*CreateShippingResponse, error) {
	out := new(CreateShippingResponse)
	err := c.cc.Invoke(ctx, "/yoshikishibata.courier.example.api.shipping.v1.Shipping/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shippingClient) Status(ctx context.Context, in *ShippingStatusRequest, opts ...grpc.CallOption) (*ShippingStatusResponse, error) {
	out := new(ShippingStatusResponse)
	err := c.cc.Invoke(ctx, "/yoshikishibata.courier.example.api.shipping.v1.Shipping/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShippingServer is the server API for Shipping service.
// All implementations must embed UnimplementedShippingServer
// for forward compatibility
type ShippingServer interface {
	// 商品の発送処理を行います。
	Create(context.Context, *CreateShippingRequest) (*CreateShippingResponse, error)
	// 商品の発送ステータスを返します。
	Status(context.Context, *ShippingStatusRequest) (*ShippingStatusResponse, error)
	mustEmbedUnimplementedShippingServer()
}

// UnimplementedShippingServer must be embedded to have forward compatible implementations.
type UnimplementedShippingServer struct {
}

func (UnimplementedShippingServer) Create(context.Context, *CreateShippingRequest) (*CreateShippingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedShippingServer) Status(context.Context, *ShippingStatusRequest) (*ShippingStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedShippingServer) mustEmbedUnimplementedShippingServer() {}

// UnsafeShippingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShippingServer will
// result in compilation errors.
type UnsafeShippingServer interface {
	mustEmbedUnimplementedShippingServer()
}

func RegisterShippingServer(s grpc.ServiceRegistrar, srv ShippingServer) {
	s.RegisterService(&Shipping_ServiceDesc, srv)
}

func _Shipping_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateShippingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShippingServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yoshikishibata.courier.example.api.shipping.v1.Shipping/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShippingServer).Create(ctx, req.(*CreateShippingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shipping_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShippingStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShippingServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yoshikishibata.courier.example.api.shipping.v1.Shipping/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShippingServer).Status(ctx, req.(*ShippingStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shipping_ServiceDesc is the grpc.ServiceDesc for Shipping service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shipping_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "yoshikishibata.courier.example.api.shipping.v1.Shipping",
	HandlerType: (*ShippingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Shipping_Create_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Shipping_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shipping/v1/shipping.proto",
}
