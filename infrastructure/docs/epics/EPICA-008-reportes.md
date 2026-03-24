# EPICA-008: Reports

## 📌 General Information
- ID: EPICA-008
- State: Completed
- Priority: Medium
- Start Date: 2026-03-23
- Target Date: 2026-07-15
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Provide comprehensive reporting capabilities for sales, inventory, and business metrics. The system must generate detailed reports, support multiple formats, and enable data export for analysis.

**What problem does it solve?**
- Manual report generation is time-consuming
- Lack of real-time business visibility
- No standardized report formats
- Limited data export capabilities

**What value does it generate?**
- Real-time business insights
- Automated report generation
- Data-driven decision making
- Export capabilities for external analysis

---

## 👥 Stakeholders

- End User: Store managers, accountants, business owners
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Reports module provides:

1. **Sales Reports**: Daily, weekly, monthly, custom range
2. **Product Reports**: Top sellers, slow movers, margins
3. **Inventory Reports**: Stock levels, movements, valuations
4. **Employee Reports**: Sales by employee, commissions
5. **Financial Reports**: Summary, by payment method
6. **Export**: PDF and Excel formats

---

## 📦 Scope

### Included:
- Sales reports (daily, period, by product, by employee)
- Inventory reports (stock levels, movements, valuation)
- Customer reports (top customers, purchase history)
- Cash drawer reports
- Profit and margin reports
- PDF export
- Excel export
- Report scheduling (future)

### Not Included:
- Advanced BI/analytics
- Real-time dashboards
- Custom report builder
- Data visualization tools

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-REP-001 | Sales Report by Period | ✅ Implemented |
| HU-REP-002 | Sales Report by Product | ✅ Implemented |
| HU-REP-003 | Sales Report by Employee | ✅ Implemented |
| HU-REP-004 | Inventory Status Report | ✅ Implemented |
| HU-REP-005 | Movement History Report | ✅ Implemented |
| HU-REP-006 | Export Report to PDF | ✅ Implemented |
| HU-REP-007 | Export Report to Excel | ✅ Implemented |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- Reports filter by tenant and branch (if applicable)
- Date ranges are inclusive
- All monetary values in COP
- Tax amounts are included in totals
- Deleted records are excluded from reports
- Export respects current filter criteria

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Report Types

**SalesSummaryReport**
- Total sales amount
- Total transactions
- Average ticket
- Sales by payment method
- Sales by hour/day

**ProductSalesReport**
- Product name and SKU
- Units sold
- Revenue
- Margin
- Return rate

**InventoryReport**
- Current stock levels
- Stock value
- Low stock items
- Items without movement

**EmployeeSalesReport**
- Employee name
- Total sales
- Commission earned
- Transaction count

---

## 📊 Success Metrics

- Report generation time < 5 seconds
- Export generation time < 10 seconds
- 100% data accuracy
- Pagination support for large datasets

---

## 🚧 Risks

- Large datasets affecting query performance
- Report export memory usage
- Complex date range queries
- Multi-tenant report isolation

---

## 🧩 Report Parameters

| Report | Parameters |
|--------|------------|
| Sales by Period | start_date, end_date, branch_id (optional), group_by (day/week/month) |
| Sales by Product | start_date, end_date, category_id (optional), branch_id (optional) |
| Sales by Employee | start_date, end_date, branch_id (optional) |
| Inventory Status | branch_id (optional), category_id (optional), stock_filter (all/low/zero) |
| Movement History | start_date, end_date, product_id (optional), movement_type (optional) |

---

## 📡 API Endpoints

### Sales Reports
```
GET    /reports/sales                    → Reporte de ventas por período
GET    /reports/sales/products           → Reporte de ventas por producto
GET    /reports/sales/employees          → Reporte de ventas por empleado
```

### Inventory Reports
```
GET    /reports/inventory                → Estado de inventario
GET    /reports/inventory/movements      → Historial de movimientos
```

### Export Reports
```
POST   /reports/:type/export/pdf         → Exportar reporte a PDF
POST   /reports/:type/export/excel       → Exportar reporte a Excel
```

### Query Parameters
- `start_date` - Fecha inicio (formato: YYYY-MM-DD)
- `end_date` - Fecha fin (formato: YYYY-MM-DD)
- `branch_id` - Filtrar por sucursal (opcional)
- `category_id` - Filtrar por categoría (opcional)
- `product_id` - Filtrar por producto (opcional)
- `group_by` - Agrupar por: day, week, month

---

## 📁 Module Structure

```
modules/reports/
├── domain.go     # Entity, Repository & Service interfaces
├── repository.go # Repository implementation
├── service.go    # Service implementation
├── handler.go    # HTTP handlers
└── routes.go     # Route registration
```

---

## Resumen

- **Total de HU**: 7
- **Completadas**: 7
- **Pendientes**: 0
- **Módulo implementado**: reports
