# HU-COMM-002 - Calculate Commissions on Sale

## 📌 General Information
- ID: HU-COMM-002
- Epic: EPICA-013
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** system
**I want to** automatically calculate commissions when a sale is completed
**So that** employee commissions are tracked

---

## 📡 Technical Requirements

This is an automated process triggered by sale completion.

### Event: sale.completed
```json
{
  "event": "sale.completed",
  "data": {
    "sales_order_id": 101,
    "user_id": 5,
    "total": 150000,
    "profit_margin": 50000
  }
}
```

### Response
```json
{
  "commissions_created": [
    {
      "employee_id": 10,
      "commission_type": "PERCENTAGE_SALE",
      "commission_rate": 5,
      "commission_amount": 7500,
      "status": "PENDING"
    }
  ]
}
```

---

## 📎 Dependencies

- EPICA-013: Commissions Epic
