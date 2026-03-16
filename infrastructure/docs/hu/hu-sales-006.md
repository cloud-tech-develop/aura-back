# HU-SALES-006 - Sales Reporting

## 📌 General Information
- ID: HU-SALES-006
- Epic: EPIC-SALES-001
- Priority: High
- State: Backlog
- Progress: 0%
- Author: QA Engineer Aura POS
- Date: 2026-03-15

---

## 👤 User Story

**As a** store manager or business owner
**I want to** generate sales reports and view analytics
**So that** I can make informed business decisions based on sales data

---

## 🧠 Functional Description

The system must provide comprehensive sales reporting capabilities including:
- Daily, weekly, monthly, and yearly sales summaries
- Sales by product, category, or brand
- Sales by employee/cashier
- Sales by payment method
- Tax reports
- Top selling products
- Sales trends and comparisons

Reports must be filterable by date range, branch, and other criteria.

---

## ✅ Acceptance Criteria

### Scenario 1: Daily sales summary
- Given that I want to view today's sales
- When I request the daily sales report
- Then the system must return:
  - Total sales amount
  - Number of transactions
  - Average transaction value
  - Sales by payment method
  - Tax totals

### Scenario 2: Sales by product
- Given that I want to analyze product performance
- When I request product sales report
- Then the system must return:
  - List of products sold
  - Quantities sold per product
  - Revenue per product
  - Top performing products

### Scenario 3: Sales by employee
- Given that I want to evaluate cashier performance
- When I request employee sales report
- Then the system must return:
  - Sales per employee
  - Number of transactions per employee
  - Average transaction value per employee

### Scenario 4: Tax report
- Given that I need tax documentation
- When I request a tax report for a period
- Then the system must return:
  - Total sales with tax breakdown
  - Tax amounts by rate
  - Taxable and non-taxable sales

### Scenario 5: Sales trends
- Given that I want to analyze sales patterns
- When I request trend analysis
- Then the system must return:
  - Sales comparison between periods
  - Growth/decline percentages
  - Seasonal patterns

---

## ❌ Error Cases

- Date range too large must return error 400
- Invalid date format must return error 400
- Accessing data from other companies must return error 403
- Report generation timeout must be handled gracefully

---

## 🔐 Business Rules

- Reports are scoped to the user's company and branch
- Users can only view reports for their authorized branches
- Report data is read-only
- Large date ranges may require background generation
- Sensitive data (like individual transaction details) requires manager permissions

---

## 🗄️ Database Schema (PostgreSQL)

No new tables required. Reports will be generated using existing tables:
- sales_order
- sales_order_item
- payment
- invoice
- product
- category
- branch
- usuario

### Query Examples

**Daily Sales Summary:**
```sql
SELECT 
    DATE(so.created_at) as sale_date,
    COUNT(DISTINCT so.id) as transaction_count,
    SUM(so.total) as total_sales,
    SUM(so.tax_total) as total_tax,
    COUNT(DISTINCT so.user_id) as cashiers_count
FROM sales_order so
WHERE so.empresa_id = :empresaId
  AND DATE(so.created_at) = :saleDate
  AND so.status IN ('PAID', 'COMPLETED')
GROUP BY DATE(so.created_at);
```

**Product Sales Report:**
```sql
SELECT 
    p.id as product_id,
    p.sku,
    p.name as product_name,
    c.name as category_name,
    SUM(roi.quantity) as total_quantity,
    SUM(roi.total) as total_revenue,
    AVG(roi.unit_price) as avg_price
FROM sales_order_item roi
JOIN product p ON roi.product_id = p.id
JOIN category c ON p.category_id = c.id
WHERE p.empresa_id = :empresaId
  AND roi.created_at BETWEEN :startDate AND :endDate
GROUP BY p.id, p.sku, p.name, c.name
ORDER BY total_revenue DESC;
```

---

## 🎨 UI/UX Considerations

- Interactive charts and graphs
- Date range picker with presets
- Export to Excel/PDF functionality
- Report customization and saving
- Dashboard widgets for key metrics

---

## 📡 Technical Requirements

### DTOs (Java)

**SalesSummaryDto:**
- total_sales (BigDecimal)
- transaction_count (Integer)
- average_transaction (BigDecimal)
- total_tax (BigDecimal)
- payment_method_breakdown (Map<String, BigDecimal>)

**ProductSalesDto:**
- product_id (Long)
- sku (String)
- product_name (String)
- category_name (String)
- total_quantity (Integer)
- total_revenue (BigDecimal)
- avg_price (BigDecimal)

**EmployeeSalesDto:**
- user_id (Long)
- user_name (String)
- total_sales (BigDecimal)
- transaction_count (Integer)
- average_transaction (BigDecimal)

### Services

- `SalesReportService` with methods:
  - `getDailySalesSummary(Date date, Integer branchId)`
  - `getProductSalesReport(DateRange dateRange, Integer categoryId)`
  - `getEmployeeSalesReport(DateRange dateRange)`
  - `getTaxReport(DateRange dateRange)`
  - `getSalesTrends(DateRange dateRange, String period)`

### Controllers

- `SalesReportController` with endpoints:
  - `GET /api/reports/sales/daily` - Daily sales summary
  - `POST /api/reports/sales/product` - Product sales report
  - `POST /api/reports/sales/employee` - Employee sales report
  - `POST /api/reports/tax` - Tax report
  - `POST /api/reports/sales/trends` - Sales trends

---

## 🚫 Out of Scope

- Advanced BI analytics
- Predictive analytics
- Real-time dashboard
- Integration with external BI tools
- Custom report builder

---

## 📎 Dependencies

- HU-SALES-003: Sale/Order Creation
- HU-SALES-004: Payment Processing
- HU-SALES-005: Invoice Generation
- Existing Branch, Product, and User entities
