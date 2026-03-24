# EPICA-006: Payment Processing

## 📌 General Information
- ID: EPICA-006
- State: Completed
- Priority: High
- Start Date: 2026-03-23
- Target Date: 2026-05-30
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Handle all payment processing for sales, including multiple payment methods, change calculation, payment splitting, and transaction history. This module integrates with the Sales and Invoicing modules to complete the payment lifecycle.

**What problem does it solve?**
- Supports diverse payment methods (cash, cards, transfers)
- Handles partial payments
- Calculates change for cash payments
- Maintains complete payment audit trail

**What value does it generate?**
- Flexible payment options for customers
- Accurate cash management
- Complete financial records
- Reconciliation support

---

## 👥 Stakeholders

- End User: Cashiers, store managers, accountants
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Payment module handles:

1. **Payment Methods**: Cash, debit card, credit card, bank transfer, credit
2. **Payment Recording**: Link payments to sales orders
3. **Change Calculation**: Calculate change for cash payments
4. **Partial Payments**: Split payments across methods
5. **Transaction History**: Complete payment audit trail

---

## 📦 Scope

### Included:
- Single payment processing
- Multiple payment methods per sale
- Change calculation for cash
- Partial payment support
- Payment method configuration
- Transaction lookup by sale
- Payment cancellation/refund

### Not Included:
- Payment gateway integration
- Card terminal integration
- Online payment processing
- Payment plan management

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-PAY-001 | Process Single Payment | ✅ Completed |
| HU-PAY-002 | Process Multiple Payment Methods | ✅ Completed |
| HU-PAY-003 | Calculate Change | ✅ Completed |
| HU-PAY-004 | Record Partial Payment | ✅ Completed |
| HU-PAY-005 | Cancel Payment | ✅ Completed |
| HU-PAY-006 | View Payment History | ✅ Completed |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- Payments must equal or exceed sale total (or create credit)
- Change is only applicable for cash payments
- Payment methods: CASH, DEBIT_CARD, CREDIT_CARD, BANK_TRANSFER, CREDIT, VOUCHER, CHECK
- Credit payments create accounts receivable
- Payments are linked to sales orders and cash drawer shifts
- All payments require user and branch identification
- Cancelled payments require authorization

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: payment**
```sql
CREATE TABLE IF NOT EXISTS payment (
    id BIGSERIAL PRIMARY KEY,
    payment_type VARCHAR(20) NOT NULL DEFAULT 'SALE' CHECK (payment_type IN ('SALE', 'PURCHASE', 'ACCOUNT_RECEIVABLE', 'ACCOUNT_PAYABLE')),
    reference_id BIGINT NOT NULL,
    reference_type VARCHAR(50) NOT NULL,
    payment_method VARCHAR(30) NOT NULL CHECK (payment_method IN ('CASH', 'DEBIT_CARD', 'CREDIT_CARD', 'BANK_TRANSFER', 'CREDIT', 'VOUCHER', 'CHECK')),
    amount DECIMAL(12,2) NOT NULL,
    reference_number VARCHAR(100),
    bank_name VARCHAR(100),
    card_type VARCHAR(20) CHECK (card_type IN ('CREDIT', 'DEBIT')),
    card_last_digits VARCHAR(4),
    authorization_code VARCHAR(50),
    change_amount DECIMAL(12,2) DEFAULT 0,
    cash_drawer_id BIGINT REFERENCES cash_drawer(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    enterprise_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'COMPLETED' CHECK (status IN ('COMPLETED', 'PENDING', 'CANCELLED', 'REFUNDED')),
    cancelled_at TIMESTAMPTZ,
    cancelled_by BIGINT REFERENCES public.users(id),
    cancellation_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);
```

**Table: payment_transaction**
```sql
CREATE TABLE IF NOT EXISTS payment_transaction (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL REFERENCES payment(id),
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('CHARGE', 'REFUND', 'CHARGEBACK')),
    amount DECIMAL(12,2) NOT NULL,
    previous_balance DECIMAL(12,2),
    new_balance DECIMAL(12,2),
    processor_reference VARCHAR(100),
    processor_response TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Table: cash_drawer**
```sql
CREATE TABLE IF NOT EXISTS cash_drawer (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    enterprise_id BIGINT NOT NULL,
    opening_balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    closing_balance DECIMAL(12,2),
    cash_in DECIMAL(12,2) NOT NULL DEFAULT 0,
    cash_out DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED')),
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ,
    notes TEXT
);
```

**Table: cash_movement**
```sql
CREATE TABLE IF NOT EXISTS cash_movement (
    id BIGSERIAL PRIMARY KEY,
    cash_drawer_id BIGINT NOT NULL REFERENCES cash_drawer(id),
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('IN', 'OUT')),
    amount DECIMAL(12,2) NOT NULL,
    description TEXT,
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 📊 Success Metrics

- Payment processing time < 500ms
- Payment accuracy 100%
- Refund processing time < 1 second
- Zero payment reconciliation errors

---

## 🚧 Risks

- Payment data security (PCI compliance consideration)
- Concurrent payment processing
- Network failures during payment
- Cash drawer reconciliation discrepancies

---

## 📡 API Endpoints

### Payments
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/payments` | Process single payment |
| POST | `/payments/batch` | Process multiple payments |
| GET | `/payments` | List payments with filters |
| GET | `/payments/:id` | Get payment by ID |
| GET | `/payments/reference/:referenceType/:referenceId` | Get payments by reference |
| POST | `/payments/:id/cancel` | Cancel payment |

### Cash Drawers
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/cash-drawers` | Open cash drawer |
| GET | `/cash-drawers` | List cash drawers |
| GET | `/cash-drawers/open` | Get open drawer for user/branch |
| GET | `/cash-drawers/:id` | Get drawer by ID |
| POST | `/cash-drawers/:id/close` | Close drawer |
| POST | `/cash-drawers/:id/cash-in` | Add cash to drawer |
| POST | `/cash-drawers/:id/cash-out` | Remove cash from drawer |

### Query Parameters
- `page` - Page number
- `limit` - Items per page
- `method` - Filter by payment method
- `status` - Filter by payment status
- `reference_id` - Filter by reference ID
- `user_id` - Filter by user (cash drawers)
- `branch_id` - Filter by branch

---

## 📁 Module Structure

```
modules/payments/
├── domain.go     # Entity, Repository & Service interfaces
├── service.go    # Repository & Service implementation
├── handler.go    # HTTP handlers
└── routes.go     # Route registration
```

### Migration Updated
- `tenant/migrations/tenant/000005_payments.up.sql` - Enhanced with payment types, transactions, cash movements

---

## Resumen

- **Total de HU**: 6
- **Completadas**: 6
- **Pendientes**: 0
- **Módulo implementado**: payments (with enhanced features)
