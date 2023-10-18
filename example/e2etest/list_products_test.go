// Copyright © 2023 Yoshiki Shibata. All rights reserved.

//go:build e2e

package e2etest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YoshikiShibata/courier/example/api/shop_v1"
	"github.com/YoshikiShibata/courier/example/api/warehouse_v1"
	"github.com/YoshikiShibata/courier/tid"
)

func TestListProducts_XXX(t *testing.T) {
	t.Parallel()

	client := newShopClient(t)
	ctx := context.Background()

	// Warehouseの在庫準備
	const numOfProducts = 100
	warehouseProducts := make([]*warehouse_v1.Product, numOfProducts)
	for i := 0; i < numOfProducts; i++ {
		warehouseProducts[i] = &warehouse_v1.Product{
			Id:                uuid.NewString(),
			Name:              fmt.Sprintf("商品-%03d", i),
			Price:             uint32(i*10 + 100),
			QuantityAvailable: uint32(i % 5), // 5個おきに、在庫無し
		}
	}

	// フェイクサービスの設定
	tid, ctx := tid.New(ctx)
	fakeWarehouseServer.SetListProductsResponse(tid,
		&warehouse_v1.ListProductsResponse{
			Products:      warehouseProducts,
			NextPageToken: "",
		}, nil)

	// 実行
	res, err := client.ListProducts(ctx, &shop_v1.ListProductsRequest{
		NumOfProducts: numOfProducts,
		PageToken:     "",
	})

	// 検査
	if status.Code(err) != codes.OK {
		t.Fatalf("ListProducts failed: %v", err)
	}

	const wantNumOfProducts = (numOfProducts * 4) / 5

	if len(res.Products) != wantNumOfProducts {
		t.Fatalf("len(res.Products) is %d, want %d", len(res.Products), wantNumOfProducts)
	}

	// 返されたProductsの確認
	// ...
}
