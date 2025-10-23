.EXPORT_ALL_VARIABLES:
OUT_DIR := ./_output
BIN_DIR := ./bin
DIST_DIR := ./dist

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -w -s -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

$(shell mkdir -p $(OUT_DIR) $(BIN_DIR) $(DIST_DIR))

.DEFAULT_GOAL := help

##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: install
install: ## Download dependencies
	go mod download
	go mod verify
	go mod tidy

.PHONY: run
run: ## Run the application locally
	go run .

.PHONY: run-dev
run-dev: ## Run with docker-compose
	docker-compose up --build

.PHONY: seed
seed: ## Seed the database
	go run database/seeders/seeders.go

##@ Testing

.PHONY: test
test: ## Run unit tests with coverage
	go test -race -coverprofile=$(OUT_DIR)/coverage.out ./...

.PHONY: integration-test
integration-test: ## Run integration tests
	go test -race -tags=integration -coverprofile=$(OUT_DIR)/coverage.out ./...

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	go test -v -race -coverprofile=$(OUT_DIR)/coverage.out ./...

.PHONY: coverage
coverage: test ## Show test coverage in browser
	go tool cover -html=$(OUT_DIR)/coverage.out

##@ Code Quality

.PHONY: lint
lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run --timeout=5m

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	goimports -w .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: security
security: ## Run security scan
	@which gosec > /dev/null || (echo "gosec not installed. Run: go install github.com/securego/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec -fmt json -out $(OUT_DIR)/gosec-report.json ./...

##@ Building

.PHONY: build
build: ## Build binary for current platform
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -trimpath -o $(BIN_DIR)/lexicon-bo-api .

.PHONY: build-linux
build-linux: ## Build binary for Linux (amd64)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -trimpath -o $(DIST_DIR)/lexicon-bo-api-linux-amd64 .

.PHONY: build-darwin
build-darwin: ## Build binary for macOS (amd64)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -trimpath -o $(DIST_DIR)/lexicon-bo-api-darwin-amd64 .

.PHONY: build-windows
build-windows: ## Build binary for Windows (amd64)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -trimpath -o $(DIST_DIR)/lexicon-bo-api-windows-amd64.exe .

.PHONY: build-all
build-all: build-linux build-darwin build-windows ## Build binaries for all platforms

##@ Docker

.PHONY: docker-build
docker-build: ## Build production Docker image
	./build-production.sh

.PHONY: docker-run
docker-run: ## Run Docker container locally
	docker run --rm -p 8080:8080 --env-file .env lexicon-bo-api:latest

.PHONY: docker-test
docker-test: docker-build docker-run ## Build and run Docker container

##@ Cleaning

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(OUT_DIR) $(BIN_DIR) $(DIST_DIR)
	go clean -cache -testcache -modcache

.PHONY: clean-docker
clean-docker: ## Remove Docker images
	docker rmi lexicon-bo-api:latest || true

##@ CI/CD

.PHONY: ci
ci: lint test build ## Run CI checks locally

.PHONY: pre-commit
pre-commit: fmt lint test ## Run before committing

.PHONY: version
version: ## Show version information
	@echo "Version:    $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
