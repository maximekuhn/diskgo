.DEFAULT_GOAL := help

fmtCheck:
	gofmt -l -s .

fmt:
	gofmt -l -s -w .

build: 
	go build -v -o bin/server cmd/server/main.go
	go build -v -o bin/client cmd/client/main.go

build_all: build
	# Raspberry Pi
	GOOS=linux GOARCH=arm GOARM=7 go build -v -o bin/rpi/server cmd/server/main.go
	GOOS=linux GOARCH=arm GOARM=7 go build -v -o bin/rpi/client cmd/client/main.go

	# MacOS (ARM)
	GOOS=darwin GOARCH=amd64 go build -v -o bin/macos/arm/server cmd/server/main.go
	GOOS=darwin GOARCH=amd64 go build -v -o bin/macos/arm/client cmd/client/main.go


clean:
	go clean
	rm -rf ./bin

lint:
	golangci-lint run

test:
	go test -tags=integration -v ./...

help:
	@echo "Available targets"
	@echo "\tfmtCheck         - Check if code is correctly formatted"
	@echo "\tfmt              - Format code"
	@echo "\tbuild            - Build all binaries for current platform"
	@echo "\tbuild_all        - Build all binaries for multiple platforms"
	@echo "\tclean            - Clean up bin/ directory"
	@echo "\tlint             - Run linter"
	@echo "\ttest             - Run all tests"

