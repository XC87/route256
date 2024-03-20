ifeq ($(OS),Windows_NT)
    EXE = .exe
endif

PHONY: .protoc-generate
.protoc-generate:
	protoc \
	-I ${PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=../bin/protoc-gen-go$(EXE) \
	--go_out pkg/${PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=../bin/protoc-gen-go-grpc$(EXE) \
	--go-grpc_out pkg/${PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=../bin/protoc-gen-validate$(EXE) \
	--validate_out="lang=go,paths=source_relative:pkg/${PROTO_PATH}" \
	--plugin=protoc-gen-grpc-gateway=../bin/protoc-gen-grpc-gateway$(EXE) \
	--grpc-gateway_out pkg/${PROTO_PATH} \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=../bin/protoc-gen-openapiv2$(EXE) \
	${SWAGGER_PATH} \
	--openapiv2_opt logtostderr=true \
	${PROTO_PATH}${PROTO_FILE}
	go mod tidy

.PHONY: .buf-generate
.buf-generate:
	../bin/buf$(EXE) mod update
	../bin/buf$(EXE) generate



