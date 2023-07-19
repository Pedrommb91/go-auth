.DEFAULT_GOAL := help

export LC_ALL=en_US.UTF-8

-include .env
export $(shell sed 's/=.*//' .env)

ifndef GOOSE_DBSTRING
	export GOOSE_DBSTRING=postgresql://auth:strong-pw@localhost:5432/auth
endif

ifndef GOOSE_DRIVER
	export GOOSE_DRIVER=postgres
endif

.PHONY: help
help: ## Help command
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: run
mod-download:
	go mod download

install-dependencies: 
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4
	go install github.com/vektra/mockery/v2@v2.20.0
	go install github.com/pressly/goose/v3/cmd/goose@latest

check-formatting: ## check formatting with linter
	golangci-lint run

generate: install-dependencies mod-download
	go generate ./...

run: generate
	go run cmd/app/main.go

tests: generate
	go clean -testcache
	go test -v ./... -cover

lint: generate
	golangci-lint run --tests=0 ./...

build: generate 
	go build -o bin/app cmd/app/main.go

local-postgres: ## run local postgres container
	docker-compose up -d postgres

local-down: ## run local postgres container
	docker-compose down

migrate: ## run database migrations
	goose -dir migrations up

migrate-rollback: ## run database migrations
	goose -dir migrations down