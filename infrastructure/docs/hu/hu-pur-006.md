# HU-PUR-006 - Supplier Account Summary

## 📌 General Information
- ID: HU-PUR-006
- Epic: EPICA-010
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** accountant
**I want to** view supplier account summary
**So that** I can manage accounts payable per supplier

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/suppliers/{supplierId}/account-summary
```

### Response (200 OK)
```json
{
  "data": {
    "supplier_id": 50,
    "supplier_name": "Global Supplies Inc.",
    "total_purchases": 15000000,
    "total_paid": 10000000,
    "total_pending": 5000000,
    "last_purchase_date": "2026-03-23",
    "purchases_count": 15
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-010: Purchases Epic
