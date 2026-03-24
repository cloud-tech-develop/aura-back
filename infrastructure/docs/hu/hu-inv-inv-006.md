# HU-INV-INV-006 - Invoice Search and Filter

## 📌 General Information
- ID: HU-INV-INV-006
- Epic: EPICA-007
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** accountant
**I want to** search and filter invoices
**So that** I can find specific invoices quickly

---

## 🧠 Functional Description

The system must provide search and filtering for invoices by number, customer, date range, and status.

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/invoices
```

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| invoice_number | string | No | Search by number |
| customer_id | int | No | Filter by customer |
| status | string | No | Filter by status |
| start_date | date | No | Start date filter |
| end_date | date | No | End date filter |
| page | int | No | Page number |
| limit | int | No | Items per page |

### Response (200 OK)
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "invoice_number": "INV-0001",
        "invoice_date": "2026-03-23",
        "customer_name": "Acme Corporation",
        "total": 119000,
        "status": "ISSUED"
      }
    ],
    "pagination": {...}
  },
  "success": true
}
```

---

## 📎 Dependencies

- EPICA-007: Invoicing Epic
