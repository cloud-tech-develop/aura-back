# HU-CART-008 - Convert Cart to Sale Order

## 📌 General Information
- ID: HU-CART-008
- Epic: EPICA-005
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** convert a cart to a sales order
**So that** I can finalize the sale and proceed to payment

---

## 🧠 Functional Description

The system must convert the cart to a sales order, update inventory (if configured), and prepare for payment processing.

---

## ✅ Acceptance Criteria

### Scenario 1: Convert to sales order
- Given that a cart has items and sufficient stock
- When I convert it to sales order
- Then:
  - Sales order is created with sequential number
  - All items and pricing preserved
  - Cart status changes to CONVERTED
  - Reference to sales order stored

### Scenario 2: Inventory reserved
- Given that auto-reserve is enabled
- When cart converts to order
- Then inventory quantities are reserved

---

## ❌ Error Cases

- Empty cart returns 400
- Insufficient stock returns 400
- Cart already converted returns 400
- Cart not OPEN returns 400

---

## 🔐 Business Rules

- Sales order number: sequential per branch
- Status: PENDING_PAYMENT
- Inventory reservation depends on configuration
- All items and totals preserved

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/carts/{cartId}/convert
```

### Response (201 Created)
```json
{
  "data": {
    "cart_id": 1,
    "sales_order": {
      "id": 1,
      "order_number": "SO-2026-0001",
      "status": "PENDING_PAYMENT",
      "customer": {
        "id": 123,
        "name": "Acme Corporation"
      },
      "branch": {
        "id": 1,
        "name": "Main Store"
      },
      "items_count": 3,
      "subtotal": 145000,
      "tax_total": 27550,
      "grand_total": 172550,
      "created_at": "2026-03-23T10:30:00Z"
    },
    "status": "CONVERTED",
    "converted_at": "2026-03-23T10:30:00Z"
  },
  "success": true,
  "message": "Cart converted to sales order"
}
```

---

## 📎 Dependencies

- EPICA-005: Cart Epic
- EPICA-004: Inventory (for reservation)
