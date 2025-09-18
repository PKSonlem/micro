# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based microservices application with a clean architecture pattern. The service provides an API for managing houses, flats, and user authentication with role-based access control (client/moderator).

## Key Commands

### Development
```bash
cd app

# Install dependencies and tools
make install

# Run the application with Docker Compose
make up                 # Start all services (app + PostgreSQL)
make down              # Stop all services

# Testing
make test              # Run unit tests
make test-integration  # Run integration tests with Testcontainers

# Code quality
make fmt               # Format code with go fmt and golangci-lint --fix
make lint              # Run golangci-lint checks
make lint-fix          # Auto-fix linting issues

# Code generation
make generate          # Run go generate (for mocks and OpenAPI server)
```

### CI/CD Commands (GitLab CI)
```bash
# These run automatically in CI pipeline on merge requests:
cd app
gofmt -l .                              # Check formatting
golangci-lint run                       # Lint code
go build -v ./...                       # Build all packages
go test ./... -v --tags=integration    # Run integration tests
```

## Architecture

### Directory Structure
- `app/` - Main application directory
  - `api/` - OpenAPI specification (api.yaml)
  - `cmd/` - Application entry point and dependency injection container
  - `internal/` - Core business logic following clean architecture:
    - `entity/` - Domain models
    - `handler/` - HTTP handlers and middleware
    - `repository/` - Data access layer
    - `service/` - Business logic layer
    - `usecase/` - Application use cases
    - `generated/` - Auto-generated OpenAPI server code
  - `migrations/` - Database migration files
  - `pkg/` - Reusable packages
  - `tests/integration/` - Integration tests using Testcontainers
  - `deploy/local/` - Docker Compose setup and environment configs

### Key Technologies
- **Framework**: Standard library HTTP server with OpenAPI code generation
- **Database**: PostgreSQL with sqlx and Squirrel for query building
- **Testing**: Testcontainers for integration tests, gomock for unit tests
- **Authentication**: JWT-based with role-based access (client/moderator)
- **Validation**: go-playground/validator
- **Monitoring**: Prometheus metrics
- **Migrations**: golang-migrate

### Environment Configuration
The application uses environment variables loaded from `deploy/local/.env`:
- PostgreSQL connection settings (POSTGRES_*)
- Server configuration (SERVER_HOST, SERVER_PORT)

For integration tests, copy `deploy/local/.test.env.example` to `deploy/local/.test.env`.

### Code Generation
The project uses code generation for:
1. OpenAPI server stubs from `api/api.yaml` → `internal/generated/server.go`
2. Mock interfaces using gomock (via `go generate ./...`)

Run `make generate` after modifying API specs or interfaces with mock annotations.