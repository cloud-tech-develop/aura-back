package reports

import (
	"context"
	"time"
)

// SalesReport represents a sales summary report
type SalesReport struct {
	Date          time.Time `json:"date"`
	TotalSales    float64   `json:"total_sales"`
	TotalOrders   int       `json:"total_orders"`
	TotalItems    int       `json:"total_items"`
	AverageTicket float64   `json:"average_ticket"`
	TaxTotal      float64   `json:"tax_total"`
	DiscountTotal float64   `json:"discount_total"`
}

// ProductSalesReport represents product performance data
type ProductSalesReport struct {
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	SKU          string  `json:"sku"`
	QuantitySold int     `json:"quantity_sold"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalProfit  float64 `json:"total_profit"`
}

// PaymentMethodReport represents payment method breakdown
type PaymentMethodReport struct {
	Method       string  `json:"method"`
	Count        int     `json:"count"`
	TotalAmount  float64 `json:"total_amount"`
	Percentage   float64 `json:"percentage"`
}

// DailySalesReport represents sales by day
type DailySalesReport struct {
	Date         time.Time `json:"date"`
	TotalSales   float64   `json:"total_sales"`
	OrderCount   int       `json:"order_count"`
}

// TopCustomer represents top customer by sales
type TopCustomer struct {
	CustomerID   int64   `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	TotalSales   float64 `json:"total_sales"`
	OrderCount   int     `json:"order_count"`
}

// ReportFilters filters for report generation
type ReportFilters struct {
	StartDate    *time.Time
	EndDate      *time.Time
	BranchID     *int64
	ProductID    *int64
	CategoryID   *int64
	UserID       *int64
	Limit        int
	GroupBy      string
}

// SalesReportResponse represents the response for HU-REP-001
type SalesReportResponse struct {
	ReportPeriod ReportPeriod    `json:"report_period"`
	Summary      SalesSummary    `json:"summary"`
	ByPaymentMethod []PaymentMethodBreakdown `json:"by_payment_method"`
	ByDay        []DailySalesData  `json:"by_day"`
	TopProducts  []TopProductSales `json:"top_products"`
	Pagination   Pagination        `json:"pagination"`
}

// ReportPeriod represents the date range
type ReportPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Days      int    `json:"days"`
}

// SalesSummary represents sales summary metrics
type SalesSummary struct {
	TotalSales            float64   `json:"total_sales"`
	TransactionCount      int       `json:"transaction_count"`
	AverageTicket         float64   `json:"average_ticket"`
	PreviousPeriodTotal   float64   `json:"previous_period_total"`
	GrowthPercentage      float64   `json:"growth_percentage"`
}

// PaymentMethodBreakdown represents payment method data
type PaymentMethodBreakdown struct {
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Count         int     `json:"count"`
	Percentage    float64 `json:"percentage"`
}

// DailySalesData represents sales data per day
type DailySalesData struct {
	Date   string  `json:"date"`
	Total  float64 `json:"total"`
	Count  int     `json:"count"`
}

// TopProductSales represents top selling product
type TopProductSales struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitsSold   int     `json:"units_sold"`
	Revenue     float64 `json:"revenue"`
}

// Pagination represents pagination info
type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

// ProductSalesReportResponse represents the response for HU-REP-002
type ProductSalesReportResponse struct {
	Summary   ProductSalesSummary `json:"summary"`
	Products  []ProductSalesData   `json:"products"`
	Pagination Pagination          `json:"pagination"`
}

// ProductSalesSummary represents product sales summary
type ProductSalesSummary struct {
	TotalProductsSold int     `json:"total_products_sold"`
	TotalRevenue      float64 `json:"total_revenue"`
}

// ProductSalesData represents product sales details
type ProductSalesData struct {
	ProductID       int64   `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductSKU      string  `json:"product_sku"`
	CategoryName    string  `json:"category_name"`
	UnitsSold       int     `json:"units_sold"`
	Revenue         float64 `json:"revenue"`
	Cost            float64 `json:"cost"`
	Profit          float64 `json:"profit"`
	MarginPercentage float64 `json:"margin_percentage"`
}

// EmployeeSalesReportResponse represents the response for HU-REP-003
type EmployeeSalesReportResponse struct {
	Employees []EmployeeSalesData `json:"employees"`
	Summary   EmployeeSalesSummary `json:"summary"`
}

// EmployeeSalesSummary represents employee sales summary
type EmployeeSalesSummary struct {
	TotalEmployees int     `json:"total_employees"`
	TotalRevenue   float64 `json:"total_revenue"`
}

// EmployeeSalesData represents employee sales details
type EmployeeSalesData struct {
	EmployeeID      int64   `json:"employee_id"`
	EmployeeName    string  `json:"employee_name"`
	UserName        string  `json:"user_name"`
	TotalSales      int     `json:"total_sales"`
	TotalRevenue    float64 `json:"total_revenue"`
	AverageTicket   float64 `json:"average_ticket"`
	CommissionEarned float64 `json:"commission_earned"`
}

// InventoryReportResponse represents the response for HU-REP-004
type InventoryReportResponse struct {
	Summary   InventorySummary    `json:"summary"`
	Products  []InventoryProduct `json:"products"`
	Pagination Pagination        `json:"pagination"`
}

// InventorySummary represents inventory summary
type InventorySummary struct {
	TotalProducts    int     `json:"total_products"`
	TotalStockValue   float64 `json:"total_stock_value"`
	LowStockCount     int     `json:"low_stock_count"`
	OutOfStockCount   int     `json:"out_of_stock_count"`
}

// InventoryProduct represents inventory product details
type InventoryProduct struct {
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	SKU          string  `json:"sku"`
	CategoryName string  `json:"category_name"`
	Quantity     int     `json:"quantity"`
	MinStock     int     `json:"min_stock"`
	UnitCost     float64 `json:"unit_cost"`
	StockValue   float64 `json:"stock_value"`
	Status       string  `json:"status"`
}

// MovementReportResponse represents the response for HU-REP-005
type MovementReportResponse struct {
	Summary   MovementSummary     `json:"summary"`
	Movements []MovementData      `json:"movements"`
	Pagination Pagination         `json:"pagination"`
}

// MovementSummary represents movement summary
type MovementSummary struct {
	TotalEntries int `json:"total_entries"`
	TotalExits   int `json:"total_exits"`
	NetChange    int `json:"net_change"`
}

// MovementData represents movement details
type MovementData struct {
	ID             int64   `json:"id"`
	ProductName    string  `json:"product_name"`
	MovementType   string  `json:"movement_type"`
	MovementReason string  `json:"movement_reason"`
	Quantity       int     `json:"quantity"`
	UserName       string  `json:"user_name"`
	CreatedAt      string  `json:"created_at"`
}

// ExportRequest represents export request body
type ExportRequest struct {
	StartDate     string `json:"start_date" binding:"required"`
	EndDate       string `json:"end_date" binding:"required"`
	BranchID      *int64 `json:"branch_id"`
	Title         string `json:"title"`
	IncludeLogo   bool   `json:"include_logo"`
	IncludeSummary bool  `json:"include_summary"`
}

// Repository interface for reports
type Repository interface {
	GetSalesSummary(ctx context.Context, enterpriseID int64, filters ReportFilters) (*SalesReport, error)
	GetProductSales(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]ProductSalesReport, error)
	GetPaymentMethodBreakdown(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]PaymentMethodReport, error)
	GetDailySales(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]DailySalesReport, error)
	GetTopCustomers(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]TopCustomer, error)

	// New methods for EPICA-008
	GetSalesByPeriod(ctx context.Context, enterpriseID int64, filters ReportFilters) (SalesReportResponse, error)
	GetSalesByProduct(ctx context.Context, enterpriseID int64, filters ReportFilters) (ProductSalesReportResponse, error)
	GetSalesByEmployee(ctx context.Context, enterpriseID int64, filters ReportFilters) (EmployeeSalesReportResponse, error)
	GetInventoryStatus(ctx context.Context, enterpriseID int64, filters ReportFilters) (InventoryReportResponse, error)
	GetMovementHistory(ctx context.Context, enterpriseID int64, filters ReportFilters) (MovementReportResponse, error)
}

// Service interface for reports business logic
type Service interface {
	GetSalesSummary(ctx context.Context, enterpriseID int64, filters ReportFilters) (*SalesReport, error)
	GetProductSales(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]ProductSalesReport, error)
	GetPaymentMethodBreakdown(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]PaymentMethodReport, error)
	GetDailySales(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]DailySalesReport, error)
	GetTopCustomers(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]TopCustomer, error)

	// New methods for EPICA-008
	GetSalesByPeriod(ctx context.Context, enterpriseID int64, filters ReportFilters) (SalesReportResponse, error)
	GetSalesByProduct(ctx context.Context, enterpriseID int64, filters ReportFilters) (ProductSalesReportResponse, error)
	GetSalesByEmployee(ctx context.Context, enterpriseID int64, filters ReportFilters) (EmployeeSalesReportResponse, error)
	GetInventoryStatus(ctx context.Context, enterpriseID int64, filters ReportFilters) (InventoryReportResponse, error)
	GetMovementHistory(ctx context.Context, enterpriseID int64, filters ReportFilters) (MovementReportResponse, error)
	ExportToPDF(ctx context.Context, enterpriseID int64, reportType string, req ExportRequest) ([]byte, error)
	ExportToExcel(ctx context.Context, enterpriseID int64, reportType string, req ExportRequest) ([]byte, error)
}
