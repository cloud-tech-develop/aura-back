# HU-COMM-001 - Configure Commission Rules

## 📌 General Information
- ID: HU-COMM-001
- Epic: EPICA-013
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** configure commission rules
**So that** I can set up commission calculations for employees

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/commission-rules
```

### Request
```json
{
  "name": "Sales Commission - Standard",
  "commission_type": "PERCENTAGE_SALE",
  "employee_id": 10,
  "value": 5,
  "start_date": "2026-01-01",
  "end_date": "2026-12-31"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "name": "Sales Commission - Standard",
    "commission_type": "PERCENTAGE_SALE",
    "value": 5,
    "employee_name": "John Doe",
    "is_active": true,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-013: Commissions Epic
