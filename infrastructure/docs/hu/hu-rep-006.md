# HU-REP-006 - Export Report to PDF

## 📌 General Information
- ID: HU-REP-006
- Epic: EPICA-008
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** export reports to PDF
**So that** I can share printed reports with stakeholders

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/reports/{reportType}/export/pdf
```

### Request
```json
{
  "start_date": "2026-03-01",
  "end_date": "2026-03-23",
  "branch_id": 1,
  "title": "Sales Report March 2026",
  "include_logo": true
}
```

### Response
```
Content-Type: application/pdf
Content-Disposition: attachment; filename="sales-report-2026-03.pdf"
```

Binary PDF data

---

## 📎 Dependencies

- EPICA-008: Reports Epic
