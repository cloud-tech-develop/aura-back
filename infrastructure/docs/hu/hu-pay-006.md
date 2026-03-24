# HU-PAY-006 - View Payment History

## 📌 General Information
- ID: HU-PAY-006
- Epic: EPICA-006
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view payment history for an order
**So that** I can track payment records and troubleshoot issues

---

## 🧠 Functional Description

The system must provide payment history for any reference (sale, purchase, account).

---

## ✅ Acceptance Criteria

### Scenario 1: View order payments
- Given that an order has multiple payments
- When I query payment history
- Then all payments are listed with details

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/payments/order/{orderId}
```

### Response (200 OK)
```json
{
  "data": {
    "order_id": 101,
    "order_total": 500000,
    "total_paid": 350000,
    "remaining": 150000,
    "status": "PARTIAL",
    "payments": [
      {
        "id": 1,
        "payment_method": "CASH",
        "amount": 200000,
        "status": "COMPLETED",
        "shift_id": 5,
        "user_name": "John Doe",
        "created_at": "2026-03-20T10:00:00Z"
      },
      {
        "id": 2,
        "payment_method": "CARD",
        "amount": 150000,
        "card_last_digits": "4567",
        "authorization_code": "AUTH123",
        "status": "COMPLETED",
        "shift_id": 7,
        "user_name": "Jane Smith",
        "created_at": "2026-03-22T14:30:00Z"
      }
    ]
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-006: Payments Epic
