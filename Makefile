.PHONY: proto-gen test lint install build clean

###############################################################################
###                                Protobuf                                 ###
###############################################################################

proto-gen:
	@echo "Generating protobuf files..."
	@buf generate

proto-lint:
	@buf lint

###############################################################################
###                                  Build                                  ###
###############################################################################

build:
	@echo "Building..."
	@go build -o bin/kudosd ./cmd/kudosd

install:
	@echo "Installing..."
	@go install ./cmd/kudosd

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf build/

###############################################################################
###                                  Tests                                  ###
###############################################################################

test:
	@echo "Running tests..."
	@go test -v ./x/kudos/...

test-cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.txt -covermode=atomic ./x/kudos/...

###############################################################################
###                                 Linting                                 ###
###############################################################################

lint:
	@echo "Running linters..."
	@golangci-lint run --timeout 5m

format:
	@echo "Formatting code..."
	@gofmt -s -w .
	@goimports -w .

###############################################################################
###                                  Help                                   ###
###############################################################################

help:
	@echo "Available targets:"
	@echo "  proto-gen     - Generate protobuf files"
	@echo "  proto-lint    - Lint protobuf files"
	@echo "  build         - Build the binary"
	@echo "  install       - Install the binary"
	@echo "  test          - Run tests"
	@echo "  test-cover    - Run tests with coverage"
	@echo "  lint          - Run linters"
	@echo "  format        - Format code"
	@echo "  clean         - Clean build artifacts"
