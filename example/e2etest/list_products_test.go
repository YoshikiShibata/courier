// Copyright © 2023 Yoshiki Shibata. All rights reserved.

//go:build e2e

package e2etest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/YoshikiShibata/courier/example/api/shop_v1"
	"github.com/YoshikiShibata/courier/example/api/warehouse_v1"
	"github.com/YoshikiShibata/courier/tid"
)

func TestListProductInventories_InvalidArgumentNumOfProducts(
	t *testing.T,
) {
	t.Parallel()

	client := newShopClient(t)
	ctx := context.Background()

	// Calling and executing ListProductInventories.
	_, err := client.ListProductInventories(ctx,
		&shop_v1.ListProductInventoriesRequest{
			NumOfProducts: 0, // invalid
			PageToken:     "",
		})

	// Inspecting the result
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("ListProductInventories returned %v, want InvalidArgument", err)
	}
}

func TestListProductInventories_InvalidPageToken(t *testing.T) {
	t.Parallel()

	client := newShopClient(t)
	tid, ctx := tid.New(context.Background())

	// Configuring responses from the fake Warehouse service.
	fakeWarehouseServer.SetListProductInventoriesResponse(tid,
		nil, status.Error(codes.InvalidArgument, "page_token"))

	// Calling and executing ListProductInventories.
	_, err := client.ListProductInventories(ctx,
		&shop_v1.ListProductInventoriesRequest{
			NumOfProducts: 100,
			PageToken:     uuid.NewString(), // invalid
		})

	// Inspecting the result
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("ListProductInventories returned %v, want InvalidArgument", err)
	}
}

func TestListProductInventories_CanceledAndDeadlineExceeded(t *testing.T) {
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
		fakeWarehouseServer.SetListProductInventoriesResponse(tid,
			nil, status.Error(tc.code, tc.msg))

		// Calling and executing ListProductInventories.
		_, err := client.ListProductInventories(ctx,
			&shop_v1.ListProductInventoriesRequest{
				NumOfProducts: 100,
				PageToken:     uuid.NewString(), // invalid
			})

		// Inspecting the result
		if status.Code(err) != tc.code {
			t.Errorf("ListProductInventories returned %v, want %v",
				err, tc.code)
		}
	}
}

func TestListProductInventories_Normal(t *testing.T) {
	t.Parallel()

	client := newShopClient(t)
	tid, ctx := tid.New(context.Background())

	// Preparing inventory at the Warehouse.
	const numOfProducts = 10
	warehouseProducts := make([]*warehouse_v1.ProductInventory,
		numOfProducts)
	for i := 0; i < numOfProducts; i++ {
		warehouseProducts[i] = &warehouse_v1.ProductInventory{
			Number:            uuid.NewString(),
			Name:              fmt.Sprintf("商品-%03d", i),
			Price:             uint32(i*10 + 100),
			QuantityAvailable: uint32(i % 5), // Out of stock every 5 products
		}
	}

	// Configuring responses from the fake Warehouse service.
	fakeWarehouseServer.SetListProductInventoriesResponse(tid,
		&warehouse_v1.ListProductInventoriesResponse{
			ProductInventories: warehouseProducts,
			NextPageToken:      "",
		}, nil)

	// Calling and executing ListProductInventories.
	res, err := client.ListProductInventories(ctx,
		&shop_v1.ListProductInventoriesRequest{
			NumOfProducts: numOfProducts,
			PageToken:     "",
		})

	// Inspecting the results
	if status.Code(err) != codes.OK {
		t.Fatalf("ListProductInventories failed: %v", err)
	}

	// Inspecting each returned ProductInventory
	wantProductInventories := make([]*shop_v1.ProductInventory, 0, len(warehouseProducts))
	for _, p := range warehouseProducts {
		wantProductInventories = append(wantProductInventories,
			&shop_v1.ProductInventory{
				Number:            p.Number,
				Name:              p.Name,
				Price:             p.Price,
				QuantityAvailable: p.QuantityAvailable,
			})
	}

	opts := cmpopts.IgnoreUnexported(shop_v1.ProductInventory{})
	if diff := cmp.Diff(res.ProductInventories, wantProductInventories, opts); diff != "" {
		t.Errorf("(-want, +got)\n%s", diff)
	}
}
