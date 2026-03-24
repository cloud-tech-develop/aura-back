# HU-SHR-002 - Configure Shrinkage Reasons

## 📌 General Information
- ID: HU-SHR-002
- Epic: EPICA-011
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** configure shrinkage reasons
**So that** I can classify inventory losses

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/shrinkage-reasons
```

### Request
```json
{
  "code": "DAMAGED",
  "name": "Damaged Products",
  "description": "Products damaged during handling",
  "requires_authorization": true,
  "authorization_threshold": 500000
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "code": "DAMAGED",
    "name": "Damaged Products",
    "requires_authorization": true,
    "authorization_threshold": 500000,
    "is_active": true
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-011: Shrinkage Epic
