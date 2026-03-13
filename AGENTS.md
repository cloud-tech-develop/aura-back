# AGENTS.md - Agent Coding Guidelines for aura-back

This document provides guidelines for AI agents working on this codebase.

## Project Overview

- **Language**: Go 1.26.1
- **Database**: PostgreSQL with `lib/pq` driver
- **Migrations**: `golang-migrate/v4` with embedded SQL files
- **HTTP**: Standard library `net/http` (no framework)
- **Architecture**: Multi-tenant SaaS with schema-per-tenant pattern

## Project Structure

```
aura-back/
├── cmd/api/main.go           # Application entry point
├── internal/
│   ├── db/db.go              # Database connection pool
│   └── migration/migrator.go # Migration logic (unused, tenant pkg has newer impl)
├── tenant/
│   ├── manager.go            # Tenant creation & migration management
│   ├── middleware.go         # X-Tenant header validation
│   └── migrations/           # Embedded SQL migrations
└── go.mod
```

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

### Formatting & Vetting
```bash
go fmt ./...
go vet ./...
go mod tidy
```

### Migrations (using golang-migrate CLI)
```bash
migrate create -ext sql -dir tenant/migrations -seq migration_name
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
  )
  ```

### Naming Conventions
- **Types/Functions**: PascalCase (e.g., `Manager`, `NewManager`)
- **Variables/Fields**: camelCase (e.g., `db`, `slug`)
- **Constants**: PascalCase or camelCase with prefix (e.g., `TenantKey`)
- **Packages**: short, lowercase, no underscores
- **SQL Schemas**: lowercase with underscores (e.g., `empresa_uno`)

### Error Handling
- Always check errors immediately after calls
- Wrap errors with `fmt.Errorf("context: %w", err)` for proper error chains
- Return errors rather than logging unless it's the final entry point
- Use sentinel errors only when necessary

### Context Usage
- Pass `context.Context` as first parameter to functions that may timeout/cancel
- Use `context.WithValue` for request-scoped values (like tenant slug)
- Check for cancellation: `if ctx.Err() != nil { return ctx.Err() }`

### Database Operations
- Use prepared statements or parameterized queries to prevent SQL injection
- Always close rows with `defer rows.Close()` after Query/QueryRow
- Use transactions for multi-step operations requiring atomicity
- Set `search_path` per-tenant for schema isolation

### HTTP Handlers
- Read tenant from context, not directly from header (middleware already extracted it)
- Return proper HTTP status codes (200, 201, 400, 404, 500)
- Log errors appropriately before returning

### Multi-Tenant Pattern
- Each tenant gets a PostgreSQL schema named after their slug
- Tenant slug must match regex: `^[a-z0-9_]+$` (lowercase, numbers, underscores)
- Use `X-Tenant` header to identify tenant per request
- All tenant-specific queries run against their schema via search_path

### Code Organization
- Keep business logic in `internal/` or `tenant/` packages
- Entry point in `cmd/`
- Use `embed.FS` for embedding static files (migrations)
- Avoid global state; pass dependencies explicitly

### Testing (when adding tests)
- Test files: `*_test.go` in same package
- Use table-driven tests for multiple cases
- Mock database operations with interfaces
- Test name format: `Test<Function>_<Scenario>`

## Database Credentials

**Do NOT commit credentials.** Use environment variables or `.env`:
```go
dsn := os.Getenv("DATABASE_URL")
// or
dsn := "postgresql://user:pass@host/db?sslmode=require"
```
