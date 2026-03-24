# HU-CART-007 - Convert Cart to Quotation

## 📌 General Information
- ID: HU-CART-007
- Epic: EPICA-005
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** convert a cart to a quotation
**So that** I can create a quote for the customer to accept later

---

## 🧠 Functional Description

The system must convert the cart to a quotation with validity period, preserving all pricing and details.

---

## ✅ Acceptance Criteria

### Scenario 1: Convert to quotation
- Given that a cart has items
- When I convert it to quotation with valid_until date
- Then:
  - Cart status changes to CONVERTED
  - Reference to quotation created
  - valid_until is set

---

## ❌ Error Cases

- Empty cart returns 400
- Cart already converted returns 400
- Cart not OPEN returns 400

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/carts/{cartId}/convert-to-quotation
```

### Request
```json
{
  "valid_until": "2026-03-30T23:59:59Z",
  "notes": "Quote valid for 7 days"
}
```

### Response (201 Created)
```json
{
  "data": {
    "cart_id": 1,
    "quotation_id": 1,
    "quotation_number": "QT-2026-0001",
    "status": "CONVERTED",
    "valid_until": "2026-03-30T23:59:59Z",
    "grand_total": 172550,
    "converted_at": "2026-03-23T10:30:00Z"
  },
  "success": true,
  "message": "Cart converted to quotation"
}
```

---

## 📎 Dependencies

- EPICA-005: Cart Epic
