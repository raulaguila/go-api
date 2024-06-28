COMPOSE_COMMAND = docker compose --env-file configs/.env

.PHONY: all
all: help
help:
	@echo "Usage:"
	@echo "	make [COMMAND]"
	@echo "	make help\n"
	@echo "Commands: \n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: ## Create environment file
	@chmod +x configs/env.sh && configs/env.sh && mv .env configs/

.PHONY: compose-up
compose-up: ## Run 'docker compose up -d' to create and start containers
	@${COMPOSE_COMMAND} up -d

.PHONY: compose-build
compose-build: ## Run 'docker compose up -d --build' to create and start containers
	@${COMPOSE_COMMAND} up -d --build

.PHONY: compose-down
compose-down: ## Run 'docker compose down' to stop and remove containers and networks
	@${COMPOSE_COMMAND} down

.PHONY: compose-remove
compose-remove: ## Run 'docker compose down -v --remove-orphans' to stop and remove containers, networks and volumes
	@echo -n "All registered data and volumes will be deleted, are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	@${COMPOSE_COMMAND} down -v --remove-orphans

.PHONY: compose-exec
compose-exec: ## Run 'docker compose exec -it backend bash' to access container bash
	@${COMPOSE_COMMAND} exec -it backend bash

.PHONY: compose-log
compose-log: ## Run 'docker compose logs -f backend' to show container logger
	@${COMPOSE_COMMAND} logs -f backend

.PHONY: compose-top
compose-top: ## Run 'docker compose top' to display containers processes
	@${COMPOSE_COMMAND} top

.PHONY: compose-stats
compose-stats: ## Run 'docker compose stats' to display containers stats
	@${COMPOSE_COMMAND} stats

.PHONY: go-run
go-run: ## Run application from source code
	@go run cmd/go-api/go-api.go

.PHONY: go-build
go-build: ## Build the application from source code
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o backend cmd/go-api/go-api.go

.PHONY: go-test
go-test: ## Run tests and generate coverage report
	@go run gotest.tools/gotestsum@latest -f testname -- ./... -race -count=1 -coverprofile=/tmp/coverage.out -covermode=atomic
	@go tool cover -html=/tmp/coverage.out

.PHONY: go-benchmark
go-benchmark: ## Benchmark code performance
	@go test ./... -benchmem -bench=. -run=^Benchmark_$

.PHONY: go-lint
go-lint: ## Run lint checks
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1 run ./... --fix

.PHONY: go-audit
go-audit: ## Conduct quality checks
	@go mod verify
	@go vet ./...
	@go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./...

.PHONY: go-swag
go-swag: ## Update swagger files
	@go run github.com/swaggo/swag/cmd/swag@v1.16.3 init -g cmd/go-api/go-api.go --parseDependency

.PHONY: go-format
go-format: ## Fix code format issues
	@go run mvdan.cc/gofumpt@latest -w -l .

.PHONY: go-tidy
go-tidy: ## Clean and tidy dependencies
	@go mod tidy -v
