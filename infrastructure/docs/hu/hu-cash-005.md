# HU-CASH-005 - Perform Cash Reconciliation

## 📌 General Information
- ID: HU-CASH-005
- Epic: EPICA-009
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** perform cash reconciliation (arqueo)
**So that** I can verify cash drawer accuracy

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/cash-shifts/{shiftId}/reconcile
```

### Response (200 OK)
```json
{
  "data": {
    "shift_id": 1,
    "opening_amount": 200000,
    "closing_amount": 500000,
    "sales_total": 450000,
    "expected_cash": 500000,
    "difference": 0,
    "by_payment_method": {
      "CASH": {
        "sales": 300000,
        "movements_in": 50000,
        "movements_out": 100000,
        "expected": 450000
      }
    },
    "status": "RECONCILED"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-009: Cash Drawer Epic
