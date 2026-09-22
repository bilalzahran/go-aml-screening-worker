BINARY_NAME=cdaq-event-worker
GO=go
GOFLAGS=-v

VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

.PHONY: all build clean test lint run fmt deps docker-up docker-down setup-env help

all: clean lint test build

## build: Build the worker binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/worker

## clean: Remove build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

## test: Run all tests with race detector
test:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

## test-short: Run tests without integration tests
test-short:
	$(GO) test -v -short ./...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## fmt: Format and vet
fmt:
	$(GO) fmt ./...
	$(GO) vet ./...

## setup-env: Create .env from .env.example (copy only if .env doesn't exist)
setup-env:
	@if [ ! -f .env ]; then \
		echo "Creating .env from .env.example..."; \
		cp .env.example .env; \
		echo "⚠️  .env created. Please edit it with your actual API keys and secrets."; \
	else \
		echo ".env already exists, skipping."; \
	fi

## run: Build and run locally with environment variables
run: build setup-env
	@echo "Loading environment variables from .env..."
	@set -a && . ./.env && set +a && ./bin/$(BINARY_NAME)

## deps: Download and tidy dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

## env: Export environment variables from .env
env: setup-env
	@echo "Exporting environment variables from .env..."
	@set -a && . ./.env && set +a && echo "✓ Environment variables loaded"

## docker-up: Start local dev infrastructure
docker-up:
	docker compose up -d

## docker-down: Stop local dev infrastructure
docker-down:
	docker compose down

## docker-build: Build Docker image
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

## help: Show this help
help:
	@echo "Usage: make [target]"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
