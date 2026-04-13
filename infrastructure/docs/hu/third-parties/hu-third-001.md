# HU-THIRD-001 - Create Third Party

## 📌 General Information
- ID: HU-THIRD-001
- Epic: EPICA-003
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** create third parties (clients, suppliers, employees)
**So that** I can manage business relationships in the system

---

## 🧠 Functional Description

The system must allow creating third parties with complete information including document details, contact information, and classification. All third parties are stored in the tenant schema and must have unique document numbers per type.

---

## ✅ Acceptance Criteria

### Scenario 1: Create a new client
- Given that I am logged in as a manager
- When I create a new client with:
  - Document Type: NIT
  - Document Number: 901234567-8
  - Name: "Acme Corporation S.A.S."
  - Email: "contact@acme.com"
  - Phone: "+573001234567"
  - Third Party Type: CLIENT
- Then the third party must be saved with:
  - Unique document validation
  - Status: ACTIVE
  - Timestamps (created_at)
  - Tenant scope applied

### Scenario 2: Create a new supplier
- Given that I am logged in as a manager
- When I create a new supplier with:
  - Document Type: NIT
  - Document Number: 800123456-9
  - Name: "Global Supplies Inc."
  - Tax Regime: RESPONSIBLE
- Then the supplier must be saved successfully
- And the supplier can be used in purchase orders

### Scenario 3: Create an employee
- Given that I am logged in as an admin
- When I create a new employee linked to a user account
- Then the third party type must be EMPLOYEE
- And the user_id field must be populated

---

## ❌ Error Cases

- Duplicate document number must return error 400
- Invalid email format must return error 400
- Invalid document type must return error 400
- Missing required fields must return validation errors
- User already linked to another employee must return error 400

---

## 🔐 Business Rules

- Document numbers must be unique per tenant and type
- Email must follow RFC 5322 standard
- Third party type must be one of: CLIENT, SUPPLIER, EMPLOYEE
- Employees must be linked to a user account
- Soft delete preserves all data

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/third-parties
```

### Method: POST

### Request
```json
{
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
  "is_taxpayer": true,
  "tax_regime": "RESPONSIBLE",
  "credit_limit": 5000000,
  "credit_days": 30
}
```

### Response (201 Created)
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
    "third_party_type": "CLIENT",
    "status": "ACTIVE",
    "credit_limit": 5000000,
    "credit_days": 30,
    "created_at": "2026-03-23T10:00:00Z"
  },
  "success": true,
  "message": "Third party created successfully"
}
```

### Error Responses
- 400: Validation errors
- 409: Duplicate document number
- 500: Server error

---

## 🧪 Testing Criteria

### Unit Tests
- Test document validation
- Test email format validation
- Test unique constraint
- Test business rule enforcement

### Integration Tests
- Test complete creation flow
- Test tenant isolation
- Test concurrent creation

---

## 📎 Dependencies

- EPICA-003: Third Parties Epic
- Existing user authentication
- Tenant context from JWT
