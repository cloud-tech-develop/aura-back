# HU-TRANS-004 - Receive Transfer

## 📌 General Information
- ID: HU-TRANS-004
- Epic: EPICA-012
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** inventory clerk
**I want to** receive a transfer
**So that** inventory is added to destination branch

---

## 📡 Technical Requirements

### Endpoint
```
PATCH /api/transfers/{transferId}/receive
```

### Request
```json
{
  "items": [
    {
      "product_id": 101,
      "received_quantity": 48
    }
  ],
  "notes": "2 units damaged in transit"
}
```

### Response (200 OK)
```json
{
  "data": {
    "transfer_id": 1,
    "status": "PARTIAL",
    "inventory_added": true,
    "received_at": "2026-03-23T16:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-012: Transfers Epic
