# HU-INV-INV-005 - Cancel Invoice

## 📌 General Information
- ID: HU-INV-INV-005
- Epic: EPICA-007
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** cancel invoices
**So that** I can correct errors while maintaining audit trail

---

## 🧠 Functional Description

The system must allow invoice cancellation with mandatory credit note reference and authorization.

---

## ✅ Acceptance Criteria

### Scenario 1: Cancel invoice
- Given that an invoice exists with status ISSUED
- When I cancel with credit_note_reference
- Then:
  - Invoice status = CANCELLED
  - Credit note created or referenced
  - Cancellation logged

---

## ❌ Error Cases

- Already cancelled returns 400
- Missing credit note reference returns 400
- Insufficient permissions returns 403

---

## 📡 Technical Requirements

### Endpoint
```
DELETE /api/invoices/{invoiceId}
```

### Request
```json
{
  "cancellation_reason": "Customer requested cancellation",
  "credit_note_number": "NC-0001"
}
```

### Response (200 OK)
```json
{
  "data": {
    "invoice_id": 1,
    "invoice_number": "INV-0001",
    "status": "CANCELLED",
    "cancelled_at": "2026-03-23T14:00:00Z",
    "cancelled_by": "Manager Name",
    "credit_note_reference": "NC-0001",
    "cancellation_reason": "Customer requested cancellation"
  },
  "success": true,
  "message": "Invoice cancelled"
}
```

---

## 📎 Dependencies

- EPICA-007: Invoicing Epic
