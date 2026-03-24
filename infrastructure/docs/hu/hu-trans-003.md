# HU-TRANS-003 - Ship Transfer

## 📌 General Information
- ID: HU-TRANS-003
- Epic: EPICA-012
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** inventory clerk
**I want to** ship a transfer
**So that** inventory is deducted from origin branch

---

## 📡 Technical Requirements

### Endpoint
```
PATCH /api/transfers/{transferId}/ship
```

### Request
```json
{
  "items": [
    {
      "product_id": 101,
      "shipped_quantity": 50
    }
  ],
  "notes": "All items shipped"
}
```

### Response (200 OK)
```json
{
  "data": {
    "transfer_id": 1,
    "status": "SHIPPED",
    "inventory_deducted": true,
    "shipped_at": "2026-03-23T14:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-012: Transfers Epic
