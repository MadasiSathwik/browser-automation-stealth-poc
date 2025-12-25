.PHONY: build run demo clean test install deps

# Build the application
build:
	@echo "Building browser-automation..."
	go build -o browser-automation cmd/main.go
	@echo "Build complete: ./browser-automation"

# Run with default config
run: build
	@echo "Running browser-automation..."
	./browser-automation --config config.yaml

# Run in demo mode
demo: build
	@echo "Running in DEMO mode..."
	./browser-automation --demo

# Run with safe mode explicitly enabled
safe: build
	@echo "Running with SAFE mode..."
	./browser-automation --safe --config config.yaml

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts and database
clean:
	@echo "Cleaning build artifacts..."
	rm -f browser-automation
	rm -f automation.db
	rm -f automation.db-journal
	@echo "Clean complete"

# Clean and rebuild
rebuild: clean build

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Show help
help:
	@echo "Browser Automation POC - Makefile Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Build and run with default config"
	@echo "  make demo     - Run in demo mode (prints actions)"
	@echo "  make safe     - Run with safe mode enabled"
	@echo "  make deps     - Install Go dependencies"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove build artifacts and database"
	@echo "  make rebuild  - Clean and rebuild"
	@echo "  make fmt      - Format Go code"
	@echo "  make lint     - Lint code (requires golangci-lint)"
	@echo "  make help     - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make demo                    # Run demo mode"
	@echo "  DEBUG=true make run          # Run with debug logging"
	@echo "  make build && ./browser-automation --config my.yaml"
