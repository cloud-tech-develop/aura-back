# HU-COMM-005 - Commission Reports

## 📌 General Information
- ID: HU-COMM-005
- Epic: EPICA-013
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** generate commission reports
**So that** I can analyze commission payouts

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/reports/commissions
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | date | No | Start date |
| end_date | date | No | End date |
| period | string | No | Settlement period |

### Response (200 OK)
```json
{
  "data": {
    "summary": {
      "total_commissions": 300000,
      "pending": 75000,
      "settled": 225000
    },
    "by_employee": [
      {
        "employee_id": 10,
        "employee_name": "John Doe",
        "total_commissions": 100000,
        "pending": 25000,
        "settled": 75000
      }
    ]
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-013: Commissions Epic
