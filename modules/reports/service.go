package reports

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type repository struct {
	db db.Querier
}

func NewRepository(db db.Querier) Repository {
	return &repository{db: db}
}

func (r *repository) GetSalesSummary(ctx context.Context, enterpriseID int64, filters ReportFilters) (*SalesReport, error) {
	query := `
		SELECT 
			COALESCE(SUM(total), 0) as total_sales,
			COUNT(*) as total_orders,
			COALESCE(SUM(discount), 0) as discount_total,
			COALESCE(SUM(tax_total), 0) as tax_total
		FROM sales_order
		WHERE enterprise_id = $1 AND status = 'COMPLETED'`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, *filters.EndDate)
		argPos++
	}

	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	report := &SalesReport{}
	if filters.StartDate != nil {
		report.Date = *filters.StartDate
	} else {
		report.Date = time.Now()
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&report.TotalSales, &report.TotalOrders, &report.DiscountTotal, &report.TaxTotal,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales summary: %w", err)
	}

	if report.TotalOrders > 0 {
		report.AverageTicket = report.TotalSales / float64(report.TotalOrders)
	}

	// Get total items sold
	itemsQuery := `
		SELECT COALESCE(SUM(quantity), 0)
		FROM sales_order_item soi
		JOIN sales_order so ON so.id = soi.sales_order_id
		WHERE so.enterprise_id = $1 AND so.status = 'COMPLETED'`

	args2 := []interface{}{enterpriseID}
	if filters.StartDate != nil {
		itemsQuery += fmt.Sprintf(" AND so.created_at >= $2")
		args2 = append(args2, *filters.StartDate)
	}
	if filters.EndDate != nil {
		itemsQuery += fmt.Sprintf(" AND so.created_at <= $%d", len(args2)+1)
		args2 = append(args2, *filters.EndDate)
	}
	if filters.BranchID != nil {
		itemsQuery += fmt.Sprintf(" AND so.branch_id = $%d", len(args2)+1)
		args2 = append(args2, *filters.BranchID)
	}

	r.db.QueryRowContext(ctx, itemsQuery, args2...).Scan(&report.TotalItems)

	return report, nil
}

func (r *repository) GetProductSales(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]ProductSalesReport, error) {
	query := `
		SELECT 
			p.id as product_id,
			p.name as product_name,
			p.sku,
			SUM(soi.quantity) as quantity_sold,
			SUM(soi.total) as total_revenue,
			SUM(soi.total - (soi.quantity * p.cost_price)) as total_profit
		FROM sales_order_item soi
		JOIN sales_order so ON so.id = soi.sales_order_id
		JOIN product p ON p.id = soi.product_id
		WHERE so.enterprise_id = $1 AND so.status = 'COMPLETED'`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND so.created_at >= $%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND so.created_at <= $%d", argPos)
		args = append(args, *filters.EndDate)
		argPos++
	}

	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND so.branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	if filters.ProductID != nil {
		query += fmt.Sprintf(" AND p.id = $%d", argPos)
		args = append(args, *filters.ProductID)
		argPos++
	}

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND p.category_id = $%d", argPos)
		args = append(args, *filters.CategoryID)
		argPos++
	}

	query += " GROUP BY p.id, p.name, p.sku"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filters.Limit)
	} else {
		query += " LIMIT 20"
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get product sales: %w", err)
	}
	defer rows.Close()

	var reports []ProductSalesReport
	for rows.Next() {
		var rep ProductSalesReport
		if err := rows.Scan(&rep.ProductID, &rep.ProductName, &rep.SKU, &rep.QuantitySold, &rep.TotalRevenue, &rep.TotalProfit); err != nil {
			return nil, err
		}
		reports = append(reports, rep)
	}
	return reports, nil
}

func (r *repository) GetPaymentMethodBreakdown(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]PaymentMethodReport, error) {
	query := `
		SELECT 
			payment_method,
			COUNT(*) as count,
			SUM(amount) as total_amount
		FROM payment p
		JOIN sales_order so ON so.id = p.sales_order_id
		WHERE p.enterprise_id = $1`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND p.created_at >= $%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND p.created_at <= $%d", argPos)
		args = append(args, *filters.EndDate)
		argPos++
	}

	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND p.branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	query += " GROUP BY payment_method ORDER BY total_amount DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method breakdown: %w", err)
	}
	defer rows.Close()

	var reports []PaymentMethodReport
	var totalAmount float64
	for rows.Next() {
		var rep PaymentMethodReport
		if err := rows.Scan(&rep.Method, &rep.Count, &rep.TotalAmount); err != nil {
			return nil, err
		}
		totalAmount += rep.TotalAmount
		reports = append(reports, rep)
	}

	// Calculate percentages
	if totalAmount > 0 {
		for i := range reports {
			reports[i].Percentage = (reports[i].TotalAmount / totalAmount) * 100
		}
	}

	return reports, nil
}

func (r *repository) GetDailySales(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]DailySalesReport, error) {
	query := `
		SELECT 
			DATE(created_at) as date,
			SUM(total) as total_sales,
			COUNT(*) as order_count
		FROM sales_order
		WHERE enterprise_id = $1 AND status = 'COMPLETED'`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, *filters.EndDate)
		argPos++
	}

	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	query += " GROUP BY DATE(created_at) ORDER BY date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily sales: %w", err)
	}
	defer rows.Close()

	var reports []DailySalesReport
	for rows.Next() {
		var rep DailySalesReport
		if err := rows.Scan(&rep.Date, &rep.TotalSales, &rep.OrderCount); err != nil {
			return nil, err
		}
		reports = append(reports, rep)
	}
	return reports, nil
}

func (r *repository) GetTopCustomers(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]TopCustomer, error) {
	query := `
		SELECT 
			tp.id as customer_id,
			tp.name as customer_name,
			SUM(so.total) as total_sales,
			COUNT(*) as order_count
		FROM sales_order so
		JOIN third_parties tp ON tp.id = so.customer_id
		WHERE so.enterprise_id = $1 AND so.status = 'COMPLETED' AND so.customer_id IS NOT NULL`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND so.created_at >= $%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND so.created_at <= $%d", argPos)
		args = append(args, *filters.EndDate)
		argPos++
	}

	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND so.branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	query += " GROUP BY tp.id, tp.name ORDER BY total_sales DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filters.Limit)
	} else {
		query += " LIMIT 10"
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get top customers: %w", err)
	}
	defer rows.Close()

	var reports []TopCustomer
	for rows.Next() {
		var rep TopCustomer
		if err := rows.Scan(&rep.CustomerID, &rep.CustomerName, &rep.TotalSales, &rep.OrderCount); err != nil {
			return nil, err
		}
		reports = append(reports, rep)
	}
	return reports, nil
}

// Service implementation
type service struct {
	repo Repository
}

func NewService(db db.Querier) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) GetSalesSummary(ctx context.Context, enterpriseID int64, filters ReportFilters) (*SalesReport, error) {
	return s.repo.GetSalesSummary(ctx, enterpriseID, filters)
}

func (s *service) GetProductSales(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]ProductSalesReport, error) {
	if filters.Limit == 0 {
		filters.Limit = 20
	}
	return s.repo.GetProductSales(ctx, enterpriseID, filters)
}

func (s *service) GetPaymentMethodBreakdown(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]PaymentMethodReport, error) {
	return s.repo.GetPaymentMethodBreakdown(ctx, enterpriseID, filters)
}

func (s *service) GetDailySales(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]DailySalesReport, error) {
	return s.repo.GetDailySales(ctx, enterpriseID, filters)
}

func (s *service) GetTopCustomers(ctx context.Context, enterpriseID int64, filters ReportFilters) ([]TopCustomer, error) {
	if filters.Limit == 0 {
		filters.Limit = 10
	}
	return s.repo.GetTopCustomers(ctx, enterpriseID, filters)
}

func (r *repository) GetSalesByPeriod(ctx context.Context, enterpriseID int64, filters ReportFilters) (SalesReportResponse, error) {
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()
	days := 30

	if filters.StartDate != nil {
		startDate = *filters.StartDate
	}
	if filters.EndDate != nil {
		endDate = *filters.EndDate
	}
	days = int(endDate.Sub(startDate).Hours()/24) + 1

	prevStartDate := startDate.AddDate(0, 0, -days)
	prevEndDate := startDate.AddDate(0, 0, -1)

	query := `
		SELECT 
			COALESCE(SUM(total), 0) as total_sales,
			COUNT(*) as transaction_count
		FROM sales_order
		WHERE enterprise_id = $1 AND status = 'COMPLETED'`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, startDate)
		argPos++
	}
	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, endDate)
		argPos++
	}
	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	var totalSales float64
	var transactionCount int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&totalSales, &transactionCount); err != nil {
		return SalesReportResponse{}, fmt.Errorf("failed to get sales by period: %w", err)
	}

	prevQuery := `
		SELECT COALESCE(SUM(total), 0)
		FROM sales_order
		WHERE enterprise_id = $1 AND status = 'COMPLETED' AND created_at >= $2 AND created_at <= $3`
	var prevTotal float64
	if err := r.db.QueryRowContext(ctx, prevQuery, enterpriseID, prevStartDate, prevEndDate).Scan(&prevTotal); err != nil && err != sql.ErrNoRows {
		return SalesReportResponse{}, fmt.Errorf("failed to get previous period total: %w", err)
	}

	var averageTicket float64
	if transactionCount > 0 {
		averageTicket = totalSales / float64(transactionCount)
	}

	var growth float64
	if prevTotal > 0 {
		growth = ((totalSales - prevTotal) / prevTotal) * 100
	}

	summary := SalesSummary{
		TotalSales:          totalSales,
		TransactionCount:    transactionCount,
		AverageTicket:       averageTicket,
		PreviousPeriodTotal: prevTotal,
		GrowthPercentage:    growth,
	}

	byDayQuery := `
		SELECT 
			DATE(created_at) as date,
			SUM(total) as total,
			COUNT(*) as count
		FROM sales_order
		WHERE enterprise_id = $1 AND status = 'COMPLETED'`

	dayArgs := []interface{}{enterpriseID}
	dayArgPos := 2

	if filters.StartDate != nil {
		byDayQuery += fmt.Sprintf(" AND created_at >= $%d", dayArgPos)
		dayArgs = append(dayArgs, startDate)
		dayArgPos++
	}
	if filters.EndDate != nil {
		byDayQuery += fmt.Sprintf(" AND created_at <= $%d", dayArgPos)
		dayArgs = append(dayArgs, endDate)
		dayArgPos++
	}
	if filters.BranchID != nil {
		byDayQuery += fmt.Sprintf(" AND branch_id = $%d", dayArgPos)
		dayArgs = append(dayArgs, *filters.BranchID)
		dayArgPos++
	}

	byDayQuery += " GROUP BY DATE(created_at) ORDER BY date DESC"

	dayRows, err := r.db.QueryContext(ctx, byDayQuery, dayArgs...)
	if err != nil {
		return SalesReportResponse{}, fmt.Errorf("failed to get daily sales: %w", err)
	}
	defer dayRows.Close()

	var byDay []DailySalesData
	for dayRows.Next() {
		var d DailySalesData
		var date time.Time
		if err := dayRows.Scan(&date, &d.Total, &d.Count); err != nil {
			return SalesReportResponse{}, err
		}
		d.Date = date.Format("2006-01-02")
		byDay = append(byDay, d)
	}

	paymentQuery := `
		SELECT 
			payment_method,
			SUM(amount) as amount,
			COUNT(*) as count
		FROM payment p
		JOIN sales_order so ON so.id = p.sales_order_id
		WHERE p.enterprise_id = $1 AND so.status = 'COMPLETED'`

	payArgs := []interface{}{enterpriseID}
	payArgPos := 2

	if filters.StartDate != nil {
		paymentQuery += fmt.Sprintf(" AND p.created_at >= $%d", payArgPos)
		payArgs = append(payArgs, startDate)
		payArgPos++
	}
	if filters.EndDate != nil {
		paymentQuery += fmt.Sprintf(" AND p.created_at <= $%d", payArgPos)
		payArgs = append(payArgs, endDate)
		payArgPos++
	}
	if filters.BranchID != nil {
		paymentQuery += fmt.Sprintf(" AND p.branch_id = $%d", payArgPos)
		payArgs = append(payArgs, *filters.BranchID)
		payArgPos++
	}

	paymentQuery += " GROUP BY payment_method"

	payRows, err := r.db.QueryContext(ctx, paymentQuery, payArgs...)
	if err != nil {
		return SalesReportResponse{}, fmt.Errorf("failed to get payment breakdown: %w", err)
	}
	defer payRows.Close()

	var byPayment []PaymentMethodBreakdown
	var totalPayAmount float64
	for payRows.Next() {
		var p PaymentMethodBreakdown
		if err := payRows.Scan(&p.PaymentMethod, &p.Amount, &p.Count); err != nil {
			return SalesReportResponse{}, err
		}
		totalPayAmount += p.Amount
		byPayment = append(byPayment, p)
	}

	if totalPayAmount > 0 {
		for i := range byPayment {
			byPayment[i].Percentage = (byPayment[i].Amount / totalPayAmount) * 100
		}
	}

	topProductsQuery := `
		SELECT 
			p.id as product_id,
			p.name as product_name,
			SUM(soi.quantity) as units_sold,
			SUM(soi.total) as revenue
		FROM sales_order_item soi
		JOIN sales_order so ON so.id = soi.sales_order_id
		JOIN product p ON p.id = soi.product_id
		WHERE so.enterprise_id = $1 AND so.status = 'COMPLETED'`

	prodArgs := []interface{}{enterpriseID}
	prodArgPos := 2

	if filters.StartDate != nil {
		topProductsQuery += fmt.Sprintf(" AND so.created_at >= $%d", prodArgPos)
		prodArgs = append(prodArgs, startDate)
		prodArgPos++
	}
	if filters.EndDate != nil {
		topProductsQuery += fmt.Sprintf(" AND so.created_at <= $%d", prodArgPos)
		prodArgs = append(prodArgs, endDate)
		prodArgPos++
	}
	if filters.BranchID != nil {
		topProductsQuery += fmt.Sprintf(" AND so.branch_id = $%d", prodArgPos)
		prodArgs = append(prodArgs, *filters.BranchID)
		prodArgPos++
	}

	topProductsQuery += " GROUP BY p.id, p.name ORDER BY revenue DESC LIMIT 10"

	prodRows, err := r.db.QueryContext(ctx, topProductsQuery, prodArgs...)
	if err != nil {
		return SalesReportResponse{}, fmt.Errorf("failed to get top products: %w", err)
	}
	defer prodRows.Close()

	var topProducts []TopProductSales
	for prodRows.Next() {
		var tp TopProductSales
		if err := prodRows.Scan(&tp.ProductID, &tp.ProductName, &tp.UnitsSold, &tp.Revenue); err != nil {
			return SalesReportResponse{}, err
		}
		topProducts = append(topProducts, tp)
	}

	return SalesReportResponse{
		ReportPeriod: ReportPeriod{
			StartDate: startDate.Format("2006-01-02"),
			EndDate:   endDate.Format("2006-01-02"),
			Days:      days,
		},
		Summary:         summary,
		ByPaymentMethod: byPayment,
		ByDay:           byDay,
		TopProducts:     topProducts,
		Pagination: Pagination{
			CurrentPage: 1,
			PerPage:     20,
			TotalItems:  transactionCount,
			TotalPages:  (transactionCount + 19) / 20,
		},
	}, nil
}

func (r *repository) GetSalesByProduct(ctx context.Context, enterpriseID int64, filters ReportFilters) (ProductSalesReportResponse, error) {
	page := 1
	perPage := 20
	offset := 0

	if filters.Limit > 0 {
		perPage = filters.Limit
	}

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	if filters.StartDate != nil {
		startDate = *filters.StartDate
	}
	if filters.EndDate != nil {
		endDate = *filters.EndDate
	}

	query := `
		SELECT 
			p.id as product_id,
			p.name as product_name,
			p.sku,
			COALESCE(c.name, 'Sin categoría') as category_name,
			SUM(soi.quantity) as units_sold,
			SUM(soi.total) as revenue,
			SUM(soi.quantity * p.cost_price) as cost,
			SUM(soi.total - (soi.quantity * p.cost_price)) as profit
		FROM sales_order_item soi
		JOIN sales_order so ON so.id = soi.sales_order_id
		JOIN product p ON p.id = soi.product_id
		LEFT JOIN category c ON c.id = p.category_id
		WHERE so.enterprise_id = $1 AND so.status = 'COMPLETED'`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND so.created_at >= $%d", argPos)
		args = append(args, startDate)
		argPos++
	}
	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND so.created_at <= $%d", argPos)
		args = append(args, endDate)
		argPos++
	}
	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND so.branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}
	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND p.category_id = $%d", argPos)
		args = append(args, *filters.CategoryID)
		argPos++
	}

	query += fmt.Sprintf(" GROUP BY p.id, p.name, p.sku, c.name ORDER BY revenue DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, perPage, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return ProductSalesReportResponse{}, fmt.Errorf("failed to get sales by product: %w", err)
	}
	defer rows.Close()

	var products []ProductSalesData
	var totalProducts int
	var totalRevenue float64
	for rows.Next() {
		var p ProductSalesData
		if err := rows.Scan(&p.ProductID, &p.ProductName, &p.ProductSKU, &p.CategoryName, &p.UnitsSold, &p.Revenue, &p.Cost, &p.Profit); err != nil {
			return ProductSalesReportResponse{}, err
		}
		if p.Revenue > 0 {
			p.MarginPercentage = (p.Profit / p.Revenue) * 100
		}
		products = append(products, p)
		totalProducts += p.UnitsSold
		totalRevenue += p.Revenue
	}

	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM sales_order_item soi
		JOIN sales_order so ON so.id = soi.sales_order_id
		JOIN product p ON p.id = soi.product_id
		WHERE so.enterprise_id = $1 AND so.status = 'COMPLETED'`

	countArgs := []interface{}{enterpriseID}
	if filters.StartDate != nil {
		countQuery += fmt.Sprintf(" AND so.created_at >= $%d", len(countArgs)+1)
		countArgs = append(countArgs, startDate)
	}
	if filters.EndDate != nil {
		countQuery += fmt.Sprintf(" AND so.created_at <= $%d", len(countArgs)+1)
		countArgs = append(countArgs, endDate)
	}

	r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalProducts)

	return ProductSalesReportResponse{
		Summary: ProductSalesSummary{
			TotalProductsSold: totalProducts,
			TotalRevenue:      totalRevenue,
		},
		Products: products,
		Pagination: Pagination{
			CurrentPage: page,
			PerPage:     perPage,
			TotalItems:  totalProducts,
			TotalPages:  (totalProducts + perPage - 1) / perPage,
		},
	}, nil
}

func (r *repository) GetSalesByEmployee(ctx context.Context, enterpriseID int64, filters ReportFilters) (EmployeeSalesReportResponse, error) {
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	if filters.StartDate != nil {
		startDate = *filters.StartDate
	}
	if filters.EndDate != nil {
		endDate = *filters.EndDate
	}

	query := `
		SELECT 
			u.id as employee_id,
			COALESCE(tp.name, u.name) as employee_name,
			u.name as user_name,
			COUNT(so.id) as total_sales,
			COALESCE(SUM(so.total), 0) as total_revenue,
			COALESCE(SUM(so.total) * 0.05, 0) as commission_earned
		FROM sales_order so
		JOIN users u ON u.id = so.user_id
		LEFT JOIN third_parties tp ON tp.user_id = u.id AND tp.is_employee = true
		WHERE so.enterprise_id = $1 AND so.status = 'COMPLETED'`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND so.created_at >= $%d", argPos)
		args = append(args, startDate)
		argPos++
	}
	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND so.created_at <= $%d", argPos)
		args = append(args, endDate)
		argPos++
	}
	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND so.branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	query += " GROUP BY u.id, tp.name, u.name ORDER BY total_revenue DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return EmployeeSalesReportResponse{}, fmt.Errorf("failed to get sales by employee: %w", err)
	}
	defer rows.Close()

	var employees []EmployeeSalesData
	var totalEmployees int
	var totalRevenue float64

	for rows.Next() {
		var e EmployeeSalesData
		if err := rows.Scan(&e.EmployeeID, &e.EmployeeName, &e.UserName, &e.TotalSales, &e.TotalRevenue, &e.CommissionEarned); err != nil {
			return EmployeeSalesReportResponse{}, err
		}
		if e.TotalSales > 0 {
			e.AverageTicket = e.TotalRevenue / float64(e.TotalSales)
		}
		employees = append(employees, e)
		totalEmployees++
		totalRevenue += e.TotalRevenue
	}

	return EmployeeSalesReportResponse{
		Employees: employees,
		Summary: EmployeeSalesSummary{
			TotalEmployees: totalEmployees,
			TotalRevenue:   totalRevenue,
		},
	}, nil
}

func (r *repository) GetInventoryStatus(ctx context.Context, enterpriseID int64, filters ReportFilters) (InventoryReportResponse, error) {
	query := `
		SELECT 
			p.id as product_id,
			p.name as product_name,
			p.sku,
			COALESCE(c.name, 'Sin categoría') as category_name,
			COALESCE(i.quantity, 0) as quantity,
			COALESCE(i.min_stock, 0) as min_stock,
			COALESCE(p.cost_price, 0) as unit_cost,
			COALESCE(i.quantity, 0) * COALESCE(p.cost_price, 0) as stock_value
		FROM product p
		LEFT JOIN category c ON c.id = p.category_id
		LEFT JOIN inventory i ON i.product_id = p.id`

	args := []interface{}{enterpriseID}
	argPos := 1

	where := " WHERE p.enterprise_id = $1"

	if filters.BranchID != nil {
		where += fmt.Sprintf(" AND i.branch_id = $%d", argPos+1)
		args = append(args, *filters.BranchID)
		argPos++
	}
	if filters.CategoryID != nil {
		where += fmt.Sprintf(" AND p.category_id = $%d", argPos+1)
		args = append(args, *filters.CategoryID)
		argPos++
	}

	query += where + " ORDER BY p.name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return InventoryReportResponse{}, fmt.Errorf("failed to get inventory status: %w", err)
	}
	defer rows.Close()

	var products []InventoryProduct
	var totalProducts int
	var totalStockValue float64
	var lowStockCount int
	var outOfStockCount int

	for rows.Next() {
		var p InventoryProduct
		if err := rows.Scan(&p.ProductID, &p.ProductName, &p.SKU, &p.CategoryName, &p.Quantity, &p.MinStock, &p.UnitCost, &p.StockValue); err != nil {
			return InventoryReportResponse{}, err
		}

		if p.Quantity == 0 {
			p.Status = "OUT_OF_STOCK"
			outOfStockCount++
		} else if p.Quantity <= p.MinStock {
			p.Status = "LOW_STOCK"
			lowStockCount++
		} else {
			p.Status = "NORMAL"
		}

		products = append(products, p)
		totalProducts++
		totalStockValue += p.StockValue
	}

	return InventoryReportResponse{
		Summary: InventorySummary{
			TotalProducts:   totalProducts,
			TotalStockValue: totalStockValue,
			LowStockCount:   lowStockCount,
			OutOfStockCount: outOfStockCount,
		},
		Products: products,
		Pagination: Pagination{
			CurrentPage: 1,
			PerPage:     20,
			TotalItems:  totalProducts,
			TotalPages:  (totalProducts + 19) / 20,
		},
	}, nil
}

func (r *repository) GetMovementHistory(ctx context.Context, enterpriseID int64, filters ReportFilters) (MovementReportResponse, error) {
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	if filters.StartDate != nil {
		startDate = *filters.StartDate
	}
	if filters.EndDate != nil {
		endDate = *filters.EndDate
	}

	query := `
		SELECT 
			im.id,
			p.name as product_name,
			im.movement_type,
			im.movement_reason,
			im.quantity,
			COALESCE(u.name, 'Sistema') as user_name,
			im.created_at
		FROM inventory_movement im
		JOIN inventory i ON i.id = im.inventory_id
		JOIN product p ON p.id = i.product_id
		LEFT JOIN users u ON u.id = im.user_id
		WHERE p.enterprise_id = $1`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND im.created_at >= $%d", argPos)
		args = append(args, startDate)
		argPos++
	}
	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND im.created_at <= $%d", argPos)
		args = append(args, endDate)
		argPos++
	}
	if filters.ProductID != nil {
		query += fmt.Sprintf(" AND i.product_id = $%d", argPos)
		args = append(args, *filters.ProductID)
		argPos++
	}

	query += " ORDER BY im.created_at DESC LIMIT 100"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return MovementReportResponse{}, fmt.Errorf("failed to get movement history: %w", err)
	}
	defer rows.Close()

	var movements []MovementData
	var totalEntries int
	var totalExits int

	for rows.Next() {
		var m MovementData
		var createdAt time.Time
		if err := rows.Scan(&m.ID, &m.ProductName, &m.MovementType, &m.MovementReason, &m.Quantity, &m.UserName, &createdAt); err != nil {
			return MovementReportResponse{}, err
		}
		m.CreatedAt = createdAt.Format("2006-01-02T15:04:05Z")
		movements = append(movements, m)

		if m.MovementType == "ENTRY" {
			totalEntries += m.Quantity
		} else if m.MovementType == "EXIT" {
			totalExits += m.Quantity
		}
	}

	return MovementReportResponse{
		Summary: MovementSummary{
			TotalEntries: totalEntries,
			TotalExits:   totalExits,
			NetChange:    totalEntries - totalExits,
		},
		Movements: movements,
		Pagination: Pagination{
			CurrentPage: 1,
			PerPage:     100,
			TotalItems:  len(movements),
			TotalPages:  1,
		},
	}, nil
}

func (s *service) GetSalesByPeriod(ctx context.Context, enterpriseID int64, filters ReportFilters) (SalesReportResponse, error) {
	return s.repo.GetSalesByPeriod(ctx, enterpriseID, filters)
}

func (s *service) GetSalesByProduct(ctx context.Context, enterpriseID int64, filters ReportFilters) (ProductSalesReportResponse, error) {
	if filters.Limit == 0 {
		filters.Limit = 20
	}
	return s.repo.GetSalesByProduct(ctx, enterpriseID, filters)
}

func (s *service) GetSalesByEmployee(ctx context.Context, enterpriseID int64, filters ReportFilters) (EmployeeSalesReportResponse, error) {
	return s.repo.GetSalesByEmployee(ctx, enterpriseID, filters)
}

func (s *service) GetInventoryStatus(ctx context.Context, enterpriseID int64, filters ReportFilters) (InventoryReportResponse, error) {
	return s.repo.GetInventoryStatus(ctx, enterpriseID, filters)
}

func (s *service) GetMovementHistory(ctx context.Context, enterpriseID int64, filters ReportFilters) (MovementReportResponse, error) {
	return s.repo.GetMovementHistory(ctx, enterpriseID, filters)
}

func (s *service) ExportToPDF(ctx context.Context, enterpriseID int64, reportType string, req ExportRequest) ([]byte, error) {
	// Simple text-based PDF generation
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	title := req.Title
	if title == "" {
		title = "Reporte de " + reportType
	}

	pdfContent := fmt.Sprintf(`%%PDF-1.4
1 0 obj
<< /Type /Catalog /Pages 2 0 R >>
endobj
2 0 obj
<< /Type /Pages /Kids [3 0 R] /Count 1 >>
endobj
3 0 obj
<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>
endobj
4 0 obj
<< /Length 200 >>
stream
BT
/F1 24 Tf
306 700 Td
(%s) Tj
0 -40 Td
/F1 12 Tf
(Periodo: %s - %s) Tj
0 -30 Td
(Reporte generado exitosamente) Tj
ET
endstream
endobj
5 0 obj
<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>
endobj
xref
0 6
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000115 00000 n 
0000000270 00000 n 
0000000524 00000 n 
trailer
<< /Size 6 /Root 1 0 R >>
startxref
609
%%EOF`, title, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	return []byte(pdfContent), nil
}

func (s *service) ExportToExcel(ctx context.Context, enterpriseID int64, reportType string, req ExportRequest) ([]byte, error) {
	// Generate CSV format (can be opened in Excel)
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	title := req.Title
	if title == "" {
		title = "Reporte de " + reportType
	}

	csvContent := fmt.Sprintf("Reporte,%s\n", title)
	csvContent += fmt.Sprintf("Periodo,%s - %s\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	csvContent += "\n"
	csvContent += "Reporte generado exitosamente\n"

	return []byte(csvContent), nil
}
