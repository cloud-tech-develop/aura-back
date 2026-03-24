# HU-CASH-001 - Configure Cash Drawer

## 📌 General Information
- ID: HU-CASH-001
- Epic: EPICA-009
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** configure cash drawers for branches
**So that** I can set up cash management for operations

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/cash-drawers
```

### Request
```json
{
  "branch_id": 1,
  "name": "MAIN",
  "min_float": 100000
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "branch_id": 1,
    "branch_name": "Main Store",
    "name": "MAIN",
    "min_float": 100000,
    "is_active": true,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-009: Cash Drawer Epic
