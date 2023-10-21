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

func TestListProducts_InvalidArgumentNumOfProducts(t *testing.T) {
	t.Parallel()

	client := newShopClient(t)
	ctx := context.Background()

	// Calling and executing ListProducts.
	_, err := client.ListProducts(ctx,
		&shop_v1.ListProductsRequest{
			NumOfProducts: 0, // invalid
			PageToken:     "",
		})

	// Inspecting the result
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("ListProducts returned %v, want InvalidArgument", err)
	}
}

func TestListProducts_InvalidPageToken(t *testing.T) {
	t.Parallel()

	client := newShopClient(t)
	tid, ctx := tid.New(context.Background())

	// Configuring responses from the fake Warehouse service.
	fakeWarehouseServer.SetListProductsResponse(tid,
		nil, status.Error(codes.InvalidArgument, "page_token"))

	// Calling and executing ListProducts.
	_, err := client.ListProducts(ctx,
		&shop_v1.ListProductsRequest{
			NumOfProducts: 100,
			PageToken:     uuid.NewString(), // invalid
		})

	// Inspecting the result
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("ListProducts returned %v, want InvalidArgument", err)
	}
}

func TestListProducts_CanceledAndDeadlineExceeded(t *testing.T) {
	t.Parallel()

	client := newShopClient(t)
	tid, ctx := tid.New(context.Background())

	for _, tc := range []struct {
		code codes.Code
		msg  string
	}{
		{
			code: codes.Canceled,
			msg:  "canceled",
		},
		{
			code: codes.DeadlineExceeded,
			msg:  "deadline exceeded",
		},
	} {
		// Configuring responses from the fake Warehouse service.
		fakeWarehouseServer.SetListProductsResponse(tid,
			nil, status.Error(tc.code, tc.msg))

		// Calling and executing ListProducts.
		_, err := client.ListProducts(ctx,
			&shop_v1.ListProductsRequest{
				NumOfProducts: 100,
				PageToken:     uuid.NewString(), // invalid
			})

		// Inspecting the result
		if status.Code(err) != tc.code {
			t.Errorf("ListProducts returned %v, want %v", err, tc.code)
		}
	}
}

func TestListProducts_XXX(t *testing.T) {
	t.Parallel()

	client := newShopClient(t)
	tid, ctx := tid.New(context.Background())

	// Preparing inventory at the Warehouse.
	const numOfProducts = 100
	warehouseProducts := make([]*warehouse_v1.Product, numOfProducts)
	for i := 0; i < numOfProducts; i++ {
		warehouseProducts[i] = &warehouse_v1.Product{
			Number:            uuid.NewString(),
			Name:              fmt.Sprintf("商品-%03d", i),
			Price:             uint32(i*10 + 100),
			QuantityAvailable: uint32(i % 5), // Out of stock every 5 products
		}
	}

	// Configuring responses from the fake Warehouse service.
	fakeWarehouseServer.SetListProductsResponse(tid,
		&warehouse_v1.ListProductsResponse{
			Products:      warehouseProducts,
			NextPageToken: "",
		}, nil)

	// Calling and executing ListProducts.
	res, err := client.ListProducts(ctx, &shop_v1.ListProductsRequest{
		NumOfProducts: numOfProducts,
		PageToken:     "",
	})

	// Inspecting the results
	if status.Code(err) != codes.OK {
		t.Fatalf("ListProducts failed: %v", err)
	}

	const wantNumOfProducts = (numOfProducts * 4) / 5
	if len(res.Products) != wantNumOfProducts {
		t.Fatalf("len(res.Products) is %d, want %d", len(res.Products), wantNumOfProducts)
	}

	// Inspecting each returned Product individually.
	// ...
}
