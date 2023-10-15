// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package server

import (
	"context"
	"fmt"
	"net"

	"github.com/YoshikiShibata/courier/example/api/shipping_v1"
	"github.com/YoshikiShibata/courier/example/api/shop_v1"
	"github.com/YoshikiShibata/courier/example/api/warehouse_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCServer struct {
	server *grpc.Server
	port   int
}

func NewGRPCServer(
	port int,
	shippingClient shipping_v1.ShippingClient,
	warehouseClient warehouse_v1.WarehouseClient,
	grpcOpts ...grpc.ServerOption,
) (
	*GRPCServer,
	error,
) {
	gsrv := grpc.NewServer(grpcOpts...)
	hsrv := health.NewServer()
	hsrv.SetServingStatus("Shop", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(gsrv, hsrv)

	srv, err := newShopServer(shippingClient, warehouseClient)
	if err != nil {
		return nil, err
	}
	shop_v1.RegisterShopServer(gsrv, srv)

	return &GRPCServer{
		server: gsrv,
		port:   port,
	}, nil
}

func (s *GRPCServer) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("net.List failed: %w", err)
	}
	return s.server.Serve(l)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

var _ shop_v1.ShopServer = (*shopServer)(nil)

type shopServer struct {
	shop_v1.UnimplementedShopServer

	shippingClient  shipping_v1.ShippingClient
	warehouseClient warehouse_v1.WarehouseClient
}

func newShopServer(
	shippingClient shipping_v1.ShippingClient,
	warehouseClient warehouse_v1.WarehouseClient,
) (srv *shopServer, err error) {
	srv = &shopServer{
		shippingClient:  shippingClient,
		warehouseClient: warehouseClient,
	}

	return srv, nil
}

func (s *shopServer) ListProducts(
	ctx context.Context,
	req *shop_v1.ListProductRequest,
) (*shop_v1.ListProductResponse, error) {

	return nil, nil
}

func (s *shopServer) GetProduct(
	ctx context.Context,
	req *shop_v1.GetProductRequest,
) (*shop_v1.GetProductResponse, error) {

	return nil, nil
}

func (s *shopServer) CreateOrder(
	ctx context.Context,
	req *shop_v1.CreateOrderRequest,
) (*shop_v1.CreateOrderResponse, error) {
	return nil, nil
}

func (s *shopServer) GetOrderStatus(
	ctx context.Context,
	req *shop_v1.GetOrderStatusRequest,
) (*shop_v1.GetOrderStatusResponse, error) {
	return nil, nil
}
