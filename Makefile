# Simple Makefile for a Go project

all: compose

generate-sql:
	sqlc generate

lint:
	@echo "*formating*"
	@docker run --rm -v .:/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint fmt
	@echo "*running lint*"
	@docker run --rm -v .:/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run

compose:
	@docker compose --profile="*" down --remove-orphans && \
	docker compose --profile="*" up --build --watch

.PHONY: all compose

