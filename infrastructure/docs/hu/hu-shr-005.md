# HU-SHR-005 - Cancel Shrinkage

## 📌 General Information
- ID: HU-SHR-005
- Epic: EPICA-011
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** cancel shrinkage records
**So that** I can correct mistakes

---

## 📡 Technical Requirements

### Endpoint
```
DELETE /api/shrinkage/{shrinkageId}
```

### Request
```json
{
  "reason": "Products found undamaged after inspection"
}
```

### Response (200 OK)
```json
{
  "data": {
    "shrinkage_id": 1,
    "status": "CANCELLED",
    "inventory_restored": true,
    "cancelled_at": "2026-03-23T14:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-011: Shrinkage Epic
