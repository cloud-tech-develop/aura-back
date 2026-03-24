# HU-REP-007 - Export Report to Excel

## 📌 General Information
- ID: HU-REP-007
- Epic: EPICA-008
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** accountant
**I want to** export reports to Excel
**So that** I can perform further analysis in spreadsheets

---

## 📡 Technical Requirements

### Endpoint
```
POST /api/reports/{reportType}/export/excel
```

### Request
```json
{
  "start_date": "2026-03-01",
  "end_date": "2026-03-23",
  "branch_id": 1,
  "include_summary": true
}
```

### Response
```
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="sales-report-2026-03.xlsx"
```

Binary Excel data

---

## 📎 Dependencies

- EPICA-008: Reports Epic
