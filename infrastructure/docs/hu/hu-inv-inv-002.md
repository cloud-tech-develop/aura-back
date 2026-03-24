# HU-INV-INV-002 - Manual Invoice Creation

## 📌 General Information
- ID: HU-INV-INV-002
- Epic: EPICA-007
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** create invoices manually
**So that** I can invoice sales without auto-generation

---

## 🧠 Functional Description

The system must allow manual invoice creation with custom details when auto-generation is disabled.

---

## ✅ Acceptance Criteria

### Scenario 1: Create manual invoice
- Given that I have invoice prefix configured
- When I create a manual invoice with items
- Then the invoice is created with next sequential number

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/invoices
```

### Request
```json
{
  "invoice_type": "SALE",
  "prefix_id": 1,
  "customer_id": 123,
  "branch_id": 1,
  "invoice_date": "2026-03-23",
  "items": [
    {
      "product_id": 101,
      "product_name": "Wireless Mouse",
      "quantity": 2,
      "unit_price": 50000,
      "discount_amount": 0,
      "tax_rate": 19
    }
  ],
  "notes": "Manual invoice"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "invoice_number": "INV-0002",
    "status": "ISSUED",
    "total": 119000,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-007: Invoicing Epic
