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

seed: ## Import master data
	go run ./cmd/tools/masterdata -action=import -dry-run=false -dir=../data/master

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf uploads/products/*

.DEFAULT_GOAL := help
