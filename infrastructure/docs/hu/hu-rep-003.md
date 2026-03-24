# HU-REP-003 - Sales Report by Employee

## 📌 General Information
- ID: HU-REP-003
- Epic: EPICA-008
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view sales by employee
**So that** I can track employee performance and calculate commissions

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/reports/sales/employees
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | date | Yes | Start date |
| end_date | date | Yes | End date |
| branch_id | int | No | Filter by branch |

### Response (200 OK)
```json
{
  "data": {
    "employees": [
      {
        "employee_id": 10,
        "employee_name": "John Doe",
        "user_name": "john.doe",
        "total_sales": 25,
        "total_revenue": 7500000,
        "average_ticket": 300000,
        "commission_earned": 375000
      }
    ],
    "summary": {
      "total_employees": 5,
      "total_revenue": 45000000
    }
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-008: Reports Epic
