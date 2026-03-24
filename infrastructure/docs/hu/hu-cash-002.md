# HU-CASH-002 - Open Cash Shift

## 📌 General Information
- ID: HU-CASH-002
- Epic: EPICA-009
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** open a cash shift
**So that** I can start processing sales and cash movements

---

## 🧠 Functional Description

The system must allow cashiers to open shifts with an initial cash amount (float). Each user can only have one active shift at a time. The shift tracks all cash movements until closure.

---

## ✅ Acceptance Criteria

### Scenario 1: Open shift with initial float
- Given that I am logged in as a cashier
- And I do not have an active shift
- When I open a shift with opening_amount=200000
- Then the shift is created with:
  - Status: OPEN
  - Opening amount: 200000
  - User from JWT
  - Current timestamp

### Scenario 2: Reject when shift already open
- Given that I have an active shift
- When I attempt to open another shift
- Then I must receive error 400
- And my existing shift is returned

---

## ❌ Error Cases

- User already has active shift returns error 400
- Invalid opening_amount (negative) returns error 400
- Missing branch_id returns error 400

---

## 🔐 Business Rules

- One active shift per user at a time
- Opening amount can be zero
- Shift is linked to branch from user context
- Opening timestamp is recorded automatically

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/cash-shifts
```

### Method: POST

### Request
```json
{
  "opening_amount": 200000,
  "notes": "Initial float for morning shift"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "cash_drawer_id": 1,
    "branch_id": 1,
    "branch_name": "Main Store",
    "user_id": 5,
    "user_name": "John Doe",
    "opening_amount": 200000,
    "status": "OPEN",
    "opened_at": "2026-03-23T08:00:00Z",
    "created_at": "2026-03-23T08:00:00Z"
  },
  "success": true,
  "message": "Shift opened successfully"
}
```

### Error Response (409 Conflict)
```json
{
  "error": "Active shift exists",
  "data": {
    "active_shift_id": 1,
    "opened_at": "2026-03-23T08:00:00Z"
  },
  "success": false
}
```

---

## 🧪 Testing Criteria

### Unit Tests
- Test shift uniqueness validation
- Test opening amount validation
- Test status transitions

### Integration Tests
- Test concurrent shift open attempts
- Test branch context
- Test user context

---

## 📎 Dependencies

- EPICA-009: Cash Drawer Epic
- Existing cash drawer configuration
- Existing user authentication
