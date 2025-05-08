SHELL:=bash

.DEFAULT_GOAL := help

# Variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -w -s -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

##@ Help
help: ## Show this help message
	@clear
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[0;33m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""

##@ Build
.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod verify

.PHONY: build
build: deps ## Build the application (development)
	rm -rf build && mkdir build
	CGO_ENABLED=0 go build -o build/img_generator -v -ldflags "$(LDFLAGS)" ./cmd

.PHONY: build-prod
build-prod: deps ## Build the application (production)
	rm -rf build && mkdir build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/img_generator -v -ldflags "$(LDFLAGS)" -trimpath ./cmd

##@ Run
.PHONY: run
run: ## Run the application
	go run cmd/main.go

##@ Development
.PHONY: clean
clean: ## Clean build directory
	rm -rf build/

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: lint-install
lint-install: ## Install golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

.PHONY: lint
lint: ## Run linter with full output
	golangci-lint run --verbose --timeout=5m --max-same-issues=0 --max-issues-per-linter=0

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
