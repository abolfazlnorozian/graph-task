# Graph Task Manager Service

A simple **Task Manager ** implemented in **Go (Gin)** with PostgreSQL, designed as part of the GRAPH Golang Developer assessment.

## 🚀 Features

- RESTful CRUD API for managing tasks
- Gin framework for HTTP routing
- PostgreSQL for data persistence
- Swagger / OpenAPI documentation
- Docker & docker-compose support
- Basic observability (metrics & tracing)
- Migration-based database schema management
- Test-ready architecture (unit + integration)


## 📁 Project Structure
```text
graph-task-service/
├── cmd/server/        # Application entry point
│   └── main.go
├── internal/
│   ├── config         # Environment & app configuration
│   ├── domain         # Core domain models
│   ├── handler        # HTTP handlers (Gin)
│   ├── middleware     # Custom middlewares
│   ├── observability  # Metrics & tracing
│   ├── repository    # Data access layer (Postgres)
│   ├── router         # Route definitions
│   └── service        # Business logic
├── migration/         # SQL migrations
│   └── 001_create_tasks.sql
├── docs/              # Swagger docs
│   ├── swagger.json
│   ├── swagger.yaml
│   └── docs.go
├── docker-compose.yml
├── Dockerfile
├── .env
├── example.env
└── go.mod

## ▶️ Run Locally

### Prerequisites
- Go **1.24+**
- Docker & Docker Compose
- PostgreSQL (via Docker)

### Start database
```bash
docker compose up -d postgres

go run ./cmd/server/main.go

### API Documentation (Swagger)

- http://localhost:8080/swagger/index.html

## Regenerate Swagger docs
```bash
 swag init -g cmd/server/main.go



## 📊 Observability

- Prometheus metrics exposed (tasks_count, request_latency_histogram, requests_total)
- Basic tracing implemented for HTTP requests


## ✅ Testing

- Designed for **TDD**
- Supports unit tests with mocked repositories
- Integration tests with PostgreSQL
- Target: **≥ 70% test coverage**

## 📝 Notes

- Database schema is managed via SQL migrations
- `RUN_MIGRATIONS` can be disabled after first run
- Clean architecture with clear separation of concerns
