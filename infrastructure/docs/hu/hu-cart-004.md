# HU-CART-004 - Remove Item from Cart

## 📌 General Information
- ID: HU-CART-004
- Epic: EPICA-005
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** remove items from the cart
**So that** I can correct mistakes before finalizing the sale

---

## 🧠 Functional Description

The system must allow removing items from the cart with proper totals recalculation.

---

## ✅ Acceptance Criteria

### Scenario 1: Remove item
- Given that a cart has items
- When I remove an item
- Then the item is deleted
- And totals are recalculated

---

## ❌ Error Cases

- Item not found returns 404
- Cart not OPEN returns 400

---

## 📡 Technical Requirements

### Endpoint
```
DELETE /api/carts/{cartId}/items/{itemId}
```

### Response (200 OK)
```json
{
  "data": {
    "cart_id": 1,
    "subtotal": 150000,
    "tax_total": 28500,
    "grand_total": 178500,
    "items_count": 2,
    "removed_item_id": 1
  },
  "success": true,
  "message": "Item removed from cart"
}
```

---

## 📎 Dependencies

- EPICA-005: Cart Epic
