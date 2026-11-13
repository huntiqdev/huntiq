.PHONY: help build run test migrate seed clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/api ./cmd/api

run: ## Run the application
	go run ./cmd/api/main.go

test: ## Run tests
	go test ./... -v

migrate: ## Run database migrations
	./scripts/migrate.sh

seed: ## Seed database with sample data
	go run ./cmd/seed/main.go

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf uploads/products/*

.DEFAULT_GOAL := help
