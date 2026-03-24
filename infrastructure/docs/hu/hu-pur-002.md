# HU-PUR-002 - Receive Goods

## 📌 General Information
- ID: HU-PUR-002
- Epic: EPICA-010
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** inventory clerk
**I want to** receive goods from a purchase order
**So that** inventory is updated with new stock

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/purchase-orders/{orderId}/receive
```

### Request
```json
{
  "items": [
    {
      "product_id": 101,
      "quantity_received": 95,
      "batch_number": "LOT-2026-001",
      "expiration_date": "2027-03-23"
    }
  ],
  "notes": "5 units damaged"
}
```

### Response (201 Created)
```json
{
  "data": {
    "purchase_order_id": 1,
    "purchase_id": 1,
    "purchase_number": "PUR-2026-0001",
    "status": "PARTIAL",
    "items_received": 95,
    "inventory_updated": true,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-010: Purchases Epic
