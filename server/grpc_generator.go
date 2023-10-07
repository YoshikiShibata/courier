// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package server

import (
	"html/template"
	"log"
	"os"
	"reflect"
)

type RPCMethod struct {
	Name        string
	ReqTypeName string
	ResTypeName string
}

// TemplateConfig defines paramters to templates.
type TemplateConfig struct {
	ServerPackageName      string
	ProtoPackageImportPath string // Import Path
	ProtoPackageImportName string // Import Name
	ServiceName            string
	ServiceClientType      reflect.Type
	Methods                []*RPCMethod
}

// Generate generates a fake server code to stdout.
func Generate(config *TemplateConfig) {
	rpcMethods := make([]*RPCMethod, 0, config.ServiceClientType.NumMethod())
	for i := 0; i < config.ServiceClientType.NumMethod(); i++ {
		method := config.ServiceClientType.Method(i)
		methodType := method.Type
		rpcMethods = append(rpcMethods, &RPCMethod{
			Name:        method.Name,
			ReqTypeName: methodType.In(2).String(),
			ResTypeName: methodType.Out(0).String(),
		})
	}

	config.Methods = rpcMethods

	fakeserver, err := template.New("fakeserver").Parse(fakeServerTemplate)
	if err != nil {
		log.Fatal(err)
	}

	if err := fakeserver.Execute(os.Stdout, config); err != nil {
		log.Fatal(err)
	}
}
