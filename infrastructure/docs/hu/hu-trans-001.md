# HU-TRANS-001 - Create Transfer Request

## 📌 General Information
- ID: HU-TRANS-001
- Epic: EPICA-012
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** create transfer requests between branches
**So that** I can redistribute inventory

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/transfers
```

### Request
```json
{
  "origin_branch_id": 1,
  "destination_branch_id": 2,
  "items": [
    {
      "product_id": 101,
      "quantity": 50
    }
  ],
  "notes": "Stock redistribution"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "transfer_number": "TRF-2026-0001",
    "origin_branch_name": "Main Store",
    "destination_branch_name": "Branch 2",
    "status": "PENDING",
    "items_count": 1,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-012: Transfers Epic
