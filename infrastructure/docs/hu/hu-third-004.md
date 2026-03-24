# HU-THIRD-004 - Update Third Party

## 📌 General Information
- ID: HU-THIRD-004
- Epic: EPICA-003
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** update third party information
**So that** I can maintain accurate business relationship data

---

## 🧠 Functional Description

The system must allow updating third party information. Only provided fields are updated (partial update). Document number changes are not allowed to prevent data integrity issues.

---

## ✅ Acceptance Criteria

### Scenario 1: Update contact information
- Given that a third party exists
- When I update phone and email
- Then only those fields are changed
- And updated_at is set to current timestamp

### Scenario 2: Update credit information
- Given that a client exists
- When I update credit_limit and credit_days
- Then the new values are saved

---

## ❌ Error Cases

- Changing document number returns error 400
- Invalid email format returns error 400
- Non-existent third party returns 404

---

## 🔐 Business Rules

- Partial update: only provided fields are modified
- Document number cannot be changed
- updated_at timestamp is automatically set

---

## 📡 Technical Requirements

### Endpoint
```
PUT /api/third-parties/{id}
```

### Request
```json
{
  "phone": "+573001111111",
  "email": "newemail@acme.com",
  "credit_limit": 10000000,
  "credit_days": 45
}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "document_type": "NIT",
    "document_number": "901234567-8",
    "name": "Acme Corporation",
    "phone": "+573001111111",
    "email": "newemail@acme.com",
    "credit_limit": 10000000,
    "credit_days": 45,
    "updated_at": "2026-03-23T14:00:00Z"
  },
  "success": true,
  "message": "Third party updated successfully"
}
```

---

## 📎 Dependencies

- EPICA-003: Third Parties Epic
