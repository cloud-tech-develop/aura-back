# HU-THIRD-006 - Search Third Parties

## 📌 General Information
- ID: HU-THIRD-006
- Epic: EPICA-003
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** quickly search third parties by name, document, or email
**So that** I can find customers or suppliers during sales or purchases

---

## 🧠 Functional Description

The system must provide fast search functionality for third parties with autocomplete-style results. Search returns matches from multiple fields.

---

## ✅ Acceptance Criteria

### Scenario 1: Search by name
- Given that third parties exist
- When I search for "acme"
- Then all third parties containing "acme" in name are returned

### Scenario 2: Search by document
- Given that third parties exist
- When I search by document number
- Then exact match is returned

### Scenario 3: Search by email
- Given that third parties exist
- When I search by email
- Then matching results are returned

---

## 🔐 Business Rules

- Only ACTIVE third parties in search results
- Minimum 3 characters for search
- Results limited to 20 items
- Search is case-insensitive

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/third-parties/search
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| q | string | Yes | Search query (min 3 chars) |
| type | string | No | Filter by type |
| limit | int | No | Max results (default: 20) |

### Response (200 OK)
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "document_type": "NIT",
        "document_number": "901234567-8",
        "name": "Acme Corporation S.A.S.",
        "email": "contact@acme.com",
        "phone": "+573001234567",
        "third_party_type": "CLIENT"
      }
    ],
    "count": 1
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-003: Third Parties Epic
