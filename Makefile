# Variables
BINARY_NAME=api
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
GOLANGCI_LINT_VERSION=latest

.PHONY: all build clean test coverage lint run docker-build docker-run help

# Default target
all: lint test build

## Build: Compile the binary
build:
	@echo "Building..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/api

## Clean: Remove binary and coverage files
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME) coverage.txt

## Test: Run unit tests
test:
	@echo "Running tests..."
	go test -v -race ./...

## Coverage: Run tests with coverage report
coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.txt ./...
	go tool cover -func=coverage.txt

## Lint: Run golangci-lint
lint:
	@echo "Linting..."
	golangci-lint run

## Run: Run the application locally
run:
	@echo "Running application..."
	go run ./cmd/api

## Docker Build: Build the Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t wherego:latest .

## Docker Run: Run the Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 wherego:latest

## Help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-15s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
