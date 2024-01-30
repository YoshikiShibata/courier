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
	fakepublisher_v1 "github.com/YoshikiShibata/courier/example/fakeservers/publisher_v1"
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

	fakePublisherServer *fakepublisher_v1.PublisherServer
)

func e2eCoverage(m *testing.M) (exitCode int) {
	// フェイクサービス用のgRPCサーバ—
	grpcServer := fakegrpc.NewGRPCServer()
	server := grpcServer.Server()
	port := grpcServer.Port()

	// Shippingフェイクサービスの作成
	fakeShippingServer = fakeshipping_v1.NewShippingServer(server)
	// Warehouseフェイクサービスの作成
	fakeWarehouseServer = fakewarehouse_v1.NewWarehouseServer(server)
	// Publisherフェイクサービスの作成
	fakePublisherServer = fakepublisher_v1.NewPublisherServer(server)

	// フェイクサービスの環境変数作成（と表示）
	shippingSvcEnv :=
		fmt.Sprintf("SHIPPING_SERVICE_ADDR=localhost:%s", port)
	warehouseSvcEnv :=
		fmt.Sprintf("WAREHOUSE_SERVICE_ADDR=localhost:%s", port)
	publisherSvcEnv :=
		fmt.Sprintf("PUBSUB_EMULATOR_HOST=localhost:%s", port)
	log.Printf("%s", shippingSvcEnv)
	log.Printf("%s", warehouseSvcEnv)
	log.Printf("%s", publisherSvcEnv)

	// フェイクサービスのgRPCサーバーの起動
	go func() { grpcServer.Serve() }()

	// テスト対象サービスのgRPCサーバーのポート割り当て
	ports := dynaport.Get(1)
	shopGRPCPort = ports[0]
	log.Printf("GRPC Port for Shop service is %d", shopGRPCPort)

	// テスト対象サービスの起動用Config設定
	config := &courier.Config{
		MakeDir:           "..",
		MakeBuildTarget:   "coverage_build",
		ServiceBinaryPath: "bin/shop_server",
		ServiceName:       "Shop",
		GRPCPort:          shopGRPCPort,
		Verbose:           testing.Verbose(),
		CoverageDir:       "coverage",
	}
	// 環境変数の設定
	config.Envs = append(config.Envs,
		fmt.Sprintf("GRPC_SERVER_PORT=%d", config.GRPCPort),
		shippingSvcEnv,
		warehouseSvcEnv,
		publisherSvcEnv,
	)

	// テスト対象サービスの起動
	terminateService, err := courier.InvokeService(config)
	if err != nil {
		log.Printf("courier.InvokeService failed: %v", err)
		return 1
	}

	// テスト関数の実行
	exitCode = m.Run()

	// テスト対象サービスの終了
	if err := terminateService(); err != nil {
		log.Printf("terminateService failed: %v", err)
		return 1
	}

	return exitCode
}
