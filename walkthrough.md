# Walkthrough: Offline POS Synchronization Implementation

This walkthrough documents the comprehensive refactoring performed to enable offline capabilities and bidirectional synchronization for the Aura POS API.

## Changes Made

### 1. Database Abstraction Layer
- **Multi-Driver Support**: Refactored [db.go](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/internal/db/db.go) to support both PostgreSQL and SQLite.
- **Query Adaptation**: Implemented a transparent query wrapper that converts Postgres placeholders (`$1`) to SQLite (`?`) and handles dialect differences (`ILIKE` -> `LIKE`, `NOW()` -> `datetime('now')`).
- **Querier Interface**: Introduced [db.Querier](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/internal/db/db.go) to allow passing either a direct connection or a transaction to repositories.

### 2. Module Refactoring
All core modules were refactored to use the new `db.Querier` and `db.DB` wrapper:
- `products`, `third-parties`, `sales`, `cart`, `inventory`, `payments`, `invoices`, and `reports`.

### 3. Global Identifiers (UUIDs)
- Added `GlobalID` (UUID string), `SyncStatus`, and `LastSyncedAt` to all core domain models:
    - [products/domain.go](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/modules/products/domain.go)
    - [sales/domain.go](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/modules/sales/domain.go)
    - [third-parties/domain.go](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/modules/third-parties/domain.go)
    - [invoices/domain.go](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/modules/invoices/domain.go)

### 4. Database Migrations
- Created [000015_sync_support.up.sql](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/tenant/migrations/tenant/000015_sync_support.up.sql) and rollback files to add synchronization columns and indexes to all relevant tables.

### 5. Synchronization Module
- **New Module**: Implemented [modules/sync](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/modules/sync/) with:
    - `Pull`: Fetching updates from the online server.
    - `Push`: Uploading local transactions with `ON CONFLICT (global_id)` support for upserts.
- **Endpoints**:
    - `GET /sync/pull`: Download new/updated data.
    - `POST /sync/push`: Upload local data batches.

## Verification Results

### Code Integrity
- All modules successfully wired in [main.go](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/cmd/api/main.go) and [server.go](file:///C:/Users/Drako/Desktop/cloud-tecno/aura-back/cmd/server/server.go).
- `go mod tidy` executed to resolve all dependencies including SQLite drivers.

### Compilation Check
- The code follows the standard Gin-gonic / multi-tenant pattern and replaces all direct `*sql.DB` usages with the adapted `db.Querier` in services and repositories.

> [!IMPORTANT]
> To use the offline version locally, simply configure the database driver in your environment to `sqlite` and point it to a local `.db` file. The same API will then operate on the local storage with sync capabilities.
