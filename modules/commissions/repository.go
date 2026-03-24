package commissions

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// ─── CommissionRule ────────────────────────────────────────────────────────────

func (r *repository) CreateRule(ctx context.Context, rule *CommissionRule) (int64, error) {
	rule.CreatedAt = time.Now()
	rule.IsActive = true

	return rule.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO commission_rule (name, commission_type, employee_id, product_id, category_id,
			value, min_sale_amount, start_date, end_date, is_active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		rule.Name, rule.CommissionType, rule.EmployeeID, rule.ProductID, rule.CategoryID,
		rule.Value, rule.MinSaleAmount, rule.StartDate, rule.EndDate, rule.IsActive, rule.CreatedAt,
	).Scan(&rule.ID)
}

func (r *repository) GetRuleByID(ctx context.Context, id int64) (*CommissionRule, error) {
	var rule CommissionRule
	var employeeID, productID, categoryID sql.NullInt64
	var startDate, endDate, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, commission_type, employee_id, product_id, category_id,
			value, min_sale_amount, start_date, end_date, is_active, created_at, updated_at
		 FROM commission_rule WHERE id = $1`, id,
	).Scan(&rule.ID, &rule.Name, &rule.CommissionType, &employeeID, &productID, &categoryID,
		&rule.Value, &rule.MinSaleAmount, &startDate, &endDate, &rule.IsActive, &rule.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if employeeID.Valid {
		rule.EmployeeID = &employeeID.Int64
	}
	if productID.Valid {
		rule.ProductID = &productID.Int64
	}
	if categoryID.Valid {
		rule.CategoryID = &categoryID.Int64
	}
	if startDate.Valid {
		rule.StartDate = &startDate.Time
	}
	if endDate.Valid {
		rule.EndDate = &endDate.Time
	}
	if updatedAt.Valid {
		rule.UpdatedAt = &updatedAt.Time
	}
	return &rule, nil
}

func (r *repository) UpdateRule(ctx context.Context, rule *CommissionRule) error {
	now := time.Now()
	rule.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE commission_rule SET name=$1, value=$2, min_sale_amount=$3, is_active=$4, updated_at=$5 WHERE id=$6`,
		rule.Name, rule.Value, rule.MinSaleAmount, rule.IsActive, now, rule.ID)
	return err
}

func (r *repository) DeleteRule(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE commission_rule SET is_active = FALSE WHERE id = $1`, id)
	return err
}

func (r *repository) ListRules(ctx context.Context, activeOnly bool) ([]CommissionRule, error) {
	query := `SELECT id, name, commission_type, employee_id, product_id, category_id,
			value, min_sale_amount, start_date, end_date, is_active, created_at, updated_at
		 FROM commission_rule`
	if activeOnly {
		query += " WHERE is_active = TRUE"
	}
	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []CommissionRule
	for rows.Next() {
		var rule CommissionRule
		var employeeID, productID, categoryID sql.NullInt64
		var startDate, endDate, updatedAt sql.NullTime
		if err := rows.Scan(&rule.ID, &rule.Name, &rule.CommissionType, &employeeID, &productID, &categoryID,
			&rule.Value, &rule.MinSaleAmount, &startDate, &endDate, &rule.IsActive, &rule.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if employeeID.Valid {
			rule.EmployeeID = &employeeID.Int64
		}
		if productID.Valid {
			rule.ProductID = &productID.Int64
		}
		if categoryID.Valid {
			rule.CategoryID = &categoryID.Int64
		}
		if startDate.Valid {
			rule.StartDate = &startDate.Time
		}
		if endDate.Valid {
			rule.EndDate = &endDate.Time
		}
		if updatedAt.Valid {
			rule.UpdatedAt = &updatedAt.Time
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *repository) GetApplicableRules(ctx context.Context, employeeID, productID, categoryID *int64, saleAmount float64) ([]CommissionRule, error) {
	query := `SELECT id, name, commission_type, employee_id, product_id, category_id,
			value, min_sale_amount, start_date, end_date, is_active, created_at, updated_at
		 FROM commission_rule
		 WHERE is_active = TRUE AND min_sale_amount <= $1`
	args := []interface{}{saleAmount}
	argIndex := 2

	query += " AND (employee_id IS NULL OR employee_id = $%d)"
	query = fmt.Sprintf(query, argIndex)
	args = append(args, employeeID)
	argIndex++

	if productID != nil {
		query += fmt.Sprintf(" AND (product_id IS NULL OR product_id = $%d)", argIndex)
		args = append(args, *productID)
		argIndex++
	}
	if categoryID != nil {
		query += fmt.Sprintf(" AND (category_id IS NULL OR category_id = $%d)", argIndex)
		args = append(args, *categoryID)
		argIndex++
	}

	query += " ORDER BY product_id DESC NULLS LAST, category_id DESC NULLS LAST, employee_id DESC NULLS LAST"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []CommissionRule
	for rows.Next() {
		var rule CommissionRule
		var empID, prodID, catID sql.NullInt64
		var startDate, endDate, updatedAt sql.NullTime
		if err := rows.Scan(&rule.ID, &rule.Name, &rule.CommissionType, &empID, &prodID, &catID,
			&rule.Value, &rule.MinSaleAmount, &startDate, &endDate, &rule.IsActive, &rule.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if empID.Valid {
			rule.EmployeeID = &empID.Int64
		}
		if prodID.Valid {
			rule.ProductID = &prodID.Int64
		}
		if catID.Valid {
			rule.CategoryID = &catID.Int64
		}
		if startDate.Valid {
			rule.StartDate = &startDate.Time
		}
		if endDate.Valid {
			rule.EndDate = &endDate.Time
		}
		if updatedAt.Valid {
			rule.UpdatedAt = &updatedAt.Time
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// ─── Commission ───────────────────────────────────────────────────────────────

func (r *repository) CreateCommission(ctx context.Context, c *Commission) (int64, error) {
	c.CreatedAt = time.Now()
	c.Status = StatusPending

	return c.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO commission (sales_order_id, employee_id, branch_id, rule_id, sale_amount,
			profit_margin, commission_type, commission_rate, commission_amount, status, settlement_period, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`,
		c.SalesOrderID, c.EmployeeID, c.BranchID, c.RuleID, c.SaleAmount,
		c.ProfitMargin, c.CommissionType, c.CommissionRate, c.CommissionAmount, c.Status, c.SettlementPeriod, c.Notes, c.CreatedAt,
	).Scan(&c.ID)
}

func (r *repository) GetCommissionByID(ctx context.Context, id int64) (*Commission, error) {
	var c Commission
	var ruleID sql.NullInt64
	var profitMargin sql.NullFloat64
	var settledAt sql.NullTime
	var settledBy sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, sales_order_id, employee_id, branch_id, rule_id, sale_amount,
			profit_margin, commission_type, commission_rate, commission_amount, status,
			settled_at, settled_by, settlement_period, notes, created_at
		 FROM commission WHERE id = $1`, id,
	).Scan(&c.ID, &c.SalesOrderID, &c.EmployeeID, &c.BranchID, &ruleID, &c.SaleAmount,
		&profitMargin, &c.CommissionType, &c.CommissionRate, &c.CommissionAmount, &c.Status,
		&settledAt, &settledBy, &c.SettlementPeriod, &c.Notes, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	if ruleID.Valid {
		c.RuleID = &ruleID.Int64
	}
	if profitMargin.Valid {
		c.ProfitMargin = &profitMargin.Float64
	}
	if settledAt.Valid {
		c.SettledAt = &settledAt.Time
	}
	if settledBy.Valid {
		c.SettledBy = &settledBy.Int64
	}
	return &c, nil
}

func (r *repository) ListCommissions(ctx context.Context, employeeID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Commission, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	baseQuery := "FROM commission WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if employeeID != nil {
		baseQuery += fmt.Sprintf(" AND employee_id = $%d", argIndex)
		args = append(args, *employeeID)
		argIndex++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if startDate != nil {
		baseQuery += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		baseQuery += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT id, sales_order_id, employee_id, branch_id, rule_id, sale_amount, profit_margin, commission_type, commission_rate, commission_amount, status, settled_at, settled_by, settlement_period, notes, created_at %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var commissions []Commission
	for rows.Next() {
		var c Commission
		var ruleID sql.NullInt64
		var profitMargin sql.NullFloat64
		var settledAt sql.NullTime
		var settledBy sql.NullInt64
		if err := rows.Scan(&c.ID, &c.SalesOrderID, &c.EmployeeID, &c.BranchID, &ruleID, &c.SaleAmount,
			&profitMargin, &c.CommissionType, &c.CommissionRate, &c.CommissionAmount, &c.Status,
			&settledAt, &settledBy, &c.SettlementPeriod, &c.Notes, &c.CreatedAt); err != nil {
			return nil, 0, err
		}
		if ruleID.Valid {
			c.RuleID = &ruleID.Int64
		}
		if profitMargin.Valid {
			c.ProfitMargin = &profitMargin.Float64
		}
		if settledAt.Valid {
			c.SettledAt = &settledAt.Time
		}
		if settledBy.Valid {
			c.SettledBy = &settledBy.Int64
		}
		commissions = append(commissions, c)
	}
	return commissions, total, nil
}

func (r *repository) SettleCommissions(ctx context.Context, ids []int64, settledBy int64, period string, notes string) error {
	if len(ids) == 0 {
		return nil
	}

	// Build placeholders for IDs
	placeholders := ""
	args := []interface{}{settledBy, period, notes}
	for i, id := range ids {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+4)
		args = append(args, id)
	}

	query := fmt.Sprintf(`UPDATE commission SET status = 'SETTLED', settled_at = NOW(), 
		settled_by = $1, settlement_period = $2, notes = $3 
		WHERE id IN (%s) AND status = 'PENDING'`, placeholders)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *repository) CancelCommission(ctx context.Context, id int64, notes string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE commission SET status = 'CANCELLED', notes = $1 WHERE id = $2 AND status = 'PENDING'`,
		notes, id)
	return err
}

// ─── Reporting ────────────────────────────────────────────────────────────────

func (r *repository) GetCommissionSummary(ctx context.Context, employeeID *int64, startDate, endDate *time.Time) ([]CommissionSummary, error) {
	query := `
		SELECT c.employee_id, tp.name as employee_name,
			COALESCE(SUM(c.sale_amount), 0) as total_sales,
			COALESCE(SUM(c.commission_amount), 0) as total_commissions,
			COALESCE(SUM(CASE WHEN c.status = 'PENDING' THEN c.commission_amount ELSE 0 END), 0) as pending_amount,
			COALESCE(SUM(CASE WHEN c.status = 'SETTLED' THEN c.commission_amount ELSE 0 END), 0) as settled_amount,
			COUNT(*) as sales_count
		FROM commission c
		JOIN third_parties tp ON tp.id = c.employee_id
		WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if employeeID != nil {
		query += fmt.Sprintf(" AND c.employee_id = $%d", argIndex)
		args = append(args, *employeeID)
		argIndex++
	}
	if startDate != nil {
		query += fmt.Sprintf(" AND c.created_at >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		query += fmt.Sprintf(" AND c.created_at <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	query += " GROUP BY c.employee_id, tp.name ORDER BY total_commissions DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []CommissionSummary
	for rows.Next() {
		var s CommissionSummary
		if err := rows.Scan(&s.EmployeeID, &s.EmployeeName, &s.TotalSales, &s.TotalCommissions,
			&s.PendingAmount, &s.SettledAmount, &s.SalesCount); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}

func (r *repository) GetCommissionTotals(ctx context.Context, employeeID *int64, startDate, endDate *time.Time) (totalSales, totalCommissions, pendingAmount, settledAmount float64, err error) {
	query := `
		SELECT COALESCE(SUM(sale_amount), 0), COALESCE(SUM(commission_amount), 0),
			COALESCE(SUM(CASE WHEN status = 'PENDING' THEN commission_amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'SETTLED' THEN commission_amount ELSE 0 END), 0)
		FROM commission WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if employeeID != nil {
		query += fmt.Sprintf(" AND employee_id = $%d", argIndex)
		args = append(args, *employeeID)
		argIndex++
	}
	if startDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&totalSales, &totalCommissions, &pendingAmount, &settledAmount)
	return
}
