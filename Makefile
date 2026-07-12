APP_NAME=booking-app
IMAGE_NAME=$(APP_NAME)
IMAGE_TAG=latest
DB_URL=postgres://root:root@localhost:5432/go_bookingdb?sslmode=disable

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
	go test -v -cover ./...

lint:
	go vet ./...

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker compose up --build

docker-down:
	docker compose down

migrate-up:
	migrate -path migrations \
	-database "$(DB_URL)" up

migrate-down:
	migrate -path migrations \
	-database "$(DB_URL)" down 1