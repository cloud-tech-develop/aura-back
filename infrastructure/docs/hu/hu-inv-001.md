# HU-INV-001 - Query Branch Inventory

## 📌 General Information
- ID: HU-INV-001
- Epic: EPICA-004
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** query inventory by branch
**So that** I can see current stock levels and make informed decisions

---

## 🧠 Functional Description

The system must provide inventory queries filtered by branch, with pagination, sorting, and filtering capabilities. All inventory data is scoped to the tenant and shows real-time stock levels.

---

## ✅ Acceptance Criteria

### Scenario 1: Query all inventory for current branch
- Given that I am logged in as a manager
- When I query inventory for my branch
- Then I must receive paginated results with:
  - Product information
  - Current quantity
  - Reserved quantity
  - Available quantity
  - Minimum stock level

### Scenario 2: Filter by category
- Given that inventory exists
- When I filter by category_id
- Then only products in that category are returned

### Scenario 3: Filter low stock items
- Given that inventory exists
- When I filter by stock_filter=low
- Then only products where quantity <= min_stock are returned

### Scenario 4: Search by product name
- Given that inventory exists
- When I search by "wireless"
- Then products with "wireless" in name are returned

---

## ❌ Error Cases

- Invalid branch_id returns 404
- Invalid pagination parameters use defaults
- Deleted products are excluded from results

---

## 🔐 Business Rules

- All queries filter by tenant from JWT
- Only active branches are accessible
- Soft-deleted products are excluded
- Pagination default: page=1, limit=20

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/inventory
```

### Method: GET

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| branch_id | int | No | Filter by branch (defaults to user's branch) |
| category_id | int | No | Filter by category |
| stock_filter | string | No | Filter: all, low, out, normal |
| search | string | No | Search by product name |
| page | int | No | Page number (default: 1) |
| limit | int | No | Items per page (default: 20, max: 100) |
| sort | string | No | Sort field (default: product_name) |
| order | string | No | Sort order: asc, desc (default: asc) |

### Response (200 OK)
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "product_id": 101,
        "product_name": "Wireless Mouse",
        "product_sku": "WM-001",
        "branch_id": 1,
        "branch_name": "Main Store",
        "category_name": "Electronics",
        "quantity": 50,
        "reserved_quantity": 5,
        "available_quantity": 45,
        "min_stock": 10,
        "location": "A1-B2"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 20,
      "total_items": 150,
      "total_pages": 8
    }
  },
  "success": true
}
```

---

## 🧪 Testing Criteria

### Unit Tests
- Test pagination calculation
- Test stock filter logic
- Test search filtering

### Integration Tests
- Test tenant isolation
- Test branch access control
- Test concurrent inventory updates

---

## 📎 Dependencies

- EPICA-004: Inventory Epic
- Existing products module
- Existing branch module
