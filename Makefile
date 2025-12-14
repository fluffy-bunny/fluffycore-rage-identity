# Makefile for fluffycore-rage-identity
# Root makefile to build all Go-App WASM applications

# Variables
MANAGEMENT_DIR := example/go-app/management
OIDC_LOGIN_DIR := example/go-app/oidc-login

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@echo ""
	@echo "  build-wasm           Build all WASM applications"
	@echo "  build-management     Build management WASM app"
	@echo "  build-oidc-login     Build oidc-login WASM app"
	@echo "  generate-static      Generate all static files"
	@echo "  clean-wasm           Clean all WASM build artifacts"
	@echo "  help                 Show this help message"
	@echo ""

# Build all WASM applications
.PHONY: build-wasm
build-wasm: build-management build-oidc-login ## Build all WASM applications

# Build management WASM app
.PHONY: build-management
build-management: ## Build management WASM app
	@echo "Building management WASM app..."
	@cd $(MANAGEMENT_DIR) && $(MAKE) build-wasm
	@echo "✅ Management WASM build complete!"

# Build oidc-login WASM app
.PHONY: build-oidc-login
build-oidc-login: ## Build oidc-login WASM app
	@echo "Building oidc-login WASM app..."
	@cd $(OIDC_LOGIN_DIR) && $(MAKE) build-wasm
	@echo "✅ OIDC-Login WASM build complete!"

# Generate all static files
.PHONY: generate-static
generate-static: ## Generate all static files for both apps
	@echo "Generating static files for all apps..."
	@cd $(MANAGEMENT_DIR) && $(MAKE) generate-static
	@cd $(OIDC_LOGIN_DIR) && $(MAKE) generate-static
	@echo "✅ All static files generated!"

# Clean WASM build artifacts
.PHONY: clean-wasm
clean-wasm: ## Clean all WASM build artifacts
	@echo "Cleaning WASM build artifacts..."
	@cd $(MANAGEMENT_DIR) && $(MAKE) clean 2>/dev/null || true
	@cd $(OIDC_LOGIN_DIR) && $(MAKE) clean 2>/dev/null || true
	@echo "✅ WASM artifacts cleaned!"

# Build the Go server
.PHONY: build-server
build-server: ## Build the Go server
	@echo "Building Go server..."
	go build -o server ./cmd/server
	@echo "✅ Server build complete!"

# Build everything
.PHONY: all
all: build-wasm build-server ## Build all WASM apps and the server
	@echo "✅ Full build complete!"
