build:
	go build -o bin/cart-server ./cmd/cart-server

run:
	go run ./cmd/cart-server/main.go

cover:
	go test -v ./internal/service -covermode count -coverprofile=coverage.out.tmp
	grep -vE "_mock.go" coverage.out.tmp > coverage.out
	go tool cover -html=coverage.out
	del coverage.out.tmp

linters:
	@echo Running gocyclo linter...
	-gocyclo -top 20 -over 5 -ignore "_test|_mock|testdata" .

	@echo Running gocognit linter...
	-gocognit -top 20 -over 5 -ignore "_test|_mock|testdata" .



