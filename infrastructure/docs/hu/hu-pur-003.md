# HU-PUR-003 - Record Purchase Payment

## 📌 General Information
- ID: HU-PUR-003
- Epic: EPICA-010
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** accountant
**I want to** record payments to suppliers
**So that** I can track accounts payable

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/purchase-payments
```

### Request
```json
{
  "purchase_id": 1,
  "payment_method": "BANK_TRANSFER",
  "amount": 2000000,
  "reference_number": "TRF-456789",
  "notes": "First payment"
}
```

### Response (201 Created)
```json
{
  "data": {
    "payment_id": 1,
    "purchase_id": 1,
    "amount": 2000000,
    "purchase_total": 3570000,
    "total_paid": 2000000,
    "remaining": 1570000,
    "status": "PARTIAL"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-010: Purchases Epic
