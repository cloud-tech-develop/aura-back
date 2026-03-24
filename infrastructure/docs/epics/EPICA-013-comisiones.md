# EPICA-013: Commissions

## 📌 General Information
- ID: EPICA-013
- State: Completed
- Priority: Low
- Start Date: 2026-03-23
- Target Date: 2026-09-15
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Manage sales commissions for employees, including commission configuration, automatic calculation per sale, commission settlement, and reporting. This module motivates sales performance and ensures accurate compensation.

**What problem does it solve?**
- Manual commission calculation
- No commission visibility
- Payment tracking issues
- Commission disputes

**What value does it generate?**
- Automatic commission calculation
- Transparent commission records
- Accurate compensation
- Performance motivation

---

## 👥 Stakeholders

- End User: Store managers, accountants, sales employees
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Commissions module handles:

1. **Configuration**: Set commission rules per employee/product
2. **Calculation**: Auto-calculate on sale completion
3. **Tracking**: View commission history
4. **Settlement**: Mark commissions as paid
5. **Reporting**: Commission reports and summaries

---

## 📦 Scope

### Included:
- Commission rules configuration
- Commission calculation per sale
- Commission history
- Settlement (mark as paid)
- Commission reports

### Not Included:
- Complex tiered commissions
- Team commissions
- Commission advances
- Integration with payroll systems

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-COMM-001 | Configure Commission Rules | ✅ Implemented |
| HU-COMM-002 | Calculate Commissions on Sale | ✅ Implemented |
| HU-COMM-003 | View Commission History | ✅ Implemented |
| HU-COMM-004 | Settle Commissions | ✅ Implemented |
| HU-COMM-005 | Commission Reports | ✅ Implemented |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- Commission types: PERCENTAGE_SALE, PERCENTAGE_MARGIN, FIXED_AMOUNT
- Commissions are calculated at sale completion
- Commission can be per employee or per product
- Settlement marks commission as paid
- Commission history is immutable
- Deleted sales do not affect settled commissions

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: commission_rule**
```sql
CREATE TABLE commission_rule (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    commission_type VARCHAR(30) NOT NULL CHECK (commission_type IN ('PERCENTAGE_SALE', 'PERCENTAGE_MARGIN', 'FIXED_AMOUNT')),
    employee_id BIGINT REFERENCES third_party(id),
    product_id BIGINT REFERENCES product(id),
    category_id BIGINT REFERENCES category(id),
    value DECIMAL(12,2) NOT NULL,
    min_sale_amount DECIMAL(12,2) DEFAULT 0,
    start_date DATE,
    end_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_commission_rule_employee ON commission_rule(employee_id) WHERE employee_id IS NOT NULL;
CREATE INDEX idx_commission_rule_product ON commission_rule(product_id) WHERE product_id IS NOT NULL;
CREATE INDEX idx_commission_rule_category ON commission_rule(category_id) WHERE category_id IS NOT NULL;
CREATE INDEX idx_commission_rule_active ON commission_rule(is_active, start_date, end_date);
```

**Table: commission**
```sql
CREATE TABLE commission (
    id BIGSERIAL PRIMARY KEY,
    sales_order_id BIGINT NOT NULL REFERENCES sales_order(id),
    employee_id BIGINT NOT NULL REFERENCES third_party(id),
    branch_id BIGINT NOT NULL REFERENCES branch(id),
    rule_id BIGINT REFERENCES commission_rule(id),
    sale_amount DECIMAL(12,2) NOT NULL,
    profit_margin DECIMAL(12,2),
    commission_type VARCHAR(30) NOT NULL,
    commission_rate DECIMAL(12,2) NOT NULL,
    commission_amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'SETTLED', 'CANCELLED')),
    settled_at TIMESTAMP,
    settled_by BIGINT REFERENCES public.users(id),
    settlement_period VARCHAR(20),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_commission_employee ON commission(employee_id);
CREATE INDEX idx_commission_order ON commission(sales_order_id);
CREATE INDEX idx_commission_status ON commission(status);
CREATE INDEX idx_commission_created ON commission(created_at);
CREATE INDEX idx_commission_period ON commission(settlement_period) WHERE settlement_period IS NOT NULL;
```

---

## 📊 Success Metrics

- Commission calculation accuracy 100%
- Calculation time < 100ms per sale
- Settlement accuracy 100%
- Zero commission disputes

---

## 🚧 Risks

- Commission rule priority conflicts
- Sale cancellation affecting pending commissions
- Large commission history performance
- Complex margin calculations

---

## 📡 API Endpoints

### Commission Rules
```
POST   /commissions/rules              → Crear regla de comisión (HU-COMM-001)
PUT    /commissions/rules/:id          → Actualizar regla
DELETE /commissions/rules/:id          → Eliminar regla
GET    /commissions/rules              → Listar reglas
```

### Commission Calculation
```
POST   /commissions/calculate          → Calcular comisiones (HU-COMM-002)
```

### Commission History
```
GET    /commissions                    → Listar historial de comisiones (HU-COMM-003)
GET    /commissions/:id                → Obtener comisión
```

### Commission Settlement
```
POST   /commissions/settle             → Liquidar comisiones (HU-COMM-004)
```

### Reports
```
GET    /commissions/report/summary     → Resumen de comisiones por empleado (HU-COMM-005)
GET    /commissions/report/totals      → Totales de comisiones
```

### Query Parameters
- `page` - Número de página
- `limit` - Items por página
- `status` - Filtrar por estado: PENDING, SETTLED, CANCELLED
- `employee_id` - Filtrar por empleado
- `active` - Filtrar activos (true/false)
- `start_date` - Fecha inicio
- `end_date` - Fecha fin

---

## 📁 Module Structure

```
modules/commissions/
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
- **Módulo implementado**: commissions
