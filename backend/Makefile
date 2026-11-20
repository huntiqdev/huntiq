.PHONY: help build run test migrate seed clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build ./cmd/api

run: ## Run the application
	go run ./cmd/api/main.go

test: ## Run tests
	go test ./... -v

migrate: ## Run database migrations
	go run ./cmd/api -migrate-only

seed: ## Verify or import master data (loccatdata)
	go run ./cmd/api -seed-only

clean: ## Clean runtime uploads
	rm -rf uploads/products/*

.DEFAULT_GOAL := help
