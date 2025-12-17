.PHONY: help init proto build run test clean docker lint fmt vet install-tools

# Default target
.DEFAULT_GOAL := help

# Variables
APP_NAME := go-microservice-template
CMD_DIR := ./cmd/server
BIN_DIR := ./bin
PROTO_DIR := ./api/proto/v1
SWAGGER_DIR := ./docs/swagger

# Colors for output
COLOR_RESET := \033[0m
COLOR_BLUE := \033[34m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m

# Help target
help: ## Show this help message
	@echo "$(COLOR_BLUE)Available targets:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-15s$(COLOR_RESET) %s\n", $$1, $$2}'

init: ## Initialize project dependencies
	@echo "$(COLOR_BLUE)Installing Go dependencies...$(COLOR_RESET)"
	go mod download
	go mod tidy
	@echo "$(COLOR_GREEN)Dependencies installed$(COLOR_RESET)"

install-tools: ## Install required tools (buf, protoc plugins)
	@echo "$(COLOR_BLUE)Installing required tools...$(COLOR_RESET)"
	@echo "Installing buf..."
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@echo "Installing protoc-gen-go..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@echo "Installing protoc-gen-go-grpc..."
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Installing protoc-gen-grpc-gateway..."
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@echo "Installing protoc-gen-openapiv2..."
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@echo "$(COLOR_GREEN)All tools installed$(COLOR_RESET)"

proto: ## Generate code from proto files
	@echo "$(COLOR_BLUE)Generating code from proto files...$(COLOR_RESET)"
	@buf lint
	@buf generate
	@echo "$(COLOR_GREEN)Proto files generated$(COLOR_RESET)"

build: proto ## Build the application
	@echo "$(COLOR_BLUE)Building $(APP_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)
	@echo "$(COLOR_GREEN)Build complete: $(BIN_DIR)/$(APP_NAME)$(COLOR_RESET)"

run: build ## Build and run the application
	@echo "$(COLOR_BLUE)Starting $(APP_NAME)...$(COLOR_RESET)"
	@$(BIN_DIR)/$(APP_NAME) config/config.yaml

run-dev: proto ## Run without building (using go run)
	@echo "$(COLOR_BLUE)Running in development mode...$(COLOR_RESET)"
	@go run $(CMD_DIR)/main.go config/config.yaml

test: ## Run tests
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	@go test -v -race -cover ./...
	@echo "$(COLOR_GREEN)Tests complete$(COLOR_RESET)"

test-coverage: ## Run tests with coverage report
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)Coverage report generated: coverage.html$(COLOR_RESET)"

lint: ## Run linters
	@echo "$(COLOR_BLUE)Running linters...$(COLOR_RESET)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "$(COLOR_YELLOW)golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(COLOR_RESET)"; \
	fi

fmt: ## Format code
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	@go fmt ./...
	@echo "$(COLOR_GREEN)Code formatted$(COLOR_RESET)"

vet: ## Run go vet
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	@go vet ./...
	@echo "$(COLOR_GREEN)Vet complete$(COLOR_RESET)"

clean: ## Clean build artifacts
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@find . -name "*.pb.go" -type f -delete
	@find . -name "*.pb.gw.go" -type f -delete
	@find $(SWAGGER_DIR) -name "*.swagger.json" -type f -delete
	@echo "$(COLOR_GREEN)Clean complete$(COLOR_RESET)"

docker-build: ## Build Docker image
	@echo "$(COLOR_BLUE)Building Docker image...$(COLOR_RESET)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(COLOR_GREEN)Docker image built$(COLOR_RESET)"

docker-run: ## Run Docker container
	@echo "$(COLOR_BLUE)Running Docker container...$(COLOR_RESET)"
	@docker run -p 8080:8080 -p 9090:9090 $(APP_NAME):latest

mod-tidy: ## Tidy go modules
	@echo "$(COLOR_BLUE)Tidying go modules...$(COLOR_RESET)"
	@go mod tidy
	@echo "$(COLOR_GREEN)Modules tidied$(COLOR_RESET)"

mod-vendor: ## Vendor dependencies
	@echo "$(COLOR_BLUE)Vendoring dependencies...$(COLOR_RESET)"
	@go mod vendor
	@echo "$(COLOR_GREEN)Dependencies vendored$(COLOR_RESET)"

swagger: proto ## Open Swagger UI (requires server to be running)
	@echo "$(COLOR_BLUE)Opening Swagger UI...$(COLOR_RESET)"
	@open http://localhost:8080/swagger/ || xdg-open http://localhost:8080/swagger/ || echo "Please open http://localhost:8080/swagger/ in your browser"

grpcurl-list: ## List gRPC services (requires grpcurl and running server)
	@echo "$(COLOR_BLUE)Listing gRPC services...$(COLOR_RESET)"
	@grpcurl -plaintext localhost:9090 list

grpcurl-describe: ## Describe gRPC service (requires grpcurl and running server)
	@echo "$(COLOR_BLUE)Describing UserService...$(COLOR_RESET)"
	@grpcurl -plaintext localhost:9090 describe api.v1.UserService

all: clean install-tools init proto build test ## Run all build steps

check: fmt vet lint test ## Run all checks
