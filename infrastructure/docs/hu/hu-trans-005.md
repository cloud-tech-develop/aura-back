# HU-TRANS-005 - Cancel Transfer

## 📌 General Information
- ID: HU-TRANS-005
- Epic: EPICA-012
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** cancel transfer requests
**So that** I can stop unwanted transfers

---

## 📡 Technical Requirements

### Endpoint
```
DELETE /api/transfers/{transferId}
```

### Request
```json
{
  "reason": "Product no longer needed at destination"
}
```

### Response (200 OK)
```json
{
  "data": {
    "transfer_id": 1,
    "status": "CANCELLED",
    "inventory_restored": true,
    "cancelled_at": "2026-03-23T14:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-012: Transfers Epic
