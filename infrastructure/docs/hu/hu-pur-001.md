# HU-PUR-001 - Create Purchase Order

## 📌 General Information
- ID: HU-PUR-001
- Epic: EPICA-010
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** create purchase orders to suppliers
**So that** I can request products for inventory replenishment

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/purchase-orders
```

### Request
```json
{
  "supplier_id": 50,
  "branch_id": 1,
  "expected_date": "2026-03-30",
  "items": [
    {
      "product_id": 101,
      "quantity": 100,
      "unit_cost": 30000,
      "discount_amount": 0,
      "tax_rate": 19
    }
  ],
  "notes": "Urgent order"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "order_number": "PO-2026-0001",
    "supplier_name": "Global Supplies Inc.",
    "status": "PENDING",
    "total": 3570000,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-010: Purchases Epic
