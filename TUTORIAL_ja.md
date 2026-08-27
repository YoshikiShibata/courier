# チュートリアル: `server/fakegrpc` と `tid` パッケージの使い方

このチュートリアルでは、courierに含まれる次の2つのパッケージの使い方を、
`example/` ディレクトリの実例とともに説明します。

- `github.com/YoshikiShibata/courier/server/fakegrpc`
  gRPCの依存先サービスを **フェイクサーバー** として立ち上げ、
  テストコードから任意のレスポンスを設定するためのライブラリ。
- `github.com/YoshikiShibata/courier/tid`
  E2Eテストが並列に走っても互いのフェイク設定が混線しないようにするための
  **Testing ID(TID)** を発行し、gRPCメタデータで伝搬させるためのライブラリ。

`example/` ディレクトリは、Shopサービス(テスト対象)がShipping/Warehouse/Publisherの
3つのgRPCサービスに依存する構成になっており、この2パッケージの使い方の完全な
リファレンス実装になっています。

---

## 全体像

E2Eテストは次のような構成で動きます。

```
                          ┌───────────────────────────────┐
                          │  E2Eテスト(list_products_test) │
                          └────────────┬──────────────────┘
                                       │ 1) tid.New で TIDを生成
                                       │ 2) fakeXxxServer.SetXxxResponse(tid, ...)
                                       │ 3) client.Xxx(ctx, req)
                                       ▼
                          ┌───────────────────────────────┐
                          │  テスト対象サービス (Shop)     │
                          │  tid.NewGRPCHeaderPropagator  │
                          │  でTIDをoutgoingへ伝搬        │
                          └────────────┬──────────────────┘
                                       │ 依存gRPC呼び出し(TIDヘッダ付き)
                                       ▼
                          ┌───────────────────────────────┐
                          │ 単一のフェイクgRPCサーバー     │
                          │  fakegrpc.NewGRPCServer()     │
                          │  ├─ Shipping fake             │
                          │  ├─ Warehouse fake            │
                          │  └─ Publisher fake            │
                          │  TIDでレスポンスを引き当てる   │
                          └───────────────────────────────┘
```

ポイントは、**フェイクサーバーがプロセス内でグローバルに1つ**でも、
テスト毎にユニークな `tid` をキーにしてレスポンスを登録できるため、
E2Eテストを `t.Parallel()` で並列実行しても衝突しないことです。

---

## 1. `tid` パッケージ

`tid.go` はわずか2つのAPIだけです。

### 1.1 `tid.New(ctx) (string, context.Context)`

TIDを新しく発行し、それをoutgoing gRPCメタデータ (ヘッダ名 `courier-testing-id`)
にセットしたコンテキストを返します。既にコンテキストにTIDが含まれていれば、
そちらを再利用します。

テスト側では毎テストの先頭でこれを呼びます。

```go
// example/e2etest/list_products_test.go より
tid, ctx := tid.New(context.Background())

fakeWarehouseServer.SetListProductInventoriesResponse(tid,
    &warehouse_v1.ListProductInventoriesResponse{ /* ... */ }, nil)

res, err := client.ListProductInventories(ctx,
    &shop_v1.ListProductInventoriesRequest{ /* ... */ })
```

- `tid` (文字列) は「今回のテストのためのレスポンス」を登録するときのキー。
- `ctx` はTIDヘッダ付きなので、`client.ListProductInventories(ctx, ...)` に渡すと
  Shopサービスまで伝搬します。

### 1.2 `tid.NewAlways(ctx) (string, context.Context)`

`New` と同じシグネチャですが、コンテキストに既存TIDが含まれていても **常に
新しいTIDを発行して差し替える** バージョンです。既存TIDは outgoing メタデータから
上書きされます。

以下のようなケースで使います。

- ループ内でテストシナリオを繰り返し、イテレーションごとに独立したTIDで
  レスポンスを設定したい場合。
- 親のセットアップで既にTID入りコンテキストが作られているが、サブテストで
  別のTIDに切り替えたい場合。

```go
_, ctx := tid.New(context.Background())

for i := 0; i < 3; i++ {
    // 各イテレーションで新しいTIDを発行
    tid, ctx := tid.NewAlways(ctx)
    fakeWarehouseServer.SetShipProductResponse(tid, resp[i], nil)
    // ... client を ctx で呼ぶ ...
}
```

通常のテストは `New` で十分です。`NewAlways` は「既存の ctx を再利用しつつ
TIDだけ差し替えたい」ときの明示的な選択肢として使ってください。

### 1.3 `tid.Extract(ctx) string`

incoming gRPCメタデータからTIDを取り出します。これはフェイクサーバー側で
「どのテストのレスポンスを返すべきか」を判断するために `GRPCServerStub`
が内部で使っているので、通常テストコードから直接呼ぶ必要はありません。

### 1.4 `tid.NewGRPCHeaderPropagator() grpc.UnaryServerInterceptor`

**テスト対象サービス**に組み込むUnaryインターセプターです。
incomingメタデータのTIDを、そのままoutgoingコンテキストに詰め替えて
下流のgRPC呼び出しにも伝搬させます。

```go
// example/main.go より
grpcOpts := []grpc.ServerOption{
    grpc_middleware.WithUnaryServerChain(tid.NewGRPCHeaderPropagator()),
}

grpcServer, err := server.NewGRPCServer(
    env.GRPCServerPort,
    shippingClient,
    warehouseClient,
    grpcOpts...,
)
```

これを入れておかないと、テスト → Shop まではTIDが届いていても、
Shop → Shipping/Warehouse への呼び出しでTIDが落ちてしまい、
フェイクサーバーがレスポンスを引き当てられず panic します
(`Response for :RpcName has not been set yet`)。

---

## 2. `server/fakegrpc` パッケージ

このパッケージは大きく2つの機能を提供します。

- **フェイクサーバーのコード生成** (`Generate` / `TemplateConfig`)
- **フェイクサーバーの実行時プリミティブ** (`GRPCServer`, `GRPCServerStub`,
  `HandleRequest`, `RequestCapture`, `ResponseBuilder`)

### 2.1 フェイクサーバーの生成

各依存gRPCサービスごとに、生成スクリプトを1本用意します。

```go
// example/scripts/shipping_v1/main.go
func main() {
    fakegrpc.Generate(&fakegrpc.TemplateConfig{
        ServerPackageName:      "fakeshipping_v1",
        ProtoPackageImportPath: "github.com/YoshikiShibata/courier/example/api/shipping_v1",
        ProtoPackageImportName: "v1",
        ServiceName:            "Shipping",
        ServiceClientType:      reflect.TypeOf(shipping_v1.NewShippingClient(nil)),
    })
}
```

`TemplateConfig` のフィールドは次の通りです。

| フィールド | 意味 |
| --- | --- |
| `ServerPackageName` | 生成するGoファイルの `package` 宣言 |
| `ProtoPackageImportPath` | protobufサービスの実装が入っているパッケージのimportパス |
| `ProtoPackageImportName` | そのパッケージのimportエイリアス名 (テンプレート内で使われる) |
| `ServiceName` | サービス名 (例: `Shipping`)。関数名の接頭辞になる |
| `ServiceClientType` | `NewXxxClient(nil)` を `reflect.TypeOf` に通した型情報。ここからRPC一覧が抽出される |

生成された `server.go` は標準出力に吐かれるので、Makefileでは次のように使います。

```makefile
# example/Makefile より
SERVICES = shipping_v1 warehouse_v1 publisher_v1

genfakes: $(SERVICES)

$(SERVICES):
	go run scripts/$@/main.go > fakeservers/$@/server.go
	gofmt -w fakeservers/$@/server.go
	goimports -w fakeservers/$@/server.go
```

生成物には、例えば `Shipping` サービスに `Create` / `Status` の2つのRPCがあれば、
以下のAPIが含まれます (`example/fakeservers/shipping_v1/server.go` 参照)。

- `NewShippingServer(grpcServer *grpc.Server) *ShippingServer`
  内部で `grpcServer` にサービス実装をregisterします。
- `(*ShippingServer).SetCreateResponse(tid, res, err)`
  この `tid` から `Create` が呼ばれたときに固定レスポンスを返す。
- `(*ShippingServer).SetCreateResponseCreator(tid, func(ctx, req) (res, err))`
  同上を関数として設定する。リクエスト内容に応じてレスポンスを変えたい場合や、
  リクエストをキャプチャしたい場合に使う。
- `(*ShippingServer).ClearAllResponses(tid)`
  その `tid` に紐づく全RPCのレスポンス設定を消す。

### 2.2 フェイクサーバーを立てる

`GRPCServer` は空きポートを自動割り当てして gRPCサーバーを組み立てるだけの
薄いラッパーです。

```go
// example/e2etest/main_test.go より
grpcServer := fakegrpc.NewGRPCServer()
server := grpcServer.Server()
port := grpcServer.Port()

// 1つのgRPCサーバーに、3種類のフェイクサービスを同居させる
fakeShippingServer  = fakeshipping_v1.NewShippingServer(server)
fakeWarehouseServer = fakewarehouse_v1.NewWarehouseServer(server)
fakePublisherServer = fakepublisher_v1.NewPublisherServer(server)

go func() { grpcServer.Serve() }()
```

依存サービスの接続先には、このフェイクサーバーのアドレスを環境変数で渡します。

```go
shippingSvcEnv  := fmt.Sprintf("SHIPPING_SERVICE_ADDR=localhost:%s", port)
warehouseSvcEnv := fmt.Sprintf("WAREHOUSE_SERVICE_ADDR=localhost:%s", port)
publisherSvcEnv := fmt.Sprintf("PUBSUB_EMULATOR_HOST=localhost:%s", port)
```

Pub/SubエミュレータもgRPCサービスなので同じ仕組みでフェイク化できます
(`example/scripts/publisher_v1/main.go`)。

### 2.3 テストからレスポンスを設定する

パターン1: **固定レスポンス**を返したい場合。

```go
tid, ctx := tid.New(context.Background())

fakeWarehouseServer.SetListProductInventoriesResponse(tid,
    &warehouse_v1.ListProductInventoriesResponse{
        ProductInventories: warehouseProducts,
        NextPageToken:      "",
    }, nil)
```

パターン2: **エラーだけ**を返したい場合。

```go
fakeWarehouseServer.SetListProductInventoriesResponse(tid,
    nil, status.Error(codes.InvalidArgument, "page_token"))
```

パターン3: **同一RPCで複数回呼ばれ、回ごとに別のレスポンス**を返したい場合、
`fakegrpc.ResponseBuilder` を使います。

```go
rb := &fakegrpc.ResponseBuilder[
    *warehouse_v1.ShipProductRequest,
    *warehouse_v1.ShipProductResponse,
]{}
rb.Append(&warehouse_v1.ShipProductResponse{ /* 1回目 */ }, nil)
rb.Append(nil, status.Error(codes.Unavailable, "retry"))
rb.Append(&warehouse_v1.ShipProductResponse{ /* 3回目 */ }, nil)

fakeWarehouseServer.SetShipProductResponseCreator(tid, rb.Creator)
```

`Creator` は呼び出し順に `Append` した通りのレスポンスを返します
(内部は `atomic.AddUint32` なので並列呼び出しにも安全)。

パターン4: **リクエスト内容を検査したい**場合、`fakegrpc.RequestCapture` を使います。

```go
rc := fakegrpc.NewRequestCapture[
    *warehouse_v1.ShipProductRequest,
    *warehouse_v1.ShipProductResponse,
](&warehouse_v1.ShipProductResponse{ /* 返したいレスポンス */ }, nil)

fakeWarehouseServer.SetShipProductResponseCreator(tid, rc.Creator)

// ... テスト対象を実行 ...

got := rc.Get() // Shopが送ってきた ShipProductRequest を検査できる
```

### 2.4 テスト後のクリーンアップ

必要に応じてテスト終了時に `ClearAllResponses(tid)` を呼びますが、
実際にはTIDはテストごとにユニークなので、放置しても他テストに影響しません。
繰り返し実行するループの中で同じTIDを使い回すようなケースでのみ有用です。

---

## 3. E2Eテストを組み立てる手順まとめ

`example/` を参考に、自プロジェクトに同じ仕組みを入れる手順は次の通りです。

### 3.1 テスト対象サービス側 (本番コード)

`grpc.NewServer` を作るときに `tid.NewGRPCHeaderPropagator()` を
Unaryインターセプターとして組み込むだけです (`example/main.go`)。

```go
import (
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    "github.com/YoshikiShibata/courier/tid"
)

grpcOpts := []grpc.ServerOption{
    grpc_middleware.WithUnaryServerChain(tid.NewGRPCHeaderPropagator()),
}
grpcServer := grpc.NewServer(grpcOpts...)
```

### 3.2 フェイクサーバーの生成スクリプトを用意

依存サービスごとに `scripts/<svc>/main.go` を1本ずつ書きます
(`example/scripts/*/main.go` がそのまま雛形になります)。

Makefileで `go run` して `fakeservers/<svc>/server.go` に吐き出し、
`gofmt` / `goimports` を掛けます。

### 3.3 E2Eテストの `TestMain`

`example/e2etest/main_test.go` の構成をそのまま流用できます。

1. `fakegrpc.NewGRPCServer()` でフェイク用gRPCサーバーを起動。
2. 各 `fakeXxxServer` をそこにregister。
3. 依存サービスのアドレス環境変数をフェイクサーバーに向ける。
4. テスト対象サービス本体を子プロセスとして起動する
   (この起動には `courier.InvokeService` を利用)。
5. `m.Run()` でテスト本体を実行。

### 3.4 各テストの流儀

```go
func TestXxx(t *testing.T) {
    t.Parallel()                              // 並列OK
    client := newShopClient(t)
    tid, ctx := tid.New(context.Background()) // このテスト用のTID発行

    // 依存サービス側のレスポンスを設定
    fakeWarehouseServer.SetListProductInventoriesResponse(tid, resp, nil)

    // テスト対象を呼ぶ (ctx にTIDが埋まっている)
    got, err := client.ListProductInventories(ctx, req)

    // 検証
    // ...
}
```

これで、フェイクサーバーは各テストが登録した「そのテスト専用のレスポンス」を
`TID` で正確に引き当てて返してくれます。

---

## 4. よくあるハマりどころ

- **`Response for <tid>:<RPC> has not been set yet` で panic する**
  - テスト対象サービス側で `tid.NewGRPCHeaderPropagator()` を組み込み忘れている。
  - `tid.New` で作った `ctx` ではなく、素の `context.Background()` を
    gRPCクライアント呼び出しに渡してしまっている。
  - あるRPCのレスポンス設定 (`SetXxxResponse`) を忘れている。
- **並列テストで別テストのレスポンスが混ざる**
  - `tid.New` を呼ばずに、共通の空文字TIDにレスポンスを設定してしまっている。
    必ずテスト毎に `tid.New` すること。
- **フェイクサーバーの生成物が古い**
  - protobuf定義を変えたら `make genfakes` を実行して
    `fakeservers/<svc>/server.go` を再生成する。
