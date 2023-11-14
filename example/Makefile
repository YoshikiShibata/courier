COURIER=github.com/YoshikiShibata/courier
COVER_PKGS=$(COURIER)/example,$(COURIER)/example/internal/...

all: compile genfakes

compile:
	protoc proto/*/v1//*.proto \
		--go_out=./api \
		--go-grpc_out=./api \
		--go_opt=paths=import \
		--go-grpc_opt=paths=import \
		--proto_path=.

genfakes:
	go run scripts/shipping_v1/main.go > fakeservers/shipping_v1/server.go
	gofmt -w fakeservers/shipping_v1/server.go
	goimports -w fakeservers/shipping_v1/server.go
	go run scripts/warehouse_v1/main.go > fakeservers/warehouse_v1/server.go
	gofmt -w fakeservers/warehouse_v1/server.go
	goimports -w fakeservers/warehouse_v1/server.go
	go run scripts/publisher_v1/main.go > fakeservers/publisher_v1/server.go
	gofmt -w fakeservers/publisher_v1/server.go
	goimports -w fakeservers/publisher_v1/server.go

# 手作業で指定することは想定されていません。
build:
	rm -fr coverage
	mkdir coverage
	go build -race -cover -o bin/shop_server -coverpkg=$(COVER_PKGS) .

clean:
	rm -fr coverage
	rm -fr bin/shop_server profile.txt
