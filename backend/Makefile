.PHONY: run dev build test test-unit test-integration lint migrate-up migrate-down migrate-test-up migrate-test-down tidy swag install-tools db-test-up db-test-down

# --- Dev DB (main) ---
DB_USER     ?= root
DB_PASSWORD ?= secret
DB_HOST     ?= localhost
DB_PORT     ?= 3306
DB_NAME     ?= expense_tracker

# --- Test DB defaults (mysql_test container on port 3307) ---
TEST_DB_USER     ?= root
TEST_DB_PASSWORD ?= secret
TEST_DB_HOST     ?= localhost
TEST_DB_PORT     ?= 3307
TEST_DB_NAME     ?= expense_tracker_test

# Run the server normally
run:
	go run ./cmd/api

# Run with Air live-reload (requires: make install-tools)
dev:
	air

# Build production binary
build:
	go build -o bin/api ./cmd/api

# Generate Swagger docs (requires: make install-tools)
# Re-run this whenever you change handler annotations.
swag:
	swag init -g cmd/api/main.go --output docs

test:
	go test ./... -v -race -count=1

test-unit:
	go test ./tests/unit/... -v -race

# Runs integration tests against the test DB container (start it first with: make db-test-up)
test-integration:
	env $(shell grep -v '^#' .env.test | xargs) go test ./tests/integration/... -v -race

# Start the test MySQL container (detached, healthcheck-ready)
db-test-up:
	docker compose -f ../docker-compose.yml up -d --wait mysql_test

# Stop and remove the test MySQL container
db-test-down:
	docker compose -f ../docker-compose.yml stop mysql_test
	docker compose -f ../docker-compose.yml rm -f mysql_test

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" up

migrate-down:
	migrate -path migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" down

migrate-test-up:
	migrate -path migrations -database "mysql://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@tcp($(TEST_DB_HOST):$(TEST_DB_PORT))/$(TEST_DB_NAME)" up

migrate-test-down:
	migrate -path migrations -database "mysql://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@tcp($(TEST_DB_HOST):$(TEST_DB_PORT))/$(TEST_DB_NAME)" down

tidy:
	go mod tidy

# Install development tools (air + swag)
install-tools:
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
