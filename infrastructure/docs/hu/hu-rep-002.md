# HU-REP-002 - Sales Report by Product

## 📌 General Information
- ID: HU-REP-002
- Epic: EPICA-008
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view sales breakdown by product
**So that** I can identify top sellers and underperformers

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/reports/sales/products
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | date | Yes | Start date |
| end_date | date | Yes | End date |
| category_id | int | No | Filter by category |
| branch_id | int | No | Filter by branch |
| sort_by | string | No | units, revenue, margin |
| page | int | No | Page number |

### Response (200 OK)
```json
{
  "data": {
    "summary": {
      "total_products_sold": 1500,
      "total_revenue": 45000000
    },
    "products": [
      {
        "product_id": 101,
        "product_name": "Wireless Mouse",
        "product_sku": "WM-001",
        "category_name": "Electronics",
        "units_sold": 150,
        "revenue": 7500000,
        "cost": 4500000,
        "profit": 3000000,
        "margin_percentage": 40
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
