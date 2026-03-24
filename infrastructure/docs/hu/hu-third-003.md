# HU-THIRD-003 - Get Third Party Details

## 📌 General Information
- ID: HU-THIRD-003
- Epic: EPICA-003
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** view third party details
**So that** I can see complete information about a business relationship

---

## 🧠 Functional Description

The system must return complete third party information including all contact details, classification, and metadata.

---

## ✅ Acceptance Criteria

### Scenario 1: Get existing third party
- Given that a third party exists
- When I request details by ID
- Then complete information is returned

### Scenario 2: Get non-existent third party
- Given that no third party exists with the ID
- When I request details
- Then 404 error is returned

---

## 🔐 Business Rules

- Returns 404 for non-existent records
- Soft-deleted records not accessible
- All tenant-scoped data

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/third-parties/{id}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "document_type": "NIT",
    "document_number": "901234567-8",
    "dv": "8",
    "name": "Acme Corporation S.A.S.",
    "trade_name": "Acme",
    "email": "contact@acme.com",
    "phone": "+573001234567",
    "mobile": "+573009876543",
    "address": "Calle 100 #45-67",
    "city": "Bogotá",
    "department": "Cundinamarca",
    "country": "COLOMBIA",
    "third_party_type": "CLIENT",
    "status": "ACTIVE",
    "is_taxpayer": true,
    "tax_regime": "RESPONSIBLE",
    "credit_limit": 5000000,
    "credit_days": 30,
    "customer_rating": 5,
    "observations": "Preferred customer",
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:00:00Z"
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-003: Third Parties Epic
