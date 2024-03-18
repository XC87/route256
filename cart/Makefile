LOCAL_BIN := $(CURDIR)/bin

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

lint:
	@echo Running gocyclo linter...
	-$(LOCAL_BIN)/gocyclo -top 20 -over 5 -ignore "_test|_mock|testdata" .

	@echo Running gocognit linter...
	-$(LOCAL_BIN)/gocognit -top 20 -over 5 -ignore "_test|_mock|testdata" .

install:
	set GOBIN=$(LOCAL_BIN) && set GOPROXY=direct && go install github.com/uudashr/gocognit/cmd/gocognit@latest
	set GOBIN=$(LOCAL_BIN) && set GOPROXY=direct && go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

