// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package server

// This file contains formats to generate a fake server.

const (
	packageFmt = `// This code is auto-generated. DO NOT EDIT.

// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package %s
`

	importFmt = `
import (
    "context"

    "google.golang.org/grpc"

    v1 "github.com/x-asia/api-go/%s"
    "github.com/x-asia/kauche-app/lib/go/e2eframework/fakeserver"
)
`
	prepareFmt = `
type fake%[1]sServerImpl struct {
    v1.%[1]sServer
    srvStub *fakeserver.ServerStub
}

type %[1]sServer struct {
    impl    *fake%[1]sServerImpl
    srvStub *fakeserver.ServerStub
}

func New%[1]sServer(
    grpcServer *grpc.Server,
) *%[1]sServer {
    srvStub := &fakeserver.ServerStub{}
    srvImpl := &fake%[1]sServerImpl{
        srvStub: srvStub,
    }
    v1.Register%[1]sServer(grpcServer, srvImpl)
    return &%[1]sServer{
        impl:    srvImpl,
        srvStub: srvStub,
    }
}

func (s *%[1]sServer) ClearAllResponses(
    tid string,
) {
    s.srvStub.ClearAllResponses(tid)
}

`

	methodFmt = `
func (s *fake%[1]sServerImpl) %[2]s(
    ctx context.Context,
    req %[3]s,
) (%[4]s, error) {
    return fakeserver.HandleRequest[%[5]s](s.srvStub, ctx, "%[2]s", req)
}

func (s *%[1]sServer) Set%[2]sResponse(
    tid string,
    res %[4]s,
    err error,
) {
    s.srvStub.SetResponse(tid, "%[2]s", res, err)
}

func (s *%[1]sServer) Set%[2]sResponseCreator(
    tid string,
    creator func(ctx context.Context, req %[3]s) (%[4]s, error)) {

    s.srvStub.SetResponseCreator(
        tid, "%[2]s",
        func(ctx context.Context, req any) (any, error) {
            res, err := creator(ctx, req.(%[3]s))
            return res, err
        },
    )
}

`
)
