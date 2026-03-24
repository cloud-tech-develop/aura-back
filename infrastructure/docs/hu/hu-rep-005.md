# HU-REP-005 - Movement History Report

## 📌 General Information
- ID: HU-REP-005
- Epic: EPICA-008
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view inventory movement history
**So that** I can analyze stock changes over time

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/reports/inventory/movements
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | date | Yes | Start date |
| end_date | date | Yes | End date |
| product_id | int | No | Filter by product |
| movement_type | string | No | ENTRY, EXIT |
| movement_reason | string | No | SALE, PURCHASE, etc. |

### Response (200 OK)
```json
{
  "data": {
    "summary": {
      "total_entries": 500,
      "total_exits": 450,
      "net_change": 50
    },
    "movements": [
      {
        "id": 1,
        "product_name": "Wireless Mouse",
        "movement_type": "ENTRY",
        "movement_reason": "PURCHASE",
        "quantity": 100,
        "user_name": "Jane Smith",
        "created_at": "2026-03-23T09:00:00Z"
      }
    ],
    "pagination": {...}
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-008: Reports Epic
