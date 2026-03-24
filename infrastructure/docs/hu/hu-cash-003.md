# HU-CASH-003 - Close Cash Shift

## 📌 General Information
- ID: HU-CASH-003
- Epic: EPICA-009
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** close my cash shift
**So that** I can end my work session and reconcile cash

---

## ✅ Acceptance Criteria

### Scenario 1: Close shift with authorization
- Given that a shift is OPEN
- When I close with:
  - closing_amount: 500000
  - authorized_by: manager_id
- Then:
  - Shift status = CLOSED
  - Expected amount calculated
  - Difference calculated

---

## 📡 Technical Requirements

### Endpoint
```
PATCH /api/cash-shifts/{shiftId}/close
```

### Request
```json
{
  "closing_amount": 500000,
  "notes": "End of day shift"
}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "status": "CLOSED",
    "opening_amount": 200000,
    "closing_amount": 500000,
    "expected_amount": 495000,
    "difference": 5000,
    "closed_at": "2026-03-23T18:00:00Z",
    "closed_by": "John Doe",
    "sales_total": 450000,
    "cash_in_total": 300000,
    "cash_out_total": 50000
  },
  "success": true,
  "message": "Shift closed successfully"
}
```

---

## 📎 Dependencies

- EPICA-009: Cash Drawer Epic
- HU-CASH-002: Open Cash Shift
