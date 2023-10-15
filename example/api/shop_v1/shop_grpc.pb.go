// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: proto/shop/v1/shop.proto

package shop_v1

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

// ShopClient is the client API for Shop service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShopClient interface {
	ListProducts(ctx context.Context, in *ListProductsRequest, opts ...grpc.CallOption) (*ListProductsResponse, error)
	GetProduct(ctx context.Context, in *GetProductRequest, opts ...grpc.CallOption) (*GetProductResponse, error)
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
	GetOrderStatus(ctx context.Context, in *GetOrderStatusRequest, opts ...grpc.CallOption) (*GetOrderStatusResponse, error)
}

type shopClient struct {
	cc grpc.ClientConnInterface
}

func NewShopClient(cc grpc.ClientConnInterface) ShopClient {
	return &shopClient{cc}
}

func (c *shopClient) ListProducts(ctx context.Context, in *ListProductsRequest, opts ...grpc.CallOption) (*ListProductsResponse, error) {
	out := new(ListProductsResponse)
	err := c.cc.Invoke(ctx, "/yoshikishibata.courier.example.api.shop.v1.Shop/ListProducts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shopClient) GetProduct(ctx context.Context, in *GetProductRequest, opts ...grpc.CallOption) (*GetProductResponse, error) {
	out := new(GetProductResponse)
	err := c.cc.Invoke(ctx, "/yoshikishibata.courier.example.api.shop.v1.Shop/GetProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shopClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/yoshikishibata.courier.example.api.shop.v1.Shop/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shopClient) GetOrderStatus(ctx context.Context, in *GetOrderStatusRequest, opts ...grpc.CallOption) (*GetOrderStatusResponse, error) {
	out := new(GetOrderStatusResponse)
	err := c.cc.Invoke(ctx, "/yoshikishibata.courier.example.api.shop.v1.Shop/GetOrderStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShopServer is the server API for Shop service.
// All implementations must embed UnimplementedShopServer
// for forward compatibility
type ShopServer interface {
	ListProducts(context.Context, *ListProductsRequest) (*ListProductsResponse, error)
	GetProduct(context.Context, *GetProductRequest) (*GetProductResponse, error)
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrderStatus(context.Context, *GetOrderStatusRequest) (*GetOrderStatusResponse, error)
	mustEmbedUnimplementedShopServer()
}

// UnimplementedShopServer must be embedded to have forward compatible implementations.
type UnimplementedShopServer struct {
}

func (UnimplementedShopServer) ListProducts(context.Context, *ListProductsRequest) (*ListProductsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProducts not implemented")
}
func (UnimplementedShopServer) GetProduct(context.Context, *GetProductRequest) (*GetProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProduct not implemented")
}
func (UnimplementedShopServer) CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedShopServer) GetOrderStatus(context.Context, *GetOrderStatusRequest) (*GetOrderStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderStatus not implemented")
}
func (UnimplementedShopServer) mustEmbedUnimplementedShopServer() {}

// UnsafeShopServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShopServer will
// result in compilation errors.
type UnsafeShopServer interface {
	mustEmbedUnimplementedShopServer()
}

func RegisterShopServer(s grpc.ServiceRegistrar, srv ShopServer) {
	s.RegisterService(&Shop_ServiceDesc, srv)
}

func _Shop_ListProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProductsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShopServer).ListProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yoshikishibata.courier.example.api.shop.v1.Shop/ListProducts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShopServer).ListProducts(ctx, req.(*ListProductsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shop_GetProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShopServer).GetProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yoshikishibata.courier.example.api.shop.v1.Shop/GetProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShopServer).GetProduct(ctx, req.(*GetProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shop_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShopServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yoshikishibata.courier.example.api.shop.v1.Shop/CreateOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShopServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shop_GetOrderStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrderStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShopServer).GetOrderStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/yoshikishibata.courier.example.api.shop.v1.Shop/GetOrderStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShopServer).GetOrderStatus(ctx, req.(*GetOrderStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shop_ServiceDesc is the grpc.ServiceDesc for Shop service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shop_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "yoshikishibata.courier.example.api.shop.v1.Shop",
	HandlerType: (*ShopServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListProducts",
			Handler:    _Shop_ListProducts_Handler,
		},
		{
			MethodName: "GetProduct",
			Handler:    _Shop_GetProduct_Handler,
		},
		{
			MethodName: "CreateOrder",
			Handler:    _Shop_CreateOrder_Handler,
		},
		{
			MethodName: "GetOrderStatus",
			Handler:    _Shop_GetOrderStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shop/v1/shop.proto",
}
