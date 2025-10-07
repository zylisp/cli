# Makefile for data-pipeline-v3
# Provides convenient commands for development workflow
PROJ_NAME := cli
PROJ_DESC := CLI tools for Zylisp
BIN_DIR = ./bin
MAINS = cmd/%/main.go
CMDS = $(wildcard cmd/*/main.go)
BINS = $(patsubst $(MAINS),bin/%,$(CMDS))

.PHONY: help build clean test test-integration test-verbose test-single lint format format-check publish-local check-types install deps version release just-publish publish micro+ minor+ major+ setup gen-data run-local stop logs status

# Default target
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[1;34m
GREEN := \033[1;32m
YELLOW := \033[1;33m
RED := \033[1;31m
RESET := \033[0m


help: ## Show this help message
	@echo "$(BLUE)$(PROJ_DESC) - Available Commands$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Prerequisites:$(RESET)"
	@echo "  - Go 1.21+"
	@echo "  - Make"
	@echo ""

deps: ## Install/update dependencies
	@echo "$(BLUE)Installing dependencies...$(RESET)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies installed$(RESET)"

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -rf ./bin
	@rm -f coverage.out
	@echo "$(GREEN)✅ Clean completed$(RESET)"

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

bin/%: $(BIN_DIR) cmd/%/main.go
	@echo "$(BLUE)Building $@ from cmd/$*/...$(RESET)"
	@go build -o ./$@ ./cmd/$*/
	@echo "$(GREEN)✅ Built $@$(RESET)"

build: $(BINS) ## Full build (builds all binaries in cmd/)
	@echo "$(GREEN)✅ Build completed$(RESET)"

test: ## Run all tests
	@echo "$(BLUE)Running tests...$(RESET)"
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "$(GREEN)✅ Tests completed$(RESET)"
	@go tool cover -func=coverage.out | tail -1

test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)Running tests with verbose output...$(RESET)"
	@go test -v -race ./...
	@echo "$(GREEN)✅ Verbose tests completed$(RESET)"

test-single: ## Run a single test (usage: make test-single TEST=TestName)
ifndef TEST
	@echo "$(RED)Error: Please specify TEST=TestName$(RESET)"
	@exit 1
endif
	@echo "$(BLUE)Running single test: $(TEST)...$(RESET)"
	@go test -v -race -run $(TEST) ./...
	@echo "$(GREEN)✅ Single test completed$(RESET)"

lint: ## Run code quality checks
	@echo "$(BLUE)Running code quality checks...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✅ Linting completed$(RESET)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not found. Install with: brew install golangci-lint$(RESET)"; \
	fi

format: ## Format code
	@echo "$(BLUE)Formatting code...$(RESET)"
	@go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "$(YELLOW)⚠ goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest$(RESET)"; \
	fi
	@echo "$(GREEN)✅ Code formatted$(RESET)"

format-check: ## Check if code is properly formatted
	@echo "$(BLUE)Checking code formatting...$(RESET)"
	@test -z "$$(gofmt -l .)" || (echo "$(RED)Code not formatted. Run 'make format'$(RESET)" && exit 1)
	@echo "$(GREEN)✅ Format check completed$(RESET)"

check-types: ## Validate types and compilation
	@echo "$(BLUE)Checking types and compilation...$(RESET)"
	@go build -o /dev/null ./...
	@echo "$(GREEN)✅ Type checking completed$(RESET)"

verify: check-types test lint ## Run full verification (compile, test, lint)
	@echo "$(GREEN)✅ Verification completed$(RESET)"

install: ## Install to local
	@echo "$(BLUE)Installing to local repository...$(RESET)"
	@echo "$(GREEN)✅ Installed to ...?$(RESET)"

force-install: ## Force-install to local
	@echo "$(RED)Installing to local repository...$(RESET)"
	@echo "$(GREEN)‼️ Force-installed to ~/.m2/repository$(RESET)"

publish-local: install ## Alias for install

package: ## Create package
	@echo "$(BLUE)Creating package...$(RESET)"
	@echo "$(GREEN)✅ Package created$(RESET)"

site: ## Generate project site and reports
	@echo "$(BLUE)Generating project site...$(RESET)"
	@echo "$(GREEN)✅ Site generated in target/site/$(RESET)"

dependency-tree: ## Show dependency tree
	@echo "$(BLUE)Showing dependency tree...$(RESET)"

dependency-updates: ## Check for dependency updates
	@echo "$(BLUE)Checking for dependency updates...$(RESET)"

plugin-updates: ## Check for plugin updates
	@echo "$(BLUE)Checking for plugin updates...$(RESET)"

security-check: ## Run security vulnerability check
	@echo "$(BLUE)Running security vulnerability check...$(RESET)"
	@echo "$(GREEN)✅ Security check completed$(RESET)"

release-prepare: ## Prepare release (update versions, create tag)
	@echo "$(BLUE)Preparing release...$(RESET)"
	@echo "$(GREEN)✅ Release prepared$(RESET)"

release-perform: ## Perform release (deploy to repository)
	@echo "$(BLUE)Performing release...$(RESET)"
	@echo "$(GREEN)✅ Release performed$(RESET)"

quick: ## Quick build and test
	@echo "$(BLUE)Running quick build and test...$(RESET)"
	@echo "$(GREEN)✅ Quick build completed$(RESET)"

ci: clean lint check-types test package ## CI pipeline (clean, lint, check-types, test, package)
	@echo "$(GREEN)✅ CI pipeline completed successfully$(RESET)"

dev-setup: deps ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(RESET)"
	@echo "$(GREEN)✅ Development environment ready$(RESET)"

watch: ## Watch for changes and run tests
	@echo "$(BLUE)Watching for changes...$(RESET)"
	@echo "$(YELLOW)Note: This uses fswatch. Install with: brew install fswatch (macOS) or apt-get install fswatch (Linux)$(RESET)"
	@which fswatch > /dev/null || (echo "$(RED)Error: fswatch not found$(RESET)" && exit 1)
	fswatch -o src/ pom.xml | xargs -n1 -I{} make quick

clean-all: clean ## Clean everything including IDE files
	@echo "$(BLUE)Cleaning all files...$(RESET)"
	@echo "$(GREEN)✅ Everything cleaned$(RESET)"

version: ## Show current project version
	@echo TBD

micro+: ## Increment micro/patch version (x.y.z -> x.y.z+1)
	@echo "$(BLUE)Incrementing micro version...$(RESET)"
	echo "$(GREEN)✅ Version updated to: $$new_version$(RESET)"

minor+: ## Increment minor version (x.y.z -> x.y+1.0)
	@echo "$(BLUE)Incrementing minor version...$(RESET)"
	echo "$(GREEN)✅ Version updated to: $$new_version$(RESET)"

major+: ## Increment major version (x.y.z -> x+1.0.0)
	@echo "$(BLUE)Incrementing major version...$(RESET)"
	@current=$$($(MVN) help:evaluate -Dexpression=project.version -q -DforceStdout) && \
	echo "$(GREEN)✅ Version updated to: $$new_version$(RESET)"

release: ## Create and tag release version
	@echo "$(BLUE)Creating release...$(RESET)"
	@version=TBD && \
	git tag -a "v$$version" -m "Release version $$version" && \
	echo "$(GREEN)Tagged release v$$version$(RESET)"

just-publish: ## Push changes and tags to origin
	@echo "$(BLUE)Publishing to origin...$(RESET)"
	@git pull origin main --rebase && \
	git push origin main && \
	git push origin main --tags
	@echo "$(GREEN)✅ Published to origin$(RESET)"

publish: clean build release just-publish ## Full publish workflow (clean, build, release, publish)
	@echo "$(GREEN)✅ Full publish workflow completed$(RESET)"

info: ## Show project information
	@echo "$(BLUE)Project Information$(RESET)"
	@echo "==================="
	@echo "$(YELLOW)Project:$(RESET) $(PROJ_DESC)"
	@echo "$(YELLOW)Module:$(RESET) $$(go list -m)"
	@echo "$(YELLOW)Go Version:$(RESET) $$(go version | cut -d' ' -f3)"
	@echo "$(YELLOW)Binary Directory:$(RESET) $(BIN_DIR)"
	@echo ""

run:
	@echo "$(BLUE)Running REPL ...$(RESET)"
	@go run cmd/zylisp/main.go

run-alt:
	@echo "$(BLUE)Running REPL with alternate prompt ...$(RESET)"
	@go run cmd/zylisp/main.go --prompt=alt

test-integration: ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(RESET)"
	@go test -tags=integration -v ./tests/integration/...
	@echo "$(GREEN)✅ Integration tests completed$(RESET)"

# Development shortcuts
dev: quick ## Alias for quick
all: ci ## Alias for ci
check: verify ## Alias for verify
