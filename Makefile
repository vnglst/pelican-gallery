# Development and build commands

.PHONY: install build run dev clean test fmt lint help

# Install dependencies and tools
install:
	@echo "🔧 Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "📦 Installing Tailwind CSS standalone binary..."
	@mkdir -p bin
	@if [ ! -f bin/tailwindcss ]; then \
		case "$(shell uname -s)-$(shell uname -m)" in \
			"Darwin-arm64") \
				curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64 -o bin/tailwindcss; \
				;; \
			"Darwin-x86_64") \
				curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-x64 -o bin/tailwindcss; \
				;; \
			"Linux-x86_64") \
				curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o bin/tailwindcss; \
				;; \
			"Linux-arm64") \
				curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64 -o bin/tailwindcss; \
				;; \
			*) \
				echo "❌ Unsupported platform. Install Tailwind CSS manually"; \
				exit 1; \
				;; \
			esac; \
		if [ ! -f bin/tailwindcss ]; then \
			echo "❌ Failed to download Tailwind CSS binary. Please check your internet connection or download manually."; \
			exit 1; \
		fi; \
		chmod +x bin/tailwindcss; \
	fi
	@echo "✅ Installation complete!"

# Build for production
build:
	@echo "🔨 Building for production..."
	@if [ ! -f bin/tailwindcss ]; then echo "Run 'make install' first"; exit 1; fi
	@./bin/tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify
	@CGO_ENABLED=0 GO_ENV=production go build -o bin/server main.go
	@echo "✅ Build complete! Binary: bin/server"

# Run the built application
run:
	@if [ ! -f bin/server ]; then echo "Run 'make build' first"; exit 1; fi
	@echo "� Starting server..."
	@./bin/server

# Development server with hot reload
dev:
	@echo "🔥 Starting development server..."
	@if [ ! -f bin/tailwindcss ]; then echo "Run 'make install' first"; exit 1; fi
	@echo "   • Go changes: automatic restart"
	@echo "   • Templates/CSS: live reload"
	@echo "   • Stop with Ctrl+C"
	@echo ""
	@# Build CSS initially
	@echo "🎨 Building CSS initially..."
	@./bin/tailwindcss -i ./static/css/input.css -o ./static/css/output.css
	@# Start both processes with proper cleanup
	@echo "📁 Starting Tailwind CSS watcher..."
	@echo "🚀 Starting Go development server..."
	@trap 'echo "Stopping processes..."; kill $$(jobs -p) 2>/dev/null' EXIT; \
	./bin/tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch=always & \
	if command -v air >/dev/null 2>&1; then \
		GO_ENV=development air; \
	else \
		echo "💡 Install Air for better hot reload: go install github.com/air-verse/air@latest"; \
		GO_ENV=development go run main.go; \
	fi

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/ tmp/
	@echo "✅ Clean complete!"

# Test the application
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Format code
fmt:
	@echo "📝 Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run

# Show help
help:
	@echo "🎨 Pelican Gallery - Make Commands"
	@echo ""
	@echo "Main Commands:"
	@echo "  install  Install dependencies and tools"
	@echo "  build    Build for production"
	@echo "  run      Run the built application"
	@echo "  dev      Development server with hot reload"
	@echo ""
	@echo "Utility Commands:"
	@echo "  clean    Clean build artifacts"
	@echo "  test     Run tests"
	@echo "  fmt      Format Go code"
	@echo "  lint     Lint Go code"
	@echo "  help     Show this help"
	@echo ""
	@echo "Typical Workflows:"
	@echo "  Development: make install && make dev"
	@echo "  Production:  make install && make build && make run"
