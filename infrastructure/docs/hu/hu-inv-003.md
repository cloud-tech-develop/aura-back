# HU-INV-003 - View Product Kardex

## 📌 General Information
- ID: HU-INV-003
- Epic: EPICA-004
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view the kardex (stock ledger) for a product
**So that** I can see complete movement history and stock tracking

---

## 🧠 Functional Description

The system must provide a complete movement history for each product, showing all entries and exits with dates, reasons, and balances.

---

## ✅ Acceptance Criteria

### Scenario 1: View kardex for a product
- Given that inventory movements exist for a product
- When I request the kardex
- Then all movements are returned in chronological order
- With previous and new balance for each movement

### Scenario 2: Filter kardex by date range
- Given that movements exist
- When I filter by date range
- Then only movements within the range are returned

### Scenario 3: Filter by movement type
- Given that movements exist
- When I filter by type=EXIT
- Then only exit movements are returned

---

## 🔐 Business Rules

- Movements ordered by created_at descending
- Date range maximum: 365 days
- All tenant-scoped data
- Soft-deleted movements excluded

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/inventory/{productId}/kardex
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | date | No | Start date filter |
| end_date | date | No | End date filter |
| movement_type | string | No | ENTRY, EXIT, ADJUSTMENT |
| movement_reason | string | No | SALE, PURCHASE, SHRINKAGE, etc. |
| page | int | No | Page number (default: 1) |
| limit | int | No | Items per page (default: 50) |

### Response (200 OK)
```json
{
  "data": {
    "product_id": 101,
    "product_name": "Wireless Mouse",
    "current_balance": 150,
    "items": [
      {
        "id": 3,
        "movement_type": "EXIT",
        "movement_reason": "SALE",
        "quantity": 2,
        "previous_balance": 152,
        "new_balance": 150,
        "reference_type": "sales_order",
        "reference_id": 301,
        "user_name": "John Doe",
        "created_at": "2026-03-23T14:30:00Z"
      },
      {
        "id": 2,
        "movement_type": "ENTRY",
        "movement_reason": "PURCHASE",
        "quantity": 100,
        "previous_balance": 52,
        "new_balance": 152,
        "batch_number": "LOT-2026-001",
        "reference_type": "purchase",
        "reference_id": 201,
        "user_name": "Jane Smith",
        "created_at": "2026-03-22T09:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 50,
      "total_items": 45
    }
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-004: Inventory Epic
