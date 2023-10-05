// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package server

import (
	"fmt"
	"reflect"
)

type rpcMethod struct {
	name        string
	reqTypeName string
	resTypeName string
}

// FormatConfig defines paramters to formats.
type FormatConfig struct {
	PackageName  string
	ProtoPackage string
	ServiceName  string
	ServiceType  reflect.Type
}

// Generate generates a fake server code to stdout.
func Generate(config *FormatConfig) {
	rpcMethods := make([]*rpcMethod, 0, config.ServiceType.NumMethod())
	for i := 0; i < config.ServiceType.NumMethod(); i++ {
		method := config.ServiceType.Method(i)
		methodType := method.Type
		rpcMethods = append(rpcMethods, &rpcMethod{
			name:        method.Name,
			reqTypeName: methodType.In(2).String(),
			resTypeName: methodType.Out(0).String(),
		})
	}

	fmt.Printf(packageFmt, config.PackageName)
	fmt.Printf(importFmt, config.ProtoPackage)
	fmt.Printf(prepareFmt, config.ServiceName)

	for _, method := range rpcMethods {
		fmt.Printf(methodFmt,
			config.ServiceName,     // [1]
			method.name,            // [2]
			method.reqTypeName,     // [3]
			method.resTypeName,     // [4]
			method.resTypeName[1:], // [5]
		)
	}
}
