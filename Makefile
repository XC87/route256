LOCAL_BIN := $(CURDIR)/bin

build-all:
	cd cart && make build
	cd loms && make build
	cd notifier && make build

run-all:
	docker-compose up --force-recreate --build

.PHONY: .bin-deps
.install-bin-deps:
	@echo Installing binary dependencies...
	set GOBIN=$(LOCAL_BIN)&& go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    set GOBIN=$(LOCAL_BIN)&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	set GOBIN=$(LOCAL_BIN)&& go install github.com/bufbuild/buf/cmd/buf@v1.21.0 && \
	set GOBIN=$(LOCAL_BIN)&& go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	set GOBIN=$(LOCAL_BIN)&& go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 && \
	set GOBIN=$(LOCAL_BIN)&& go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1 && \
	set GOBIN=$(LOCAL_BIN)&& go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5
	@echo Installing linters
	set GOBIN=$(LOCAL_BIN)&& go install github.com/uudashr/gocognit/cmd/gocognit@latest
	set GOBIN=$(LOCAL_BIN)&& go install github.com/fzipp/gocyclo/cmd/gocyclo@latest