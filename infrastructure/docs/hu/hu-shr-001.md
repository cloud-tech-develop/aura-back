# HU-SHR-001 - Register Shrinkage

## 📌 General Information
- ID: HU-SHR-001
- Epic: EPICA-011
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** register shrinkage (mermas)
**So that** I can document inventory losses

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/shrinkage
```

### Request
```json
{
  "branch_id": 1,
  "reason_id": 1,
  "shrinkage_date": "2026-03-23",
  "items": [
    {
      "product_id": 101,
      "quantity": 5,
      "unit_cost": 30000,
      "reason_detail": "Damaged packaging"
    }
  ],
  "notes": "Found damaged products in storage"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "shrinkage_number": "SHR-2026-0001",
    "status": "PENDING",
    "total_value": 150000,
    "inventory_updated": true,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-011: Shrinkage Epic
