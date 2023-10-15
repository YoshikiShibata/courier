// Copyright © 2023 Yoshiki Shibata. All rights reserved.

//go:build e2e

package e2etest

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/YoshikiShibata/courier/server/fakegrpc"
	"github.com/travisjeffery/go-dynaport"

	"github.com/YoshikiShibata/courier"
	fakeshipping_v1 "github.com/YoshikiShibata/courier/example/fakeservers/shipping_v1"
	fakewarehouse_v1 "github.com/YoshikiShibata/courier/example/fakeservers/warehouse_v1"
)

func TestMain(m *testing.M) {
	flag.Parse()

	exitCode := e2eCoverage(m)
	os.Exit(exitCode)
}

var (
	fakeShippingServer  *fakeshipping_v1.ShippingServer
	fakeWarehouseServer *fakewarehouse_v1.WarehouseServer
	shopGRPCPort        int
)

func e2eCoverage(m *testing.M) (exitCode int) {
	ports := dynaport.Get(1)
	shopGRPCPort = ports[0]
	log.Printf("GRPC Port for Shop service is %d", shopGRPCPort)

	config := &courier.Config{
		MakeDir:           "..",
		MakeBuildTarget:   "build",
		ServiceBinaryPath: "bin/shop_server",
		ServiceName:       "Shop",
		GRPCPort:          shopGRPCPort,
		Verbose:           testing.Verbose(),
		CoverageDir:       "coverage",
	}

	grpcServer := fakegrpc.NewGRPCServer()

	fakeShippingServer = fakeshipping_v1.NewShippingServer(grpcServer.Server())
	fakeWarehouseServer = fakewarehouse_v1.NewWarehouseServer(grpcServer.Server())

	shippingSvcEnv := fmt.Sprintf("SHIPPING_SERVICE_ADDR=localhost:%s", grpcServer.Port())
	log.Printf("%s", shippingSvcEnv)
	warehouseSvcEnv := fmt.Sprintf("WAREHOUSE_SERVICE_ADDR=localhost:%s", grpcServer.Port())
	log.Printf("%s", warehouseSvcEnv)

	config.Envs = append(config.Envs,
		fmt.Sprintf("GRPC_SERVER_PORT=%d", config.GRPCPort),
		shippingSvcEnv,
		warehouseSvcEnv,
	)

	go func() { grpcServer.Serve() }()

	terminateService, err := courier.InvokeService(config)
	if err != nil {
		log.Printf("courier.InvokeService failed: %v", err)
		return 1
	}

	exitCode = m.Run()

	if err := terminateService(); err != nil {
		log.Printf("terminateService failed: %v", err)
		return 1
	}

	return exitCode
}
