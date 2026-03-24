# HU-PUR-005 - View Purchase History

## 📌 General Information
- ID: HU-PUR-005
- Epic: EPICA-010
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** accountant
**I want to** view purchase history
**So that** I can track supplier transactions

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/purchases
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| supplier_id | int | No | Filter by supplier |
| status | string | No | Filter by status |
| start_date | date | No | Start date |
| end_date | date | No | End date |

### Response (200 OK)
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "purchase_number": "PUR-2026-0001",
        "supplier_name": "Global Supplies Inc.",
        "purchase_date": "2026-03-23",
        "total": 3570000,
        "paid_amount": 2000000,
        "pending_amount": 1570000,
        "status": "PARTIAL"
      }
    ],
    "pagination": {...}
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-010: Purchases Epic
