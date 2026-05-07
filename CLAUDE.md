# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run development server
go run ./cmd/api/main.go

# Build
go build ./cmd/api/main.go

# Run all tests
go test ./...

# Run tests for a specific module
go test -v ./modules/catalog/products/...

# Run a single test
go test -v -run TestService_Create_ValidSlugFormat ./modules/catalog/products/...

# Format and vet
go fmt ./...
go vet ./...

# Create a new migration
migrate create -ext sql -dir tenant/migrations/tenant -seq migration_name

# Docker
docker-compose up --build
```

## Architecture Overview

This is a multi-tenant POS (Point of Sale) backend with dual Postgres/SQLite support for offline mode.

### Entry Points

- `cmd/api/main.go` — loads `.env`, initializes DB and migrations, wires all modules, generates the offline binary (`static/bin/aura-pos-offline.exe`) at startup when in Postgres mode.
- `cmd/server/server.go` — sets up the Gin router, CORS, auth middleware, and calls each module's `Register()` function.

### Database Strategy

**Multi-tenant via Postgres schemas:**
- `public` schema holds shared tables: `enterprises`, `users`, `roles`, `user_roles`, `plans`.
- Each enterprise gets its own schema (e.g., `empresa_uno`) with tenant-specific tables: products, sales, inventory, etc.
- All tenant queries use `db.WithSchema(q, slug)` to scope the schema.

**Dual-mode (Postgres / SQLite):**
- `DATABASE_URL` present → Postgres mode (port 8081)
- `DATABASE_URL` absent → SQLite offline mode using `aura_pos.db` (port 8091)
- `internal/db/` provides a `Querier` interface that transparently adapts Postgres queries to SQLite: `$1` → `?`, `ILIKE` → `LIKE`, `NOW()` → `CURRENT_TIMESTAMP`, removes schema prefixes.

### Module Pattern

Every feature module follows the same layout:

```
modules/<domain>/<feature>/
  domain.go       — entity structs, Repository/Service interfaces, event constants
  service.go      — business logic, returns interface not struct
  repository.go   — DB queries via db.Querier (compatible with both drivers)
  handler.go      — HTTP request/response types + Gin handler methods
  routes.go       — Register(public, protected, handler) function
  logger.go       — subscribes to domain events for audit logging
  *_test.go       — unit tests with testify mocks and go-sqlmock
```

New modules are wired in `cmd/api/main.go` (dependency injection) and `cmd/server/server.go` (route registration).

### Authentication & Middleware Stack

Request flow:
```
CORS → AuthMiddleware (JWT + IP validation) → Tenant resolution (slug from JWT) → Handler
```

JWT claims include: `user_id`, `enterprise_id`, `slug`, `email`, `roles`, `role_level`, `ip`. The token IP must match the client IP on every request. Role levels: 0=SUPERADMIN, 1=ADMIN, 2=SUPERVISOR, 3+=USER.

Public routes (no auth): `POST /login`, `POST /enterprises`, `GET /enterprises/:slug`.

### Event Bus & Online↔Offline Sync

**Local events** (audit logs, internal) use the in-memory bus in both modes.

**Cross-server sync** uses RabbitMQ (`USE_RABBITMQ=true`), topic exchange `aura.{env}`.

Routing key convention: `{entity}.{direction}.{action}.{slug}`

```
Online → Offline:  product.offline.created.empresa_uno  → queue aura.{env}.product.offline.created.empresa_uno
Offline → Online:  product.online.created.empresa_uno   → queue aura.{env}.product.online.created
                                                           binding: product.online.created.*  (wildcard, all tenants)
```

**Online server** (`b.tenant == ""`): routing key is built by extracting `tenant_slug` from the event payload. Subscribes with `.*` wildcard to receive from all tenants in one queue.

**Offline instance** (`b.tenant == slug`): routing key automatically appends the slug. Subscribes to exact tenant queue.

**Activation sequence** (offline):
1. App starts with in-memory bus (no slug yet).
2. User calls `POST /offline/ping` with JWT.
3. `ActivateRabbitMQ(slug)` creates a tenant-scoped RabbitMQ bus, subscribes catalog handlers, calls `Start()`.
4. `SetSyncBus(rb)` is called on catalog services so subsequent create/update/delete operations publish to RabbitMQ with the correct slug in the routing key.

### Shared Packages

- `shared/response/` — HTTP response builders: `response.OK`, `response.Created`, `response.BadRequest`, `response.Conflict`, etc.
- `shared/errors/` — sentinel domain errors used for HTTP status mapping.
- `shared/domain/vo/` — value objects: `Email`, `Document`, `DateTime`.
- `shared/events/` — `EventBus` interface and base event type.

### Migrations

SQL files live in `tenant/migrations/{public,tenant,offline}/`, embedded in the binary via Go's `embed` package, and auto-applied on startup. Use the `migrate` CLI to generate new migration files.

## Key Environment Variables

```
DATABASE_URL=postgresql://...     # omit for SQLite offline mode
PORT=8081
JWT_SECRET=...
URL_PROD=http://localhost:8081
USE_RABBITMQ=true                 # optional
RABBITMQ_HOST=...
RABBITMQ_DEFAULT_USER=...
RABBITMQ_DEFAULT_PASS=...
APP_ENV=dev                       # or "prod"
```

## Testing Conventions

- Test naming: `Test<Component>_<Method>_<Scenario>` (e.g., `TestService_Create_DuplicateSlug`).
- Repository tests use `go-sqlmock`; service tests use testify mocks of the repository interface.
- Each module's mock lives alongside its tests.
