# HU-CASH-004 - Record Cash Movement

## 📌 General Information
- ID: HU-CASH-004
- Epic: EPICA-009
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** record cash movements (drops, withdrawals)
**So that** I can maintain accurate cash drawer balance

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/cash-movements
```

### Request
```json
{
  "shift_id": 1,
  "movement_type": "IN",
  "reason": "DROPS",
  "amount": 100000,
  "notes": "Safe drop"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "shift_id": 1,
    "movement_type": "IN",
    "reason": "DROPS",
    "amount": 100000,
    "created_at": "2026-03-23T14:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-009: Cash Drawer Epic
