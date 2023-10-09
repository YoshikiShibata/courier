compile:
	protoc proto/*/v1//*.proto \
		--go_out=./api \
		--go-grpc_out=./api \
		--go_opt=paths=import \
		--go-grpc_opt=paths=import \
		--proto_path=.

genfakes:
	go run cmd/logserver/main.go > fakeservers/v1/fakelogserver.go
	gofmt -w fakeservers/v1/fakelogserver.go
