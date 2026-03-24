# HU-INV-002 - Register Inventory Movement

## 📌 General Information
- ID: HU-INV-002
- Epic: EPICA-004
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** register inventory movements
**So that** I can track stock changes with proper documentation

---

## 🧠 Functional Description

The system must record all inventory movements with complete audit trail including movement type, reason, quantities, and user information.

---

## ✅ Acceptance Criteria

### Scenario 1: Register entry movement
- Given that inventory exists for a product
- When I register an ENTRY movement with:
  - Quantity: 100
  - Reason: PURCHASE
  - Reference: purchase_id
- Then inventory quantity increases
- And movement record is created with previous/new balance

### Scenario 2: Register exit movement
- Given that inventory has sufficient stock
- When I register an EXIT movement
- Then inventory quantity decreases
- And movement is recorded

### Scenario 3: Reject negative stock
- Given that inventory has 10 units
- When I try to register an EXIT of 15 units
- Then the movement is rejected with error 400

---

## ❌ Error Cases

- Insufficient stock returns error 400
- Invalid movement type returns error 400
- Invalid reason returns error 400
- Missing required fields returns validation error

---

## 🔐 Business Rules

- Movements are immutable after creation
- Stock cannot go negative
- Every movement requires: type, reason, quantity, user
- Reference links to source document when applicable

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/inventory/movements
```

### Request
```json
{
  "product_id": 101,
  "branch_id": 1,
  "movement_type": "ENTRY",
  "movement_reason": "PURCHASE",
  "quantity": 100,
  "batch_number": "LOT-2026-001",
  "expiration_date": "2027-03-23",
  "reference_id": 501,
  "reference_type": "purchase",
  "notes": "Received from supplier"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "inventory_id": 101,
    "product_name": "Wireless Mouse",
    "movement_type": "ENTRY",
    "movement_reason": "PURCHASE",
    "quantity": 100,
    "previous_balance": 50,
    "new_balance": 150,
    "batch_number": "LOT-2026-001",
    "reference_type": "purchase",
    "reference_id": 501,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true,
  "message": "Movement registered successfully"
}
```

---

## 📎 Dependencies

- EPICA-004: Inventory Epic
- Existing products module
