# HU-PAY-004 - Record Partial Payment

## 📌 General Information
- ID: HU-PAY-004
- Epic: EPICA-006
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** record partial payments for credit sales
**So that** customers can pay in installments

---

## 🧠 Functional Description

The system must record partial payments against orders, track remaining balance, and update order status when fully paid.

---

## ✅ Acceptance Criteria

### Scenario 1: First partial payment
- Given that an order has total 500000
- When I record a payment of 200000
- Then remaining = 300000
- And order status = PARTIAL

### Scenario 2: Final payment completes order
- Given that an order has remaining 200000
- When I record a payment of 200000
- Then remaining = 0
- And order status = PAID

---

## 🔐 Business Rules

- Partial payments create/update accounts receivable
- Order status: PENDING_PAYMENT, PARTIAL, PAID
- Payment history maintained

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/payments
```

### Request
```json
{
  "payment_type": "SALE",
  "reference_id": 101,
  "reference_type": "sales_order",
  "payment_method": "CASH",
  "amount": 200000,
  "shift_id": 5,
  "notes": "First installment"
}
```

### Response (201 Created)
```json
{
  "data": {
    "payment_id": 1,
    "amount": 200000,
    "order_total": 500000,
    "total_paid": 200000,
    "remaining": 300000,
    "order_status": "PARTIAL"
  },
  "success": true,
  "message": "Partial payment recorded"
}
```

---

## 📎 Dependencies

- EPICA-006: Payments Epic
