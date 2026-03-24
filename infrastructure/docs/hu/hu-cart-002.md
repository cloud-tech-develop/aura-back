# HU-CART-002 - Add Item to Cart

## 📌 General Information
- ID: HU-CART-002
- Epic: EPICA-005
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** add items to the cart
**So that** I can build the sale with products and quantities

---

## 🧠 Functional Description

The system must add products to the cart with quantity, apply correct pricing (including customer-specific prices), and update cart totals.

---

## ✅ Acceptance Criteria

### Scenario 1: Add product to cart
- Given that a cart exists and is OPEN
- When I add a product with quantity 2
- Then the item is added with:
  - Product price from catalog
  - Calculated tax
  - Updated cart totals

### Scenario 2: Add item with customer pricing
- Given that a cart has a customer with specific pricing
- When I add a product
- Then customer price is used instead of catalog price

### Scenario 3: Update quantity if product exists
- Given that a cart has the product already
- When I add the same product again
- Then quantities are summed

---

## ❌ Error Cases

- Product not found returns 404
- Insufficient stock returns 400
- Cart is not OPEN returns 400
- Invalid quantity returns 400

---

## 🔐 Business Rules

- Price comes from: customer price > volume price > catalog price
- Tax calculated based on product tax_rate
- Stock validation before adding
- If item exists, quantities are summed

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/carts/{cartId}/items
```

### Request
```json
{
  "product_id": 101,
  "quantity": 2,
  "notes": "Customer requested gift wrap"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "cart_id": 1,
    "product_id": 101,
    "product_name": "Wireless Mouse",
    "product_sku": "WM-001",
    "quantity": 2,
    "unit_price": 50000,
    "discount_amount": 0,
    "tax_rate": 19,
    "tax_amount": 19000,
    "line_total": 119000,
    "cart": {
      "subtotal": 100000,
      "discount_total": 0,
      "tax_total": 19000,
      "grand_total": 119000,
      "items_count": 1
    }
  },
  "success": true,
  "message": "Item added to cart"
}
```

---

## 📎 Dependencies

- EPICA-005: Cart Epic
- HU-CART-001: Create Cart
