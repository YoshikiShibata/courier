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
	go run scripts/warehouse_v1/main.go > fakeservers/warehouse_v1/server.go
	gofmt -w fakeservers/warehouse_v1/server.go
