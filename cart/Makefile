include ../mk/vendor-proto.mk
include ../mk/proto_compile.mk

build:
	go build -o bin/cart-server ./cmd/cart-server

run:
	go run ./cmd/cart-server/main.go

cover:
	go test -v ./internal/service ./internal/repository -covermode count -coverprofile=coverage.out.tmp
	grep -vE "_mock.go" coverage.out.tmp > coverage.out
	go tool cover -html=coverage.out
	del coverage.out.tmp

test:
	go test -v -parallel 45 -race  ./internal/repository
	go test -v -parallel 45 -race  ./internal/service
test_intergration:
	go test -v -parallel 45 -race  ./tests/main_integration_test.go

lint:
	@echo Running gocyclo linter...
	-$(LOCAL_BIN)/gocyclo -top 20 -over 5 -ignore "_test|_mock|testdata" .

	@echo Running gocognit linter...
	-$(LOCAL_BIN)/gocognit -top 20 -over 5 -ignore "_test|_mock|testdata" .

PROTO_PATH:="api/v1"
PROTO_FILE:="/loms.proto"
# компиляция proto файлов
protoc: .protoc-generate

# Вендоринг внешних proto файлов
vendor-proto: .vendor-rm  vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options

