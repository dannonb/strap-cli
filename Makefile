# Microservice Bootstrapper Build Configuration

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Binary name
BINARY_NAME = strap

# Go build flags
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE) -s -w"

# Build directories
BUILD_DIR = build
DIST_DIR = dist

# Supported platforms
PLATFORMS = \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

.PHONY: all build clean test lint install build-all release help

# Default target
all: build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

# Build for all platforms
build-all: clean
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_NAME=$(BINARY_NAME); \
		if [ $$OS = "windows" ]; then OUTPUT_NAME=$(BINARY_NAME).exe; fi; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$$OS-$$ARCH/$$OUTPUT_NAME ./cmd; \
		if [ $$? -ne 0 ]; then \
			echo "Failed to build for $$OS/$$ARCH"; \
			exit 1; \
		fi; \
	done

# Create release packages
release: build-all
	@echo "Creating release packages..."
	@cd $(DIST_DIR) && for dir in */; do \
		platform=$$(basename "$$dir"); \
		echo "Packaging $$platform..."; \
		if echo "$$platform" | grep -q "windows"; then \
			zip -r "$(BINARY_NAME)-$$platform.zip" "$$dir"; \
		else \
			tar -czf "$(BINARY_NAME)-$$platform.tar.gz" "$$dir"; \
		fi; \
	done
	@echo "Release packages created in $(DIST_DIR)/"

# Install locally
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, running go vet instead..."; \
		go vet ./...; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)

# Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build for current platform"
	@echo "  build-all  - Build for all supported platforms"
	@echo "  release    - Create release packages for all platforms"
	@echo "  install    - Install binary locally (requires sudo)"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  clean      - Clean build artifacts"
	@echo "  version    - Show version information"
	@echo "  help       - Show this help message"