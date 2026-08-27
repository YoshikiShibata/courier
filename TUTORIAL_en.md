# Tutorial: Using the `server/fakegrpc` and `tid` Packages

This tutorial walks through the two core packages of courier, using the
sample project under `example/` as the running reference.

- `github.com/YoshikiShibata/courier/server/fakegrpc`
  A library for standing up gRPC dependencies as **fake servers** and
  configuring their responses from test code.
- `github.com/YoshikiShibata/courier/tid`
  A library that issues a **Testing ID (TID)** and propagates it through
  gRPC metadata so that parallel E2E tests do not interfere with each
  other's fake-server configurations.

The `example/` directory contains a Shop service (the system under test)
that depends on three gRPC services — Shipping, Warehouse, and Publisher —
and serves as a complete reference for both packages.

---

## Big Picture

An E2E test in this framework wires up as follows:

```
                          ┌───────────────────────────────┐
                          │  E2E test (list_products_test)│
                          └────────────┬──────────────────┘
                                       │ 1) generate TID via tid.New
                                       │ 2) fakeXxxServer.SetXxxResponse(tid, ...)
                                       │ 3) client.Xxx(ctx, req)
                                       ▼
                          ┌───────────────────────────────┐
                          │  Service under test (Shop)    │
                          │  tid.NewGRPCHeaderPropagator  │
                          │  copies TID onto outgoing ctx │
                          └────────────┬──────────────────┘
                                       │ downstream gRPC call (with TID header)
                                       ▼
                          ┌───────────────────────────────┐
                          │ Single fake gRPC server       │
                          │  fakegrpc.NewGRPCServer()     │
                          │  ├─ Shipping fake             │
                          │  ├─ Warehouse fake            │
                          │  └─ Publisher fake            │
                          │  Look up response by TID      │
                          └───────────────────────────────┘
```

The key idea: even though the fake server is a **single process-wide
instance**, every test registers responses under its own unique `tid`,
so tests can safely run with `t.Parallel()` without stepping on each
other.

---

## 1. The `tid` Package

`tid.go` exposes only two APIs.

### 1.1 `tid.New(ctx) (string, context.Context)`

Generates a fresh TID and returns a context with that TID attached as an
outgoing gRPC metadata header named `courier-testing-id`. If the given
context already contains a TID, the existing one is reused.

Test code calls this at the top of each test:

```go
// from example/e2etest/list_products_test.go
tid, ctx := tid.New(context.Background())

fakeWarehouseServer.SetListProductInventoriesResponse(tid,
    &warehouse_v1.ListProductInventoriesResponse{ /* ... */ }, nil)

res, err := client.ListProductInventories(ctx,
    &shop_v1.ListProductInventoriesRequest{ /* ... */ })
```

- The `tid` string is the key you use to register "the response for this
  test".
- The `ctx` carries the TID header, so passing it to
  `client.ListProductInventories(ctx, ...)` makes the TID reach the Shop
  service.

### 1.2 `tid.NewAlways(ctx) (string, context.Context)`

Same signature as `New`, but **always issues a fresh TID even if the
context already contains one**. Any existing TID in the outgoing metadata
is replaced.

Typical use cases:

- Loops that run a scenario repeatedly and need an independent TID per
  iteration for distinct response setups.
- A setup step has already attached a TID to the context, but a sub-test
  wants to switch to a different one.

```go
_, ctx := tid.New(context.Background())

for i := 0; i < 3; i++ {
    // Issue a fresh TID for every iteration
    tid, ctx := tid.NewAlways(ctx)
    fakeWarehouseServer.SetShipProductResponse(tid, resp[i], nil)
    // ... exercise the client with ctx ...
}
```

Ordinary tests should keep using `New`. Reach for `NewAlways` only when
you want to reuse an existing context but refresh the TID.

### 1.3 `tid.Extract(ctx) string`

Extracts a TID from incoming gRPC metadata. `GRPCServerStub` uses this
internally to decide "which test's response should I return?", so you
normally do not need to call it from test code.

### 1.4 `tid.NewGRPCHeaderPropagator() grpc.UnaryServerInterceptor`

A unary interceptor to install into the **service under test**. It reads
the TID from incoming metadata and re-attaches it to the outgoing context
of downstream gRPC calls.

```go
// from example/main.go
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

Without this interceptor the TID reaches the Shop service but is dropped
when Shop calls Shipping/Warehouse. The fake server then cannot look up
a matching response and panics with
`Response for :RpcName has not been set yet`.

---

## 2. The `server/fakegrpc` Package

This package provides two capabilities:

- **Fake-server code generation** (`Generate` / `TemplateConfig`)
- **Runtime primitives for a fake server** (`GRPCServer`, `GRPCServerStub`,
  `HandleRequest`, `RequestCapture`, `ResponseBuilder`)

### 2.1 Generating a Fake Server

Write one generator script per gRPC dependency:

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

`TemplateConfig` fields:

| Field | Meaning |
| --- | --- |
| `ServerPackageName` | The `package` declaration of the generated file |
| `ProtoPackageImportPath` | Import path of the generated protobuf/gRPC package |
| `ProtoPackageImportName` | Import alias used inside the template |
| `ServiceName` | The service name (e.g. `Shipping`); used as a prefix for the generated helpers |
| `ServiceClientType` | Result of `reflect.TypeOf(NewXxxClient(nil))`; RPC list is discovered from this type |

`Generate` writes the fake-server source to stdout, so the Makefile
redirects it into place:

```makefile
# from example/Makefile
SERVICES = shipping_v1 warehouse_v1 publisher_v1

genfakes: $(SERVICES)

$(SERVICES):
	go run scripts/$@/main.go > fakeservers/$@/server.go
	gofmt -w fakeservers/$@/server.go
	goimports -w fakeservers/$@/server.go
```

For example, when the `Shipping` service has RPCs `Create` and `Status`,
the generated file (`example/fakeservers/shipping_v1/server.go`) exposes:

- `NewShippingServer(grpcServer *grpc.Server) *ShippingServer`
  Registers the fake implementation on `grpcServer`.
- `(*ShippingServer).SetCreateResponse(tid, res, err)`
  Returns a fixed response for `Create` when it is called under this
  `tid`.
- `(*ShippingServer).SetCreateResponseCreator(tid, func(ctx, req) (res, err))`
  Same as above but as a function — useful when the response depends on
  the request, or when you want to capture the incoming request.
- `(*ShippingServer).ClearAllResponses(tid)`
  Removes every RPC response registered under this `tid`.

### 2.2 Starting the Fake Server

`GRPCServer` is a thin wrapper that assigns a free port and builds a
gRPC server.

```go
// from example/e2etest/main_test.go
grpcServer := fakegrpc.NewGRPCServer()
server := grpcServer.Server()
port := grpcServer.Port()

// Host three fake services on a single gRPC server
fakeShippingServer  = fakeshipping_v1.NewShippingServer(server)
fakeWarehouseServer = fakewarehouse_v1.NewWarehouseServer(server)
fakePublisherServer = fakepublisher_v1.NewPublisherServer(server)

go func() { grpcServer.Serve() }()
```

Point the service under test at the fake server via environment variables:

```go
shippingSvcEnv  := fmt.Sprintf("SHIPPING_SERVICE_ADDR=localhost:%s", port)
warehouseSvcEnv := fmt.Sprintf("WAREHOUSE_SERVICE_ADDR=localhost:%s", port)
publisherSvcEnv := fmt.Sprintf("PUBSUB_EMULATOR_HOST=localhost:%s", port)
```

The Pub/Sub emulator is a gRPC service too, so it can be faked with the
same mechanism (`example/scripts/publisher_v1/main.go`).

### 2.3 Registering Responses from Tests

Pattern 1 — return a **fixed response**:

```go
tid, ctx := tid.New(context.Background())

fakeWarehouseServer.SetListProductInventoriesResponse(tid,
    &warehouse_v1.ListProductInventoriesResponse{
        ProductInventories: warehouseProducts,
        NextPageToken:      "",
    }, nil)
```

Pattern 2 — return **an error only**:

```go
fakeWarehouseServer.SetListProductInventoriesResponse(tid,
    nil, status.Error(codes.InvalidArgument, "page_token"))
```

Pattern 3 — return **different responses across repeated calls** to the
same RPC, using `fakegrpc.ResponseBuilder`:

```go
rb := &fakegrpc.ResponseBuilder[
    *warehouse_v1.ShipProductRequest,
    *warehouse_v1.ShipProductResponse,
]{}
rb.Append(&warehouse_v1.ShipProductResponse{ /* 1st call */ }, nil)
rb.Append(nil, status.Error(codes.Unavailable, "retry"))
rb.Append(&warehouse_v1.ShipProductResponse{ /* 3rd call */ }, nil)

fakeWarehouseServer.SetShipProductResponseCreator(tid, rb.Creator)
```

`Creator` returns the appended responses in order. It uses
`atomic.AddUint32` internally and is safe against concurrent invocations.

Pattern 4 — **inspect the request** the service under test sent, using
`fakegrpc.RequestCapture`:

```go
rc := fakegrpc.NewRequestCapture[
    *warehouse_v1.ShipProductRequest,
    *warehouse_v1.ShipProductResponse,
](&warehouse_v1.ShipProductResponse{ /* the response to return */ }, nil)

fakeWarehouseServer.SetShipProductResponseCreator(tid, rc.Creator)

// ... exercise the service under test ...

got := rc.Get() // Inspect the ShipProductRequest that Shop sent
```

### 2.4 Cleanup Between Tests

You *can* call `ClearAllResponses(tid)` after a test if you want, but
because each test gets a unique TID, leaving old entries in place does
not affect other tests. `ClearAllResponses` is only useful when a single
loop reuses the same TID across iterations.

---

## 3. Wiring Up an E2E Test — Step by Step

Follow the pattern of `example/` when adopting this framework in your
own project.

### 3.1 Service Under Test (Production Code)

Just install `tid.NewGRPCHeaderPropagator()` as a unary interceptor when
you build `grpc.NewServer` (`example/main.go`):

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

### 3.2 Generator Scripts

Write one `scripts/<svc>/main.go` per dependency
(`example/scripts/*/main.go` are ready-made templates).

The Makefile runs each via `go run`, writes the output to
`fakeservers/<svc>/server.go`, and formats it with `gofmt` / `goimports`.

### 3.3 `TestMain` for E2E

Reuse the structure of `example/e2etest/main_test.go`:

1. `fakegrpc.NewGRPCServer()` — start the fake gRPC server.
2. Register each `fakeXxxServer` on it.
3. Point the service-under-test's dependency env vars at the fake server.
4. Launch the service under test as a child process
   (`courier.InvokeService`).
5. Run the tests with `m.Run()`.

### 3.4 Per-test Skeleton

```go
func TestXxx(t *testing.T) {
    t.Parallel()                              // safe to run in parallel
    client := newShopClient(t)
    tid, ctx := tid.New(context.Background()) // TID for this test

    // Configure downstream responses
    fakeWarehouseServer.SetListProductInventoriesResponse(tid, resp, nil)

    // Exercise the service under test (ctx carries the TID)
    got, err := client.ListProductInventories(ctx, req)

    // Assertions
    // ...
}
```

The fake server will use the `TID` to route to exactly the response this
test registered.

---

## 4. Common Pitfalls

- **Panic: `Response for <tid>:<RPC> has not been set yet`**
  - The service under test forgot to install
    `tid.NewGRPCHeaderPropagator()`.
  - The test used a bare `context.Background()` for the gRPC call
    instead of the ctx returned by `tid.New`.
  - You forgot to call `SetXxxResponse` for one of the RPCs that gets
    invoked.
- **Parallel tests see each other's responses**
  - You skipped `tid.New` and registered responses under an empty TID.
    Always call `tid.New` per test.
- **Stale generated fake servers**
  - When you change the `.proto`, run `make genfakes` to regenerate
    `fakeservers/<svc>/server.go`.
