# HU-CART-003 - Update Cart Item

## 📌 General Information
- ID: HU-CART-003
- Epic: EPICA-005
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** update item quantities in the cart
**So that** I can adjust the sale before finalizing

---

## 🧠 Functional Description

The system must allow updating item quantities and recalculate totals. Removing an item (setting quantity to 0) deletes the line.

---

## ✅ Acceptance Criteria

### Scenario 1: Update quantity
- Given that a cart has items
- When I update an item quantity to 5
- Then the new quantity is saved
- And totals are recalculated

### Scenario 2: Remove item (quantity = 0)
- Given that a cart has items
- When I set quantity to 0
- Then the item is removed from cart
- And totals are recalculated

---

## ❌ Error Cases

- Item not found returns 404
- Cart not OPEN returns 400
- Invalid quantity returns 400
- Insufficient stock returns 400

---

## 🔐 Business Rules

- Quantity must be positive integer
- Zero quantity removes the item
- Totals recalculated on any change
- Stock validation on quantity increase

---

## 📡 Technical Requirements

### Endpoint
```
PUT /api/carts/{cartId}/items/{itemId}
```

### Request
```json
{
  "quantity": 5
}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "cart_id": 1,
    "product_name": "Wireless Mouse",
    "quantity": 5,
    "unit_price": 50000,
    "tax_amount": 47500,
    "line_total": 297500,
    "cart": {
      "subtotal": 250000,
      "tax_total": 47500,
      "grand_total": 297500,
      "items_count": 1
    }
  },
  "success": true,
  "message": "Cart item updated"
}
```

---

## 📎 Dependencies

- EPICA-005: Cart Epic
- HU-CART-002: Add Item to Cart
