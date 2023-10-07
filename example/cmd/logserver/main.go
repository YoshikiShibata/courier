package main

import (
	"reflect"

	log_v1 "github.com/YoshikiShibata/courier/example/api/v1"
	"github.com/YoshikiShibata/courier/server"
)

func main() {
	server.Generate(&server.TemplateConfig{
		ServerPackageName:      "fakelogv1",
		ProtoPackageImportPath: "github.com/YoshikiShibata/example/api/v1",
		ProtoPackageImportName: "log_v1",
		ServiceName:            "Log",
		ServiceClientType:      reflect.TypeOf(log_v1.NewLogClient(nil)),
	})
}
