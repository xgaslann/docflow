.PHONY: help dev build test clean lint fmt package deploy

# Variables
BINARY_NAME=docflow-server
BUILD_DIR=build
BIN_DIR=bin

# Default target
help:
	@echo "DocFlow - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev          Start both backend and frontend"
	@echo "  make dev-backend  Start backend only"
	@echo "  make dev-frontend Start frontend only"
	@echo ""
	@echo "Build:"
	@echo "  make build        Build everything"
	@echo "  make build-backend  Build backend binary"
	@echo "  make build-frontend Build frontend"
	@echo "  make package      Create release package"
	@echo ""
	@echo "Testing:"
	@echo "  make test         Run all tests"
	@echo "  make test-backend Run backend tests"
	@echo "  make test-frontend Run frontend tests"
	@echo "  make coverage     Run tests with coverage"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint         Run linters"
	@echo "  make fmt          Format code"
	@echo ""
	@echo "Deployment:"
	@echo "  make docker       Build Docker images"
	@echo "  make docker-up    Start Docker containers"
	@echo "  make docker-down  Stop Docker containers"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        Remove build artifacts"
	@echo "  make deps         Install dependencies"

# Development
dev:
	@echo "Starting development servers..."
	@./scripts/build.sh help 2>/dev/null || true
	@(cd backend && go run cmd/server/main.go &)
	@(cd frontend && npm run dev)

dev-backend:
	cd backend && go run cmd/server/main.go

dev-frontend:
	cd frontend && npm run dev

# Build
build: build-backend build-frontend
	@echo "Build complete!"

build-backend:
	@echo "Building backend..."
	@mkdir -p $(BIN_DIR)
	cd backend && go build -ldflags="-s -w" -o ../$(BIN_DIR)/$(BINARY_NAME) cmd/server/main.go
	@echo "Backend built: $(BIN_DIR)/$(BINARY_NAME)"

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm ci && npm run build
	@echo "Frontend built: frontend/dist/"

package:
	@echo "Creating release package..."
	./scripts/build.sh package

# Testing
test: test-backend test-frontend
	@echo "All tests passed!"

test-backend:
	@echo "Running backend tests..."
	cd backend && go test -v ./...

test-frontend:
	@echo "Running frontend tests..."
	cd frontend && npm run test:run

coverage:
	@echo "Running tests with coverage..."
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	cd frontend && npm run test:coverage

# Code Quality
lint: lint-backend lint-frontend

lint-backend:
	@echo "Linting backend..."
	cd backend && go vet ./...

lint-frontend:
	@echo "Linting frontend..."
	cd frontend && npm run lint

fmt:
	@echo "Formatting code..."
	cd backend && go fmt ./...

# Docker
docker:
	docker-compose build

docker-up:
	docker-compose up -d
	@echo "DocFlow running at http://localhost"

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Dependencies
deps:
	@echo "Installing dependencies..."
	cd backend && go mod tidy
	cd frontend && npm install

# Clean
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -rf $(BIN_DIR)
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf backend/temp
	rm -rf backend/output
	rm -f backend/coverage.out
	rm -f backend/coverage.html
	@echo "Clean complete!"
