package main

import (
	"reflect"

	"github.com/YoshikiShibata/courier/example/api/shipping_v1"
	"github.com/YoshikiShibata/courier/server/fakegrpc"
)

func main() {
	fakegrpc.Generate(&fakegrpc.TemplateConfig{
		ServerPackageName:      "fakeshipping_v1",
		ProtoPackageImportPath: "github.com/YoshikiShibata/courier/example/api/shipping_v1",
		ProtoPackageImportName: "v1",
		ServiceName:            "Shipping",
		ServiceClientType:      reflect.TypeOf(shipping_v1.NewShippingClient(nil)),
	})
}
