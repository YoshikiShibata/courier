// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package server

import (
	"fmt"
	"reflect"
)

type rpcMethod struct {
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
}

// Generate generates a fake server code to stdout.
func Generate(config *TemplateConfig) {
	rpcMethods := make([]*rpcMethod, 0, config.ServiceClientType.NumMethod())
	for i := 0; i < config.ServiceClientType.NumMethod(); i++ {
		method := config.ServiceClientType.Method(i)
		methodType := method.Type
		rpcMethods = append(rpcMethods, &rpcMethod{
			Name:        method.Name,
			ReqTypeName: methodType.In(2).String(),
			ResTypeName: methodType.Out(0).String(),
		})
	}

	fmt.Printf(packageFmt, config.ServerPackageName)
	fmt.Printf(importFmt, config.ProtoPackageImportPath)
	fmt.Printf(prepareFmt, config.ServiceName)

	for _, method := range rpcMethods {
		fmt.Printf(methodFmt,
			config.ServiceName,     // [1]
			method.Name,            // [2]
			method.ReqTypeName,     // [3]
			method.ResTypeName,     // [4]
			method.ResTypeName[1:], // [5]
		)
	}
}
