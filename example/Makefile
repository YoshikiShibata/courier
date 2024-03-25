COURIER=github.com/YoshikiShibata/courier
COVER_PKGS=$(COURIER)/example,$(COURIER)/example/internal/...
VALIDATE_PROTO_PATH=$(HOME)/tools/protoc-gen-validate/
SERVICES = shipping_v1 warehouse_v1 publisher_v1

all: compile genfakes

.PHONY: compile
compile:
	protoc proto/*/v1//*.proto \
		-I $(VALIDATE_PROTO_PATH) \
		--go_out=./api \
		--go-grpc_out=./api \
		--go_opt=paths=import \
		--go-grpc_opt=paths=import \
		--validate_out="lang=go:./api" \
		--proto_path=.

.PHONY: genfakes
genfakes: $(SERVICES)

$(SERVICES):
	go run scripts/$@/main.go > fakeservers/$@/server.go
	gofmt -w fakeservers/$@/server.go
	goimports -w fakeservers/$@/server.go


# 手作業で指定することは想定されていません。
coverage_build:
	rm -fr coverage
	mkdir coverage
	go build -race -cover -o bin/shop_server -coverpkg=$(COVER_PKGS) .

clean:
	rm -fr coverage
	rm -fr bin/shop_server profile.txt
