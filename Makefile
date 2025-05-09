# Variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOLINT := golangci-lint
BINARY_NAME := archiver

.DEFAULT_GOAL := help

.PHONY: all build test lint clean help

all: build

# Build the binary
build:
	mkdir -p build
	$(GOBUILD) -o build/$(BINARY_NAME) -v ./cmd/archiver

# Run tests with verbose output
test:
	$(GOTEST) -v ./...

# Run linter (golangci-lint or staticcheck)
lint:
	@which staticcheck > /dev/null || (echo "staticcheck not found, please install it"; exit 1)
	staticcheck ./...
	# @which $(GOLINT) > /dev/null || (echo "golangci-lint not found, please install it"; exit 1)
	# $(GOLINT) run

# Clean build artifacts
clean:
	rm -f build/$(BINARY_NAME)

# Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build    - Build the Go binary"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linter (golangci-lint)"
	@echo "  clean    - Remove build artifacts"
	@echo "  help     - Show this help message"
