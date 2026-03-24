# HU-INV-INV-004 - Generate Invoice PDF

## 📌 General Information
- ID: HU-INV-INV-004
- Epic: EPICA-007
- Priority: High
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** cashier
**I want to** generate invoice PDF
**So that** I can provide printed invoices to customers

---

## 🧠 Functional Description

The system must generate a PDF document for the invoice with all required fiscal information.

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/invoices/{invoiceId}/pdf
```

### Response
```
Content-Type: application/pdf
Content-Disposition: attachment; filename="INV-0001.pdf"
```

Binary PDF data

---

## 📎 Dependencies

- EPICA-007: Invoicing Epic
