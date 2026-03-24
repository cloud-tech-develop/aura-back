# HU-COMM-004 - Settle Commissions

## 📌 General Information
- ID: HU-COMM-004
- Epic: EPICA-013
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** settle (mark as paid) commissions
**So that** employees can receive their earnings

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/commissions/settle
```

### Request
```json
{
  "employee_id": 10,
  "period": "2026-03",
  "commission_ids": [1, 2, 3],
  "notes": "Monthly settlement March 2026"
}
```

### Response (200 OK)
```json
{
  "data": {
    "settled_count": 3,
    "total_amount": 75000,
    "employee_name": "John Doe",
    "period": "2026-03",
    "settled_at": "2026-03-23T10:00:00Z",
    "settled_by": "Manager Name"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-013: Commissions Epic
