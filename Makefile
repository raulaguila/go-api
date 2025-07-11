DOCKER_BASE ?= $(or ${base}, built)
DOCKER_SERVICE ?= $(or ${service}, backend)
COMPOSE_COMMAND := BASE=${DOCKER_BASE} docker compose --env-file configs/.env -f build/docker/compose.yml
GOBUILD_COMMAND := CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s"

define clean_dangling_images
	@echo "\033[1;33m🧹 Cleaning up dangling Docker images...\n\033[0m"
	@if test -n "$$(docker images -f "dangling=true" -q)"; then \
		docker rmi $$(docker images -f "dangling=true" -q); \
	fi > /dev/null
endef

.PHONY: all
all: help
help: ## Display available commands and their descriptions
	@echo "\033[1;36mUsage:\033[0m"
	@echo "  make [COMMAND]\n"
	@echo "\033[1;36mExample:\033[0m"
	@echo "  make build\n"
	@echo "\033[1;36mCommands:\033[0m\n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: ## Create environment file
	@echo "\033[1;34m⚙️  Initializing environment setup...\033[0m"
	@chmod +x configs/env.sh && configs/env.sh && mv .env configs/
	@echo "\033[1;32m✅ Environment file successfully created and configured!\033[0m\n"

.PHONY: test
test: ## Run tests and generate coverage report
	@echo "\033[1;34m🔍 Running tests...\033[0m"
	@-go install github.com/axw/gocov/gocov@v1.2.1
	@-go install github.com/matm/gocov-html/cmd/gocov-html@v1.4.0
	@-go clean -testcache
	@-go test ./... -coverprofile cover.out
	@-go tool cover -html=cover.out
	@gocov convert cover.out | gocov-html -t kit > report.html
	@echo "\033[1;32m✅ Tests completed!\033[0m\n"
	@-x-www-browser ./report.html

.PHONY: run
run: ## Run application from source code
	@echo "\033[1;36m▶️  Running the application...\033[0m"
	@go run cmd/backend/backend.go
	@echo "\033[1;32m✅ Application stopped.\033[0m"

.PHONY: build
build: ## Build the all applications from source code
	@echo "\033[1;34m🚀 Building application...\033[0m"
	@${GOBUILD_COMMAND} -o binbackend cmd/backend/backend.go
	@${GOBUILD_COMMAND} -o bingenerator cmd/generator/generator.go
	@echo "\033[1;32m✅ Build completed successfully!\033[0m\n"

.PHONY: swag
swag: ## Update swagger files
	@echo "\033[1;34m📄 Updating Swagger API documentation...\033[0m"
	@go run github.com/swaggo/swag/cmd/swag@v1.16.3 init -g cmd/backend/backend.go --parseDependency 2>&1 > /dev/null
	@echo "\033[1;32m✅ Swagger files updated successfully.\033[0m\n"

.PHONY: format
format: ## Fix code format issues
	@echo "\033[1;33m📝 Formatting code to fix style issues...\033[0m"
	@go run mvdan.cc/gofumpt@latest -w -l . 2>&1 > /dev/null
	@echo "\033[1;32m✅ Code formatting complete! All issues fixed.\033[0m\n"

.PHONY: tidy
tidy: ## Clean and tidy dependencies
	@echo "\033[1;33m🔧 Cleaning and tidying Go dependencies...\033[0m"
	@go mod tidy -v 2>&1 > /dev/null
	@echo "\033[1;32m✅ Dependencies cleaned and tidied successfully.\033[0m\n"

.PHONY: lint
lint: ## Run lint checks
	@echo "\033[1;33m🔍 Running lint checks on the code...\033[0m"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1 run ./... --fix
	@echo "\033[1;32m✅ Linting complete! Code issues fixed where possible.\033[0m\n"

.PHONY: audit
audit: ## Conduct quality checks
	@echo "\033[1;33m🔎 Running code audit...\033[0m"
	@go mod verify 2>&1 > /dev/null
	@go vet ./... 2>&1 > /dev/null
	@go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./... 2>&1 > /dev/null
	@echo "\033[1;32m✅ Code audit finished!\033[0m\n"

.PHONY: benchmark
benchmark: ## Benchmark code performance
	@echo "\033[1;35m⚡ Running benchmarks...\033[0m"
	@go test ./... -benchmem -bench=. -run=^Benchmark_$ 2>&1 > /dev/null
	@echo "\033[1;32m✅ Benchmark completed!\033[0m\n"

### Docker compose commands  ---------------------------------------------

.PHONY: compose-up
compose-up: ## Create and start containers
	@echo "\033[1;34m🚀 Starting Docker containers...\033[0m"
	@${COMPOSE_COMMAND} up -d 2>&1 > /dev/null
	@echo "\033[1;32m✅ Containers are up and running!\033[0m\n"

.PHONY: compose-build
compose-build: ## Build, create and start containers
	@echo "\033[1;34m🚢 Building and starting Docker containers...\033[0m"
	@${COMPOSE_COMMAND} up -d --build 2>&1 > /dev/null
	@echo "\033[1;32m✅ Containers are up and running!\033[0m\n"
	@$(clean_dangling_images)
	
.PHONY: compose-down
compose-down: ## Stop and remove containers and networks
	@echo "\033[1;33m🛑 Stopping and removing containers...\033[0m"
	@${COMPOSE_COMMAND} down 2>&1 > /dev/null
	@echo "\033[1;32m✅ Containers stopped.\033[0m\n"

.PHONY: compose-clean
compose-clean: ## Clear dangling Docker images
	$(clean_dangling_images)

.PHONY: compose-remove
compose-remove: ## Stop and remove containers, networks and volumes
	@echo "\033[1;31m⚠️  WARNING: This will permanently delete all containers, networks, and VOLUMES!\033[0m"
	@echo -n "\033[1;31m❌ All data will be lost. Are you sure? [y/N] \033[0m" && read ans && [ $${ans:-N} = y ]
	@echo "\033[1;33m\n🛑 Stopping and removing all Docker resources...\033[0m"
	@${COMPOSE_COMMAND} down -v --remove-orphans 2>&1 > /dev/null
	@echo "\033[1;32m✅ Containers, networks, and volumes removed successfully.\033[0m\n"

.PHONY: compose-exec
compose-exec: ## Access container bash
	@echo "\033[1;34m🔑 Accessing the container shell...\033[0m"
	@${COMPOSE_COMMAND} exec -it ${DOCKER_SERVICE} bash

.PHONY: compose-log
compose-log: ## Show container logger
	@echo "\033[1;34m📜 Fetching container logs...\033[0m"
	@${COMPOSE_COMMAND} logs -f ${DOCKER_SERVICE}

.PHONY: compose-top
compose-top: ## Display containers processes
	@echo "\033[1;34m🔍 Displaying container processes...\033[0m"
	@${COMPOSE_COMMAND} top

.PHONY: compose-stats
compose-stats: ## Display containers stats
	@echo "\033[1;36m📊 Showing container statistics...\033[0m"
	@${COMPOSE_COMMAND} stats
