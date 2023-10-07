// This code is auto-generated. DO NOT EDIT.

package fakelogv1

import (
	"context"

	"google.golang.org/grpc"

	"github.com/YoshikiShibata/courier/server"
	log_v1 "github.com/YoshikiShibata/example/api/v1"
)

type fakeLogServerImpl struct {
	log_v1.LogServer
	srvStub *server.ServerStub
}

type LogServer struct {
	impl    *fakeLogServerImpl
	srvStub *server.ServerStub
}

func NewLogServer(
	grpcServer *grpc.Server,
) *LogServer {
	srvStub := &server.ServerStub{}
	srvImpl := &fakeLogServerImpl{
		srvStub: srvStub,
	}
	log_v1.RegisterLogServer(grpcServer, srvImpl)
	return &LogServer{
		impl:    srvImpl,
		srvStub: srvStub,
	}
}

func (s *LogServer) ClearAllResponses(
	tid string,
) {
	s.srvStub.ClearAllResponses(tid)
}

func (s *fakeLogServerImpl) Consume(
	ctx context.Context,
	req *log_v1.ConsumeRequest,
) (*log_v1.ConsumeResponse, error) {
	return server.HandleRequest[*log_v1.ConsumeResponse](s.srvStub, ctx, "Consume", req)
}

func (s *LogServer) SetConsumeResponse(
	tid string,
	res *log_v1.ConsumeResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "Consume", res, err)
}

func (s *LogServer) SetConsumeResponseCreator(
	tid string,
	creator func(ctx context.Context, req *log_v1.ConsumeRequest) (*log_v1.ConsumeResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "Consume",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*log_v1.ConsumeRequest))
			return res, err
		},
	)
}

func (s *fakeLogServerImpl) Produce(
	ctx context.Context,
	req *log_v1.ProduceRequest,
) (*log_v1.ProduceResponse, error) {
	return server.HandleRequest[*log_v1.ProduceResponse](s.srvStub, ctx, "Produce", req)
}

func (s *LogServer) SetProduceResponse(
	tid string,
	res *log_v1.ProduceResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "Produce", res, err)
}

func (s *LogServer) SetProduceResponseCreator(
	tid string,
	creator func(ctx context.Context, req *log_v1.ProduceRequest) (*log_v1.ProduceResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "Produce",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*log_v1.ProduceRequest))
			return res, err
		},
	)
}
