.DEFAULT_GOAL := help
COMPOSE = docker compose -f ../docker-compose.yml $(if $(wildcard .env),--env-file .env)

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

setup: ## First-time setup: .env, deps, dev tools
	@test -f .env || cp .env.example .env
	go mod tidy
	go install github.com/air-verse/air@latest

db-up: ## Start dev + test Postgres (Docker) and wait until healthy
	$(COMPOSE) up -d --wait

db-down: ## Stop databases (dev data kept)
	$(COMPOSE) down

db-reset: ## Stop databases and DELETE dev data
	$(COMPOSE) down -v

db-psql: ## Open psql on the dev database
	$(COMPOSE) exec postgres psql -U $${DB_USER:-app} -d $${DB_NAME:-app}

run: ## Run the API (migrations apply on startup)
	go run ./cmd/server

dev: ## Run with hot reload (air)
	air

build: ## Build binary to bin/server
	go build -o bin/server ./cmd/server

test: ## Run all tests (DB tests skip if test DB is down)
	go test -count=1 ./...

test-v: ## Verbose tests with race detector
	go test -count=1 -race -v ./...

cover: ## Test coverage report
	go test -count=1 -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | tail -1

lint: ## go vet + gofmt check
	go vet ./...
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed:"; gofmt -l .; exit 1)

fmt: ## Format code
	gofmt -w .

new-migration: ## Create migration files: make new-migration name=create_items
	@test -n "$(name)" || (echo "usage: make new-migration name=create_items"; exit 1)
	@n=$$(printf "%04d" $$(( $$(ls migrations/*.up.sql 2>/dev/null | wc -l) + 1 ))); \
	touch migrations/$${n}_$(name).up.sql migrations/$${n}_$(name).down.sql; \
	echo "created migrations/$${n}_$(name).{up,down}.sql"

.PHONY: help setup db-up db-down db-reset db-psql run dev build test test-v cover lint fmt new-migration
