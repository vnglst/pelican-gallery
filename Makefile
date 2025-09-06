# Development and build commands

.PHONY: install build run dev clean

# Install dependencies
install:
	go mod tidy
	go mod download

# Build the application
build:
	GO_ENV=production go build -o bin/server main.go

# Run the built application
run: build
	./bin/server

# Run development server with hot reload
dev:
	@echo "Starting development server with hot reload..."
	@echo "Go code changes will trigger automatic restarts"
	@echo "Templates and static files will be read from disk for live updates"
	@if command -v air >/dev/null 2>&1; then \
		GO_ENV=development air; \
	elif [ -f "/Users/vnglst/go/bin/air" ]; then \
		GO_ENV=development /Users/vnglst/go/bin/air; \
	else \
		echo "Air not found. Install it with: go install github.com/air-verse/air@latest"; \
		echo "Falling back to basic go run..."; \
		GO_ENV=development go run main.go; \
	fi

# Run development server without hot reload (fallback)
dev-simple:
	@echo "Starting development server (simple mode)..."
	@echo "Templates and static files will be read from disk for live updates"
	@echo "Restart manually for Go code changes"
	GO_ENV=development go run main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Test the application
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  install    - Install dependencies"
	@echo "  build      - Build the application for production"
	@echo "  run        - Build and run the application"
	@echo "  dev        - Run development server with hot reload (requires Air)"
	@echo "  dev-simple - Run development server without hot reload"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  help       - Show this help message"
