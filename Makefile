.PHONY: help dev build clean docker-up docker-down migrate install-tools

help: ## Show this help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## Install development tools
	go install github.com/air-verse/air@latest
	go install entgo.io/ent/cmd/ent@latest

dev: ## Start development server with hot reload
	air

build: ## Build the application
	go build -o bin/server ./cmd/server

clean: ## Clean build artifacts
	rm -rf bin/ tmp/

docker-up: ## Start PostgreSQL with Docker Compose
	docker-compose up -d

docker-down: ## Stop Docker Compose services
	docker-compose down

migrate: ## Run database migrations
	go run ./cmd/server

test: ## Run tests
	go test ./...

generate: ## Generate Ent code
	go generate ./ent

deps: ## Download dependencies
	go mod tidy
	go mod download

.DEFAULT_GOAL := help