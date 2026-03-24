# HU-CART-005 - Apply Discounts

## 📌 General Information
- ID: HU-CART-005
- Epic: EPICA-005
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** apply discounts to cart items or the entire cart
**So that** I can offer promotions and price adjustments

---

## 🧠 Functional Description

The system must support item-level and cart-level discounts with percentage or fixed amount types.

---

## ✅ Acceptance Criteria

### Scenario 1: Apply item discount
- Given that a cart has an item
- When I apply a 10% discount to the item
- Then the discount_amount is calculated
- And line_total is updated

### Scenario 2: Apply cart discount
- Given that a cart has items
- When I apply a fixed $10000 cart discount
- Then discount_total is updated
- And grand_total reflects the discount

### Scenario 3: Apply percentage cart discount
- Given that a cart has subtotal of 200000
- When I apply 5% cart discount
- Then discount_total = 10000
- And grand_total = 190000 + taxes

---

## ❌ Error Cases

- Discount exceeds item/cart total returns 400
- Invalid discount type returns 400
- Cart not OPEN returns 400

---

## 🔐 Business Rules

- Discount types: PERCENTAGE, FIXED
- Percentage: 0-100%
- Fixed: cannot exceed line or cart total
- Discount calculated before tax
- Maximum discount limits may apply

---

## 📡 Technical Requirements

### Endpoint (Item Discount)
```
PATCH /api/carts/{cartId}/items/{itemId}/discount
```

### Request
```json
{
  "discount_type": "PERCENTAGE",
  "discount_value": 10
}
```

### Endpoint (Cart Discount)
```
PATCH /api/carts/{cartId}/discount
```

### Request
```json
{
  "discount_type": "FIXED",
  "discount_value": 10000,
  "reason": "Customer loyalty"
}
```

### Response (200 OK)
```json
{
  "data": {
    "cart_id": 1,
    "discount_type": "FIXED",
    "discount_value": 10000,
    "discount_total": 10000,
    "subtotal": 200000,
    "grand_total": 226000,
    "items_count": 3
  },
  "success": true,
  "message": "Discount applied"
}
```

---

## 📎 Dependencies

- EPICA-005: Cart Epic
