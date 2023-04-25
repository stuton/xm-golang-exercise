include config/local.env

PATH_TO_MIGRATIONS=migrations
MIGRATE_VERSION=4.15.2
MIGRATE_COMMAND := docker run -v $(shell pwd)/$(PATH_TO_MIGRATIONS):/migrations --network host migrate/migrate:v$(MIGRATE_VERSION) -path=/migrations/ -database "$(DATABASE_URI)"

GOLANGCI_LINT_VERSION=1.52.2
GOLANGCI_LINT_COMMAND := docker run -t --rm -v $(shell pwd):/app -v ~/.cache/golangci-lint/v$(GOLANGCI_LINT_VERSION):/root/.cache -w /app golangci/golangci-lint:v$(GOLANGCI_LINT_VERSION) golangci-lint run -v

.PHONY: migrate
migrate:
	@echo "Running all new database migrations..."
	$(MIGRATE_COMMAND) up

.PHONY: migrate-down
migrate-down:
	@echo "Reverting database to the last migration step..."
	$(MIGRATE_COMMAND) down 1

.PHONY: migrate-new
migrate-new:
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE_COMMAND) create -ext sql -dir /migrations/ $${name// /_}

.PHONY: migrate-reset
migrate-reset:
	@echo "Resetting database..."
	$(MIGRATE_COMMAND) drop -f
	@echo "Running all new database migrations..."
	$(MIGRATE_COMMAND) up

run:
	@echo "Running docker containers..."
	@docker-compose up --build

down:
	@echo "Shutdown docker containers..."
	@docker-compose down -v

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	@echo "Linting project..."
	$(GOLANGCI_LINT_COMMAND)

.PHONY: tidy
tidy:
	go mod tidy
