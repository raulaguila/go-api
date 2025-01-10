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

.PHONY: compose-build-services
compose-build-services: ## Run 'docker compose --env-file configs/.env --profile services up -d --build' to create and start services containers
	@#BUILDKIT_PROGRESS=plain ${COMPOSE_COMMAND} --profile services up -d --build
	@${COMPOSE_COMMAND} --profile services up -d --build

.PHONY: compose-build-source
compose-build-source: ## Run 'docker compose --env-file configs/.env --profile services --profile source up -d --build' to create and start containers from source code
	@${COMPOSE_COMMAND} --profile services --profile source up -d --build

.PHONY: compose-build-built
compose-build-built: ## Run 'docker compose --env-file configs/.env --profile services --profile built up -d --build' to create and start containers from built
	@${COMPOSE_COMMAND} --profile services --profile built up -d --build

.PHONY: compose-down
compose-down: ## Run 'docker compose --env-file configs/.env --profile all down' to stop and remove containers and networks
	@${COMPOSE_COMMAND} --profile all down

.PHONY: compose-remove
compose-remove: ## Run 'docker compose --env-file configs/.env --profile all down -v --remove-orphans' to stop and remove containers, networks and volumes
	@echo -n "All registered data and volumes will be deleted, are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	@${COMPOSE_COMMAND} --profile all down -v --remove-orphans

.PHONY: compose-exec-built
compose-exec-built: ## Run 'docker compose --env-file configs/.env exec -it backend_built bash' to access container bash
	@${COMPOSE_COMMAND} exec -it backend_built bash

.PHONY: compose-exec-source
compose-exec-source: ## Run 'docker compose --env-file configs/.env exec -it backend_source bash' to access container bash
	@${COMPOSE_COMMAND} exec -it backend_source bash

.PHONY: compose-log-built
compose-log-built: ## Run 'docker compose --env-file configs/.env logs -f backend_built' to show container logger
	@${COMPOSE_COMMAND} logs -f backend_built

.PHONY: compose-log-source
compose-log-source: ## Run 'docker compose --env-file configs/.env logs -f backend_source' to show container logger
	@${COMPOSE_COMMAND} logs -f backend_source

.PHONY: compose-top
compose-top: ## Run 'docker compose --env-file configs/.env top' to display containers processes
	@${COMPOSE_COMMAND} top

.PHONY: compose-stats
compose-stats: ## Run 'docker compose --env-file configs/.env stats' to display containers stats
	@${COMPOSE_COMMAND} stats

.PHONY: go-run
go-run: ## Run application from source code
	@go run cmd/go-api/go-api.go

.PHONY: go-test
go-test: ## Run tests and generate coverage report
	@go install github.com/axw/gocov/gocov@v1.2.1
	@go install github.com/matm/gocov-html/cmd/gocov-html@v1.4.0
	$(eval packages:=$(shell go list ./... | grep -v github.com/raulaguila/go-api/docs))
	@gocov test $(packages) | gocov-html -t kit > report.html
	-#open -a "Google Chrome" ./report.html

.PHONY: go-build
go-build: ## Build the application from source code
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o backend cmd/go-api/go-api.go

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
