# HU-DASH-003 - View Top Products

## 📌 General Information
- ID: HU-DASH-003
- Epic: EPICA-014
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** business owner
**I want to** view top selling products on dashboard
**So that** I can identify best performers

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/dashboard/top-products
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| period | string | No | today, week, month (default: today) |
| limit | int | No | Number of products (default: 10) |

### Response (200 OK)
```json
{
  "data": {
    "period": "today",
    "products": [
      {
        "rank": 1,
        "product_id": 101,
        "product_name": "Wireless Mouse",
        "sku": "WM-001",
        "units_sold": 25,
        "revenue": 1250000
      },
      {
        "rank": 2,
        "product_id": 103,
        "product_name": "Keyboard",
        "sku": "KB-001",
        "units_sold": 18,
        "revenue": 1440000
      }
    ]
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-014: Dashboard Epic
