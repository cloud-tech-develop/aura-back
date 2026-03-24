# HU-INV-007 - Serialized Product Tracking

## 📌 General Information
- ID: HU-INV-007
- Epic: EPICA-004
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** track serialized products
**So that** I can manage individual unit tracking for warranty and tracking

---

## 🧠 Functional Description

The system must support serialized product tracking where each unit has a unique serial number, enabling individual unit management.

---

## ✅ Acceptance Criteria

### Scenario 1: Register serialized entry
- Given that a product supports serialization
- When I register an entry with serial numbers
- Then each unit is tracked individually

### Scenario 2: Track serial during sale
- Given that serialized inventory exists
- When I sell a serialized product
- Then the specific serial number is recorded

### Scenario 3: View serial history
- Given that a serialized product was sold
- When I query serial history
- Then all movements for that serial are shown

---

## 🔐 Business Rules

- Serial numbers must be unique per product
- Each serial can only be in one place at a time
- Serial history is immutable
- Serial cannot be sold twice

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/inventory/{productId}/serials
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| status | string | No | AVAILABLE, SOLD, RESERVED |
| page | int | No | Page number |
| limit | int | No | Items per page |

### Response (200 OK)
```json
{
  "data": {
    "product_id": 101,
    "product_name": "Laptop Model X",
    "total_units": 25,
    "available": 20,
    "sold": 5,
    "serials": [
      {
        "serial_number": "SN-LAP-001",
        "status": "AVAILABLE",
        "purchase_date": "2026-03-01",
        "purchase_id": 401
      },
      {
        "serial_number": "SN-LAP-002",
        "status": "SOLD",
        "purchase_date": "2026-03-01",
        "sale_date": "2026-03-15",
        "sales_order_id": 301
      }
    ],
    "pagination": {...}
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-004: Inventory Epic
