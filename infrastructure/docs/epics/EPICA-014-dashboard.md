# EPICA-014: Dashboard

## 📌 General Information
- ID: EPICA-014
- State: Backlog
- Priority: Medium
- Start Date: 2026-03-23
- Target Date: 2026-09-30
- Owner: Aura POS Backend Team
- Percentage: 0%

---

## 🎯 Business Objective

Provide a comprehensive dashboard with key business metrics, real-time KPIs, and actionable insights. The dashboard aggregates data from all modules to give business owners and managers a complete view of operations.

**What problem does it solve?**
- No unified view of business metrics
- Manual data aggregation
- Lack of real-time visibility
- No trend analysis

**What value does it generate?**
- Real-time business visibility
- Data-driven decisions
- Trend identification
- Performance monitoring

---

## 👥 Stakeholders

- End User: Business owners, store managers, accountants
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Dashboard module provides:

1. **Sales Metrics**: Today's sales, comparison with previous periods
2. **Inventory Alerts**: Low stock, expiring products
3. **Top Products**: Best sellers, slow movers
4. **Customer Metrics**: New customers, repeat customers
5. **Financial Summary**: Revenue, costs, margins
6. **Trend Charts**: Sales over time, product performance

---

## 📦 Scope

### Included:
- Today's sales summary
- Sales comparison (day, week, month)
- Top selling products
- Low stock alerts
- Recent transactions
- Quick actions
- Date range filtering

### Not Included:
- Customizable widgets
- Advanced analytics
- Data export from dashboard
- Push notifications

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-DASH-001 | View Sales Summary | ⏳ Pending |
| HU-DASH-002 | View Inventory Alerts | ⏳ Pending |
| HU-DASH-003 | View Top Products | ⏳ Pending |
| HU-DASH-004 | View Recent Activity | ⏳ Pending |
| HU-DASH-005 | Filter Dashboard by Period | ⏳ Pending |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- Dashboard data filters by tenant and user permissions
- All monetary values in COP
- Comparisons use same period from previous month/year
- Alerts consider branch-level inventory
- Recent activity limited to last 24 hours or configurable

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Dashboard Data Structures

**SalesSummary**
```go
type SalesSummary struct {
    TodaySales     decimal.Decimal
    YesterdaySales decimal.Decimal
    WeekSales      decimal.Decimal
    MonthSales     decimal.Decimal
    SalesGrowth    float64
    TransactionCount int
    AverageTicket  decimal.Decimal
}
```

**InventoryAlert**
```go
type InventoryAlert struct {
    ProductID   int64
    ProductName string
    BranchID    int64
    CurrentStock int
    MinStock    int
    AlertType   string // LOW_STOCK, EXPIRING, OUT_OF_STOCK
}
```

**TopProduct**
```go
type TopProduct struct {
    ProductID    int64
    ProductName  string
    UnitsSold    int
    Revenue      decimal.Decimal
    ProfitMargin decimal.Decimal
    Rank         int
}
```

---

## 📊 Success Metrics

- Dashboard load time < 2 seconds
- Real-time data refresh < 30 seconds
- 100% data accuracy
- 99% uptime

---

## 🚧 Risks

- Large dataset aggregation affecting performance
- Multi-branch data consolidation
- Caching strategy complexity
- Real-time updates with high concurrency

---

## 📡 API Endpoints (Planificados)

### Dashboard Metrics
```
GET    /dashboard/summary              → Resumen general del día
GET    /dashboard/sales                → Métricas de ventas
GET    /dashboard/inventory/alerts     → Alertas de inventario
GET    /dashboard/products/top         → Productos más vendidos
GET    /dashboard/activity/recent      → Actividad reciente
```

### Query Parameters
- `period` - Período: today, week, month, custom
- `start_date` - Fecha inicio (si custom)
- `end_date` - Fecha fin (si custom)
- `branch_id` - Filtrar por sucursal

---

## Resumen

- **Total de HU**: 5
- **Completadas**: 0
- **Pendientes**: 5
- **Estado**: Backlog
