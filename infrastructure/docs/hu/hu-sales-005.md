# HU-SALES-005 - Invoice Generation

## 📌 General Information

- ID: HU-SALES-005
- Epic: EPIC-SALES-001
- Priority: High
- State: Backlog
- Progress: 0%
- Author: QA Engineer Aura POS
- Date: 2026-03-15

---

## 👤 User Story

**As a** store manager or cashier
**I want to** generate invoices for completed sales
**So that** I can provide customers with proper documentation for their purchases

---

## 🧠 Functional Description

The system must generate invoices for completed sales orders. Invoices include:

- Invoice number with prefix and sequence
- Customer information
- Complete item breakdown with taxes
- Payment details
- Company information
- Legal requirements for tax documentation

Invoices are generated automatically when a sale is completed, but can also be generated manually for existing orders.

---

## ✅ Acceptance Criteria

### Scenario 1: Automatic invoice generation

- Given that a sales order is completed with payment
- When the payment is processed successfully
- Then an invoice must be automatically generated with:
  - Unique invoice number (prefix + sequence)
  - Date and time of issue
  - Customer details
  - Item list with quantities and prices
  - Tax calculations
  - Total amount
  - Payment information
- And the invoice must be linked to the sales order

### Scenario 2: Invoice number sequencing

- Given that multiple invoices exist
- When a new invoice is generated
- Then the invoice number must be sequential
- And the prefix must be configurable per branch

### Scenario 3: Manual invoice generation

- Given that a sales order exists without an invoice
- When I manually generate an invoice
- Then the invoice must be created with all required data
- And linked to the original sales order

### Scenario 4: Invoice retrieval and printing

- Given that an invoice exists
- When I search for the invoice
- Then I can view all details
- And I can print or download the invoice

### Scenario 5: Tax calculation on invoice

- Given that products have different tax rates
- When the invoice is generated
- Then taxes must be calculated per item
- And the summary must show total tax amount

---

## ❌ Error Cases

- Generating invoice for non-existent order must return error 404
- Duplicate invoice numbers must be prevented
- Invoice generation without required data must return error 400
- Accessing invoices from other companies must return error 403

---

## 🔐 Business Rules

- Invoice numbers are unique per company and branch
- Prefix is configurable per branch
- Invoices are immutable after creation
- Soft delete is used for invoice removal
- Only users with appropriate permissions can view all invoices
- Invoice generation is atomic with sales order completion

---

## 🗄️ Database Schema (PostgreSQL)

### Table: invoice

```sql
CREATE TABLE invoice (
    id BIGSERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) NOT NULL,
    prefix VARCHAR(10) NOT NULL,
    sequence BIGINT NOT NULL,
    sales_order_id BIGINT NOT NULL REFERENCES sales_order(id),
    customer_id INTEGER REFERENCES customer(id),
    user_id INTEGER NOT NULL REFERENCES usuario(id),
    branch_id INTEGER NOT NULL REFERENCES branch(id),
    enterprise_id INTEGER NOT NULL REFERENCES empresa(id),
    subtotal DECIMAL(12,2) NOT NULL,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    issue_date TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'ISSUED' CHECK (status IN ('ISSUED', 'PAID', 'CANCELLED', 'OVERDUE')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT invoice_empresa_fk FOREIGN KEY (enterprise_id) REFERENCES empresa(id),
    CONSTRAINT invoice_branch_fk FOREIGN KEY (branch_id) REFERENCES branch(id),
    CONSTRAINT invoice_user_fk FOREIGN KEY (user_id) REFERENCES usuario(id),
    CONSTRAINT invoice_sales_order_fk FOREIGN KEY (sales_order_id) REFERENCES sales_order(id),
    CONSTRAINT invoice_customer_fk FOREIGN KEY (customer_id) REFERENCES customer(id),
    CONSTRAINT invoice_number_unique UNIQUE (enterprise_id, invoice_number),
    CONSTRAINT invoice_sequence_unique UNIQUE (enterprise_id, branch_id, prefix, sequence)
);

CREATE INDEX idx_invoice_empresa ON invoice(enterprise_id);
CREATE INDEX idx_invoice_branch ON invoice(branch_id);
CREATE INDEX idx_invoice_customer ON invoice(customer_id);
CREATE INDEX idx_invoice_number ON invoice(invoice_number);
CREATE INDEX idx_invoice_status ON invoice(status);
CREATE INDEX idx_invoice_deleted_at ON invoice(deleted_at) WHERE deleted_at IS NULL;
```

### Table: invoice_item

```sql
CREATE TABLE invoice_item (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoice(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    sales_order_item_id BIGINT REFERENCES sales_order_item(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL,
    tax_amount DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT invoice_item_invoice_fk FOREIGN KEY (invoice_id) REFERENCES invoice(id) ON DELETE CASCADE,
    CONSTRAINT invoice_item_product_fk FOREIGN KEY (product_id) REFERENCES product(id),
    CONSTRAINT invoice_item_sales_order_item_fk FOREIGN KEY (sales_order_item_id) REFERENCES sales_order_item(id)
);

CREATE INDEX idx_invoice_item_invoice ON invoice_item(invoice_id);
CREATE INDEX idx_invoice_item_product ON invoice_item(product_id);
```

### Table: invoice_prefix

```sql
CREATE TABLE invoice_prefix (
    id BIGSERIAL PRIMARY KEY,
    branch_id INTEGER NOT NULL REFERENCES branch(id),
    enterprise_id INTEGER NOT NULL REFERENCES empresa(id),
    prefix VARCHAR(10) NOT NULL,
    next_sequence BIGINT NOT NULL DEFAULT 1,
    description VARCHAR(100),

    CONSTRAINT invoice_prefix_branch_fk FOREIGN KEY (branch_id) REFERENCES branch(id),
    CONSTRAINT invoice_prefix_empresa_fk FOREIGN KEY (enterprise_id) REFERENCES empresa(id),
    CONSTRAINT invoice_prefix_unique UNIQUE (enterprise_id, branch_id, prefix)
);

CREATE INDEX idx_invoice_prefix_empresa ON invoice_prefix(enterprise_id);
CREATE INDEX idx_invoice_prefix_branch ON invoice_prefix(branch_id);
```

---

## 🎨 UI/UX Considerations

- Invoice preview before printing
- Search by invoice number, customer, or date
- Invoice template customization
- Print and download options (PDF)
- Invoice status indicators

---

## 📡 Technical Requirements

### Entities (Java)

**InvoiceEntity:**

- id (Long, PK)
- invoice_number (String, unique)
- prefix (String)
- sequence (Long)
- sales_order_id (Long, FK)
- customer_id (Long, FK, nullable)
- user_id (Long, FK)
- branch_id (Long, FK)
- enterprise_id (Long, FK)
- subtotal (BigDecimal)
- discount (BigDecimal)
- tax_total (BigDecimal)
- total (BigDecimal)
- issue_date (LocalDateTime)
- due_date (LocalDateTime, nullable)
- status (String): ISSUED, PAID, CANCELLED, OVERDUE
- notes (String, TEXT)
- created_at (LocalDateTime)
- updated_at (LocalDateTime, nullable)
- deleted_at (LocalDateTime, nullable)

**InvoicePrefixEntity:**

- id (Long, PK)
- branch_id (Long, FK)
- enterprise_id (Long, FK)
- prefix (String)
- next_sequence (Long)
- description (String, nullable)

### Services

- `InvoiceService` with methods:
  - `generateFromOrder(SalesOrderDto order)`
  - `generateManual(Long orderId)`
  - `getInvoiceById(Long id)`
  - `searchInvoices(PageableDto pageable)`
  - `deleteInvoice(Long id)` (soft delete)

### Controllers

- `InvoiceController` with endpoints:
  - `POST /api/invoices` - Generate invoice (manual)
  - `GET /api/invoices/{id}` - Get invoice details
  - `GET /api/invoices/number/{number}` - Get by invoice number
  - `POST /api/invoices/page` - Paginated search
  - `DELETE /api/invoices/{id}` - Soft delete invoice
  - `GET /api/invoices/{id}/pdf` - Generate PDF

---

## 🚫 Out of Scope

- Electronic invoicing (DIAN integration)
- Complex invoice templates
- Invoice editing after issuance
- Credit notes and adjustments
- Multi-currency support

---

## 📎 Dependencies

- HU-SALES-003: Sale/Order Creation
- HU-SALES-004: Payment Processing
- Existing Branch and Empresa entities
- PDF generation library (e.g., iText, Flying Saucer)
