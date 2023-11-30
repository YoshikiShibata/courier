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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
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

func (s *shopServer) ListProductInventories(
	ctx context.Context,
	req *shop_v1.ListProductInventoriesRequest,
) (*shop_v1.ListProductInventoriesResponse, error) {
	numOfProducts := req.GetNumOfProducts()
	if numOfProducts == 0 {
		return nil, status.Error(codes.InvalidArgument, "num_of_products")
	}

	res, err := s.warehouseClient.ListProductInventories(ctx,
		&warehouse_v1.ListProductInventoriesRequest{
			NumOfProducts: numOfProducts,
			PageToken:     req.GetPageToken(),
		})

	switch status.Code(err) {
	case codes.OK:
		break
	case codes.Canceled,
		codes.DeadlineExceeded:
		return nil, err
	case codes.InvalidArgument:
		return nil, status.Error(codes.InvalidArgument, "page_token is invalid")
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	productInventories := make([]*shop_v1.ProductInventory, 0, len(res.ProductInventories))
	for _, p := range res.ProductInventories {
		productInventories = append(productInventories,
			&shop_v1.ProductInventory{
				Number:            p.Number,
				Name:              p.Name,
				Price:             p.Price,
				QuantityAvailable: p.QuantityAvailable,
			})
	}
	return &shop_v1.ListProductInventoriesResponse{
		ProductInventories: productInventories,
		NextPageToken:      res.NextPageToken,
	}, nil
}

func (s *shopServer) CreateOrder(
	ctx context.Context,
	req *shop_v1.CreateOrderRequest,
) (*shop_v1.CreateOrderResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Not Implemented Yet")
}

func (s *shopServer) GetShippingStatus(
	ctx context.Context,
	req *shop_v1.GetShippingStatusRequest,
) (*shop_v1.GetShippingStatusResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Not Implemented Yet")
}
