// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package fakegrpc

// This file contains formats to generate a fake server.
const fakeServerTemplate = `// This code is auto-generated. DO NOT EDIT.

package {{.ServerPackageName}}

import (
    "context"

    "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/YoshikiShibata/courier/server/fakegrpc"

	{{.ProtoPackageImportName}} "{{.ProtoPackageImportPath}}"
)

type fake{{.ServiceName}}ServerImpl struct {
	{{.ProtoPackageImportName}}.{{.ServiceName}}Server
    srvStub *fakegrpc.GRPCServerStub
}

type {{.ServiceName}}Server struct {
    impl    *fake{{.ServiceName}}ServerImpl
    srvStub *fakegrpc.GRPCServerStub
}

func New{{.ServiceName}}Server(
    grpcServer *grpc.Server,
) *{{.ServiceName}}Server {
    srvStub := &fakegrpc.GRPCServerStub{}
    srvImpl := &fake{{.ServiceName}}ServerImpl{
        srvStub: srvStub,
    }
	{{.ProtoPackageImportName}}.Register{{.ServiceName}}Server(grpcServer, srvImpl)
    return &{{.ServiceName}}Server{
        impl:    srvImpl,
        srvStub: srvStub,
    }
}

func (s *{{.ServiceName}}Server) ClearAllResponses(
    tid string,
) {
    s.srvStub.ClearAllResponses(tid)
}

{{range .Methods}}

func (s *fake{{$.ServiceName}}ServerImpl) {{.Name}}(
    ctx context.Context,
    req {{.ReqTypeName}},
) ({{.ResTypeName}}, error) {
    return fakegrpc.HandleRequest[{{.ResTypeName}}](s.srvStub, ctx, "{{.Name}}", req)
}

func (s *{{$.ServiceName}}Server) Set{{.Name}}Response(
    tid string,
    res {{.ResTypeName}},
    err error,
) {
    s.srvStub.SetResponse(tid, "{{.Name}}", res, err)
}

func (s *{{$.ServiceName}}Server) Set{{.Name}}ResponseCreator(
    tid string,
    creator func(ctx context.Context, req {{.ReqTypeName}}) ({{.ResTypeName}}, error)) {

    s.srvStub.SetResponseCreator(
        tid, "{{.Name}}",
        func(ctx context.Context, req any) (any, error) {
            res, err := creator(ctx, req.({{.ReqTypeName}}))
            return res, err
        },
    )
}

{{end}}
`
