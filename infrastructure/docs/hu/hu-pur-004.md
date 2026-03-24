# HU-PUR-004 - Cancel Purchase

## 📌 General Information
- ID: HU-PUR-004
- Epic: EPICA-010
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** cancel purchases
**So that** I can correct errors or handle returns

---

## 📡 Technical Requirements

### Endpoint
```
DELETE /api/purchases/{purchaseId}
```

### Request
```json
{
  "reason": "Supplier cancelled delivery"
}
```

### Response (200 OK)
```json
{
  "data": {
    "purchase_id": 1,
    "purchase_number": "PUR-2026-0001",
    "status": "CANCELLED",
    "inventory_reversed": true,
    "cancelled_at": "2026-03-23T14:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-010: Purchases Epic
