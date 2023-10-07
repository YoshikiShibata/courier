compile:
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

genfakes:
	go run cmd/logserver/main.go > fakeservers/v1/fakelogserver.go
	gofmt -w fakeservers/v1/fakelogserver.go
