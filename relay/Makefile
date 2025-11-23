.PHONY: help build run dev docker-build docker-run docker-stop docker-compose-up docker-compose-down clean test lint fmt install-air mod-tidy health-check

# Variables
BINARY_NAME=relay
MAIN_PATH=./cmd/main.go
GO=go
DOCKER_IMAGE_NAME=relay
DOCKER_IMAGE_TAG=latest
PORT?=8080

# Colors for output
BLUE=\033[0;34m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "$(BLUE)Ethereum Relay Service - Makefile Commands$(NC)"
	@echo ""
	@echo "$(GREEN)Development$(NC)"
	@echo "  make dev              - Run with hot reload using air"
	@echo "  make build            - Build the binary locally"
	@echo "  make run              - Run the compiled binary"
	@echo "  make clean            - Remove build artifacts"
	@echo ""
	@echo "$(GREEN)Docker$(NC)"
	@echo "  make docker-build     - Build production Docker image"
	@echo "  make docker-run       - Run container (requires .env file)"
	@echo "  make docker-stop      - Stop running containers"
	@echo "  make docker-compose-up    - Start with docker-compose (hot reload)"
	@echo "  make docker-compose-down  - Stop docker-compose services"
	@echo ""
	@echo "$(GREEN)Code Quality$(NC)"
	@echo "  make lint             - Run Go linter"
	@echo "  make fmt              - Format code with gofmt"
	@echo "  make vet              - Run go vet"
	@echo ""
	@echo "$(GREEN)Dependencies$(NC)"
	@echo "  make install-air      - Install air for hot reload"
	@echo "  make mod-tidy         - Tidy go modules"
	@echo "  make mod-download     - Download dependencies"
	@echo ""
	@echo "$(GREEN)Testing$(NC)"
	@echo "  make health-check     - Check if service is running"
	@echo "  make test-relay       - Send test relay request"
	@echo ""
	@echo "$(GREEN)Info$(NC)"
	@echo "  make version          - Show Go version"
	@echo "  make env              - Show environment variables from .env"
	@echo ""

## build: Build the binary locally
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	$(GO) build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Binary built successfully: ./$(BINARY_NAME)$(NC)"

## run: Run the compiled binary
run: build
	@echo "$(BLUE)Starting $(BINARY_NAME)...$(NC)"
	@./$(BINARY_NAME)

## dev: Run with hot reload using air
dev: install-air
	@echo "$(BLUE)Starting with hot reload (air)...$(NC)"
	@air -c .air.toml

## install-air: Install air for hot reload development
install-air:
	@echo "$(BLUE)Installing air...$(NC)"
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@echo "$(GREEN)✓ Air installed$(NC)"

## mod-download: Download Go dependencies
mod-download:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	$(GO) mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

## mod-tidy: Tidy go modules
mod-tidy:
	@echo "$(BLUE)Tidying modules...$(NC)"
	$(GO) mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

## docker-build: Build production Docker image
docker-build:
	@echo "$(BLUE)Building Docker image: $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)...$(NC)"
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -f Dockerfile .
	@echo "$(GREEN)✓ Docker image built successfully$(NC)"

## docker-run: Run Docker container
docker-run: docker-build
	@if [ ! -f .env ]; then \
		echo "$(RED)Error: .env file not found$(NC)"; \
		echo "Please copy .env.example to .env and fill in your values"; \
		exit 1; \
	fi
	@echo "$(BLUE)Running Docker container...$(NC)"
	docker run -p $(PORT):8080 --env-file .env $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

## docker-stop: Stop running Docker containers
docker-stop:
	@echo "$(BLUE)Stopping Docker containers...$(NC)"
	@docker ps -q --filter "ancestor=$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" | xargs -r docker stop
	@echo "$(GREEN)✓ Containers stopped$(NC)"

## docker-compose-up: Start services with docker-compose
docker-compose-up:
	@if [ ! -f .env ]; then \
		echo "$(RED)Error: .env file not found$(NC)"; \
		echo "Please copy .env.example to .env and fill in your values"; \
		exit 1; \
	fi
	@echo "$(BLUE)Starting services with docker-compose...$(NC)"
	docker-compose up

## docker-compose-down: Stop docker-compose services
docker-compose-down:
	@echo "$(BLUE)Stopping docker-compose services...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

## clean: Remove build artifacts
clean:
	@echo "$(BLUE)Cleaning up...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -rf tmp/
	@rm -f build-errors.log
	@$(GO) clean
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

## fmt: Format code with gofmt
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

## lint: Run Go linter
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run ./...
	@echo "$(GREEN)✓ Linting complete$(NC)"

## vet: Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GO) vet ./...
	@echo "$(GREEN)✓ Vet check complete$(NC)"

## version: Show Go version
version:
	@echo "$(BLUE)Go version:$(NC)"
	@$(GO) version

## env: Show environment variables from .env
env:
	@echo "$(BLUE)Environment variables from .env:$(NC)"
	@if [ -f .env ]; then \
		cat .env; \
	else \
		echo "$(RED).env file not found. Use: cp .env.example .env$(NC)"; \
	fi

## health-check: Check if service is running
health-check:
	@echo "$(BLUE)Checking service health...$(NC)"
	@curl -s http://localhost:$(PORT)/health | jq . 2>/dev/null || echo "$(RED)Service is not responding on port $(PORT)$(NC)"

## test-relay: Send test relay request
test-relay:
	@echo "$(BLUE)Sending test relay request...$(NC)"
	@curl -X POST http://localhost:$(PORT)/relay \
		-H "Content-Type: application/json" \
		-d '{"ephemeralPublicKey":"test_key_12345","burnAddress":"0x0000000000000000000000000000000000000000"}' \
		| jq . 2>/dev/null || echo "$(RED)Failed to send request$(NC)"

## setup: Initial project setup
setup: mod-download install-air
	@echo "$(BLUE)Setting up project...$(NC)"
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)Creating .env from .env.example...$(NC)"; \
		cp .env.example .env; \
		echo "$(YELLOW)Please edit .env with your configuration$(NC)"; \
	else \
		echo "$(GREEN)✓ .env already exists$(NC)"; \
	fi
	@echo "$(GREEN)✓ Setup complete$(NC)"

## all: Install dependencies, format, lint, and build
all: fmt vet mod-tidy build
	@echo "$(GREEN)✓ All checks passed and binary built$(NC)"

# Development workflow targets
.PHONY: setup all

# Quick targets
.PHONY: b
b: build

.PHONY: r
r: run

.PHONY: d
d: dev

.PHONY: dcu
dcu: docker-compose-up

.PHONY: dcd
dcd: docker-compose-down

.PHONY: hc
hc: health-check

.PHONY: tr
tr: test-relay
