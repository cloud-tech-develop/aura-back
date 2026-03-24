# HU-SHR-003 - Authorize High-Value Shrinkage

## 📌 General Information
- ID: HU-SHR-003
- Epic: EPICA-011
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** authorize high-value shrinkage
**So that** large losses are properly approved

---

## 📡 Technical Requirements

### Endpoint
```
PATCH /api/shrinkage/{shrinkageId}/authorize
```

### Request
```json
{
  "approved": true,
  "notes": "Approved after inspection"
}
```

### Response (200 OK)
```json
{
  "data": {
    "shrinkage_id": 1,
    "status": "APPROVED",
    "authorized_by": "Manager Name",
    "authorized_at": "2026-03-23T14:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-011: Shrinkage Epic
