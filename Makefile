# ─────────────────────────────────────────────────────────────────────────────
# Universal Go Backend — Makefile
#
# All database credentials are read from the environment (set via .envrc /
# direnv). Sensible defaults match the docker-compose.yml values.
# ─────────────────────────────────────────────────────────────────────────────

# ── Config ────────────────────────────────────────────────────────────────────
APP_NAME    := my-app
BINARY_DIR  := ./bin
BINARY      := $(BINARY_DIR)/api
MAIN        := ./cmd/api/main.go
MIGRATE_DIR := internal/db/migrations

APP_PORT    ?= 8080
APP_ENV     ?= development
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_USER     ?= postgres
DB_PASSWORD ?= postgres
DB_NAME     ?= myapp
DB_MAX_CONNS ?= 25
DB_MIN_CONNS ?= 5

MIGRATE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Colours for help target
RESET  := \033[0m
BOLD   := \033[1m
GREEN  := \033[32m
YELLOW := \033[33m

.PHONY: dev build run test lint lint-fix migrate-up migrate-down migrate-create \
        sqlc swagger docker-up docker-down docker-build tidy help

# ── Development ───────────────────────────────────────────────────────────────

## dev: Start the server with live reload (air). Does NOT run migrations or swagger gen.
dev:
	@APP_PORT=$(APP_PORT) APP_ENV=$(APP_ENV) \
	 DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) \
	 DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) \
	 DB_MAX_CONNS=$(DB_MAX_CONNS) DB_MIN_CONNS=$(DB_MIN_CONNS) \
	 air

## build: Generate Swagger docs then compile the binary to ./bin/api.
build:
	@echo "$(YELLOW)→ Generating Swagger docs...$(RESET)"
	@swag init -g $(MAIN) -o docs/
	@echo "$(YELLOW)→ Building binary...$(RESET)"
	@mkdir -p $(BINARY_DIR)
	@go build -ldflags="-w -s" -o $(BINARY) $(MAIN)
	@echo "$(GREEN)✓ Binary written to $(BINARY)$(RESET)"

## run: Run the compiled binary (requires make build first).
run:
	@APP_PORT=$(APP_PORT) APP_ENV=$(APP_ENV) \
	 DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) \
	 DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) \
	 DB_MAX_CONNS=$(DB_MAX_CONNS) DB_MIN_CONNS=$(DB_MIN_CONNS) \
	 $(BINARY)

## test: Run all tests with race detector.
test:
	@go test ./... -v -race

## tidy: Tidy and verify go.mod / go.sum.
tidy:
	@go mod tidy
	@go mod verify

# ── Linting ───────────────────────────────────────────────────────────────────

## lint: Run golangci-lint.
lint:
	@golangci-lint run ./...

## lint-fix: Run golangci-lint with auto-fix.
lint-fix:
	@golangci-lint run --fix ./...

# ── Database migrations ───────────────────────────────────────────────────────
# POLICY: Migrations NEVER run automatically.
# make migrate-up is the ONLY way migrations run.

## migrate-up: Apply all pending migrations.
migrate-up:
	@echo "$(YELLOW)→ Applying migrations...$(RESET)"
	@migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_URL)" up
	@echo "$(GREEN)✓ Migrations applied$(RESET)"

## migrate-down: Roll back the last applied migration.
migrate-down:
	@echo "$(YELLOW)→ Rolling back last migration...$(RESET)"
	@migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_URL)" down 1

## migrate-create name=xxx: Create a new timestamped up+down migration pair.
migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=<migration_name>" && exit 1)
	@migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)
	@echo "$(GREEN)✓ Migration files created in $(MIGRATE_DIR)$(RESET)"

# ── Code generation ───────────────────────────────────────────────────────────

## sqlc: Regenerate type-safe DB code from SQL queries.
sqlc:
	@echo "$(YELLOW)→ Running sqlc generate...$(RESET)"
	@sqlc generate
	@echo "$(GREEN)✓ sqlc output written to internal/db/generated/$(RESET)"

## swagger: Regenerate Swagger docs from Go annotations (run after adding/changing handlers).
swagger:
	@echo "$(YELLOW)→ Generating Swagger docs...$(RESET)"
	@swag init -g $(MAIN) -o docs/
	@echo "$(GREEN)✓ Docs written to docs/$(RESET)"

# ── Docker ────────────────────────────────────────────────────────────────────

## docker-up: Start the postgres service in the background (app runs locally via make dev).
docker-up:
	@docker compose up -d postgres

## docker-down: Stop and remove all containers (data volume is preserved).
docker-down:
	@docker compose down

## docker-build: Build the Docker image for the app service.
docker-build:
	@docker build -t $(APP_NAME) .

# ── Help ──────────────────────────────────────────────────────────────────────

## help: Print this help message.
help:
	@echo ""
	@echo "$(BOLD)$(APP_NAME) — available targets$(RESET)"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | awk -F': ' \
		'{ printf "  $(GREEN)%-22s$(RESET) %s\n", $$1, $$2 }'
	@echo ""
