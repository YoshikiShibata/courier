// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package main

import (
	"reflect"

	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"github.com/YoshikiShibata/courier/server/fakegrpc"
)

func main() {
	fakegrpc.Generate(&fakegrpc.TemplateConfig{
		ServerPackageName:      "publisher_v1",
		ProtoPackageImportPath: "cloud.google.com/go/pubsub/apiv1/pubsubpb",
		ProtoPackageImportName: "v1",
		ServiceName:            "Publisher",
		ServiceClientType:      reflect.TypeOf(pubsubpb.NewPublisherClient(nil)),
	})
}
