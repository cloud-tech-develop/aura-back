# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run development server (SQLite offline mode if DATABASE_URL is unset)
go run ./cmd/api/main.go

# Build
go build ./cmd/api/main.go

# Run all tests
go test ./...

# Run tests for a specific module
go test -v ./modules/catalog/brands/...

# Run a single test
go test -v -run TestService_Create_ValidSlugFormat ./modules/catalog/products/...

# Format and vet
go fmt ./...
go vet ./...

# Create a new PostgreSQL migration
migrate create -ext sql -dir tenant/migrations/public -seq migration_name
migrate create -ext sql -dir tenant/migrations/tenant -seq migration_name

# Docker Compose (Postgres + RabbitMQ)
docker-compose up --build
```

## Architecture Overview

Multi-tenant POS (Point of Sale) backend with dual Postgres/SQLite support for online↔offline sync via RabbitMQ.

### Entry Points

- `cmd/api/main.go` — loads `.env`, initializes DB + migrations, wires all modules (DI), generates the offline binary (`static/bin/aura-pos-offline.exe`) at startup in Postgres mode. Uses build-time ldflags to bake-in RabbitMQ/APP_ENV defaults into the offline binary.
- `cmd/server/server.go` — sets up Gin router with CORS, registers `POST /login`, static files, health check, and mounts all modules via `Register(public, protected, handler)`.

### Database Strategy

**Multi-tenant via Postgres schemas (schema-per-tenant pattern):**
- `public` schema: `enterprises`, `users`, `roles`, `user_roles`, `plans`.
- Each enterprise gets its own schema (e.g., `empresa_uno`) with tenant-specific tables.
- Tenant queries use `db.WithSchema(q, slug)` which sets `SET search_path TO <slug>` per-connection.

**Dual-mode (Postgres / SQLite):**
- `DATABASE_URL` set → Postgres mode (default port 8081)
- `DATABASE_URL` unset → SQLite offline mode (`aura_pos.db`, port 8091)
- `internal/db/db.go` provides a `Querier` interface that transparently adapts Postgres queries to SQLite: `$N` → `?`, `ILIKE` → `LIKE`, `NOW()` → `CURRENT_TIMESTAMP`, `TIMESTAMPTZ` → `TEXT`, `BOOLEAN` → `INTEGER`, schema prefixes removed.
- SQLite uses `modernc.org/sqlite` (pure Go, no CGO). Connection pool limited to `MaxOpenConns=1` with `PRAGMA busy_timeout = 5000`.
- Test repos use `db.NewMock(conn)` to pass `*sql.DB` from `go-sqlmock`.

### Migration System

Three migration directories, embedded via Go `embed.FS` in `tenant/manager.go`:

```
tenant/migrations/
├── public/          # Online (Postgres) .up.sql / .down.sql — shared tables
├── tenant/          # Online (Postgres) .up.sql / .down.sql — per-tenant tables
└── offline/
    ├── public/      # Offline (SQLite) .sql — native SQLite DDL for shared tables
    └── tenant/      # Offline (SQLite) .sql — native SQLite DDL for tenant tables
```

- **Online**: `golang-migrate/v4` compatible `.up.sql` / `.down.sql` files, version-tracked in `schema_migrations` table. Runs via `manager.MigratePublic()` at startup, `manager.MigrateAll()` in background for existing tenants.
- **Offline**: Plain `.sql` files per version, tracked in per-directory `schema_offline_*` tables. Runs via `manager.MigrateOffline()` at startup. No `down.sql` — SQLite migrations are forward-only.

### Module Pattern

Every feature module follows the same layout:

```
modules/<domain>/<feature>/
  domain.go       — entity structs + Repository/Service interfaces + event constants
  service.go      — business logic (unexported struct, returns interface)
  repository.go   — DB queries via db.Querier (compatible with both drivers)
  handler.go      — HTTP request/response types + Gin handler methods
  routes.go       — Register(public, protected, handler) function
  logger.go       — subscribes to domain events for audit logging
  sync.go         — sync payload types + event constructors for online↔offline sync
  sync_handler.go — handles incoming sync events from RabbitMQ (for sync-only modules)
  *_test.go       — unit tests with testify mocks and go-sqlmock
```

New modules are wired in `cmd/api/main.go` and registered in `cmd/server/server.go`.

**Two sync handler patterns exist for catalog modules:**

| Module | Sync mechanism | Event handler lives on |
|--------|---------------|----------------------|
| products | Service implements `events.EventHandler` directly | `service.go` — `Handle()` method |
| categories | Service implements `events.EventHandler` directly | `service.go` — `Handle()` method |
| brands | Service implements `events.EventHandler` directly | `service.go` — `Handle()` method |
| presentations | Separate `SyncHandler` struct | `sync_handler.go` — standalone handler |

Products, categories, and brands use the dual-role approach where the service handles both business logic and sync. Presentations use a separate `SyncHandler` struct that takes only a `Repository`.

Modules are grouped under `modules/catalog/` but each sub-module (products, brands, categories, units, presentations) is an independent Go package with its own routes, registered separately in `cmd/server/server.go`.

### Authentication & Middleware

Request flow: `CORS → AuthMiddleware (JWT + IP validation) → Tenant resolution (slug from JWT) → Handler`

JWT claims: `user_id`, `enterprise_id`, `tenant_id`, `slug`, `email`, `roles`, `role_level`, `ip`. Token IP must match client IP on every request via `X-Forwarded-For` / `X-Real-IP` / `RemoteAddr`. Role levels: 0=SUPERADMIN, 1=ADMIN, 2=SUPERVISOR, 3+=USER.

Public routes: `POST /login`, `POST /enterprises`, `GET /enterprises/:slug`.

Login flow: authenticates against `public.users` + `public.enterprises`, fetches roles from `public.roles` joined via `public.user_roles`, resolves role level (minimum level number = highest privilege), fetches third-party name from tenant schema, returns JWT + user + enterprise info.

### Event Bus Infrastructure

Two implementations of `shared/events.EventBus`:

**Memory (in-memory bus):**
- `infrastructure/messaging/memory/bus.go` — goroutine workers consuming from a buffered channel. Default: buffer 100, 5 workers. Used for audit logging and as the initial bus in offline mode before RabbitMQ activation. Subscriptions match by exact event name.

**RabbitMQ (cross-server bus):**
- `infrastructure/messaging/rabbit/bus.go` — topic exchange `aura.<APP_ENV>`, persistent delivery, JSON-encoded payloads.
- Routing keys: `{eventName}.{tenantSlug}` (e.g., `product.offline.created.empresa_uno`).
- Supports wildcard subscriptions via `.*` suffix for receiving all tenants.
- Tenant-scoped bus (`NewRabbitMQEventBusWithTenant(slug)`) auto-appends slug to routing keys.
- Validates APP_ENV to prevent malformed exchange/queue names.

### Online↔Offline Sync Architecture

**Offline activation sequence:**
1. SQLite mode starts with in-memory bus (no tenant slug yet).
2. `POST /offline/ping` fetches enterprise slug from JWT (or falls back to local SQLite).
3. **Bulk sync**: HTTP GET requests pull all data from the online server (`URL_PROD`) in parallel goroutines: enterprises, plans, users, user_roles, third_parties, categories, brands, units, products, presentations.
4. **RabbitMQ activation**: `ActivateRabbitMQ(slug)` creates a tenant-scoped RabbitMQ bus, subscribes catalog service handlers (products, categories, brands) and the presentations `SyncHandler`, then calls `Start()`.
5. `SetSyncBus(rb)` is called on catalog services so subsequent creates/updates/deletes publish to RabbitMQ with the correct slug routing key.

**Sync event routing key convention:** `{entity}.{direction}.{action}.{slug}`

- Online → Offline: `product.offline.created.empresa_uno` → queue `aura.dev.product.offline.created.empresa_uno`
- Offline → Online: `product.online.created.empresa_uno` → queue `aura.dev.product.online.created` (wildcard binding `product.online.created.*` catches all tenants)

**Conflict resolution:** Timestamp-based. If the online record has a newer `updated_at` than the offline event's timestamp, the online version wins. Conflicts are logged and silently resolved in favor of online.

### Shared Packages

- `shared/response/` — HTTP response helpers: `response.OK`, `response.Created`, `response.BadRequest`, `response.Conflict`, etc.
- `shared/errors/` — sentinel domain errors mapped to HTTP status codes.
- `shared/domain/vo/` — value objects: `Email` (with validation), `Document`, `DateTime`.
- `shared/events/` — `Event`, `EventBus`, and `EventHandler` interfaces; `BaseEvent` struct with `name`, `payload`, `timestamp`.
- `shared/domain/` — `PageResult` for paginated queries, `Common` struct for shared entity fields.
- `shared/logging/` — `LoggerHandler` that writes structured log files to a directory.

### Key Dependencies

| Package | Purpose |
|---------|---------|
| `gin-gonic/gin` v1.12 | HTTP framework |
| `golang-jwt/jwt/v5` | JWT auth with IP binding |
| `golang-migrate/v4` | Postgres migrations |
| `modernc.org/sqlite` | Pure-Go SQLite driver (no CGO) |
| `rabbitmq/amqp091-go` | RabbitMQ client for cross-server sync |
| `stretchr/testify` | Assertions + mocks |
| `DATA-DOG/go-sqlmock` | SQL mock for repository tests |
| `joho/godotenv` | .env loading |
| `lib/pq` | Postgres driver |

## Key Environment Variables

```
DATABASE_URL=postgresql://...      # omit for SQLite offline mode
DATABASE_DRIVER=postgres|sqlite    # auto-detected from DATABASE_URL
PORT=8081                          # 8091 in SQLite mode
JWT_SECRET=...
URL_PROD=http://localhost:8081     # online server URL (offline pulls data from here)
USE_RABBITMQ=true                  # enable cross-server sync
RABBITMQ_HOST=...
RABBITMQ_DEFAULT_USER=...
RABBITMQ_DEFAULT_PASS=...
APP_ENV=dev|prod                   # scopes RabbitMQ exchange name
```

## Testing Conventions

- Test naming: `Test<Component>_<Method>_<Scenario>` (e.g., `TestService_Create_DuplicateSlug`).
- Repository tests use `go-sqlmock` with `db.NewMock(mockDB)` to wrap the mock connection.
- Service tests use testify mocks of the repository interface (mock structs defined alongside tests in each module).
- Each module's mock lives in its own test files.

## Creating Repositories (DB Queries)

### Postgres vs SQLite differences to handle

Every `Repository` implementation receives `db.Querier` and must use ONLY `QueryContext`, `QueryRowContext`, `ExecContext` methods — never the raw `*sql.DB` methods. The `Querier` interface transparently adapts queries for the active driver.

**Schema prefix handling:**
```go
// Always use SchemaPrefix() when referencing tenant tables in queries
prefix := q.SchemaPrefix(tenantSlug)
query := fmt.Sprintf(`SELECT * FROM %sproducts WHERE id = ?`, prefix)
```
- Postgres: produces `"empresa_uno".products`
- SQLite: produces `` `products` `` (empty prefix, tables are flat)

**Placeholder style:** Always use `$1, $2, $3...` in queries. The `AdaptQuery` function converts them to `?` for SQLite. Never use `?` directly in Postgres queries — `lib/pq` does not support it.

**Type mapping reference:**
| Postgres type | SQLite equivalent | Notes |
|---|---|---|
| `BIGSERIAL` | `INTEGER PRIMARY KEY AUTOINCREMENT` | Adapted automatically |
| `TIMESTAMPTZ` | `TEXT` | `vo.DateTime` handles parsing |
| `BOOLEAN` | `INTEGER` (0/1) | Adapted automatically |
| `NOW()` | `CURRENT_TIMESTAMP` | Adapted automatically |
| `ILIKE` | `LIKE` | SQLite LIKE is case-insensitive by default for ASCII |
| `RETURNING id` | _(removed)_ | Use `LastInsertId()` pattern instead or query after insert |
| Schema prefix `"slug".table` | _(removed)_ | Tables in same file |

**Transaction safety:** For SQLite, `MaxOpenConns=1` prevents `database is locked` errors. When using transactions within a tenant context, pass the `Querier` from `db.WithSchema()`.

### Repository test pattern with go-sqlmock

```go
func TestRepository_Create_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    q := db.NewMock(db)
    repo := NewRepository(q)

    mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO ...`)).
        WithArgs(...).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    err := repo.Create(ctx, "test_slug", entity)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

## Implementing Online↔Offline Sync

### Module sync file structure

Every catalog module that needs real-time sync has these additional files:

```
modules/catalog/<feature>/
  sync.go          — Sync payload struct + event constructors (both directions)
  sync_handler.go  — Batch sync handler (for modules with separate sync handler) — OR —
  service.go       — Handle() method directly on the service (products, categories, brands)
```

### Sync event naming convention

Events follow the pattern: `{entity}.{direction}.{action}`

| Event constant | Direction | Example routing key |
|---|---|---|
| `EventProductOfflineCreated` | Online → Offline | `product.offline.created.empresa_uno` |
| `EventProductOnlineCreated` | Offline → Online | `product.online.created.empresa_uno` |

Each direction gets its own event constructors:
```go
NewSyncCreatedEvent(slug, entity)           → EventProductOfflineCreated (source:"online")
NewSyncCreatedEventFromOffline(slug, entity) → EventProductOnlineCreated  (source:"offline")
```

### Service sync wiring (the products reference pattern)

The products module is the reference implementation. Every catalog service that supports sync must:

1. **Store two buses** — `eventBus` (local/initial RabbitMQ) and `syncBus` (set after offline activation):
```go
type service struct {
    eventBus  events.EventBus // Local + online RabbitMQ
    syncBus   events.EventBus // Cross-server sync (nil until SetSyncBus)
    syncMu    sync.RWMutex
    isOffline bool
}
```

2. **Single `publishSync()` with fallback** — never separate publishSyncNew/Update/Delete + publishOnline:
```go
func (s *service) publishSync(event events.Event) {
    s.syncMu.RLock()
    bus := s.syncBus
    s.syncMu.RUnlock()
    if bus == nil {
        bus = s.eventBus // fallback — online mode uses eventBus directly
    }
    if bus == nil { return }
    if err := bus.Publish(event); err != nil {
        s.logger.Logf("warn: sync publish failed: %v", err)
    }
}
```

3. **Publish the correct event direction on CUD operations:**
```go
func (s *service) Create(ctx, slug string, entity *Entity) error {
    // ... validation + repo.Create ...
    if s.isOffline {
        s.publishSync(NewSyncCreatedEventFromOffline(slug, entity))
    } else {
        s.publishSync(NewSyncCreatedEvent(slug, entity))
    }
    return nil
}
```

4. **Timestamp-based conflict resolution in remote handlers:**
```go
func (s *service) handleRemoteCreate(ctx, slug string, payload map[string]interface{}) error {
    entity := entityFromPayload(payload)
    existing, err := s.repo.GetByID(ctx, slug, entity.ID)
    if err == nil {
        var eventTime time.Time
        if tStr, ok := payload["timestamp"].(string); ok {
            eventTime, _ = time.Parse(time.RFC3339, tStr)
        }
        if existing.UpdatedAt != nil && time.Time(*existing.UpdatedAt).After(eventTime) {
            return nil // local is newer, skip
        }
        entity.CreatedAt = existing.CreatedAt
        return s.repo.Update(ctx, slug, entity)
    }
    if err != sql.ErrNoRows { return err }
    return s.repo.Create(ctx, slug, entity)
}
```

5. **Implement `events.EventHandler` interface** — the service's `Handle()` dispatches inbound RabbitMQ events by event name:
```go
func (s *service) Handle(event events.Event) error {
    payload := event.GetPayload().(map[string]interface{})
    slug := payload["tenant_slug"].(string)
    switch event.GetName() {
    case EventOfflineCreated, EventOnlineCreated: return s.handleRemoteCreate(ctx, slug, payload)
    case EventOfflineUpdated, EventOnlineUpdated: return s.handleRemoteUpdate(ctx, slug, payload)
    case EventOfflineDeleted, EventOnlineDeleted: return s.handleRemoteDelete(ctx, slug, payload)
    }
    return nil
}
```

6. **Online subscribes with wildcard**, offline defers until `/offline/ping`:
```go
func (s *service) subscribeToRabbitMQEvents() {
    if s.eventBus == nil { return }
    if s.isOffline {
        s.logger.Log("Offline mode: sync deferred until /offline/ping")
        return
    }
    s.eventBus.Subscribe(EventOnlineCreated+".*", s) // wildcard = all tenants
}
```

7. **SyncBus setter** for offline activation:
```go
func (s *service) SetSyncBus(bus events.EventBus) {
    s.syncMu.Lock(); defer s.syncMu.Unlock()
    s.syncBus = bus
}
```

### Adding sync to a new module checklist

1. Create `sync.go` with `SyncPayload`, event constants, `ToSyncPayload()`, `NewSyncCreatedEvent[FromOffline]()`, `NewSyncUpdatedEvent[FromOffline]()`, `NewSyncDeletedEvent[FromOffline]()`
2. If the module uses a separate sync handler: create `sync_handler.go` with `SyncHandler` struct that implements `events.EventHandler` with batch sync methods (`*SyncFromOffline`) and timestamp conflict resolution
3. Add `isOffline`, `syncBus`, `syncMu` fields to the service struct
4. Add `SetSyncBus()`, `Handle()`, `subscribeToRabbitMQEvents()`, `publishSync()` methods
5. Add `handleRemoteCreate/Update/Delete` with timestamp conflict resolution
6. Add `entityFromPayload()` helper to extract entity from RabbitMQ map payload
7. Wire Create/Update/Delete to call `publishSync()` with correct event direction
8. Register the service's `Handle()` or `SyncHandler` in `offline.Service.ActivateRabbitMQ()`
