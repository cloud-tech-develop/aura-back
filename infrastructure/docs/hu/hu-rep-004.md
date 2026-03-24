# HU-REP-004 - Inventory Status Report

## 📌 General Information
- ID: HU-REP-004
- Epic: EPICA-008
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view inventory status report
**So that** I can monitor stock levels and valuation

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/reports/inventory
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| branch_id | int | No | Filter by branch |
| category_id | int | No | Filter by category |
| stock_filter | string | No | all, low, out, normal |

### Response (200 OK)
```json
{
  "data": {
    "summary": {
      "total_products": 500,
      "total_stock_value": 50000000,
      "low_stock_count": 25,
      "out_of_stock_count": 5
    },
    "products": [
      {
        "product_id": 101,
        "product_name": "Wireless Mouse",
        "sku": "WM-001",
        "category_name": "Electronics",
        "quantity": 150,
        "min_stock": 50,
        "unit_cost": 30000,
        "stock_value": 4500000,
        "status": "NORMAL"
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
