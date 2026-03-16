# HU-SALES-004 - Payment Processing

## 📌 General Information
- ID: HU-SALES-004
- Epic: EPIC-SALES-001
- Priority: High
- State: Backlog
- Progress: 0%
- Author: QA Engineer Aura POS
- Date: 2026-03-15

---

## 👤 User Story

**As a** cashier
**I want to** process payments for sales orders using multiple payment methods
**So that** I can complete customer transactions efficiently

---

## 🧠 Functional Description

The system must support multiple payment methods for sales orders:
- Cash
- Credit/Debit Card
- Bank Transfer
- Credit Account (for registered customers)

The system must record all payment details, calculate change for cash payments, and update the sales order status accordingly.

---

## ✅ Acceptance Criteria

### Scenario 1: Cash payment with exact amount
- Given that a sales order total is $100,000
- When I receive $100,000 cash
- Then the payment must be recorded as complete
- And the sales order status must update to PAID

### Scenario 2: Cash payment with change
- Given that a sales order total is $100,000
- When I receive $150,000 cash
- Then the payment must be recorded as $100,000
- And change of $50,000 must be calculated and displayed
- And the sales order status must update to PAID

### Scenario 3: Card payment
- Given that a customer chooses card payment
- When I process the card transaction
- Then the payment must be recorded with:
  - Card type (credit/debit)
  - Authorization code
  - Last 4 digits
  - Bank information
- And the sales order status must update to PAID

### Scenario 4: Split payment
- Given that a customer wants to pay with multiple methods
- When I record cash + card payment
- Then each payment method must be recorded separately
- And the sum must equal the order total
- And the sales order status must update to PAID

### Scenario 5: Credit account payment
- Given that the customer has a credit account
- When I process the payment as credit
- Then the payment must be recorded
- And the customer's credit balance must be updated
- And the sales order status must update to CREDIT

---

## ❌ Error Cases

- Payment amount less than order total must return error 400
- Invalid card data must return error 400
- Credit account with insufficient limit must return error 400
- Payment on already paid order must return error 400
- Card decline must be handled gracefully

---

## 🔐 Business Rules

- Cash payments require exact change calculation
- Card payments require authorization code validation
- Credit payments require customer credit limit check
- All payments must be recorded with timestamp and user
- Payment records are immutable after creation
- Cash drawer reconciliation is tracked per user/shift

---

## 🗄️ Database Schema (PostgreSQL)

### Table: payment
```sql
CREATE TABLE payment (
    id BIGSERIAL PRIMARY KEY,
    sales_order_id BIGINT NOT NULL REFERENCES sales_order(id),
    payment_method VARCHAR(50) NOT NULL CHECK (payment_method IN ('CASH', 'CARD', 'TRANSFER', 'CREDIT')),
    amount DECIMAL(12,2) NOT NULL,
    reference VARCHAR(100),
    card_type VARCHAR(20) CHECK (card_type IN ('CREDIT', 'DEBIT')),
    card_last_four VARCHAR(4),
    bank_name VARCHAR(100),
    authorization_code VARCHAR(50),
    user_id INTEGER NOT NULL REFERENCES usuario(id),
    branch_id INTEGER NOT NULL REFERENCES branch(id),
    empresa_id INTEGER NOT NULL REFERENCES empresa(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT payment_order_fk FOREIGN KEY (sales_order_id) REFERENCES sales_order(id),
    CONSTRAINT payment_user_fk FOREIGN KEY (user_id) REFERENCES usuario(id),
    CONSTRAINT payment_branch_fk FOREIGN KEY (branch_id) REFERENCES branch(id),
    CONSTRAINT payment_empresa_fk FOREIGN KEY (empresa_id) REFERENCES empresa(id),
    CONSTRAINT payment_positive_amount CHECK (amount > 0)
);

CREATE INDEX idx_payment_order ON payment(sales_order_id);
CREATE INDEX idx_payment_user ON payment(user_id);
CREATE INDEX idx_payment_empresa ON payment(empresa_id);
```

### Table: cash_drawer
```sql
CREATE TABLE cash_drawer (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES usuario(id),
    branch_id INTEGER NOT NULL REFERENCES branch(id),
    empresa_id INTEGER NOT NULL REFERENCES empresa(id),
    opening_balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    closing_balance DECIMAL(12,2),
    cash_in DECIMAL(12,2) NOT NULL DEFAULT 0,
    cash_out DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED')),
    opened_at TIMESTAMP NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMP,
    
    CONSTRAINT cash_drawer_user_fk FOREIGN KEY (user_id) REFERENCES usuario(id),
    CONSTRAINT cash_drawer_branch_fk FOREIGN KEY (branch_id) REFERENCES branch(id),
    CONSTRAINT cash_drawer_empresa_fk FOREIGN KEY (empresa_id) REFERENCES empresa(id)
);

CREATE INDEX idx_cash_drawer_user ON cash_drawer(user_id);
CREATE INDEX idx_cash_drawer_status ON cash_drawer(status);
```

---

## 🎨 UI/UX Considerations

- Payment method selection with icons
- Cash calculator with change display
- Card reader integration interface
- Payment confirmation dialog
- Receipt preview

---

## 📡 Technical Requirements

### Entities (Java)

**PaymentEntity:**
- id (Long, PK)
- sales_order_id (Long, FK)
- payment_method (String): CASH, CARD, TRANSFER, CREDIT
- amount (BigDecimal)
- reference (String, nullable)
- card_type (String): CREDIT, DEBIT, nullable
- card_last_four (String, nullable)
- bank_name (String, nullable)
- authorization_code (String, nullable)
- user_id (Long, FK)
- branch_id (Long, FK)
- empresa_id (Long, FK)
- created_at (LocalDateTime)

**CashDrawerEntity:**
- id (Long, PK)
- user_id (Long, FK)
- branch_id (Long, FK)
- empresa_id (Long, FK)
- opening_balance (BigDecimal)
- closing_balance (BigDecimal, nullable)
- cash_in (BigDecimal)
- cash_out (BigDecimal)
- status (String): OPEN, CLOSED
- opened_at (LocalDateTime)
- closed_at (LocalDateTime, nullable)

### Controllers

- `PaymentController` with endpoints:
  - `POST /api/payments` - Process payment
  - `GET /api/payments/order/{orderId}` - Get payments for order
  - `POST /api/cash-drawer/open` - Open cash drawer
  - `POST /api/cash-drawer/close` - Close cash drawer

---

## 🚫 Out of Scope

- Payment gateway integration (external)
- Advanced fraud detection
- Recurring payments
- Payment plan management

---

## 📎 Dependencies

- HU-SALES-003: Sale/Order Creation
- Existing Branch and Usuario entities
- POS hardware integration (card reader)
