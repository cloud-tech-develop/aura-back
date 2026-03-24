# HU-CART-001 - Create Cart

## 📌 General Information
- ID: HU-CART-001
- Epic: EPICA-005
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** create a shopping cart
**So that** I can start adding products for a sale or quotation

---

## 🧠 Functional Description

The system must allow creating shopping carts for sales or quotations. Carts are associated with a specific branch and user, and can optionally be linked to a customer for personalized pricing.

---

## ✅ Acceptance Criteria

### Scenario 1: Create a sale cart
- Given that I am logged in as a cashier
- When I create a new cart with type=SALE
- Then the cart must be created with:
  - Status: OPEN
  - User from JWT
  - Branch from JWT
  - Empty items array
  - Zero totals

### Scenario 2: Create a quotation cart
- Given that I am logged in as a cashier
- When I create a new cart with type=QUOTATION and valid_until
- Then the cart must be created with:
  - Status: OPEN
  - Cart type: QUOTATION
  - Expiration date set

### Scenario 3: Create cart with customer
- Given that I am logged in as a cashier
- When I create a cart linked to an existing customer
- Then the cart must include customer pricing rules

---

## ❌ Error Cases

- Invalid cart_type returns error 400
- Non-existent customer_id returns 404
- Cart creation fails if user has open cart (existing cart returned)

---

## 🔐 Business Rules

- One active cart per user at a time (optional enforcement)
- Cart expires after 24 hours by default
- Cart totals start at zero
- Customer pricing applies when customer_id is provided

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/carts
```

### Method: POST

### Request
```json
{
  "cart_type": "SALE",
  "customer_id": 123,
  "valid_until": "2026-03-30T23:59:59Z",
  "notes": "Customer request for express service"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "cart_type": "SALE",
    "customer_id": 123,
    "customer_name": "Acme Corporation",
    "branch_id": 1,
    "user_id": 5,
    "status": "OPEN",
    "subtotal": 0,
    "discount_total": 0,
    "tax_total": 0,
    "grand_total": 0,
    "items_count": 0,
    "valid_until": null,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true,
  "message": "Cart created successfully"
}
```

---

## 🧪 Testing Criteria

### Unit Tests
- Test cart type validation
- Test default values
- Test expiration logic

### Integration Tests
- Test tenant isolation
- Test user context
- Test customer linking

---

## 📎 Dependencies

- EPICA-005: Cart Epic
- Existing third parties module
- Existing user authentication
