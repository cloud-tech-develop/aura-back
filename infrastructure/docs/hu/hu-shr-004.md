# HU-SHR-004 - View Shrinkage Report

## 📌 General Information
- ID: HU-SHR-004
- Epic: EPICA-011
- Priority: Low
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view shrinkage reports
**So that** I can analyze inventory losses

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/reports/shrinkage
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | date | No | Start date |
| end_date | date | No | End date |
| reason_id | int | No | Filter by reason |

### Response (200 OK)
```json
{
  "data": {
    "summary": {
      "total_shrinkage_value": 2500000,
      "total_items": 50,
      "most_common_reason": "Damaged Products"
    },
    "items": [
      {
        "id": 1,
        "shrinkage_number": "SHR-2026-0001",
        "shrinkage_date": "2026-03-23",
        "reason_name": "Damaged Products",
        "total_value": 150000,
        "items_count": 5,
        "status": "APPROVED"
      }
    ],
    "pagination": {...}
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-011: Shrinkage Epic
