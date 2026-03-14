# AGENTS.md - Agent Coding Guidelines for aura-back

This document provides guidelines for AI agents working on this codebase.

## Project Overview

- **Language**: Go 1.26.1
- **Database**: PostgreSQL with `lib/pq` driver
- **Migrations**: `golang-migrate/v4` with embedded SQL files
- **HTTP**: Standard library `net/http` (no framework)
- **Architecture**: Multi-tenant SaaS with schema-per-tenant pattern
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **Go Version**: 1.26.1

## Project Structure

```
aura-back/
├── cmd/
│   ├── api/main.go            # Dependency wiring & entry point
│   └── server/server.go       # Router, middleware & module registration
├── infrastructure/            # Infrastructure implementations
│   └── messaging/             # Message bus implementations
│       └── memory/            # In-memory event bus
├── internal/
│   └── db/db.go               # Database connection pool
├── shared/                    # Cross-cutting concerns
│   ├── errors/                # Domain-specific errors
│   ├── events/                # Event bus interfaces & base types
│   ├── logging/               # Generic logging handlers
│   └── response/              # Standard HTTP response helpers
├── tenant/
│   ├── manager.go             # Multi-tenant logic (migrations, CRUD)
│   ├── auth.go                # JWT, Login & AuthMiddleware
│   ├── middleware.go          # Tenant middleware
│   └── migrations/            # PostgreSQL migrations (public + tenant)
├── modules/                   # Vertical Feature Modules (Self-contained)
│   └── enterprise/            # Example: Enterprise management
│       ├── domain.go          # Domain entity, Repository & Service interfaces
│       ├── service.go         # Business logic implementation
│       ├── repository.go      # PostgreSQL implementation of repository
│       ├── handler.go         # HTTP handlers
│       ├── routes.go          # Module route registration
│       └── logger.go          # Module-specific event logging
└── go.mod
```

### Module Guidelines
- Each POS feature (Sales, Products, etc.) must be its own module in `modules/`.
- Modules should be decoupled; avoid cross-module imports where possible.
- Use `shared/` for truly global logic and interfaces.
- Always use `context.Context` for repository and service methods.
- Modules should register their routes via a `RegisterRoutes` function.

## Build, Lint & Test Commands

### Running the Application
```bash
go run ./cmd/api/main.go
```

### Build
```bash
go build ./...
```

### Run Tests
```bash
go test ./...
```

### Run Single Test
```bash
go test -run TestName ./...
go test -v -run TestName ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Formatting & Vetting
```bash
go fmt ./...
go vet ./...
go mod tidy
```

### Linting (if golangci-lint is installed)
```bash
golangci-lint run
```

### Migrations (using golang-migrate CLI)
```bash
migrate create -ext sql -dir tenant/migrations -seq migration_name
migrate -path tenant/migrations -database "postgres://user:pass@localhost/db?sslmode=disable" up
migrate -path tenant/migrations -database "postgres://user:pass@localhost/db?sslmode=disable" down
```

## Code Style Guidelines

### Imports
- Standard library first, then external packages
- Grouped with blank line between groups:
  ```go
  import (
      "context"
      "fmt"
      "net/http"
      
      "github.com/golang-migrate/migrate/v4"
      "github.com/lib/pq"
      "github.com/joho/godotenv"
  )
  ```
- Internal packages use the full module path:
  ```go
  "github.com/cloud-tech-develop/aura-back/internal/db"
  ```
- Use blank lines between groups (standard, external, internal)

### Naming Conventions
- **Types/Functions**: PascalCase (e.g., `Manager`, `NewManager`)
- **Variables/Fields**: camelCase (e.g., `db`, `slug`)
- **Constants**: PascalCase or camelCase with prefix (e.g., `TenantKey`)
- **Packages**: short, lowercase, no underscores
- **SQL Schemas**: lowercase with underscores (e.g., `empresa_uno`)
- **Interface Names**: Should end with -er when they have one method (e.g., `Migrator`, `Logger`)
- **Struct Fields**: Use camelCase, export only when needed

### Error Handling
- Always check errors immediately after calls
- Wrap errors with `fmt.Errorf("context: %w", err)` for proper error chains
- Return errors rather than logging unless it's the final entry point
- Use sentinel errors only when necessary (defined in `shared/errors/`)
- Handle specific errors with `errors.Is()` and `errors.As()` when appropriate
- Return `sql.ErrNoRows` from repository methods when no data is found

### Context Usage
- Pass `context.Context` as first parameter to functions that may timeout/cancel
- Use `context.WithValue` for request-scoped values (like tenant slug)
- Check for cancellation: `if ctx.Err() != nil { return ctx.Err() }`
- Never store contexts in struct types; pass them explicitly
- Use `context.Background()` only at the top level (main, tests)

### Database Operations
- Use prepared statements or parameterized queries to prevent SQL injection
- Always close rows with `defer rows.Close()` after Query/QueryRow
- Use transactions for multi-step operations requiring atomicity
- Set `search_path` per-tenant for schema isolation
- Handle `sql.ErrNoRows` appropriately
- Use `QueryRowContext` and `ExecContext` when context is available

### HTTP Handlers
- Read tenant from context, not directly from header (middleware already extracted it)
- Return proper HTTP status codes (200, 201, 400, 404, 500)
- Log errors appropriately before returning
- Use helper functions from `shared/response/` for consistent responses
- Validate request bodies early and return 400 for invalid input
- Use `context.WithTimeout` for outgoing HTTP calls

### Multi-Tenant Pattern
- Each tenant gets a PostgreSQL schema named after their slug
- Tenant slug must match regex: `^[a-z0-9_]+$` (lowercase, numbers, underscores)
- Use `X-Tenant` header to identify tenant per request
- All tenant-specific queries run against their schema via search_path
- Tenant middleware should set the search_path and validate the slug
- Public tables (in public schema) are shared across tenants
- When creating tenants, run migrations against their schema
- Email validation: Check email uniqueness across all enterprises (in public schema)

### Code Organization
- Keep business logic in `internal/` or `tenant/` packages
- Entry point in `cmd/`
- Use `embed.FS` for embedding static files (migrations)
- Avoid global state; pass dependencies explicitly
- Group related constants in `const` blocks
- Use receiver methods appropriately (value vs pointer)
- Keep functions focused and under 50 lines when possible

### Testing (when adding tests)
- Test files: `*_test.go` in same package
- Use table-driven tests for multiple cases
- Mock database operations with interfaces
- Test name format: `Test<Function>_<Scenario>`
- Test both success and error paths
- Use `t.Parallel()` for independent tests
- Set up test data in `TestMain` when needed for the package
- Don't test private functions directly; test through public interface

### Event Bus Usage
- Use the event bus from `shared/events/` for loose coupling
- Define event types as constants
- Publishers should not know about subscribers
- Handle event processing errors gracefully
- Use synchronous processing for critical paths
- Consider dead letter queues for failed events

### Logging
- Use the logger from `shared/logging/`
- Log at appropriate levels (debug, info, warn, error)
- Include context in logs (request ID, tenant ID, etc.)
- Don't log sensitive information (passwords, tokens)
- Use structured logging when possible
- Pass logger as dependency rather than using globals

## Dependency Management
- Use `go mod tidy` regularly to clean up dependencies
- Pin versions in go.mod for reproducible builds
- Vendoring is not used; rely on module proxy
- Update dependencies with `go get -u ./...`
- Check for outdated packages with `go list -m -u all`

## Security Guidelines
- Never hardcode credentials; use environment variables
- Validate all input (both syntactic and semantic)
- Use parameterized queries to prevent SQL injection
- Set secure flags on cookies when using them
- Implement rate limiting at the middleware level
- Hash passwords using bcrypt or similar (see `tenant/auth.go`)
- Use environment-specific configuration
- Regularly update dependencies to patch vulnerabilities
- Email uniqueness validation must be enforced at service level

## Docker Guidelines (if applicable)
- Use multi-stage builds to minimize image size
- Run as non-root user in containers
- Cache dependency downloads
- Expose only necessary ports
- Use healthchecks
- Don't store secrets in images; use secrets management

## CI/CD Considerations
- Run `go vet` and `go fmt -l` in CI to catch style issues
- Run tests with race detector: `go test -race ./...`
- Build for multiple platforms if needed: `GOOS=linux GOARCH=amd64 go build`
- Tag Docker images with git SHA and version
- Use blue-green deployments for zero-downtime releases

## AGENTS.md Guidelines
- This file should be kept up to date with any architectural or process changes
- New agents should read this file before making changes to the codebase
- If there are specific rules for AI agents (e.g., Cursor rules, Copilot instructions), include them here
