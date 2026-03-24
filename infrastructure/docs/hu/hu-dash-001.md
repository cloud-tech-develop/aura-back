# HU-DASH-001 - View Sales Summary

## 📌 General Information
- ID: HU-DASH-001
- Epic: EPICA-014
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** business owner
**I want to** view sales summary on dashboard
**So that** I can see business performance at a glance

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/dashboard/sales-summary
```

### Response (200 OK)
```json
{
  "data": {
    "today": {
      "total_sales": 2500000,
      "transaction_count": 45,
      "average_ticket": 55555
    },
    "comparison": {
      "vs_yesterday": {
        "amount": 2000000,
        "change_percentage": 25.0
      },
      "vs_last_week": {
        "amount": 15000000,
        "change_percentage": 10.5
      }
    },
    "by_payment_method": {
      "CASH": 1500000,
      "CREDIT_CARD": 750000,
      "DEBIT_CARD": 250000
    },
    "generated_at": "2026-03-23T14:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-014: Dashboard Epic
