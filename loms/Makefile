include ../mk/vendor-proto.mk
include ../mk/proto_compile.mk

build:
	go build -o bin/loms-server ./cmd/loms-server

run:
	go run ./cmd/loms-server/main.go

cover:
	go test -v ./internal/service/loms -covermode count -coverprofile=coverage.out.tmp
	grep -vE "_mock.go" coverage.out.tmp > coverage.out
	go tool cover -html=coverage.out
	del coverage.out.tmp

test:
	go test -v -parallel 45 -race  ./internal/service/loms

lint:
	@echo Running gocyclo linter...
	-$(LOCAL_BIN)gocyclo -top 20 -over 5 -ignore "_test|_mock|testdata|.pb." .

	@echo Running gocognit linter...
	-$(LOCAL_BIN)gocognit -top 20 -over 5 -ignore "_test|_mock|testdata|.pb." .

sqlc:
	sqlc generate

migrate-up:
	goose -dir migrations postgres "postgresql://postgres:password@localhost:5442/loms?sslmode=disable" up
migrate-reset:
	goose -dir migrations postgres "postgresql://postgres:password@localhost:5442/loms?sslmode=disable" reset

PROTO_PATH:="api/v1"
PROTO_FILE:="/loms.proto"
SWAGGER_PATH:="--openapiv2_out=api/openapiv2"
# компиляция proto файлов
protoc: .protoc-generate

# Вендоринг внешних proto файлов
vendor-proto: .vendor-rm  vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options
