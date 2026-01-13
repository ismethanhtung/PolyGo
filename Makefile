.PHONY: all build run test clean docker swagger lint bench help

# Variables
APP_NAME := polygo
MAIN_PATH := ./cmd/server
BUILD_DIR := ./build
GO := go
GOFLAGS := -v
LDFLAGS := -ldflags="-w -s"

# Default target
all: build

## Build

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

build-linux: ## Build for Linux
	@echo "Building $(APP_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux $(MAIN_PATH)

build-darwin: ## Build for macOS
	@echo "Building $(APP_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin $(MAIN_PATH)

## Run

run: ## Run the application
	@echo "Starting $(APP_NAME)..."
	$(GO) run $(MAIN_PATH)

run-dev: ## Run with hot reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

## Test

test: ## Run tests
	@echo "Running tests..."
	$(GO) test -v ./tests/...

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	$(GO) test -v ./tests/unit/...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GO) test -v ./tests/integration/...

test-cover: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./tests/...

## Dependencies

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

tidy: ## Tidy go.mod
	$(GO) mod tidy

## Docker

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name $(APP_NAME) $(APP_NAME):latest

docker-stop: ## Stop Docker container
	docker stop $(APP_NAME) && docker rm $(APP_NAME)

docker-compose-up: ## Start with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop docker-compose
	docker-compose down

## Swagger

swagger: ## Generate Swagger documentation
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/server/main.go -o docs

## Code Quality

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt: ## Format code
	$(GO) fmt ./...
	@which goimports > /dev/null && goimports -w . || true

vet: ## Run go vet
	$(GO) vet ./...

## Clean

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	$(GO) clean

## Help

help: ## Show this help
	@echo "PolyGo - High-Performance Polymarket API Proxy"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
