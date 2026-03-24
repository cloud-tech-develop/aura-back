# HU-INV-004 - Low Stock Alerts

## 📌 General Information
- ID: HU-INV-004
- Epic: EPICA-004
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** receive low stock alerts
**So that** I can take action before stockouts occur

---

## 🧠 Functional Description

The system must identify products where current stock is at or below minimum stock level and provide alerting functionality.

---

## ✅ Acceptance Criteria

### Scenario 1: View all low stock items
- Given that products exist with stock at or below min_stock
- When I query low stock alerts
- Then all such products are returned
- With current stock vs min_stock comparison

### Scenario 2: Get critical stock items
- Given that inventory exists
- When I filter by critical level
- Then products at 50% of min_stock are highlighted

---

## 🔐 Business Rules

- Alert threshold: quantity <= min_stock
- Critical threshold: quantity <= min_stock * 0.5
- Branch-scoped results
- Excludes discontinued products

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/inventory/alerts
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| branch_id | int | No | Filter by branch |
| alert_type | string | No | all, low, critical |
| category_id | int | No | Filter by category |

### Response (200 OK)
```json
{
  "data": {
    "items": [
      {
        "product_id": 101,
        "product_name": "Wireless Mouse",
        "product_sku": "WM-001",
        "branch_id": 1,
        "branch_name": "Main Store",
        "category_name": "Electronics",
        "current_stock": 8,
        "min_stock": 10,
        "alert_type": "LOW",
        "deficit": 2
      },
      {
        "product_id": 102,
        "product_name": "USB Cable",
        "product_sku": "USB-001",
        "branch_id": 1,
        "current_stock": 3,
        "min_stock": 10,
        "alert_type": "CRITICAL",
        "deficit": 7
      }
    ],
    "summary": {
      "total_alerts": 15,
      "low_count": 10,
      "critical_count": 5
    }
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-004: Inventory Epic
