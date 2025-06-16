# Go API Boilerplate

A modern REST API boilerplate built with Go, featuring:

- **Huma** - OpenAPI 3.0 schema generation and input/output validation
- **Chi** - Lightweight and fast HTTP router
- **SQLC** - Compile-time safe SQL query generation
- **Zerolog** - High-performance structured logging
- **PostgreSQL** - Database backend
- **Docker** - Containerization support

## Features

- ✅ RESTful API endpoints with OpenAPI documentation
- ✅ Type-safe database operations with SQLC
- ✅ Database migrations support
- ✅ Request logging and error handling
- ✅ Health check endpoint with database connectivity
- ✅ Environment-based configuration
- ✅ Code quality tools (golangci-lint)
- ✅ Hot-reload development environment

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Docker (optional)

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd go-api-boilerplate
```

2. Start the application with Docker Compose:
```bash
make compose
```

The API will be available at `http://localhost:8080`

## API Documentation

Once the server is running, you can access the interactive API documentation at:

- Swagger UI: `http://localhost:8080/v1/docs`
- OpenAPI JSON: `http://localhost:8080/v1/openapi.json`

### Preview

![Documentation preview](./images/docs.png)

## Database Access

PGweb UI is available at `http://localhost:8081` for database management.

## Database Schema

The application uses a simple `visit` table:

```sql
CREATE TABLE visit (
    ip inet not null,
    name varchar not null,
    visited_at timestamptz not null,
    primary key (ip, visited_at)
);
```

## Configuration

The application can be configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://admin:password@localhost:5432/database?sslmode=disable` | PostgreSQL connection string |
| `SERVER_PORT` | `8080` | Server port |
| `SERVER_HOST` | `0.0.0.0` | Server host |
| `API_BASE_PATH` | `/v1` | API base path |
| `SERVICE_NAME` | `api` | Service name |
| `ENVIRONMENT` | `local` | Environment name |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |

## Development

### Running Tests

```bash
go test ./...
```

### Code Quality

Run linting and formatting:

```bash
make lint
```

### Building for Production

```bash
go build -o bin/api ./main.go
```

### Database Operations

Generate SQL code from queries:

```bash
make generate-sql
```

Run database migrations:

```bash
make migrate-up
```

## Docker Support

Build and run with Docker:

```bash
docker build -t go-api .
docker run -p 8080:8080 go-api
```

## Project Structure

```
go-api-boilerplate/
├── main.go                     # Application entry point
├── internal/
│   ├── app/                    # Application setup and routing
│   │   ├── app.go             # Main application struct and startup
│   │   ├── routes.go          # Route definitions
│   │   ├── handlers/          # Request handlers
│   │   │   ├── health.go      # Health check handler
│   │   │   └── greetings.go   # Greeting endpoints
│   │   ├── middleware/        # HTTP middleware
│   │   │   └── request.go     # Request middleware
│   │   └── models/            # API models/DTOs
│   │       ├── health.go      # Health check models
│   │       └── greetings.go   # Greeting models
│   ├── config/                # Configuration management
│   │   ├── app.go            # Application configuration
│   │   └── database.go       # Database configuration
│   ├── database/              # Database connection
│   │   └── database.go       # Database setup and connection
│   └── store/                 # Generated SQLC code
│       ├── db.go             # Database interface
│       ├── models.go         # Generated models
│       └── visits.sql.go     # Generated queries
├── queries/                   # SQL queries for SQLC
│   └── visits.sql            # Visit-related queries
├── migrations/                # Database migrations
│   ├── 001_init.up.sql       # Initial schema
│   └── 001_init.down.sql     # Rollback schema
├── images/                    # Documentation assets
│   └── docs.png              # API documentation screenshot
├── Dockerfile                 # Docker configuration
├── compose.yml               # Docker Compose setup
├── Makefile                  # Build and development commands
├── sqlc.yaml                 # SQLC configuration
├── .golangci.yml            # Linting configuration
├── go.mod                    # Go module definition
└── go.sum                    # Go module checksums
```

## Available Commands

- `make compose` - Start the application with Docker Compose
- `make generate-sql` - Generate Go code from SQL queries
- `make lint` - Run code formatting and linting
- `make migrate-up` - Run database migrations (when implemented)
- `make migrate-down` - Rollback database migrations (when implemented)
