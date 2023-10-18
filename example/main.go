// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/YoshikiShibata/courier/example/api/shipping_v1"
	"github.com/YoshikiShibata/courier/example/api/warehouse_v1"
	"github.com/YoshikiShibata/courier/example/internal/server"
	"github.com/YoshikiShibata/courier/tid"
)

type Env struct {
	GRPCServerPort       int    `envconfig:"GRPC_SERVER_PORT" required:"true"`
	ShippingServiceAddr  string `envconfig:"SHIPPING_SERVICE_ADDR" required:"true"`
	WarehouseServiceAddr string `envconfig:"WAREHOUSE_SERVICE_ADDR" required:"true"`
}

func main() {
	fmt.Printf("GOCOVERDIR: %s\n", os.Getenv("GOCOVERDIR"))
	var env Env
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("envconfig.Process failed: %v", err)
	}

	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	shippingCC, err := grpc.Dial(env.ShippingServiceAddr, clientOptions...)
	if err != nil {
		log.Fatalf("grpc.Dial(%s) failed: %v", env.ShippingServiceAddr, err)
	}
	defer func() { _ = shippingCC.Close() }()
	shippingClient := shipping_v1.NewShippingClient(shippingCC)

	warehouseCC, err := grpc.Dial(env.WarehouseServiceAddr, clientOptions...)
	if err != nil {
		log.Fatalf("grpc.Dial(%s) failed: %v", env.WarehouseServiceAddr, err)
	}
	defer func() { _ = warehouseCC.Close() }()
	warehouseClient := warehouse_v1.NewWarehouseClient(warehouseCC)

	grpcOpts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(tid.NewGRPCHeaderPropagator()),
	}

	grpcServer, err := server.NewGRPCServer(
		env.GRPCServerPort,
		shippingClient,
		warehouseClient,
		grpcOpts...,
	)
	if err != nil {
		log.Fatalf("server.NewGRPCServer failed: %v", err)
	}

	go func() {
		log.Printf("Start Server")
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("grpcServer.Start failed: %v", err)
		}
	}()

	quitC := make(chan os.Signal, 1)
	signal.Notify(quitC, syscall.SIGINT, syscall.SIGTERM)
	<-quitC

	grpcServer.Stop()
	log.Printf("Service gracefully shut down")
}
