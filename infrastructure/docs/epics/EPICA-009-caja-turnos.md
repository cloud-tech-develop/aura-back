# EPICA-009: Cash Drawer & Shifts

## 📌 General Information
- ID: EPICA-009
- State: Completed
- Priority: Medium
- Start Date: 2026-03-23
- Target Date: 2026-06-15
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Manage cash drawer operations, including shift management, cash movements, and cash drawer reconciliation (arqueo). This module ensures proper cash control and accountability for each user and branch.

**What problem does it solve?**
- No double shift opening
- Cash movement tracking
- Shift closing with authorization
- Cash reconciliation (arqueo)

**What value does it generate?**
- Cash accountability
- Fraud prevention
- Accurate end-of-day reporting
- Complete cash audit trail

---

## 👥 Stakeholders

- End User: Cashiers, store managers, accountants
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Cash Drawer module handles:

1. **Cash Drawer Setup**: Configure drawers per branch
2. **Shift Management**: Open and close shifts
3. **Cash Movements**: Entry/exit with reasons
4. **Arqueo (Reconciliation)**: Expected vs actual cash
5. **Shift Reports**: Summary by shift

---

## 📦 Scope

### Included:
- Cash drawer configuration per branch
- Shift open/close operations
- Cash movement recording
- Cash reconciliation (arqueo)
- Shift reports
- Movement history

### Not Included:
- Automatic cash counting
- Cash vault management
- Multi-currency handling

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-CASH-001 | Configure Cash Drawer | ✅ Implemented |
| HU-CASH-002 | Open Cash Shift | ✅ Implemented |
| HU-CASH-003 | Close Cash Shift | ✅ Implemented |
| HU-CASH-004 | Record Cash Movement | ✅ Implemented |
| HU-CASH-005 | Perform Cash Reconciliation | ✅ Implemented |
| HU-CASH-006 | View Shift Summary | ✅ Implemented |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- One active shift per user at a time
- Only one cash drawer per branch
- Shift close requires authorization (configurable role)
- Cash movements only within active shift
- Opening amount required for shift open
- Closing amount required for shift close
- All movements linked to user and shift

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: cash_drawer**
```sql
CREATE TABLE cash_drawer (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL UNIQUE REFERENCES branch(id),
    name VARCHAR(50) NOT NULL DEFAULT 'MAIN',
    is_active BOOLEAN DEFAULT TRUE,
    min_float DECIMAL(12,2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_cash_drawer_branch ON cash_drawer(branch_id);
```

**Table: cash_shift**
```sql
CREATE TABLE cash_shift (
    id BIGSERIAL PRIMARY KEY,
    cash_drawer_id BIGINT NOT NULL REFERENCES cash_drawer(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES branch(id),
    opening_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    closing_amount DECIMAL(12,2),
    expected_amount DECIMAL(12,2),
    difference DECIMAL(12,2),
    opening_notes TEXT,
    closing_notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED', 'AUDITED')),
    opened_at TIMESTAMP NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMP,
    closed_by BIGINT REFERENCES public.users(id),
    authorized_by BIGINT REFERENCES public.users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cash_shift_drawer ON cash_shift(cash_drawer_id);
CREATE INDEX idx_cash_shift_user ON cash_shift(user_id);
CREATE INDEX idx_cash_shift_status ON cash_shift(status);
CREATE INDEX idx_cash_shift_opened ON cash_shift(opened_at);
```

**Table: cash_movement**
```sql
CREATE TABLE cash_movement (
    id BIGSERIAL PRIMARY KEY,
    shift_id BIGINT NOT NULL REFERENCES cash_shift(id),
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('IN', 'OUT')),
    reason VARCHAR(30) NOT NULL CHECK (reason IN ('SALE', 'OPENING', 'CLOSING', 'EXPENSE', 'DROPS', 'WITHDRAWAL', 'ADJUSTMENT', 'REFUND')),
    amount DECIMAL(12,2) NOT NULL,
    reference_id BIGINT,
    reference_type VARCHAR(50),
    notes TEXT,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT cash_movement_amount_check CHECK (amount > 0)
);

CREATE INDEX idx_cash_movement_shift ON cash_movement(shift_id);
CREATE INDEX idx_cash_movement_type ON cash_movement(movement_type);
CREATE INDEX idx_cash_movement_created ON cash_movement(created_at);
```

---

## 📊 Success Metrics

- Zero double shift openings
- 100% shift closure rate
- Reconciliation accuracy > 99%
- Cash difference tolerance < 0.1%

---

## 🚧 Risks

- Concurrent shift operations
- Cash amount discrepancies
- Shift authorization bypass
- System crashes during open shift

---

## 📡 API Endpoints

### Cash Drawer
```
POST   /cash/drawer                  → Configurar caja (HU-CASH-001)
GET    /cash/drawer/:branchID        → Obtener caja por branch
```

### Cash Shift
```
POST   /cash/shift/open              → Abrir turno (HU-CASH-002)
POST   /cash/shift/:shiftID/close    → Cerrar turno (HU-CASH-003)
GET    /cash/shift/active            → Obtener turno activo
GET    /cash/shift/:shiftID          → Ver resumen (HU-CASH-006)
GET    /cash/shifts                  → Listar turnos
```

### Cash Movement
```
POST   /cash/movement                → Registrar movimiento (HU-CASH-004)
```

### Reconciliation
```
POST   /cash/shift/:shiftID/reconcile → Conciliar turno (HU-CASH-005)
```

### Query Parameters
- `page` - Número de página
- `limit` - Items por página
- `status` - Filtrar por estado: OPEN, CLOSED, AUDITED
- `start_date` - Fecha inicio
- `end_date` - Fecha fin

---

## 📁 Module Structure

```
modules/cash/
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
- **Módulo implementado**: cash
