// This code is auto-generated. DO NOT EDIT.

package publisher_v1

import (
	"context"

	"github.com/YoshikiShibata/courier/server/fakegrpc"
	"google.golang.org/grpc"

	v1 "cloud.google.com/go/pubsub/apiv1/pubsubpb"
)

type fakePublisherServerImpl struct {
	v1.PublisherServer
	srvStub *fakegrpc.GRPCServerStub
}

type PublisherServer struct {
	impl    *fakePublisherServerImpl
	srvStub *fakegrpc.GRPCServerStub
}

func NewPublisherServer(
	grpcServer *grpc.Server,
) *PublisherServer {
	srvStub := &fakegrpc.GRPCServerStub{}
	srvImpl := &fakePublisherServerImpl{
		srvStub: srvStub,
	}
	v1.RegisterPublisherServer(grpcServer, srvImpl)
	return &PublisherServer{
		impl:    srvImpl,
		srvStub: srvStub,
	}
}

func (s *PublisherServer) ClearAllResponses(
	tid string,
) {
	s.srvStub.ClearAllResponses(tid)
}

func (s *fakePublisherServerImpl) CreateTopic(
	ctx context.Context,
	req *v1.Topic,
) (*v1.Topic, error) {
	return fakegrpc.HandleRequest[*v1.Topic](s.srvStub, ctx, "CreateTopic", req)
}

func (s *PublisherServer) SetCreateTopicResponse(
	tid string,
	res *v1.Topic,
	err error,
) {
	s.srvStub.SetResponse(tid, "CreateTopic", res, err)
}

func (s *PublisherServer) SetCreateTopicResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.Topic) (*v1.Topic, error)) {

	s.srvStub.SetResponseCreator(
		tid, "CreateTopic",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.Topic))
			return res, err
		},
	)
}

func (s *fakePublisherServerImpl) DeleteTopic(
	ctx context.Context,
	req *v1.DeleteTopicRequest,
) (*v1.Empty, error) {
	return fakegrpc.HandleRequest[*v1.Empty](s.srvStub, ctx, "DeleteTopic", req)
}

func (s *PublisherServer) SetDeleteTopicResponse(
	tid string,
	res *v1.Empty,
	err error,
) {
	s.srvStub.SetResponse(tid, "DeleteTopic", res, err)
}

func (s *PublisherServer) SetDeleteTopicResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.DeleteTopicRequest) (*v1.Empty, error)) {

	s.srvStub.SetResponseCreator(
		tid, "DeleteTopic",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.DeleteTopicRequest))
			return res, err
		},
	)
}

func (s *fakePublisherServerImpl) DetachSubscription(
	ctx context.Context,
	req *v1.DetachSubscriptionRequest,
) (*v1.DetachSubscriptionResponse, error) {
	return fakegrpc.HandleRequest[*v1.DetachSubscriptionResponse](s.srvStub, ctx, "DetachSubscription", req)
}

func (s *PublisherServer) SetDetachSubscriptionResponse(
	tid string,
	res *v1.DetachSubscriptionResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "DetachSubscription", res, err)
}

func (s *PublisherServer) SetDetachSubscriptionResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.DetachSubscriptionRequest) (*v1.DetachSubscriptionResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "DetachSubscription",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.DetachSubscriptionRequest))
			return res, err
		},
	)
}

func (s *fakePublisherServerImpl) GetTopic(
	ctx context.Context,
	req *v1.GetTopicRequest,
) (*v1.Topic, error) {
	return fakegrpc.HandleRequest[*v1.Topic](s.srvStub, ctx, "GetTopic", req)
}

func (s *PublisherServer) SetGetTopicResponse(
	tid string,
	res *v1.Topic,
	err error,
) {
	s.srvStub.SetResponse(tid, "GetTopic", res, err)
}

func (s *PublisherServer) SetGetTopicResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.GetTopicRequest) (*v1.Topic, error)) {

	s.srvStub.SetResponseCreator(
		tid, "GetTopic",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.GetTopicRequest))
			return res, err
		},
	)
}

func (s *fakePublisherServerImpl) ListTopicSnapshots(
	ctx context.Context,
	req *v1.ListTopicSnapshotsRequest,
) (*v1.ListTopicSnapshotsResponse, error) {
	return fakegrpc.HandleRequest[*v1.ListTopicSnapshotsResponse](s.srvStub, ctx, "ListTopicSnapshots", req)
}

func (s *PublisherServer) SetListTopicSnapshotsResponse(
	tid string,
	res *v1.ListTopicSnapshotsResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "ListTopicSnapshots", res, err)
}

func (s *PublisherServer) SetListTopicSnapshotsResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.ListTopicSnapshotsRequest) (*v1.ListTopicSnapshotsResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "ListTopicSnapshots",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.ListTopicSnapshotsRequest))
			return res, err
		},
	)
}

func (s *fakePublisherServerImpl) ListTopicSubscriptions(
	ctx context.Context,
	req *v1.ListTopicSubscriptionsRequest,
) (*v1.ListTopicSubscriptionsResponse, error) {
	return fakegrpc.HandleRequest[*v1.ListTopicSubscriptionsResponse](s.srvStub, ctx, "ListTopicSubscriptions", req)
}

func (s *PublisherServer) SetListTopicSubscriptionsResponse(
	tid string,
	res *v1.ListTopicSubscriptionsResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "ListTopicSubscriptions", res, err)
}

func (s *PublisherServer) SetListTopicSubscriptionsResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.ListTopicSubscriptionsRequest) (*v1.ListTopicSubscriptionsResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "ListTopicSubscriptions",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.ListTopicSubscriptionsRequest))
			return res, err
		},
	)
}

func (s *fakePublisherServerImpl) ListTopics(
	ctx context.Context,
	req *v1.ListTopicsRequest,
) (*v1.ListTopicsResponse, error) {
	return fakegrpc.HandleRequest[*v1.ListTopicsResponse](s.srvStub, ctx, "ListTopics", req)
}

func (s *PublisherServer) SetListTopicsResponse(
	tid string,
	res *v1.ListTopicsResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "ListTopics", res, err)
}

func (s *PublisherServer) SetListTopicsResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.ListTopicsRequest) (*v1.ListTopicsResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "ListTopics",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.ListTopicsRequest))
			return res, err
		},
	)
}

func (s *fakePublisherServerImpl) Publish(
	ctx context.Context,
	req *v1.PublishRequest,
) (*v1.PublishResponse, error) {
	return fakegrpc.HandleRequest[*v1.PublishResponse](s.srvStub, ctx, "Publish", req)
}

func (s *PublisherServer) SetPublishResponse(
	tid string,
	res *v1.PublishResponse,
	err error,
) {
	s.srvStub.SetResponse(tid, "Publish", res, err)
}

func (s *PublisherServer) SetPublishResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.PublishRequest) (*v1.PublishResponse, error)) {

	s.srvStub.SetResponseCreator(
		tid, "Publish",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.PublishRequest))
			return res, err
		},
	)
}

func (s *fakePublisherServerImpl) UpdateTopic(
	ctx context.Context,
	req *v1.UpdateTopicRequest,
) (*v1.Topic, error) {
	return fakegrpc.HandleRequest[*v1.Topic](s.srvStub, ctx, "UpdateTopic", req)
}

func (s *PublisherServer) SetUpdateTopicResponse(
	tid string,
	res *v1.Topic,
	err error,
) {
	s.srvStub.SetResponse(tid, "UpdateTopic", res, err)
}

func (s *PublisherServer) SetUpdateTopicResponseCreator(
	tid string,
	creator func(ctx context.Context, req *v1.UpdateTopicRequest) (*v1.Topic, error)) {

	s.srvStub.SetResponseCreator(
		tid, "UpdateTopic",
		func(ctx context.Context, req any) (any, error) {
			res, err := creator(ctx, req.(*v1.UpdateTopicRequest))
			return res, err
		},
	)
}
