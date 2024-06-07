.DEFAULT_GOAL := help

fmtCheck:
	gofmt -l -s .

fmt:
	gofmt -l -s -w .

build: 
	go build -v -o bin/server cmd/server/main.go
	go build -v -o bin/client cmd/client/main.go

clean:
	go clean
	rm -rf ./bin

lint:
	golangci-lint run

test:
	go test -v ./...

help:
	@echo "Available targets"
	@echo "\tfmtCheck         - Check if code is correctly formatted"
	@echo "\tfmt              - Format code"
	@echo "\tbuild            - Build all binaries"
	@echo "\tclean            - Clean up bin/ directory"
	@echo "\tlint             - Run linter"
	@echo "\ttest             - Run all tests"

