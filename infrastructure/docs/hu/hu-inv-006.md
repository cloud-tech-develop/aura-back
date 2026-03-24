# HU-INV-006 - Batch Tracking

## 📌 General Information
- ID: HU-INV-006
- Epic: EPICA-004
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** track products by batch
**So that** I can manage expiration dates and batch-specific inventory

---

## 🧠 Functional Description

The system must support batch tracking for products that require it, including batch numbers, expiration dates, and FIFO/LIFO selection.

---

## ✅ Acceptance Criteria

### Scenario 1: Register batch on entry
- Given that a product supports batch tracking
- When I register an entry with batch_number and expiration_date
- Then inventory is tracked per batch

### Scenario 2: View batch inventory
- Given that batch-tracked inventory exists
- When I query batch inventory for a product
- Then I see stock per batch with expiration dates

### Scenario 3: FIFO/LIFO selection on exit
- Given that multiple batches exist for a product
- When I register an exit and specify FIFO
- Then the oldest batch is used first

---

## 🔐 Business Rules

- Batch products require batch_number on entry
- Expiration tracking required for batch products
- FIFO (First In First Out) is default selection method
- Expired batches flagged in queries

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/inventory/{productId}/batches
```

### Response (200 OK)
```json
{
  "data": {
    "product_id": 101,
    "product_name": "Wireless Mouse",
    "batches": [
      {
        "batch_number": "LOT-2026-001",
        "quantity": 50,
        "expiration_date": "2027-03-23",
        "days_until_expiry": 365,
        "status": "ACTIVE"
      },
      {
        "batch_number": "LOT-2025-002",
        "quantity": 25,
        "expiration_date": "2026-04-15",
        "days_until_expiry": 23,
        "status": "EXPIRING_SOON"
      }
    ]
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-004: Inventory Epic
