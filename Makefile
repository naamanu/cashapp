.PHONY: swagger build docker-build docker-up docker-down

# Generate Swagger documentation
swagger: swagger-user swagger-ledger

swagger-user:
	@echo "Generating User Service Swagger documentation..."
	@swag init -g cmd/user/main.go -o ./docs/user

swagger-ledger:
	@echo "Generating Ledger Service Swagger documentation..."
	@swag init -g cmd/ledger/main.go -o ./docs/ledger

# Build the application
build: build-user build-ledger

build-user:
	@echo "Building User Service..."
	@go build -o bin/user ./cmd/user

build-ledger:
	@echo "Building Ledger Service..."
	@go build -o bin/ledger ./cmd/ledger

# Build Docker image
docker-build:
	@echo "Building Docker images..."
	@docker-compose build

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

# Run the application locally
run-user:
	@echo "Running User Service..."
	@go run cmd/user/main.go

run-ledger:
	@echo "Running Ledger Service..."
	@go run cmd/ledger/main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Install swagger CLI (if not installed)
install-swagger:
	@echo "Installing swagger CLI..."
	@go install github.com/swaggo/swag/cmd/swag@latest
