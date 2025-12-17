.PHONY: help all build clean demo server client test

# Colors for output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

help: ## Show this help message
	@echo '$(CYAN)Help$(RESET)'
	@echo ''
	@echo 'Usage:'
	@echo '  $(YELLOW)make$(RESET) $(GREEN)<target>$(RESET)'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: build ## Default target

build: server client demo ## Build all binaries

server: ## Build server binary
	@echo "Building server..."
	@go build -o bin/server cmd/server/main.go
	@echo "✅ Server built: bin/server"

client: ## Build client binary
	@echo "Building client..."
	@go build -o bin/client cmd/client/main.go
	@echo "✅ Client built: bin/client"

demo: ## Build demo binary
	@echo "Building demo..."
	@go build -o bin/demo main.go
	@echo "✅ Demo built: bin/demo"

run-demo: ## Run the demo
	@echo "Running demo..."
	@go run main.go

run-server: ## Run server
	@echo "Starting server..."
	@go run cmd/server/main.go

run-client: ## Run client
	@echo "Starting client..."
	@go run cmd/client/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f server client demo server_bin client_bin
	@echo "✅ Cleaned"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint: ## Lint code
	@echo "Linting code..."
	@golint ./...

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"