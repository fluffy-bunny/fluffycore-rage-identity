# Makefile for Go WebAssembly Application
# mastodon-identity management go-app

# Variables
APP_NAME := mastodon-identity-management-go-app
WASM_BINARY := app.wasm
STATIC_OUTPUT_DIR := ../../../cmd/server/static/go-app/management/static_output
WASM_OUTPUT_PATH := $(STATIC_OUTPUT_DIR)/web/$(WASM_BINARY)
STATIC_GENERATOR_BINARY := static_generator
SERVER_BINARY := server

# Static generator options
APP_TITLE := MyApp
BASE_HREF := management

# Build versioning with ldflags
BUILD_VERSION := $(shell date +%s 2>/dev/null || powershell -Command "Get-Date -UFormat %s")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S' 2>/dev/null || powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH:mm:ss'")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/common.AppVersion=$(BUILD_VERSION) \
           -X github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/common.BuildTime=$(BUILD_TIME) \
           -X github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/common.GitCommit=$(GIT_COMMIT) \
           -X github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/common.GitBranch=$(GIT_BRANCH)

# Go build variables
GOOS := js
GOARCH := wasm

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Customization:"
	@echo "  Override variables: make generate-static APP_TITLE=\"My App\" BASE_HREF=\"/app/\""
	@echo "  Example: make generate-static APP_TITLE=\"Dashboard\" BASE_HREF=\"/admin/\" STATIC_OUTPUT_DIR=\"./dist\""

# Build targets
.PHONY: all
all: build-wasm build-static-generator build-server ## Build all components (wasm, static_generator, server)

.PHONY: build-wasm
build-wasm: ## Build the WebAssembly client
	@echo "Building WebAssembly client..."
	@mkdir -p $(STATIC_OUTPUT_DIR)
	GOOS=js GOARCH=wasm go build -ldflags "$(LDFLAGS)" -o $(WASM_OUTPUT_PATH) ./cmd/wasm
	@echo "✅ WASM build successful! Output: $(WASM_OUTPUT_PATH)"

.PHONY: build-static-generator
build-static-generator: ## Build the static_generator binary
	@echo "Building static_generator..."
	go build -ldflags "$(LDFLAGS)" -o $(STATIC_GENERATOR_BINARY) ./cmd/static_generator
	@echo "✅ static_generator build successful!"

.PHONY: build-server
build-server: ## Build the echo server binary
	@echo "Building server..."
	go build -ldflags "$(LDFLAGS)" -o $(SERVER_BINARY) ./cmd/server
	@echo "✅ server build successful!"

.PHONY: build-fixup-static-html
build-fixup-static-html: ## Build the fixup_static_html binary
	@echo "Building fixup_static_html..."
	go build -ldflags "$(LDFLAGS)" -o fixup_static_html ./cmd/fixup_static_html
	@echo "✅ fixup_static_html build successful!"

.PHONY: build-gen-images
build-gen-images: ## Build the gen_images binary
	@echo "Building gen_images..."
	go build -o gen_images ./cmd/gen_images
	@echo "✅ gen_images build successful!"

.PHONY: gen-images
gen-images: build-gen-images ## Generate PNG favicons from SVG
	@echo "Generating PNG images from SVG..."
	@cd cmd/gen_images && ../../gen_images || powershell -Command "cd cmd/gen_images; ../../gen_images"
	@echo "✅ Images generated in ./web/"

# Generate static output
.PHONY: generate-static
generate-static: gen-images ## Build and generate static files
	@echo "Cleaning static output directory..."
	@rm -rf $(STATIC_OUTPUT_DIR) 2>/dev/null || powershell -Command "Remove-Item -Path $(STATIC_OUTPUT_DIR) -Recurse -Force -ErrorAction SilentlyContinue"
	@mkdir -p $(STATIC_OUTPUT_DIR) || powershell -Command "New-Item -ItemType Directory -Path $(STATIC_OUTPUT_DIR) -Force | Out-Null"
	@echo "Syncing shared webauthn.js from example/go-app/web/..."
	@cp ../web/webauthn.js ./web/webauthn.js 2>/dev/null || powershell -Command "Copy-Item -Path ../web/webauthn.js -Destination ./web/webauthn.js -Force"
	@$(MAKE) build-wasm build-static-generator build-fixup-static-html
	@echo "Generating build_version.js in source web directory..."
	@echo "// Build version information" > ./web/build_version.js
	@echo "// This file is auto-generated during the build process" >> ./web/build_version.js
	@echo "window.BUILD_VERSION = window.BUILD_VERSION || {" >> ./web/build_version.js
	@echo "  version: '$(BUILD_VERSION)'," >> ./web/build_version.js
	@echo "  buildTime: '$(BUILD_TIME)'," >> ./web/build_version.js
	@echo "  gitCommit: '$(GIT_COMMIT)'," >> ./web/build_version.js
	@echo "  gitBranch: '$(GIT_BRANCH)'" >> ./web/build_version.js
	@echo "};" >> ./web/build_version.js
	@echo "Generating static output..."
	MSYS_NO_PATHCONV=1 ./$(STATIC_GENERATOR_BINARY) -title "$(APP_TITLE)" -output $(STATIC_OUTPUT_DIR) -base-href "/$(BASE_HREF)"
	@echo "Copying web assets..."
	@mkdir -p $(STATIC_OUTPUT_DIR)/web
	@cp -r ./web/* $(STATIC_OUTPUT_DIR)/web/ 2>/dev/null || powershell -Command "Copy-Item -Path ./web/* -Destination $(STATIC_OUTPUT_DIR)/web/ -Recurse -Force"
	@echo "Fixing up HTML files..."
	$(MSYS_NO_PATHCONV) ./fixup_static_html -dir $(STATIC_OUTPUT_DIR) -basehref "$(BASE_HREF)" -title "$(APP_TITLE)"
	@echo "Cleaning up unnecessary HTML files..."
	@cd $(STATIC_OUTPUT_DIR) && (ls *.html 2>/dev/null | grep -v "^index.html$$" | grep -v "_template.html$$" | xargs rm -f 2>/dev/null || true)
	@cd $(STATIC_OUTPUT_DIR) && powershell -Command "Get-ChildItem -Filter *.html | Where-Object { $$_.Name -ne 'index.html' -and $$_.Name -notlike '*_template.html' } | Remove-Item -Force -ErrorAction SilentlyContinue" 2>/dev/null || true
	@echo "✅ Static files generated in $(STATIC_OUTPUT_DIR)"

.PHONY: fixup-static-html
fixup-static-html: build-fixup-static-html ## Fix up HTML files in static output
	@echo "Fixing up HTML files in $(STATIC_OUTPUT_DIR)..."
	./fixup_static_html -dir $(STATIC_OUTPUT_DIR) -basehref "$(BASE_HREF)" -title "$(APP_TITLE)"

# Run targets
.PHONY: run-generator
run-generator: build-static-generator ## Run the static_generator with default options
	@echo "Running static_generator..."
	./$(STATIC_GENERATOR_BINARY) -title "$(APP_TITLE)" -output $(STATIC_OUTPUT_DIR) -base-href "$(BASE_HREF)"

.PHONY: run-server
run-server: ## Run the echo server
	@echo "Starting echo server to host $(STATIC_OUTPUT_DIR)..."
	@echo "Open your browser to http://localhost:8080"
	go run ./cmd/server

.PHONY: run
run: generate-static run-server ## Generate static files and run the server

# Clean targets
.PHONY: clean
clean: ## Clean all build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f $(STATIC_GENERATOR_BINARY) $(STATIC_GENERATOR_BINARY).exe 2>/dev/null || powershell -Command "Remove-Item -Path $(STATIC_GENERATOR_BINARY),$(STATIC_GENERATOR_BINARY).exe -ErrorAction SilentlyContinue"
	@rm -f $(SERVER_BINARY) $(SERVER_BINARY).exe 2>/dev/null || powershell -Command "Remove-Item -Path $(SERVER_BINARY),$(SERVER_BINARY).exe -ErrorAction SilentlyContinue"
	@rm -f fixup_static_html fixup_static_html.exe 2>/dev/null || powershell -Command "Remove-Item -Path fixup_static_html,fixup_static_html.exe -ErrorAction SilentlyContinue"
	@rm -rf $(STATIC_OUTPUT_DIR) 2>/dev/null || powershell -Command "Remove-Item -Path $(STATIC_OUTPUT_DIR) -Recurse -Force -ErrorAction SilentlyContinue"
	@echo "✅ Cleaned all build artifacts"

# Development targets
.PHONY: rebuild
rebuild: clean all ## Clean and rebuild all components

# Go targets
.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies updated"

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: fmt vet ## Run formatting and vetting

# Utility targets
.PHONY: check
check: deps lint test ## Run all checks (deps, lint, test)

.PHONY: info
info: ## Show project information
	@echo "Project: $(APP_NAME)"
	@echo "WASM Output: $(WASM_OUTPUT_PATH)"
	@echo "Static Generator: $(STATIC_GENERATOR_BINARY)"
	@echo "Server Binary: $(SERVER_BINARY)"
	@echo "Go Version: $$(go version)"
	@echo "WASM Build Target: $(GOOS)/$(GOARCH)"

# Development workflow
.PHONY: dev
dev: clean all ## Full development workflow (clean, build all components)
