# HU-INV-005 - Inventory Consolidation

## 📌 General Information
- ID: HU-INV-005
- Epic: EPICA-004
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** business owner
**I want to** view consolidated inventory across all branches
**So that** I can understand total stock availability for the company

---

## 🧠 Functional Description

The system must aggregate inventory data across all branches of the company, showing total stock levels per product.

---

## ✅ Acceptance Criteria

### Scenario 1: View consolidated inventory
- Given that inventory exists across multiple branches
- When I request consolidated inventory
- Then products are listed with total stock across all branches
- And branch-level breakdown is available

### Scenario 2: Filter consolidated by category
- Given that consolidated inventory exists
- When I filter by category
- Then only products in that category are shown

---

## 🔐 Business Rules

- Totals include all active branches
- Excludes soft-deleted products
- All tenant-scoped data

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/inventory/consolidated
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| category_id | int | No | Filter by category |
| stock_filter | string | No | all, low, out |
| page | int | No | Page number |
| limit | int | No | Items per page |

### Response (200 OK)
```json
{
  "data": {
    "items": [
      {
        "product_id": 101,
        "product_name": "Wireless Mouse",
        "product_sku": "WM-001",
        "category_name": "Electronics",
        "total_stock": 250,
        "total_value": 12500000,
        "low_stock": false,
        "branch_breakdown": [
          {
            "branch_id": 1,
            "branch_name": "Main Store",
            "quantity": 150
          },
          {
            "branch_id": 2,
            "branch_name": "Branch 2",
            "quantity": 100
          }
        ]
      }
    ],
    "summary": {
      "total_products": 500,
      "total_stock_value": 150000000,
      "products_low_stock": 25
    },
    "pagination": {...}
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-004: Inventory Epic
- Existing branch module
