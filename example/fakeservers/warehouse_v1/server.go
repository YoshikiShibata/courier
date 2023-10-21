// This code is auto-generated. DO NOT EDIT.

package fakewarehouse_v1

import (
	"context"

	"github.com/YoshikiShibata/courier/server/fakegrpc"
	"google.golang.org/grpc"

	v1 "github.com/YoshikiShibata/courier/example/api/warehouse_v1"
)

type fakeWarehouseServerImpl struct {
	v1.WarehouseServer
	srvStub *fakegrpc.GRPCServerStub
}

type WarehouseServer struct {
	impl    *fakeWarehouseServerImpl
	srvStub *fakegrpc.GRPCServerStub
}

func NewWarehouseServer(
	grpcServer *grpc.Server,
) *WarehouseServer {
	srvStub := &fakegrpc.GRPCServerStub{}
	srvImpl := &fakeWarehouseServerImpl{
		srvStub: srvStub,
	}
	v1.RegisterWarehouseServer(grpcServer, srvImpl)
	return &WarehouseServer{
		impl:    srvImpl,
		srvStub: srvStub,
	}
}

func (s *WarehouseServer) ClearAllResponses(
	tid string,
) {
	s.srvStub.ClearAllResponses(tid)
}

func (s *fakeWarehouseServerImpl) ListProductInventories(
	ctx context.Context,
	req *v1.ListProductInventoriesRequest,
) (*v1.ListProductInventoriesResponse, error) {
	return fakegrpc.HandleRequest[*v1.ListProductInventoriesResponse](s.srvStub, ctx, "ListProductInventories", req)
}

func (s *WarehouseServer) SetListProductInventoriesResponse(
	tid string,
	res *v1.ListProductInventoriesResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "ListProductInventories", res, err)
}

func (s *WarehouseServer) SetListProductInventoriesResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.ListProductInventoriesRequest) (*v1.ListProductInventoriesResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "ListProductInventories",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.ListProductInventoriesRequest))
			return res, err
		},
	)
}

func (s *fakeWarehouseServerImpl) ShipProduct(
	ctx context.Context,
	req *v1.ShipProductRequest,
) (*v1.ShipProductResponse, error) {
	return fakegrpc.HandleRequest[*v1.ShipProductResponse](s.srvStub, ctx, "ShipProduct", req)
}

func (s *WarehouseServer) SetShipProductResponse(
	tid string,
	res *v1.ShipProductResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "ShipProduct", res, err)
}

func (s *WarehouseServer) SetShipProductResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.ShipProductRequest) (*v1.ShipProductResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "ShipProduct",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.ShipProductRequest))
			return res, err
		},
	)
}
