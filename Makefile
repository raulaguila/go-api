COMPOSE_COMMAND = docker compose --env-file configs/.env

.PHONY: all
all: help
help: ## Display help screen
	@echo "Usage:"
	@echo "	make [COMMAND]"
	@echo "	make help\n"
	@echo "Commands: \n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: ## Create environment variables
	@chmod +x configs/env.sh && configs/env.sh && mv .env configs/

.PHONY: run
run: ## Run application
	@go run cmd/go-api/go-api.go

.PHONY: swag
swag: ## Update swagger files
	@swag init -g cmd/go-api/go-api.go --parseDependency

.PHONY: build
build: ## Build the application from source code
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o backend cmd/go-api/go-api.go

.PHONY: compose-up
compose-up: ## Run docker compose up for create and start containers
	@${COMPOSE_COMMAND} up -d

.PHONY: compose-build
compose-build: ## Run docker compose up --build for create and start containers
	@${COMPOSE_COMMAND} up -d --build

.PHONY: compose-down
compose-down: ## Run docker compose down to stop and remove containers and networks
	@${COMPOSE_COMMAND} down

.PHONY: compose-remove
compose-remove: ## Run docker compose down to stop and remove containers, networks and volumes
	@echo -n "All registered data and volumes will be deleted, are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	@${COMPOSE_COMMAND} down -v --remove-orphans

.PHONY: compose-exec
compose-exec: ## Run docker compose exec to access container bash
	@${COMPOSE_COMMAND} exec -it backend bash

.PHONY: compose-log
compose-log: ## Run docker compose logs to show container logger
	@${COMPOSE_COMMAND} logs -f backend

.PHONY: compose-top
compose-top: ## Run docker compose top to display containers processes
	@${COMPOSE_COMMAND} top

.PHONY: compose-stats
compose-stats: ## Run docker compose stats to display containers stats
	@${COMPOSE_COMMAND} stats

.PHONY: go-test
go-test: ## Run tests and show coverage
	@go test ./... -coverprofile cover.out
	@go tool cover -html=cover.out ## -o coverage.html

.PHONY: go-lint
go-lint: ## Run golang lint
	@golangci-lint run ./...
