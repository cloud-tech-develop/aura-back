# HU-CART-006 - Calculate Totals and Taxes

## 📌 General Information
- ID: HU-CART-006
- Epic: EPICA-005
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** see accurate totals with taxes calculated
**So that** I can confirm the final sale amount with the customer

---

## 🧠 Functional Description

The system must calculate subtotal, discounts, taxes (IVA, RETEICA, etc.), and grand total for the cart.

---

## ✅ Acceptance Criteria

### Scenario 1: View cart totals
- Given that a cart has items
- When I request cart details
- Then I see:
  - subtotal (sum of line totals before discounts)
  - discount_total
  - tax_exempt
  - taxable_amount
  - iva_19 (19% tax)
  - iva_5 (5% tax)
  - reteica
  - grand_total

### Scenario 2: Totals update on changes
- Given that a cart has items
- When I add/remove items or apply discounts
- Then all totals are recalculated

---

## 🔐 Business Rules

- Tax rates: 19% (general), 5% (some products), 0% (exempt)
- RETEICA withholding may apply
- All calculations in COP
- Totals recalculated on any cart modification

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/carts/{cartId}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "cart_type": "SALE",
    "customer": {
      "id": 123,
      "name": "Acme Corporation"
    },
    "items": [
      {
        "id": 1,
        "product_name": "Wireless Mouse",
        "quantity": 2,
        "unit_price": 50000,
        "discount_amount": 0,
        "tax_rate": 19,
        "tax_amount": 19000,
        "line_total": 119000
      },
      {
        "id": 2,
        "product_name": "USB Cable",
        "quantity": 3,
        "unit_price": 15000,
        "discount_amount": 0,
        "tax_rate": 19,
        "tax_amount": 8550,
        "line_total": 35850
      }
    ],
    "totals": {
      "subtotal": 145000,
      "discount_total": 0,
      "tax_exempt": 0,
      "taxable_amount": 145000,
      "iva_19": 27550,
      "iva_5": 0,
      "reteica": 0,
      "retefuente": 0,
      "grand_total": 172550
    },
    "status": "OPEN"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-005: Cart Epic
