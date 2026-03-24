# HU-COMM-003 - View Commission History

## 📌 General Information
- ID: HU-COMM-003
- Epic: EPICA-013
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view commission history for employees
**So that** I can track accumulated commissions

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/commissions
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| employee_id | int | No | Filter by employee |
| status | string | No | Filter by status |
| start_date | date | No | Start date |
| end_date | date | No | End date |

### Response (200 OK)
```json
{
  "data": {
    "summary": {
      "total_pending": 75000,
      "total_settled": 225000
    },
    "items": [
      {
        "id": 1,
        "sales_order_number": "SO-2026-0001",
        "employee_name": "John Doe",
        "sale_amount": 150000,
        "commission_type": "PERCENTAGE_SALE",
        "commission_rate": 5,
        "commission_amount": 7500,
        "status": "PENDING",
        "created_at": "2026-03-23T10:00:00Z"
      }
    ],
    "pagination": {...}
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-013: Commissions Epic
