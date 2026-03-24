# HU-DASH-004 - View Recent Activity

## 📌 General Information
- ID: HU-DASH-004
- Epic: EPICA-014
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view recent activity on dashboard
**So that** I can stay informed of latest transactions

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/dashboard/recent-activity
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| limit | int | No | Number of activities (default: 10) |

### Response (200 OK)
```json
{
  "data": {
    "activities": [
      {
        "id": 1,
        "type": "SALE",
        "description": "Sale #SO-2026-0045 completed",
        "amount": 150000,
        "user_name": "John Doe",
        "created_at": "2026-03-23T14:30:00Z"
      },
      {
        "id": 2,
        "type": "PURCHASE",
        "description": "Purchase #PUR-2026-0012 received",
        "amount": 2500000,
        "user_name": "Jane Smith",
        "created_at": "2026-03-23T14:00:00Z"
      }
    ]
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-014: Dashboard Epic
