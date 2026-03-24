# EPICA-004: Inventory Management

## 📌 General Information
- ID: EPICA-004
- State: Completed
- Priority: High
- Start Date: 2026-03-23
- Target Date: 2026-04-30
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Implement comprehensive inventory control per branch and tenant. The system must track stock movements, maintain Kardex (stock ledger), provide low stock alerts, and ensure inventory consistency across all sales, purchases, and shrinkage operations.

**What problem does it solve?**
- Real-time inventory tracking per branch
- Prevention of overselling
- Complete traceability of stock movements
- Low stock alerting

**What value does it generate?**
- Accurate stock information for decision making
- Reduced losses from stockouts or overselling
- Complete audit trail of inventory changes

---

## 👥 Stakeholders

- End User: Store managers, inventory clerks, cashiers
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Inventory module provides complete stock management:

1. **Stock Control**: Track quantity per product per branch
2. **Kardex**: Complete movement history for each product
3. **Alerts**: Low stock notifications
4. **Movements**: Entry and exit tracking with reasons
5. **Consolidation**: Multi-branch inventory view

---

## 📦 Scope

### Included:
- Stock inquiry by branch
- Stock movement registration (entry/exit)
- Kardex per product with movement details
- Low stock alerts configuration
- Inventory consolidation across branches
- Movement reasons management
- Batch tracking for products with expiration dates
- Serial number tracking for serialized products

### Not Included:
- Advanced demand forecasting
- Automatic reorder point calculations
- Barcode printing
- RFID integration

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-INV-001 | Query Branch Inventory | ✅ Completed |
| HU-INV-002 | Register Inventory Movement | ✅ Completed |
| HU-INV-003 | View Product Kardex | ✅ Completed |
| HU-INV-004 | Low Stock Alerts | ✅ Completed |
| HU-INV-005 | Inventory Consolidation | ✅ Completed |
| HU-INV-006 | Batch Tracking | ✅ Completed |
| HU-INV-007 | Serialized Product Tracking | ✅ Completed |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- All movements are immutable (cannot be modified or deleted)
- Stock cannot be negative
- Every movement requires: reason, user, quantity, previous and new balance
- Movements are linked to: sales, purchases, shrinkage, transfers
- Branches belong to a tenant
- Batch products require expiration date tracking
- Serialized products require unique serial number per unit

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: inventory**
```sql
CREATE TABLE IF NOT EXISTS inventory (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES product(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER DEFAULT 0,
    min_stock INTEGER DEFAULT 0,
    max_stock INTEGER,
    location VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    
    CONSTRAINT inventory_product_branch_unique UNIQUE (product_id, branch_id),
    CONSTRAINT inventory_quantity_check CHECK (quantity >= 0),
    CONSTRAINT inventory_reserved_check CHECK (reserved_quantity >= 0)
);
```

**Table: inventory_movement**
```sql
CREATE TABLE IF NOT EXISTS inventory_movement (
    id BIGSERIAL PRIMARY KEY,
    inventory_id BIGINT NOT NULL REFERENCES inventory(id),
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('ENTRY', 'EXIT', 'ADJUSTMENT')),
    movement_reason VARCHAR(30) NOT NULL CHECK (movement_reason IN ('SALE', 'PURCHASE', 'SHRINKAGE', 'TRANSFER_IN', 'TRANSFER_OUT', 'ADJUSTMENT', 'RETURN', 'INITIAL', 'DAMAGE', 'THEFT', 'EXPIRED')),
    quantity INTEGER NOT NULL,
    previous_balance INTEGER NOT NULL,
    new_balance INTEGER NOT NULL,
    reference_id BIGINT,
    reference_type VARCHAR(50),
    batch_number VARCHAR(50),
    serial_number VARCHAR(100),
    expiration_date DATE,
    notes TEXT,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT inventory_movement_quantity_check CHECK (quantity > 0),
    CONSTRAINT inventory_movement_balance_check CHECK (new_balance >= 0)
);
```

**Table: movement_reason**
```sql
CREATE TABLE IF NOT EXISTS movement_reason (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('ENTRY', 'EXIT', 'ADJUSTMENT')),
    requires_authorization BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 📊 Success Metrics

- Inventory accuracy > 99%
- Movement processing time < 100ms
- Low stock alert delivery < 1 minute
- Zero negative stock occurrences

---

## 🚧 Risks

- Concurrent stock updates causing race conditions
- Complex batch and serial tracking complexity
- Multi-branch inventory consistency
- Performance with large movement history

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/inventory` | List inventory with filters |
| GET | `/inventory/low-stock` | Get low stock items |
| GET | `/inventory/:productId/:branchId` | Get inventory for product/branch |
| GET | `/inventory/product/:productId` | Get inventory by product (all branches) |
| GET | `/inventory/kardex/:productId/:branchId` | Get product Kardex |
| POST | `/inventory/movements` | Register stock movement |
| GET | `/movements` | List movements |
| GET | `/movements/:id` | Get movement by ID |
| GET | `/movement-reasons` | List movement reasons |

### Query Parameters
- `page` - Page number
- `limit` - Items per page
- `branch_id` - Filter by branch
- `product_id` - Filter by product
- `low_stock` - Filter low stock items
- `movement_type` - Filter by type (ENTRY, EXIT, ADJUSTMENT)
- `reason` - Filter by reason

---

## 📁 Module Structure

```
modules/inventory/
├── domain.go     # Entity, Repository & Service interfaces
├── service.go    # Repository & Service implementation
├── handler.go    # HTTP handlers
└── routes.go     # Route registration
```

### Migration Created
- `tenant/migrations/tenant/000007_inventory.up.sql`
- `tenant/migrations/tenant/000007_inventory.down.sql`

---

## Resumen

- **Total de HU**: 7
- **Completadas**: 7
- **Pendientes**: 0
- **Módulo implementado**: inventory
