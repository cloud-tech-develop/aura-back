# EPICA-011: Shrinkage Management

## 📌 General Information
- ID: EPICA-011
- State: Completed
- Priority: Low
- Start Date: 2026-03-23
- Target Date: 2026-08-15
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Manage inventory shrinkage (mermas) due to damage, expiration, theft, or other losses. The system must track shrinkage reasons, reduce inventory accordingly, and provide reporting for loss analysis.

**What problem does it solve?**
- No tracking of inventory losses
- Unaccounted stock discrepancies
- Lack of shrinkage analysis
- No authorization for losses

**What value does it generate?**
- Accurate inventory records
- Loss analysis and prevention
- Financial loss tracking
- Proper authorization workflow

---

## 👥 Stakeholders

- End User: Store managers, inventory clerks, auditors
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Shrinkage module handles:

1. **Shrinkage Registration**: Record losses with reasons
2. **Reason Configuration**: Define valid shrinkage reasons
3. **Inventory Update**: Automatic stock reduction
4. **Authorization**: Require approval for high-value losses
5. **Reporting**: Shrinkage analysis and trends

---

## 📦 Scope

### Included:
- Shrinkage registration
- Reason configuration
- High-value authorization
- Inventory adjustment
- Shrinkage reports
- Trend analysis

### Not Included:
- Automatic shrinkage detection
- Theft prevention systems
- Insurance claim automation

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-SHR-001 | Register Shrinkage | ✅ Implemented |
| HU-SHR-002 | Configure Shrinkage Reasons | ✅ Implemented |
| HU-SHR-003 | Authorize High-Value Shrinkage | ✅ Implemented |
| HU-SHR-004 | View Shrinkage Report | ✅ Implemented |
| HU-SHR-005 | Cancel Shrinkage | ✅ Implemented |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- Shrinkage requires a reason
- Inventory is reduced upon shrinkage registration
- High-value shrinkage requires manager authorization
- Shrinkage values are tracked for reporting
- Cancellation restores inventory
- Shrinkage reasons cannot be deleted if used

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: shrinkage_reason**
```sql
CREATE TABLE shrinkage_reason (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    requires_authorization BOOLEAN DEFAULT FALSE,
    authorization_threshold DECIMAL(12,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Table: shrinkage**
```sql
CREATE TABLE shrinkage (
    id BIGSERIAL PRIMARY KEY,
    shrinkage_number VARCHAR(50) NOT NULL,
    branch_id BIGINT NOT NULL REFERENCES branch(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    reason_id BIGINT NOT NULL REFERENCES shrinkage_reason(id),
    shrinkage_date DATE NOT NULL,
    total_value DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED')),
    notes TEXT,
    authorized_by BIGINT REFERENCES public.users(id),
    authorized_at TIMESTAMP,
    cancellation_reason TEXT,
    cancelled_by BIGINT REFERENCES public.users(id),
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_shrinkage_number ON shrinkage(shrinkage_number);
CREATE INDEX idx_shrinkage_branch ON shrinkage(branch_id);
CREATE INDEX idx_shrinkage_reason ON shrinkage(reason_id);
CREATE INDEX idx_shrinkage_status ON shrinkage(status);
CREATE INDEX idx_shrinkage_date ON shrinkage(shrinkage_date);
```

**Table: shrinkage_item**
```sql
CREATE TABLE shrinkage_item (
    id BIGSERIAL PRIMARY KEY,
    shrinkage_id BIGINT NOT NULL REFERENCES shrinkage(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    batch_number VARCHAR(50),
    serial_number VARCHAR(100),
    quantity DECIMAL(10,2) NOT NULL,
    unit_cost DECIMAL(12,2) NOT NULL,
    total_value DECIMAL(12,2) NOT NULL,
    reason_detail TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shrinkage_item_shrinkage ON shrinkage_item(shrinkage_id);
```

---

## 📊 Success Metrics

- Shrinkage registration accuracy 100%
- Inventory update accuracy 100%
- Authorization compliance > 99%
- Shrinkage reporting time < 3 seconds

---

## 🚧 Risks

- Shrinkage abuse or fraud
- High-value shrinkage authorization bypass
- Batch product expiration tracking
- Shrinkage cancellation abuse

---

## 📡 API Endpoints

### Shrinkage
```
POST   /shrinkage                      → Registrar merma (HU-SHR-001)
GET    /shrinkage/:id                  → Obtener merma
GET    /shrinkage                      → Listar mermas
POST   /shrinkage/:id/authorize        → Autorizar merma (HU-SHR-003)
POST   /shrinkage/:id/cancel           → Cancelar merma (HU-SHR-005)
```

### Shrinkage Reasons
```
POST   /shrinkage/reasons              → Crear razón de merma (HU-SHR-002)
GET    /shrinkage/reasons              → Listar razones
```

### Reporting
```
GET    /shrinkage/report               → Reporte de mermas (HU-SHR-004)
```

### Query Parameters
- `page` - Número de página
- `limit` - Items por página
- `status` - Filtrar por estado: PENDING, APPROVED, REJECTED, CANCELLED
- `branch_id` - Filtrar por sucursal
- `active` - Filtrar activos (true/false)
- `start_date` - Fecha inicio
- `end_date` - Fecha fin

---

## 📁 Module Structure

```
modules/shrinkage/
├── domain.go     # Entity, Repository & Service interfaces
├── repository.go # Repository implementation
├── service.go    # Service implementation
├── handler.go    # HTTP handlers
└── routes.go     # Route registration
```

---

## Resumen

- **Total de HU**: 5
- **Completadas**: 5
- **Pendientes**: 0
- **Módulo implementado**: shrinkage
