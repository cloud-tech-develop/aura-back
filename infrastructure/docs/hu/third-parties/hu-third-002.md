# HU-THIRD-002 - List Third Parties

## 📌 General Information
- ID: HU-THIRD-002
- Epic: EPICA-003
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** list third parties with pagination and filters
**So that** I can find and manage business relationships efficiently

---

## 🧠 Functional Description

The system must provide paginated listing of third parties with filtering by type, status, and search criteria. All results are scoped to the tenant.

---

## ✅ Acceptance Criteria

### Scenario 1: List all third parties
- Given that third parties exist
- When I request the list without filters
- Then I must receive paginated results with pagination metadata

### Scenario 2: Filter by type
- Given that third parties exist
- When I filter by type=CLIENT
- Then only clients are returned

### Scenario 3: Search by name
- Given that third parties exist
- When I search by "acme"
- Then third parties with "acme" in name are returned

---

## 🔐 Business Rules

- Results filtered by tenant from JWT
- Soft-deleted records excluded
- Default pagination: page=1, limit=20

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/third-parties
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| type | string | No | Filter: CLIENT, SUPPLIER, EMPLOYEE |
| status | string | No | Filter: ACTIVE, INACTIVE |
| search | string | No | Search by name or document |
| page | int | No | Page number (default: 1) |
| limit | int | No | Items per page (default: 20) |

### Response (200 OK)
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "document_type": "NIT",
        "document_number": "901234567-8",
        "name": "Acme Corporation",
        "email": "contact@acme.com",
        "phone": "+573001234567",
        "third_party_type": "CLIENT",
        "status": "ACTIVE",
        "created_at": "2026-03-23T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 20,
      "total_items": 150,
      "total_pages": 8
    }
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-003: Third Parties Epic
