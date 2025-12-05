.PHONY: help build run test lint lint-fix fmt vet deps test-coverage analyse clean docker-build docker-run

# Variables
BINARY_NAME=gokode
CMD_PATH=./cmd/gokode
METRICS_DIR=metrics
GOLANGCI_LINT_VERSION=v1.55.2
GOPATH=$(shell go env GOPATH)
GOLANGCI_LINT=$(GOPATH)/bin/golangci-lint

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the gokode binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BINARY_NAME)"

run: build ## Build and run the gokode tool on current directory
	@echo "Running $(BINARY_NAME) analyse ."
	@./$(BINARY_NAME) analyse .

test: ## Run all tests
	@echo "Running tests..."
	@go test ./... -v

lint: deps ## Run golangci-lint
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT) run ./...

lint-fix: deps ## Run golangci-lint with --fix
	@echo "Running golangci-lint with --fix..."
	@$(GOLANGCI_LINT) run --fix ./...

fmt: ## Format code with gofmt
	@echo "Formatting code..."
	@gofmt -s -w .
	@echo "Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "go vet passed"

deps: ## Install dependencies including golangci-lint
	@echo "Checking dependencies..."
	@if ! [ -f $(GOLANGCI_LINT) ]; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	else \
		echo "golangci-lint is already installed"; \
	fi
	@go mod download
	@go mod tidy
	@echo "Dependencies ready"

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=$(METRICS_DIR)/coverage.out
	@go tool cover -html=$(METRICS_DIR)/coverage.out -o $(METRICS_DIR)/coverage.html
	@echo "Coverage report generated: $(METRICS_DIR)/coverage.html"

analyse: fmt vet lint-fix ## Run full analysis: format, vet, and lint with fixes
	@echo "Running full analysis..."
	@mkdir -p $(METRICS_DIR)
	@echo "Running go vet and writing to $(METRICS_DIR)/vet.txt..."
	@go vet ./... 2>&1 | tee $(METRICS_DIR)/vet.txt || true
	@echo "Running golangci-lint with JSON output to $(METRICS_DIR)/report.json..."
	@$(GOLANGCI_LINT) run --out-format json ./... > $(METRICS_DIR)/report.json 2>&1 || true
	@echo "Analysis complete. Reports written to $(METRICS_DIR)/"

clean: ## Clean build artifacts and metrics
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(METRICS_DIR)/*.json $(METRICS_DIR)/*.txt $(METRICS_DIR)/*.out $(METRICS_DIR)/*.html
	@echo "Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

docker-run: docker-build ## Build and run Docker container
	@echo "Running Docker container..."
	@docker run --rm $(BINARY_NAME):latest
