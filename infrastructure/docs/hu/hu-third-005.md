# HU-THIRD-005 - Change Third Party Status

## 📌 General Information
- ID: HU-THIRD-005
- Epic: EPICA-003
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** change third party status (active/inactive)
**So that** I can manage business relationship availability

---

## 🧠 Functional Description

The system must allow changing third party status between ACTIVE and INACTIVE. Inactive third parties cannot be selected in new transactions but remain visible for historical records.

---

## ✅ Acceptance Criteria

### Scenario 1: Deactivate a client
- Given that a client exists with status ACTIVE
- When I change status to INACTIVE
- Then the client status is updated
- And the client cannot be selected in new sales

### Scenario 2: Reactivate a client
- Given that a client exists with status INACTIVE
- When I change status to ACTIVE
- Then the client status is updated
- And the client can be selected in new sales

---

## ❌ Error Cases

- Invalid status value returns error 400
- Non-existent third party returns 404

---

## 🔐 Business Rules

- Valid statuses: ACTIVE, INACTIVE
- Inactive third parties excluded from active lists
- Historical records remain accessible
- Soft delete is preferred over status change

---

## 📡 Technical Requirements

### Endpoint
```
PATCH /api/third-parties/{id}/status
```

### Request
```json
{
  "status": "INACTIVE"
}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "name": "Acme Corporation",
    "status": "INACTIVE",
    "updated_at": "2026-03-23T14:00:00Z"
  },
  "success": true,
  "message": "Status updated successfully"
}
```

---

## 📎 Dependencies

- EPICA-003: Third Parties Epic
