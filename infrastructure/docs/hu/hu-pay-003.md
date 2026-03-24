# HU-PAY-003 - Calculate Change

## 📌 General Information
- ID: HU-PAY-003
- Epic: EPICA-006
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** calculate change when customer pays more than total
**So that** I can return the correct change amount

---

## 🧠 Functional Description

The system must calculate change amount when cash payment exceeds the sale total.

---

## ✅ Acceptance Criteria

### Scenario 1: Calculate change
- Given that a sales order has total 85000
- When customer pays 100000 cash
- Then change = 15000

### Scenario 2: No change needed
- Given that a sales order has total 100000
- When customer pays exactly 100000
- Then change = 0

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/payments/calculate-change
```

### Request
```json
{
  "order_total": 85000,
  "cash_received": 100000
}
```

### Response (200 OK)
```json
{
  "data": {
    "order_total": 85000,
    "cash_received": 100000,
    "change_amount": 15000,
    "can_complete": true
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-006: Payments Epic
