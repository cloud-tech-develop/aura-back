# HU-REP-001 - Sales Report by Period

## 📌 General Information
- ID: HU-REP-001
- Epic: EPICA-008
- Priority: Medium
- State: Backlog
- Progress: 0%
- Author: Aura POS Backend Team
- Date: 2026-03-23

---

## 👤 User Story

**As a** store manager
**I want to** generate sales reports by period
**So that** I can analyze sales performance over time

---

## 🧠 Functional Description

The system must generate comprehensive sales reports filtered by date range, with aggregated metrics, breakdowns by payment method, and product-level details. Reports can be paginated and exported.

---

## ✅ Acceptance Criteria

### Scenario 1: Daily sales report
- Given that I am logged in as a manager
- When I request a sales report for today
- Then I must receive:
  - Total sales amount
  - Transaction count
  - Average ticket
  - Sales by payment method
  - Top products sold

### Scenario 2: Weekly comparison
- Given that I am logged in as a manager
- When I request a weekly report with comparison
- Then I must receive:
  - Current week totals
  - Previous week totals
  - Percentage change

### Scenario 3: Filter by branch
- Given that I have access to multiple branches
- When I filter by specific branch
- Then only that branch's sales are included

---

## ❌ Error Cases

- Invalid date range returns error 400
- Future end date returns error 400
- No data returns empty results with zeros

---

## 🔐 Business Rules

- Date range maximum: 365 days
- All monetary values in COP
- Deleted records excluded
- Branch filter required for multi-branch users

---

## 📡 Technical Requirements

### Endpoint
```
GET /api/reports/sales
```

### Method: GET

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | date | Yes | Report start date (YYYY-MM-DD) |
| end_date | date | Yes | Report end date (YYYY-MM-DD) |
| branch_id | int | No | Filter by branch |
| group_by | string | No | Group by: day, week, month (default: day) |
| include_details | bool | No | Include product details (default: false) |

### Response (200 OK)
```json
{
  "data": {
    "report_period": {
      "start_date": "2026-03-16",
      "end_date": "2026-03-22",
      "days": 7
    },
    "summary": {
      "total_sales": 15000000,
      "transaction_count": 245,
      "average_ticket": 61224,
      "previous_period_total": 12000000,
      "growth_percentage": 25.0
    },
    "by_payment_method": [
      {
        "payment_method": "CASH",
        "amount": 9000000,
        "count": 147,
        "percentage": 60.0
      },
      {
        "payment_method": "CREDIT_CARD",
        "amount": 4500000,
        "count": 74,
        "percentage": 30.0
      },
      {
        "payment_method": "DEBIT_CARD",
        "amount": 1500000,
        "count": 24,
        "percentage": 10.0
      }
    ],
    "by_day": [
      {
        "date": "2026-03-16",
        "total": 2100000,
        "count": 35
      },
      {
        "date": "2026-03-17",
        "total": 1800000,
        "count": 29
      }
    ],
    "top_products": [
      {
        "product_id": 101,
        "product_name": "Wireless Mouse",
        "units_sold": 45,
        "revenue": 2025000
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 20,
      "total_items": 245,
      "total_pages": 13
    }
  },
  "success": true
}
```

---

## 🧪 Testing Criteria

### Unit Tests
- Test date range validation
- Test aggregation calculations
- Test grouping logic

### Integration Tests
- Test multi-tenant isolation
- Test date range performance
- Test large dataset handling

---

## 📎 Dependencies

- EPICA-008: Reports Epic
- Existing sales module
- Existing payment module
