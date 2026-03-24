# HU-PAY-002 - Process Multiple Payment Methods

## 📌 General Information
- ID: HU-PAY-002
- Epic: EPICA-006
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** process multiple payment methods for a single sale
**So that** customers can combine payment options

---

## 🧠 Functional Description

The system must allow splitting payments across multiple methods (e.g., cash + card) as long as total equals or exceeds sale amount.

---

## ✅ Acceptance Criteria

### Scenario 1: Split payment cash + card
- Given that a sales order has total 150000
- When I process:
  - Cash: 100000
  - Card: 50000
- Then both payments are recorded
- And total equals sale total

### Scenario 2: Combined payments exceed total
- Given that a sales order has total 100000
- When I process:
  - Cash: 60000
  - Card: 50000
- Then payments recorded with credit of 10000

---

## ❌ Error Cases

- Total payments less than order total returns 400
- Invalid payment method returns 400
- Missing shift returns 400

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/payments/split
```

### Request
```json
{
  "payment_type": "SALE",
  "reference_id": 101,
  "reference_type": "sales_order",
  "shift_id": 5,
  "payments": [
    {
      "payment_method": "CASH",
      "amount": 100000
    },
    {
      "payment_method": "CREDIT_CARD",
      "amount": 50000,
      "authorization_code": "AUTH123456",
      "card_last_digits": "4567"
    }
  ],
  "notes": "Split payment"
}
```

### Response (201 Created)
```json
{
  "data": {
    "payments": [
      {
        "id": 1,
        "payment_method": "CASH",
        "amount": 100000,
        "change_amount": 0,
        "status": "COMPLETED"
      },
      {
        "id": 2,
        "payment_method": "CREDIT_CARD",
        "amount": 50000,
        "status": "COMPLETED"
      }
    ],
    "total_paid": 150000,
    "order_total": 150000,
    "remaining": 0
  },
  "success": true,
  "message": "Split payment processed"
}
```

---

## 📎 Dependencies

- EPICA-006: Payments Epic
- HU-PAY-001: Process Single Payment
