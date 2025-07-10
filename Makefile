# Makefile for Go Auth Template

.PHONY: run-backend run-frontend dev migrate migrate-only tidy build install-deps

# Development
dev:
	@echo "Stopping Docker containers and killing processes on ports 8080 and 3000..."
	@-docker-compose down || true
	@-pkill -f "go run.*cmd/server" || true
	@-pkill -f "next dev" || true
	@-lsof -ti:8080 | xargs kill -9 || true
	@-lsof -ti:3000 | xargs kill -9 || true
	@sleep 1
	@echo "Starting development environment..."
	./scripts/start-dev.sh

run-backend:
	@echo "Starting backend server..."
	cd authserver && export $$(cat ../.env.local | sed 's/#.*//g' | xargs) && go run ./cmd/server

run-backend-air:
	@echo "Starting backend server with hot reload..."
	cd authserver && export $$(cat ../.env.local | sed 's/#.*//g' | xargs) && air -c .air.toml || (echo 'Air not installed. Run: go install github.com/air-verse/air@latest' && go run ./cmd/server)

run-frontend:
	@echo "Starting frontend server..."
	cd front-end && npm run dev

# Database
migrate:
	@echo "Running migrations..."
	cd authserver && export $$(cat ../.env.local | sed 's/#.*//g' | xargs) && go run cmd/server/main.go --migrate

migrate-only:
	@echo "Running migrate-only..."
	cd authserver && export $$(cat ../.env.local | sed 's/#.*//g' | xargs) && go run cmd/server/main.go --migrate-only

# Go commands
tidy:
	@echo "Tidying Go modules..."
	cd authserver && go mod tidy

build:
	@echo "Building backend..."
	cd authserver && go build -o bin/authserver cmd/server/main.go

# Setup
install-deps:
	@echo "Installing dependencies..."
	cd authserver && go mod tidy
	cd front-end && npm install

install-tools:
	@echo "Installing development tools..."
	go install github.com/air-verse/air@latest
	brew install act || echo "act install failed - install manually"

# Testing
test:
	@echo "Running tests..."
	cd authserver && go test ./...
	cd front-end && npm test

test-ci:
	@echo "Running CI tests locally with act..."
	act -j test || echo "act not installed. Run 'make install-tools' or install act manually"

# Docker
docker-dev:
	@echo "Starting Docker development environment..."
	docker-compose up -d

docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

# Utilities
clean:
	@echo "Cleaning build artifacts..."
	cd authserver && rm -rf bin/
	cd front-end && rm -rf .next/
	docker-compose down -v || true

stop:
	@echo "Stopping all development processes..."
	@-pkill -f "go run.*cmd/server" || true
	@-pkill -f "next dev" || true
	@-lsof -ti:8080 | xargs kill -9 || true
	@-lsof -ti:3000 | xargs kill -9 || true
	@echo "All processes stopped."

help:
	@echo "Available commands:"
	@echo "Development:"
	@echo "  dev            - Start both backend and frontend"
	@echo "  run-backend    - Start only backend server"
	@echo "  run-backend-air - Start backend with hot reload"
	@echo "  run-frontend   - Start only frontend server"
	@echo ""
	@echo "Docker:"
	@echo "  docker-dev     - Start Docker development environment"
	@echo "  docker-build   - Build Docker images"
	@echo "  docker-stop    - Stop Docker containers"
	@echo "  docker-logs    - Show Docker logs"
	@echo ""
	@echo "Database:"
	@echo "  migrate        - Run database migrations"
	@echo "  migrate-only   - Run migrations only (no server start)"
	@echo ""
	@echo "Testing & Build:"
	@echo "  test           - Run all tests"
	@echo "  test-ci        - Run CI tests locally with act"
	@echo "  build          - Build backend binary"
	@echo "  tidy           - Tidy Go modules"
	@echo ""
	@echo "Setup:"
	@echo "  install-deps   - Install all dependencies"
	@echo "  install-tools  - Install development tools (air, act)"
	@echo ""
	@echo "Utilities:"
	@echo "  clean          - Clean build artifacts"
	@echo "  stop           - Stop all development processes"
	@echo "  help           - Show this help"