.PHONY: test test-verbose build run clean

# Test commands (dependency injection prevents TUI from starting)
test:
	go test ./...

test-verbose:
	go test -v ./...

test-short:
	go test -short ./...

# Build the application
build:
	go build -o hyprland-monitor-tui .

# Run the application
run:
	go run .

# Clean build artifacts
clean:
	rm -f hyprland-monitor-tui

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run tests with race detection
test-race:
	go test -race ./...

# Run benchmarks
bench:
	go test -bench=. ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./... 