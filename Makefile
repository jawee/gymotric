# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go &
	@npm install --prefix ./frontend
	@npm run dev --prefix ./frontend
# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

db-status:
	@GOOSE_DRIVER=sqlite3 GOOSE_MIGRATION_DIR=./migrations GOOSE_DBSTRING=./db/test.db goose status

# goose sqlite3 ./foo.db status
#    goose sqlite3 ./foo.db create init sql
#    goose sqlite3 ./foo.db create add_some_column sql
#    goose sqlite3 ./foo.db create fetch_user_data go
#    goose sqlite3 ./foo.db up
#
db-create-migration:
	GOOSE_DRIVER=sqlite3 GOOSE_MIGRATION_DIR=./migrations GOOSE_DBSTRING=./db/test.db goose create sql

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch
