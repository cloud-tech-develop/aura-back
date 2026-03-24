# HU-DASH-005 - Filter Dashboard by Period

## 📌 General Information
- ID: HU-DASH-005
- Epic: EPICA-014
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** business owner
**I want to** filter dashboard data by period
**So that** I can analyze different time frames

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/dashboard?period=week&branch_id=1
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| period | string | No | today, week, month, custom |
| start_date | date | No | Custom start date |
| end_date | date | No | Custom end date |
| branch_id | int | No | Filter by branch |

### Response (200 OK)
```json
{
  "data": {
    "period": {
      "type": "week",
      "start_date": "2026-03-16",
      "end_date": "2026-03-23"
    },
    "branch": {
      "id": 1,
      "name": "Main Store"
    },
    "sales_summary": {...},
    "inventory_alerts": {...},
    "top_products": {...},
    "recent_activity": {...}
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-014: Dashboard Epic
