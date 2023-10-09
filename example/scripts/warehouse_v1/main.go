package main

import (
	"reflect"

	"github.com/YoshikiShibata/courier/example/api/warehouse_v1"
	"github.com/YoshikiShibata/courier/server"
)

func main() {
	server.Generate(&server.TemplateConfig{
		ServerPackageName:      "fakewarehouse_v1",
		ProtoPackageImportPath: "github.com/YoshikiShibata/courier/example/api/warehouse_v1",
		ProtoPackageImportName: "v1",
		ServiceName:            "Warehouse",
		ServiceClientType:      reflect.TypeOf(warehouse_v1.NewWarehouseClient(nil)),
	})
}
