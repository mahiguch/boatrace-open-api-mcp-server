.PHONY: build test lint clean help

BINARY_NAME=boatrace-mcp
BINARY_PATH=bin/$(BINARY_NAME)

help:
	@echo "Available targets:"
	@echo "  build       - Build the MCP server binary"
	@echo "  test        - Run all tests"
	@echo "  lint        - Run golangci-lint"
	@echo "  clean       - Clean build artifacts"

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) ./cmd/boatrace-mcp
	@echo "Build complete: $(BINARY_PATH)"

test:
	@echo "Running tests..."
	@go test -v -cover ./...

lint:
	@echo "Running go vet..."
	@go vet ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean
	@echo "Clean complete"
