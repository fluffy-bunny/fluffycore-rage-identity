# Makefile for fluffycore-rage-identity

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@echo ""
	@echo "  build-server         Build the Go server"
	@echo "  all                  Build everything"
	@echo "  help                 Show this help message"
	@echo ""

# Build the Go server
.PHONY: build-server
build-server: ## Build the Go server
	@echo "Building Go server..."
	go build -o server ./cmd/server
	@echo "✅ Server build complete!"

# Build everything
.PHONY: all
all: build-server ## Build the server
	@echo "✅ Full build complete!"
