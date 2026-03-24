# EPICA-012: Inventory Transfers

## 📌 General Information
- ID: EPICA-012
- State: Completed
- Priority: Low
- Start Date: 2026-03-23
- Target Date: 2026-08-30
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Manage inventory transfers between branches of the same company. The system must support transfer requests, shipment, receipt, and cancellation with proper inventory adjustments at both origin and destination.

**What problem does it solve?**
- Manual transfer tracking
- No inventory visibility between branches
- Transfer authorization
- Incomplete transfer records

**What value does it generate?**
- Centralized inventory management
- Proper stock redistribution
- Transfer audit trail
- Reduced stockouts at branches

---

## 👥 Stakeholders

- End User: Store managers, inventory clerks
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Transfers module handles:

1. **Transfer Request**: Create transfer from source branch
2. **Transfer Approval**: Authorize transfers
3. **Shipment**: Dispatch from origin
4. **Receipt**: Accept at destination
5. **Cancellation**: Cancel pending transfers

---

## 📦 Scope

### Included:
- Transfer request creation
- Stock validation
- Shipment processing
- Receipt confirmation
- Transfer cancellation
- Transfer history

### Not Included:
- Inter-company transfers
- Third-party logistics integration
- Automated routing

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-TRANS-001 | Create Transfer Request | ✅ Implemented |
| HU-TRANS-002 | Approve Transfer | ✅ Implemented |
| HU-TRANS-003 | Ship Transfer | ✅ Implemented |
| HU-TRANS-004 | Receive Transfer | ✅ Implemented |
| HU-TRANS-005 | Cancel Transfer | ✅ Implemented |
| HU-TRANS-006 | View Transfer History | ✅ Implemented |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- Transfers only between branches of same company
- Stock is deducted at origin upon shipment
- Stock is added at destination upon receipt
- Partial receipts are supported
- Cancellation restores stock at origin (if shipped)
- Pending transfers can be cancelled

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: transfer**
```sql
CREATE TABLE transfer (
    id BIGSERIAL PRIMARY KEY,
    transfer_number VARCHAR(50) NOT NULL,
    origin_branch_id BIGINT NOT NULL REFERENCES branch(id),
    destination_branch_id BIGINT NOT NULL REFERENCES branch(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'SHIPPED', 'PARTIAL', 'RECEIVED', 'CANCELLED')),
    requested_date TIMESTAMP NOT NULL DEFAULT NOW(),
    shipped_date TIMESTAMP,
    received_date TIMESTAMP,
    notes TEXT,
    shipped_by BIGINT REFERENCES public.users(id),
    received_by BIGINT REFERENCES public.users(id),
    cancellation_reason TEXT,
    cancelled_by BIGINT REFERENCES public.users(id),
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_transfer_number ON transfer(transfer_number);
CREATE INDEX idx_transfer_origin ON transfer(origin_branch_id);
CREATE INDEX idx_transfer_destination ON transfer(destination_branch_id);
CREATE INDEX idx_transfer_status ON transfer(status);
CREATE INDEX idx_transfer_dates ON transfer(requested_date);
```

**Table: transfer_item**
```sql
CREATE TABLE transfer_item (
    id BIGSERIAL PRIMARY KEY,
    transfer_id BIGINT NOT NULL REFERENCES transfer(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    requested_quantity DECIMAL(10,2) NOT NULL,
    shipped_quantity DECIMAL(10,2) NOT NULL DEFAULT 0,
    received_quantity DECIMAL(10,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT transfer_item_quantity_check CHECK (requested_quantity > 0)
);

CREATE INDEX idx_transfer_item_transfer ON transfer_item(transfer_id);
CREATE INDEX idx_transfer_item_product ON transfer_item(product_id);
```

---

## 📊 Success Metrics

- Transfer processing time < 1 second
- Inventory accuracy after transfers 100%
- Transfer completion rate > 95%
- Zero unauthorized transfers

---

## 🚧 Risks

- Concurrent transfer requests
- Stock validation race conditions
- Partial receipt tracking complexity
- Cross-branch inventory consistency

---

## 📡 API Endpoints

### Transfers
```
POST   /transfers                      → Crear traslado (HU-TRANS-001)
GET    /transfers/:id                  → Obtener traslado
GET    /transfers                      → Listar historial de traslados (HU-TRANS-006)
POST   /transfers/:id/approve          → Aprobar traslado (HU-TRANS-002)
POST   /transfers/:id/ship             → Enviar traslado (HU-TRANS-003)
POST   /transfers/:id/receive          → Recibir traslado (HU-TRANS-004)
POST   /transfers/:id/cancel           → Cancelar traslado (HU-TRANS-005)
```

### Query Parameters
- `page` - Número de página
- `limit` - Items por página
- `status` - Filtrar por estado: PENDING, APPROVED, SHIPPED, PARTIAL, RECEIVED, CANCELLED
- `origin_branch_id` - Filtrar por sucursal origen
- `destination_branch_id` - Filtrar por sucursal destino
- `start_date` - Fecha inicio
- `end_date` - Fecha fin

---

## 📁 Module Structure

```
modules/transfers/
├── domain.go     # Entity, Repository & Service interfaces
├── repository.go # Repository implementation
├── service.go    # Service implementation
├── handler.go    # HTTP handlers
└── routes.go     # Route registration
```

---

## Resumen

- **Total de HU**: 6
- **Completadas**: 6
- **Pendientes**: 0
- **Módulo implementado**: transfers
