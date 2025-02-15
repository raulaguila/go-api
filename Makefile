COMPOSE_COMMAND = docker compose --env-file configs/.env
GOBUILD_COMMAND = CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s"

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

.PHONY: build
build: ## Build the application from source code
	@${GOBUILD_COMMAND} -o backend cmd/go-api/go-api.go
	@${GOBUILD_COMMAND} -o generator cmd/generator/generator.go

.PHONY: run
run: ## Run application from source code
	@go run cmd/go-api/go-api.go

.PHONY: test
test: ## Run tests and generate coverage report
	@go install github.com/axw/gocov/gocov@v1.2.1
	@go install github.com/matm/gocov-html/cmd/gocov-html@v1.4.0
	@go clean -testcache
	@go test -coverprofile=cover.out ./...
	@go tool cover -html=cover.out
	@gocov convert cover.out | gocov-html -t kit > report.html
	@-open ./report.html

### Docker compose commands  ---------------------------------------------

.PHONY: compose-build-services
compose-build-services: ## Create and start services containers
	@#BUILDKIT_PROGRESS=plain ${COMPOSE_COMMAND} up -d --build
	@${COMPOSE_COMMAND} -f build/services.compose.yml up -d --build

.PHONY: compose-build-built
compose-build-built: ## Create and start containers from built
	@${COMPOSE_COMMAND} -f build/built.compose.yml up -d --build

.PHONY: compose-build-source
compose-build-source: ## Create and start containers from source code
	@${COMPOSE_COMMAND} -f build/source.compose.yml up -d --build

.PHONY: compose-down
compose-down: ## Stop and remove containers and networks
	@${COMPOSE_COMMAND} -f build/built.compose.yml down

.PHONY: compose-remove
compose-remove: ## Stop and remove containers, networks and volumes
	@echo -n "All registered data and volumes will be deleted, are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	@${COMPOSE_COMMAND} -f build/built.compose.yml down -v --remove-orphans

.PHONY: compose-exec
compose-exec: ## Access container bash
	@${COMPOSE_COMMAND} -f build/built.compose.yml exec -it backend bash

.PHONY: compose-log
compose-log: ## Show container logger
	@${COMPOSE_COMMAND} -f build/built.compose.yml logs -f backend

.PHONY: compose-top
compose-top: ## Display containers processes
	@${COMPOSE_COMMAND} -f build/built.compose.yml top

.PHONY: compose-stats
compose-stats: ## Display containers stats
	@${COMPOSE_COMMAND} -f build/built.compose.yml stats

###   ---------------------------------------------

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
	@go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g cmd/go-api/go-api.go --parseDependency

.PHONY: go-format
go-format: ## Fix code format issues
	@go run mvdan.cc/gofumpt@latest -w -l .

.PHONY: go-tidy
go-tidy: ## Clean and tidy dependencies
	@go mod tidy -v
