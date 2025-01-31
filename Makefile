# Variables
APP_NAME := minipack
PORT := 8585

# Targets
.PHONY: all run build test clean lint help

all: clean build

# Run the server
run:
	@echo "Starting the server on port $(PORT)..."
	go run cmd/server/main.go

# Build the server binary
build:
	@echo "Building the $(APP_NAME) binary..."
	go build -o $(APP_NAME).exe cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Clean the build artifacts
clean:
	@echo "Cleaning up..."
	del /F /Q $(APP_NAME).exe 2>nul || exit 0

# Run linter (requires golangci-lint)
lint:
	@echo "Running lint checks..."
	golangci-lint run || echo "Linting skipped: golangci-lint not found"

# Help command
help:
	@echo "Available targets:"
	@echo "  run       - Start the server"
	@echo "  build     - Build the server binary"
	@echo "  test      - Run all tests"
	@echo "  clean     - Remove build artifacts"
	@echo "  lint      - Run linting checks (requires golangci-lint)"
	@echo "  help      - Show this help message"
