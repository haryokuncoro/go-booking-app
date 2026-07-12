APP_NAME=booking-app
IMAGE_NAME=$(APP_NAME)
IMAGE_TAG=latest
DB_HOST ?= localhost
DB_PORT ?= 5433
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= bookingdb
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: help build run test docker-build docker-run docker-up docker-down migrate-up migrate-down lint

help:
	@echo "Available targets:"
	@echo "  make build          Build the Go binary"
	@echo "  make run            Run the application locally"
	@echo "  make test           Run unit tests"
	@echo "  make docker-build   Build Docker image"
	@echo "  make docker-run     Run Docker container"
	@echo "  make docker-up      Start app and dependencies with docker compose"
	@echo "  make docker-down    Stop app and dependencies"
	@echo "  make migrate-up     Apply database migrations"
	@echo "  make migrate-down   Roll back the last migration"

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

lint:
	go vet ./...

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env --name $(APP_NAME) $(IMAGE_NAME):$(IMAGE_TAG)

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

migrate-up:
	migrate -path migrations \
	-database "$(DB_URL)" up

migrate-down:
	migrate -path migrations \
	-database "$(DB_URL)" down 1