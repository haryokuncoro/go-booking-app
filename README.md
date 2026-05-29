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
* Worker Pool
* Goroutines
* Context
* Mutex
* Docker Compose

---

## Features

### Authentication

* Register User
* Login User
* JWT Authentication
* Protected Routes

### Booking

* Create Booking
* Get Booking Detail
* Get User Bookings
* Room Availability Validation
* Double Booking Prevention

### Performance

* Redis Cache
* Cache Aside Pattern
* Background Email Processing
* Worker Pool

### Observability

* Structured Logging with Zap
* Request Logging

### Concurrency

* Goroutines
* Channels
* Worker Pool
* WaitGroup
* Mutex

### Reliability

* Context Propagation
* Request Timeout
* Graceful Error Handling

---


# Prerequisites

Install:

* Go 1.26+
* Docker
* Docker Compose
* PostgreSQL Client (optional)

Verify:

```bash
go version
docker --version
docker compose version
```

---

# Installation

Clone project:

```bash
git clone https://github.com/yourname/booking-app.git

cd booking-app
```

Install dependencies:

```bash
go mod tidy
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

# Run Application

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

# API Endpoints

## Health Check

### Request

```http
GET /health
```

### Response

```json
{
  "status": "UP"
}
```

---

# Register

### Request

```http
POST /api/v1/auth/register
```

```json
{
  "name": "Haryo",
  "email": "haryo@test.com",
  "password": "123456"
}
```

### Response

```json
{
  "success": true,
  "message": "user registered"
}
```

---

# Login

### Request

```http
POST /api/v1/auth/login
```

```json
{
  "email": "haryo@test.com",
  "password": "123456"
}
```

### Response

```json
{
  "success": true,
  "message": "login success",
  "data": {
    "access_token": "jwt-token"
  }
}
```

---

# Current User

### Request

```http
GET /api/v1/me
```

Headers:

```text
Authorization: Bearer <token>
```

### Response

```json
{
  "id": 1,
  "name": "Haryo",
  "email": "haryo@test.com"
}
```

---

# Create Booking

### Request

```http
POST /api/v1/bookings
```

Headers:

```text
Authorization: Bearer <token>
```

Body:

```json
{
  "room_id": 1,
  "date": "2026-06-20"
}
```

### Response

```json
{
  "success": true,
  "message": "booking created"
}
```

---

# List My Bookings

### Request

```http
GET /api/v1/bookings
```

Headers:

```text
Authorization: Bearer <token>
```

### Response

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "room_id": 1,
      "status": "CONFIRMED"
    }
  ]
}
```

---

# Booking Detail

### Request

```http
GET /api/v1/bookings/{id}
```

Headers:

```text
Authorization: Bearer <token>
```

---

# Concurrency Features

## Goroutine

Used for background processing.

```go
go StartEmailWorker(...)
```

## Channel

Used as queue.

```go
EmailQueue <- job
```

## Worker Pool

```go
for i := 0; i < 5; i++ {
    go StartEmailWorker(i, EmailQueue)
}
```

## WaitGroup

```go
wg.Add(1)
go process()
wg.Wait()
```

## Mutex

Used to prevent concurrent booking conflicts.

```go
roomLock.Lock()
defer roomLock.Unlock()
```

---

# Cache Strategy

Cache Aside Pattern.

```text
Request
   │
   ▼
Redis
 │      │
Hit     Miss
 │      │
 ▼      ▼
Return  Database
          │
          ▼
      Cache Result
```

TTL:

```text
5 minutes
```

---

# Security

* Password Hashing (bcrypt)
* JWT Authentication
* Protected Routes
* Validation
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