.PHONY: help run build test docker-build docker-run docker-compose-up docker-compose-down clean

# Default target
help:
	@echo "PingMe API - Available Commands:"
	@echo "  make run              - Run the application locally"
	@echo "  make build            - Build the application binary"
	@echo "  make test             - Run the test suite"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run Docker container"
	@echo "  make docker-compose-up   - Start with docker-compose"
	@echo "  make docker-compose-down - Stop docker-compose"
	@echo "  make clean            - Remove build artifacts"

# Run the application
run:
	@echo "Starting PingMe API..."
	go run main.go

# Build the application
build:
	@echo "Building PingMe API..."
	go build -o pingme-api main.go
	@echo "Binary created: ./pingme-api"

# Run tests
test:
	@echo "Running test suite..."
	@./tests/api-tests.sh

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t pingme-api:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name pingme-api-container pingme-api:latest

# Start with docker-compose
docker-compose-up:
	@echo "Starting with docker-compose..."
	docker-compose up -d
	@echo "API is running at http://localhost:8080"

# Stop docker-compose
docker-compose-down:
	@echo "Stopping docker-compose..."
	docker-compose down

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f pingme-api
	@docker-compose down 2>/dev/null || true
	@docker rm -f pingme-api-container 2>/dev/null || true
	@echo "Clean complete!"