# EPICA-007: Invoicing

## 📌 General Information
- ID: EPICA-007
- State: Completed
- Priority: High
- Start Date: 2026-03-23
- Target Date: 2026-06-30
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Manage invoice generation and lifecycle for Aura POS, including automatic invoice creation from sales, invoice numbering, fiscal documentation, and basic invoice management. This module provides the foundation for electronic invoicing compliance.

**What problem does it solve?**
- Automatic invoice generation from sales
- Sequential invoice numbering
- Complete fiscal documentation
- Basic invoice management (view, cancel)

**What value does it generate?**
- Tax compliance foundation
- Automated documentation
- Customer invoice delivery capability
- Financial audit trail

---

## 👥 Stakeholders

- End User: Cashiers, accountants, store managers
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Invoicing module handles:

1. **Invoice Generation**: Auto-create from completed sales
2. **Invoice Numbering**: Sequential prefixes per branch
3. **Fiscal Data**: Tax breakdown, totals, customer info
4. **Invoice States**: DRAFT, ISSUED, SENT, VIEWED, CANCELLED
5. **PDF Generation**: Invoice document creation

---

## 📦 Scope

### Included:
- Automatic invoice from sale
- Manual invoice creation
- Invoice prefixes and sequences per branch
- Invoice PDF generation
- Invoice cancellation (credit note reference)
- Invoice search and filtering
- Customer invoice history

### Not Included:
- DIAN electronic invoice integration (EPIC-003)
- XML generation for tax authority
- Advanced credit/debit notes
- Multi-currency invoices

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-INV-001 | Auto-generate Invoice from Sale | ✅ Completed |
| HU-INV-002 | Manual Invoice Creation | ✅ Completed |
| HU-INV-003 | View Invoice Details | ✅ Completed |
| HU-INV-004 | Generate Invoice PDF | ✅ Completed |
| HU-INV-005 | Cancel Invoice | ✅ Completed |
| HU-INV-006 | Invoice Search and Filter | ✅ Completed |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- Invoice numbers are sequential per branch prefix
- Invoice numbers cannot be reused after cancellation
- Cancelled invoices require credit note reference
- Invoice totals must match linked sale
- Invoices are immutable after issue (except cancellation)
- Soft delete only for audit compliance
- Each invoice requires customer information

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: invoice_prefix**
```sql
CREATE TABLE IF NOT EXISTS invoice_prefix (
    id BIGSERIAL PRIMARY KEY,
    prefix VARCHAR(10) NOT NULL,
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    enterprise_id BIGINT NOT NULL REFERENCES public.enterprises(id),
    current_number INTEGER NOT NULL DEFAULT 0,
    resolution_number VARCHAR(50),
    resolution_date DATE,
    valid_from DATE,
    valid_until DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);
```

**Table: invoice**
```sql
CREATE TABLE IF NOT EXISTS invoice (
    id BIGSERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) NOT NULL,
    prefix_id BIGINT NOT NULL REFERENCES invoice_prefix(id),
    invoice_type VARCHAR(20) NOT NULL DEFAULT 'SALE' CHECK (invoice_type IN ('SALE', 'CREDIT_NOTE', 'DEBIT_NOTE')),
    reference_id BIGINT,
    reference_type VARCHAR(50),
    sales_order_id BIGINT REFERENCES sales_order(id),
    customer_id BIGINT NOT NULL REFERENCES third_parties(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    enterprise_id BIGINT NOT NULL REFERENCES public.enterprises(id),
    invoice_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE,
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_exempt DECIMAL(12,2) NOT NULL DEFAULT 0,
    taxable_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    iva_19 DECIMAL(12,2) NOT NULL DEFAULT 0,
    iva_5 DECIMAL(12,2) NOT NULL DEFAULT 0,
    reteica DECIMAL(12,2) NOT NULL DEFAULT 0,
    retefuente DECIMAL(12,2) NOT NULL DEFAULT 0,
    reteica_rate DECIMAL(5,2) DEFAULT 0,
    retefuente_rate DECIMAL(5,2) DEFAULT 0,
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    payment_method VARCHAR(30) DEFAULT 'CASH',
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'ISSUED', 'SENT', 'VIEWED', 'CANCELLED')),
    cancelled_at TIMESTAMPTZ,
    cancelled_by BIGINT REFERENCES public.users(id),
    cancellation_reason TEXT,
    credit_note_id BIGINT REFERENCES invoice(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
```

**Table: invoice_item**
```sql
CREATE TABLE IF NOT EXISTS invoice_item (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoice(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    product_name VARCHAR(200) NOT NULL,
    product_sku VARCHAR(50),
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 19.00,
    tax_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    line_total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Table: invoice_log**
```sql
CREATE TABLE IF NOT EXISTS invoice_log (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoice(id),
    action VARCHAR(20) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    details TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 📊 Success Metrics

- Invoice generation time < 2 seconds
- Invoice number accuracy 100%
- Tax calculation accuracy 100%
- Zero duplicate invoice numbers

---

## 🚧 Risks

- Concurrent invoice generation causing number conflicts
- Complex tax calculation for different products
- Invoice PDF generation performance
- Large invoice history queries

---

## 📡 API Endpoints

### Invoices
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/invoices/generate` | Generate invoice from sales order |
| POST | `/invoices` | Create manual invoice |
| GET | `/invoices` | List invoices with filters |
| GET | `/invoices/:id` | Get invoice by ID |
| GET | `/invoices/number/:invoiceNumber` | Get invoice by number |
| POST | `/invoices/:id/issue` | Issue invoice (change status to ISSUED) |
| POST | `/invoices/:id/cancel` | Cancel invoice |
| GET | `/invoices/:id/logs` | Get invoice audit logs |

### Invoice Prefixes
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/invoice-prefixes` | Create invoice prefix |
| GET | `/invoice-prefixes` | List invoice prefixes |

### Query Parameters
- `page` - Page number
- `limit` - Items per page
- `status` - Filter by status (DRAFT, ISSUED, SENT, VIEWED, CANCELLED)
- `type` - Filter by invoice type (SALE, CREDIT_NOTE, DEBIT_NOTE)
- `branch_id` - Filter by branch
- `customer_id` - Filter by customer
- `start_date` - Filter by start date
- `end_date` - Filter by end date

---

## 📁 Module Structure

```
modules/invoices/
├── domain.go     # Entity, Repository & Service interfaces
├── service.go    # Repository & Service implementation
├── handler.go    # HTTP handlers
└── routes.go     # Route registration
```

### Migration Updated
- `tenant/migrations/tenant/000006_invoices.up.sql` - Enhanced with fiscal data, invoice types, audit logs

---

## Resumen

- **Total de HU**: 6
- **Completadas**: 6
- **Pendientes**: 0
- **Módulo implementado**: invoices (with full fiscal support)
