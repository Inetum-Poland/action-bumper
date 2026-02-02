# Copyright (c) 2024 Inetum Poland.

.PHONY: help build test test-verbose test-coverage lint clean run

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building bumper..."
	@go build -o bin/bumper ./cmd/bumper
	@echo "Binary built: bin/bumper"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	@go test -v -race -count=1 ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping..."; \
		echo "Install: https://golangci-lint.run/usage/install/"; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

run: build ## Build and run locally (requires env vars)
	@echo "Running bumper..."
	@./bin/bumper

all: tidy fmt vet test build ## Run all checks and build

.DEFAULT_GOAL := help
