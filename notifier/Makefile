LOCAL_BIN := $(CURDIR)/../bin

build:
	go build -o bin/notifier ./cmd/notifier

run:
	go run ./cmd/notifier/main.go

lint:
	@echo Running gocyclo linter...
	-$(LOCAL_BIN)/gocyclo -top 20 -over 5 -ignore "_test|_mock|testdata" .

	@echo Running gocognit linter...
	-$(LOCAL_BIN)/gocognit -top 20 -over 5 -ignore "_test|_mock|testdata" .
