# HU-PAY-005 - Cancel Payment

## 📌 General Information
- ID: HU-PAY-005
- Epic: EPICA-006
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** cancel payments
**So that** I can correct mistakes or process refunds

---

## 🧠 Functional Description

The system must allow payment cancellation with proper authorization, updating the linked order status accordingly.

---

## ✅ Acceptance Criteria

### Scenario 1: Cancel completed payment
- Given that a payment exists
- When I cancel with reason
- Then:
  - Payment status = CANCELLED
  - Order remaining updated
  - Cancellation logged

---

## ❌ Error Cases

- Already cancelled payment returns 400
- Missing cancellation reason returns 400
- Insufficient permissions returns 403

---

## 🔐 Business Rules

- Cancellation requires reason
- Manager authorization for cancellation
- Payment history preserved
- Order status updated accordingly

---

## 📡 Technical Requirements

### Endpoint
```
DELETE /api/payments/{paymentId}
```

### Request
```json
{
  "reason": "Customer returned product"
}
```

### Response (200 OK)
```json
{
  "data": {
    "payment_id": 1,
    "original_amount": 100000,
    "status": "CANCELLED",
    "cancelled_at": "2026-03-23T14:00:00Z",
    "cancelled_by": "Manager Name",
    "order": {
      "id": 101,
      "remaining": 50000,
      "status": "PARTIAL"
    }
  },
  "success": true,
  "message": "Payment cancelled"
}
```

---

## 📎 Dependencies

- EPICA-006: Payments Epic
