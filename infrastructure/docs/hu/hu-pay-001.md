# HU-PAY-001 - Process Single Payment

## 📌 General Information
- ID: HU-PAY-001
- Epic: EPICA-006
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** process a payment for a sale
**So that** the transaction can be completed

---

## 🧠 Functional Description

The system must process single payments for sales orders, supporting multiple payment methods (cash, card, transfer). All payments are linked to a shift and require proper documentation.

---

## ✅ Acceptance Criteria

### Scenario 1: Process cash payment with exact amount
- Given that a sales order exists with total=100000
- When I process a cash payment of 100000
- Then the payment is recorded with:
  - Status: COMPLETED
  - Change: 0

### Scenario 2: Process cash payment with change
- Given that a sales order exists with total=85000
- When I process a cash payment of 100000
- Then the payment is recorded with:
  - Amount: 100000
  - Change: 15000
  - Status: COMPLETED

### Scenario 3: Process card payment
- Given that a sales order exists
- When I process a card payment with:
  - Method: CREDIT_CARD
  - Authorization code: ABC123
  - Card last digits: 4567
- Then the payment is recorded successfully

---

## ❌ Error Cases

- Payment amount less than order total returns error 400 (for non-credit)
- Invalid payment method returns error 400
- Missing shift returns error 400
- Closed shift returns error 400

---

## 🔐 Business Rules

- Cash payments require change calculation
- Card payments require authorization code
- All payments linked to active shift
- Credit payments create accounts receivable
- Payments cannot exceed configurable limits without authorization

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/payments
```

### Method: POST

### Request
```json
{
  "payment_type": "SALE",
  "reference_id": 101,
  "reference_type": "sales_order",
  "payment_method": "CASH",
  "amount": 100000,
  "shift_id": 5,
  "reference_number": null,
  "bank_name": null,
  "card_last_digits": null,
  "authorization_code": null,
  "notes": "Customer paid in cash"
}
```

### Response (201 Created)
```json
{
  "data": {
    "id": 1,
    "payment_type": "SALE",
    "reference_id": 101,
    "reference_type": "sales_order",
    "payment_method": "CASH",
    "amount": 100000,
    "change_amount": 0,
    "shift_id": 5,
    "status": "COMPLETED",
    "created_at": "2026-03-23T10:30:00Z"
  },
  "success": true,
  "message": "Payment processed successfully"
}
```

### Request (Card Payment)
```json
{
  "payment_type": "SALE",
  "reference_id": 101,
  "reference_type": "sales_order",
  "payment_method": "CREDIT_CARD",
  "amount": 100000,
  "shift_id": 5,
  "reference_number": "TXN-456789",
  "bank_name": "Bancolombia",
  "card_last_digits": "4567",
  "authorization_code": "AUTH123456"
}
```

---

## 🧪 Testing Criteria

### Unit Tests
- Test change calculation
- Test payment method validation
- Test amount validation

### Integration Tests
- Test shift integration
- Test order status update
- Test payment recording

---

## 📎 Dependencies

- EPICA-006: Payments Epic
- Existing sales orders module
- Existing cash drawer module
