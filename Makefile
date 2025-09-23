# Makefile for the pbar project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Linting
GOLANGCILINT=golangci-lint

# Binary name
BINARY_NAME=pbar

.PHONY: all test build clean fmt vet lint

all: build

# Build the binary for the current platform
build:
	go build -o pbar main.go

# Run all tests
test: 
	$(GOTEST) -v ./...

# Run tests with coverage
coverage: 
	$(GOTEST) -cover ./...

# Format the code
fmt: 
	$(GOFMT) ./...

# Run go vet
vet: 
	$(GOVET) ./...

# Run linter
lint: vet
	$(GOLANGCILINT) run ./...

# Cross-compile for all target platforms
release: 
	@./build.sh

# Clean up build artifacts
clean:
	rm -f pbar.exe
	rm -rf dist/

# Install the binary
install:
	$(GOCMD) install .

# Help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build the binary for the current platform"
	@echo "  test        Run all tests"
	@echo "  coverage    Run tests with code coverage"
	@echo "  fmt         Format the code"
	@echo "  vet         Run go vet"
	@echo "  lint        Run linter (includes vet)"
	@echo "  release     Cross-compile for all target platforms"
	@echo "  clean       Clean up build artifacts"
	@echo "  install     Install the binary"
	@echo "  help        Show this help message"
