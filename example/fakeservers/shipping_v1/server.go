// This code is auto-generated. DO NOT EDIT.

package fakeshipping_v1

import (
	"context"

	"google.golang.org/grpc"

	v1 "github.com/YoshikiShibata/courier/example/api/shipping_v1"
	"github.com/YoshikiShibata/courier/server"
)

type fakeShippingServerImpl struct {
	v1.ShippingServer
	srvStub *server.ServerStub
}

type ShippingServer struct {
	impl    *fakeShippingServerImpl
	srvStub *server.ServerStub
}

func NewShippingServer(
	grpcServer *grpc.Server,
) *ShippingServer {
	srvStub := &server.ServerStub{}
	srvImpl := &fakeShippingServerImpl{
		srvStub: srvStub,
	}
	v1.RegisterShippingServer(grpcServer, srvImpl)
	return &ShippingServer{
		impl:    srvImpl,
		srvStub: srvStub,
	}
}

func (s *ShippingServer) ClearAllResponses(
	tid string,
) {
	s.srvStub.ClearAllResponses(tid)
}

func (s *fakeShippingServerImpl) Create(
	ctx context.Context,
	req *v1.CreateShippingRequest,
) (*v1.CreateShippingResponse, error) {
	return server.HandleRequest[*v1.CreateShippingResponse](s.srvStub, ctx, "Create", req)
}

func (s *ShippingServer) SetCreateResponse(
	tid string,
	res *v1.CreateShippingResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "Create", res, err)
}

func (s *ShippingServer) SetCreateResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.CreateShippingRequest) (*v1.CreateShippingResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "Create",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.CreateShippingRequest))
			return res, err
		},
	)
}

func (s *fakeShippingServerImpl) Status(
	ctx context.Context,
	req *v1.ShippingStatusRequest,
) (*v1.ShippingStatusResponse, error) {
	return server.HandleRequest[*v1.ShippingStatusResponse](s.srvStub, ctx, "Status", req)
}

func (s *ShippingServer) SetStatusResponse(
	tid string,
	res *v1.ShippingStatusResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "Status", res, err)
}

func (s *ShippingServer) SetStatusResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.ShippingStatusRequest) (*v1.ShippingStatusResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "Status",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.ShippingStatusRequest))
			return res, err
		},
	)
}
