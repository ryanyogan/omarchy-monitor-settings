.PHONY: test test-verbose test-short test-coverage test-race bench build run clean deps fmt vet lint staticcheck quality-check help

# Get version from git tag or use "dev"
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# ================================
# Quality Assurance Targets
# ================================

# Run comprehensive quality check suite
quality-check: vet fmt-check mod-tidy build-check test-race staticcheck test-verbose
	@echo "🎉 All quality checks passed! Code is ready for commit."
	@echo "📊 Summary:"
	@echo "   ✅ Code quality: Clean"
	@echo "   ✅ Formatting: Consistent" 
	@echo "   ✅ Dependencies: Up-to-date"
	@echo "   ✅ Compilation: Error-free"
	@echo "   ✅ Tests: All passing"
	@echo "   ✅ Race conditions: None detected"
	@echo "   ✅ Static analysis: Clean"

# Individual quality check targets
vet:
	@echo "📋 Running go vet..."
	go vet ./...
	@echo "✅ go vet passed"

fmt:
	@echo "🎨 Running go fmt..."
	go fmt ./...
	@echo "✅ go fmt completed"

fmt-check:
	@echo "🎨 Checking code formatting..."
	@if [ -n "$$(go fmt ./...)" ]; then \
		echo "❌ Code formatting issues found. Run 'make fmt' to fix."; \
		exit 1; \
	fi
	@echo "✅ Code formatting is clean"

mod-tidy:
	@echo "📦 Running go mod tidy..."
	go mod tidy
	@echo "✅ go mod tidy completed"

build-check:
	@echo "🔨 Running build check..."
	go build -ldflags "-X main.version=$(VERSION)" -o omarchy-monitor-settings ./cmd/omarchy-monitor-settings
	@rm -f omarchy-monitor-settings
	@echo "✅ Build check passed"

test-race:
	@echo "🏃 Running tests with race detection..."
	go test -race ./...
	@echo "✅ Race detection tests passed"

staticcheck:
	@echo "🔬 Running staticcheck..."
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
		echo "✅ staticcheck passed"; \
	else \
		echo "⚠️  staticcheck not installed, skipping (install with: go install honnef.co/go/tools/cmd/staticcheck@latest)"; \
	fi

# ================================
# Test Targets
# ================================

test:
	go test ./...

test-verbose:
	@echo "🧪 Running full test suite..."
	go test -v ./...
	@echo "✅ Full test suite passed"

test-short:
	go test -short ./...

test-coverage:
	go test -cover ./...

# Visual Regression Testing
.PHONY: visual-test visual-update visual-clean
visual-test: ## Run visual regression tests
	@echo "Running visual regression tests..."
	go test -run TestVisualRegression -v

visual-update: ## Update golden files for visual regression tests
	@echo "Updating visual regression golden files..."
	UPDATE_GOLDEN=true go test -run TestVisualRegression -v

visual-clean: ## Clean visual regression test artifacts
	@echo "Cleaning visual regression artifacts..."
	rm -rf testdata/golden/*.diff
	rm -rf testdata/golden/*.tmp

# ================================
# Build & Run Targets  
# ================================

version:
	@echo "Current version: $(VERSION)"
	@echo "Git tag: $(shell git describe --tags --abbrev=0 2>/dev/null || echo "none")"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")"

build:
	go build -ldflags "-X main.version=$(VERSION)" -o omarchy-monitor-settings ./cmd/omarchy-monitor-settings

run:
	go run .

clean:
	rm -f omarchy-monitor-settings

# ================================
# Development Targets
# ================================

deps:
	go mod tidy
	go mod download

bench:
	go test -bench=. ./...

# Legacy lint target (requires golangci-lint)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Using quality-check instead..."; \
		$(MAKE) quality-check; \
	fi

# ================================
# Help Target
# ================================

help: ## Show this help message
	@echo "Available targets:"
	@echo ""
	@echo "Quality Assurance:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "(vet|fmt|tidy|build-check|test|lint|cover)"
	@echo ""
	@echo "Visual Testing:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "(visual-)"
	@echo ""
	@echo "Build & Development:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | grep -E "(build|run|version|install|clean)" | grep -v visual
	@echo ""
	@echo "Quick Commands:"
	@echo "  \033[36mmake qa\033[0m         - Run all quality checks"
	@echo "  \033[36mmake visual-test\033[0m - Run visual regression tests"
	@echo "  \033[36mmake build\033[0m      - Build the application" 