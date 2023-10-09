package main

import (
	"reflect"

	"github.com/YoshikiShibata/courier/example/api/shipping_v1"
	"github.com/YoshikiShibata/courier/server"
)

func main() {
	server.Generate(&server.TemplateConfig{
		ServerPackageName:      "fakeshipping_v1",
		ProtoPackageImportPath: "github.com/YoshikiShibata/courier/example/api/shipping_v1",
		ProtoPackageImportName: "shipping_v1",
		ServiceName:            "Shipping",
		ServiceClientType:      reflect.TypeOf(shipping_v1.NewShippingClient(nil)),
	})
}
