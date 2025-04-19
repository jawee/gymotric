# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."

	mjml internal/email/emails/*.mjml -o internal/email/emails/

	@CGO_ENABLED=1 GOOS=darwin go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go &
	@npm install --prefix ./frontend
	@npm run dev --prefix ./frontend

run-api:
	@go run cmd/api/main.go
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

seed:
	@go run cmd/seed/main.go

db-status:
	goose status

db-create:
	cp .default.env .env
	mkdir db
	go run cmd/goose/main.go
	go run cmd/seed/main.go

db-create-docker:
	cp .default.env .env
	mkdir db
	docker build -t goose -f Dockerfile.goose .
	docker run --rm --env-file .env -v ./db:/app/db goose

db-up:
	goose up

db-reset:
	goose reset

db-create-migration:
	goose create a sql

test-coverage:
	rm -r testresults || true
	mkdir testresults
	go test -coverprofile ./testresults/cover.out -o ./testresults ./...
	go tool cover -html=./testresults/cover.out -o ./testresults/cover.html
	echo "Coverage report is saved in testresults/cover.html"

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
