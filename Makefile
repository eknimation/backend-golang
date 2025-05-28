# Makefile for Backend Go API

.PHONY: help build up down logs restart clean test swagger

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the Docker images
	docker-compose build --no-cache

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

logs: ## Show logs from all services
	docker-compose logs -f

logs-api: ## Show logs from API service only
	docker-compose logs -f api

logs-db: ## Show logs from MongoDB service only
	docker-compose logs -f mongodb

restart: ## Restart all services
	docker-compose restart

restart-api: ## Restart API service only
	docker-compose restart api

clean: ## Remove all containers, networks, and volumes
	docker-compose down -v --rmi all --remove-orphans

clean-volumes: ## Remove all volumes (WARNING: This will delete all data)
	docker-compose down -v

status: ## Show status of all services
	docker-compose ps

shell-api: ## Open shell in API container
	docker-compose exec api sh

shell-db: ## Open MongoDB shell
	docker-compose exec mongodb mongosh -u root -p password --authenticationDatabase admin

test: ## Run tests inside the container
	docker-compose exec api go test ./...

swagger: ## Generate Swagger documentation
	docker-compose exec api swag init -g cmd/api/main.go -o docs

dev: ## Start services for development
	docker-compose up --build

dev-rebuild: ## Rebuild and start services for development
	docker-compose up --build --force-recreate

# Database commands
db-backup: ## Backup MongoDB database
	docker-compose exec mongodb mongodump --uri="mongodb://root:password@localhost:27017/backend_challenge?authSource=admin" --out=/tmp/backup
	docker cp mongodb:/tmp/backup ./backup-$(shell date +%Y%m%d_%H%M%S)

db-restore: ## Restore MongoDB database (usage: make db-restore BACKUP_PATH=./backup-folder)
	@if [ -z "$(BACKUP_PATH)" ]; then echo "Please specify BACKUP_PATH. Usage: make db-restore BACKUP_PATH=./backup-folder"; exit 1; fi
	docker cp $(BACKUP_PATH) mongodb:/tmp/restore
	docker-compose exec mongodb mongorestore --uri="mongodb://root:password@localhost:27017/backend_challenge?authSource=admin" /tmp/restore

# Monitoring
health: ## Check health of all services
	@echo "Checking API health..."
	@curl -s http://localhost:5555/health | jq . || echo "API not responding"
	@echo "\nChecking MongoDB health..."
	@docker-compose exec mongodb mongosh --eval "db.adminCommand('ping')" --quiet || echo "MongoDB not responding"

# Development helpers
install: ## Install dependencies locally
	go mod download
	go mod tidy

lint: ## Run linter
	golangci-lint run

format: ## Format code
	go fmt ./...
