# HU-INV-INV-001 - Auto-generate Invoice from Sale

## 📌 General Information
- ID: HU-INV-INV-001
- Epic: EPICA-007
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** automatically receive an invoice when completing a sale
**So that** I have proper fiscal documentation for the transaction

---

## 🧠 Functional Description

The system must automatically generate invoices when sales orders are completed. Invoice numbering follows sequential prefixes per branch, and all fiscal data is properly calculated and stored.

---

## ✅ Acceptance Criteria

### Scenario 1: Auto-generate invoice on sale completion
- Given that a sales order is completed
- When the sale finalizes
- Then an invoice must be auto-generated with:
  - Sequential invoice number (prefix+number)
  - All sale items as invoice items
  - Correct tax calculations
  - Status: ISSUED
  - Linked to sales order

### Scenario 2: Invoice number increment
- Given that the last invoice for branch had number "INV-001"
- When a new invoice is generated
- Then the new number must be "INV-002"

### Scenario 3: Invoice with multiple tax rates
- Given that a sales order has items with different tax rates (0%, 5%, 19%)
- When the invoice is generated
- Then each item must have its correct tax rate
- And totals must be broken down by tax rate

---

## ❌ Error Cases

- Duplicate invoice number returns error 500
- Missing prefix configuration returns error 500
- Invoice generation failure rolls back sale completion

---

## 🔐 Business Rules

- Invoice numbers are sequential per branch prefix
- Invoice prefixes are configured per branch
- Invoice totals must match linked sale
- Invoices are immutable after generation
- Only soft delete allowed for audit compliance
- Tax calculations follow Colombian regulations

---

## 📡 Technical Requirements

### Endpoint (Manual Invoice)
```
POST /api/invoices
```

### Method: POST

### Request
```json
{
  "invoice_type": "SALE",
  "sales_order_id": 101,
  "customer_id": 123,
  "notes": "Invoice generated at customer request"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "invoice_number": "INV-0001",
    "prefix": "INV",
    "invoice_type": "SALE",
    "sales_order_id": 101,
    "customer_id": 123,
    "customer_name": "Acme Corporation",
    "branch_id": 1,
    "invoice_date": "2026-03-23",
    "subtotal": 84000,
    "discount_total": 0,
    "tax_exempt": 0,
    "taxable_amount": 84000,
    "iva_19": 15960,
    "total": 99960,
    "status": "ISSUED",
    "items": [
      {
        "id": 1,
        "product_name": "Wireless Mouse",
        "product_sku": "WM-001",
        "quantity": 2,
        "unit_price": 42000,
        "tax_rate": 19,
        "tax_amount": 7980,
        "line_total": 49980
      }
    ],
    "created_at": "2026-03-23T10:30:00Z"
  },
  "success": true,
  "message": "Invoice generated successfully"
}
```

### Event Trigger
```json
{
  "event": "invoice.generated",
  "data": {
    "invoice_id": 1,
    "invoice_number": "INV-0001",
    "sales_order_id": 101,
    "total": 99960
  }
}
```

---

## 🧪 Testing Criteria

### Unit Tests
- Test sequential number generation
- Test tax calculation
- Test invoice item mapping

### Integration Tests
- Test auto-generation on sale completion
- Test prefix increment
- Test concurrent invoice generation

---

## 📎 Dependencies

- EPICA-007: Invoicing Epic
- Existing sales orders module
- Existing third parties module
- Invoice prefix configuration
