# Makefile for the Go HTTP Rental Server

# Define the name of the output binary
BINARY_NAME=rental-server

# Define the default database URL for development.
# This can be overridden from the command line, e.g., make run DB_URL=...
DATABASE_URL="postgres://devuser:devpassword@localhost:5432/rentaldb?sslmode=disable"

# Phony targets are targets that are not files.
.PHONY: all run test clean db-up db-down db-logs migrate-up tidy build

all: build run

# ====================================================================================
# Development Commands
# ====================================================================================

## run: Builds and runs the application.
run: build
	@echo "Running the application..."
	@DATABASE_URL=${DATABASE_URL} ./$(BINARY_NAME)

## test: Runs all tests in the project.
test:
	@echo "Running tests..."
	@go test -v ./...

## tidy: Tidies up go.mod and go.sum files.
tidy:
	@echo "Tidying go modules..."
	@go mod tidy

# ====================================================================================
# Database Commands (Docker Compose)
# ====================================================================================

## db-up: Starts the PostgreSQL database container in the background.
db-up:
	@echo "Starting PostgreSQL container..."
	@docker compose up -d

## db-down: Stops and removes the PostgreSQL database container.
db-down:
	@echo "Stopping PostgreSQL container..."
	@docker compose down

## db-logs: Tails the logs from the PostgreSQL container.
db-logs:
	@echo "Tailing PostgreSQL logs..."
	@docker compose logs -f

# ====================================================================================
# Migration Commands
# ====================================================================================

## migrate-up: Applies all 'up' database migrations.
migrate-up:
	@echo "Applying database migrations..."
	@migrate -database "${DATABASE_URL}" -path migrations up

# ====================================================================================
# Build and Cleanup Commands
# ====================================================================================

## build: Compiles the application into a binary.
build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) ./cmd/server/main.go

## clean: Removes the built binary.
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)