# HU-CASH-006 - View Shift Summary

## 📌 General Information
- ID: HU-CASH-006
- Epic: EPICA-009
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view shift summary
**So that** I can review cash operations for a specific shift

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/cash-shifts/{shiftId}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "cash_drawer": {
      "id": 1,
      "name": "MAIN"
    },
    "user": {
      "id": 5,
      "name": "John Doe"
    },
    "branch": {
      "id": 1,
      "name": "Main Store"
    },
    "opening_amount": 200000,
    "closing_amount": 500000,
    "expected_amount": 495000,
    "difference": 5000,
    "status": "CLOSED",
    "opened_at": "2026-03-23T08:00:00Z",
    "closed_at": "2026-03-23T18:00:00Z",
    "movements": [
      {
        "id": 1,
        "movement_type": "IN",
        "reason": "SALE",
        "amount": 150000,
        "created_at": "2026-03-23T10:00:00Z"
      }
    ],
    "summary": {
      "sales_count": 45,
      "sales_total": 450000,
      "cash_in_total": 350000,
      "cash_out_total": 100000
    }
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-009: Cash Drawer Epic
