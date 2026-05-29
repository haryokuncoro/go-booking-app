# Booking App API

A production-style Booking API built with Golang using Clean Architecture principles.

## Tech Stack

* Go
* Gin
* PostgreSQL
* GORM
* Redis
* JWT Authentication
* Zap Logger
* Goroutines
* Worker Pool
* Context
* Mutex
* golang-migrate
* Docker Compose

---

## Features

### Authentication

* User Registration
* User Login
* JWT Authentication
* Protected Routes

### Booking

* Create Booking
* Get Booking Detail
* List User Bookings
* Room Availability Validation
* Double Booking Prevention

### Performance

* Redis Cache
* Cache Aside Pattern
* Background Email Processing
* Worker Pool

### Observability

* Structured Logging (Zap)
* HTTP Request Logging

### Concurrency

* Goroutines
* Channels
* Worker Pool
* WaitGroup
* Mutex

### Database

* PostgreSQL
* GORM
* Versioned SQL Migrations
* Rollback Support

---


# Prerequisites

Install:

* Go 1.26+
* Docker
* Docker Compose
* golang-migrate

MacOS:

```bash
brew install golang-migrate
```

Verify:

```bash
go version
docker --version
docker compose version
migrate -version
```

---

# Environment Variables

Create `.env`

```env
APP_NAME=Booking App
APP_PORT=8080
APP_ENV=local

DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=bookingdb

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=my-super-secret-key
JWT_EXPIRE_HOUR=24
```

---

# Start Infrastructure

```bash
docker compose up -d
```

Verify:

```bash
docker ps
```

Expected:

```text
booking-postgres
booking-redis
```

---

# Database Migration

This project uses **golang-migrate**.

Schema changes are managed through versioned SQL files.

## Run Migrations

```bash
migrate \
-path migrations \
-database "postgres://admin:admin@localhost:5432/bookingdb?sslmode=disable" \
up
```

## Rollback Last Migration

```bash
migrate \
-path migrations \
-database "postgres://admin:admin@localhost:5432/bookingdb?sslmode=disable" \
down 1
```

## Create New Migration

```bash
migrate create \
-ext sql \
-dir migrations \
create_refresh_tokens_table
```

---

# Makefile Commands

Run migrations:

```bash
make migrate-up
```

Rollback:

```bash
make migrate-down
```

---

# Application Startup

1. Start PostgreSQL & Redis

```bash
docker compose up -d
```

2. Run migrations

```bash
make migrate-up
```

3. Start application

```bash
go run cmd/server/main.go
```

Expected logs:

```text
database connected
redis connected
application started
```

---

# Database Schema

## Users

| Column     | Type                |
| ---------- | ------------------- |
| id         | BIGSERIAL           |
| name       | VARCHAR(100)        |
| email      | VARCHAR(255) UNIQUE |
| password   | VARCHAR(255)        |
| created_at | TIMESTAMP           |
| updated_at | TIMESTAMP           |

## Bookings

| Column       | Type        |
| ------------ | ----------- |
| id           | BIGSERIAL   |
| user_id      | BIGINT      |
| room_id      | BIGINT      |
| booking_date | DATE        |
| status       | VARCHAR(50) |
| created_at   | TIMESTAMP   |
| updated_at   | TIMESTAMP   |

Constraint:

```sql
UNIQUE(room_id, booking_date)
```

This prevents double booking.

---

# API Endpoints

## Public

| Method | Endpoint              |
| ------ | --------------------- |
| POST   | /api/v1/auth/register |
| POST   | /api/v1/auth/login    |
| GET    | /health               |

## Protected

| Method | Endpoint              |
| ------ | --------------------- |
| GET    | /api/v1/me            |
| POST   | /api/v1/bookings      |
| GET    | /api/v1/bookings      |
| GET    | /api/v1/bookings/{id} |

---

# Concurrency Features

## Goroutine

Used for background processing.

## Channel

Used as asynchronous job queue.

## Worker Pool

5 email workers process jobs concurrently.

## WaitGroup

Used for synchronization.

## Mutex

Used to prevent race conditions during booking creation.

---

# Cache Strategy

Cache Aside Pattern.

Booking Detail Request

↓

Redis

├── Hit → Return

└── Miss

↓

PostgreSQL

↓

Cache Result

TTL:

```text
5 minutes
```

---

# Security

* Password Hashing (bcrypt)
* JWT Authentication
* Protected Routes
* Request Validation
* Request Timeout

---

# Logging

Structured logging using Zap.

Example:

```json
{
  "level":"info",
  "msg":"booking created",
  "user_id":1,
  "room_id":10
}
```

---
