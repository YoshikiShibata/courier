// Copyright © 2023 Yoshiki Shibata. All rights reserved.

//go:build e2e

package e2etest

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YoshikiShibata/courier/example/api/shop_v1"
)

func TestListProducts_XXX(t *testing.T) {
	client := newShopClient(t)
	ctx := context.Background()

	_, err := client.ListProducts(ctx, &shop_v1.ListProductsRequest{})
	if status.Code(err) != codes.OK {
		t.Fatalf("ListProducts failed: %v", err)
	}
}
