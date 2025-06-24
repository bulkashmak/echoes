# Variables
DB_URL ?= postgres://echoes:echoes@localhost:5432/echoes
MIGRATIONS_DIR = sql/schema

# Target
.PHONY: build, migrate

# Build a binary
build:
	go build

# Run all Goose UP migrations
migrate:
	goose -dir ${MIGRATIONS_DIR} postgres "${DB_URL}" up

migrate-down:
	goose -dir ${MIGRATIONS_DIR} postgres "${DB_URL}" down
