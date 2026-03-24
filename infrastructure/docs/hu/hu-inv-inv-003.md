# HU-INV-INV-003 - View Invoice Details

## 📌 General Information
- ID: HU-INV-INV-003
- Epic: EPICA-007
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** view invoice details
**So that** I can provide invoice information to customers

---

## 🧠 Functional Description

The system must return complete invoice information including all items, taxes, and customer details.

---

## ✅ Acceptance Criteria

### Scenario 1: Get invoice details
- Given that an invoice exists
- When I request details by ID
- Then complete invoice with items is returned

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/invoices/{invoiceId}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "invoice_number": "INV-0001",
    "prefix": "INV",
    "invoice_type": "SALE",
    "invoice_date": "2026-03-23",
    "customer": {
      "id": 123,
      "name": "Acme Corporation",
      "document_number": "901234567-8"
    },
    "branch": {
      "id": 1,
      "name": "Main Store"
    },
    "items": [
      {
        "product_name": "Wireless Mouse",
        "quantity": 2,
        "unit_price": 50000,
        "discount_amount": 0,
        "tax_rate": 19,
        "tax_amount": 19000,
        "line_total": 119000
      }
    ],
    "subtotal": 100000,
    "discount_total": 0,
    "tax_exempt": 0,
    "taxable_amount": 100000,
    "iva_19": 19000,
    "total": 119000,
    "status": "ISSUED",
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-007: Invoicing Epic
