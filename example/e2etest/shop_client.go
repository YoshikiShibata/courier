// Copyright © 2023 Yoshiki Shibata. All rights reserved.

//go:build e2e

package e2etest

import (
	"fmt"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/YoshikiShibata/courier/example/api/shop_v1"
)

func newShopClient(t *testing.T) shop_v1.ShopClient {
	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cc, err := grpc.Dial(fmt.Sprintf("localhost:%d", shopGRPCPort), clientOptions...)
	if err != nil {
		t.Fatalf("grpc.Dial failed: %v", err)
	}
	t.Cleanup(func() {
		_ = cc.Close()
	})

	client := shop_v1.NewShopClient(cc)
	return client
}
