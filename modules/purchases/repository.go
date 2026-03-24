package purchases

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

// ─── Purchase Order ────────────────────────────────────────────────────────────

func (r *repository) CreatePurchaseOrder(ctx context.Context, po *PurchaseOrder) (int64, error) {
	po.CreatedAt = time.Now()
	po.Status = StatusPending
	if po.OrderDate.IsZero() {
		po.OrderDate = time.Now()
	}

	query := `
		INSERT INTO purchase_order (order_number, supplier_id, branch_id, user_id, order_date, expected_date,
			status, subtotal, discount_total, tax_total, total, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		po.OrderNumber, po.SupplierID, po.BranchID, po.UserID, po.OrderDate, po.ExpectedDate,
		po.Status, po.Subtotal, po.DiscountTotal, po.TaxTotal, po.Total, po.Notes, po.CreatedAt,
	).Scan(&id)
	return id, err
}

func (r *repository) GetPurchaseOrderByID(ctx context.Context, id int64) (*PurchaseOrder, error) {
	var po PurchaseOrder
	var expectedDate, updatedAt, deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, order_number, supplier_id, branch_id, user_id, order_date, expected_date,
			status, subtotal, discount_total, tax_total, total, notes, created_at, updated_at, deleted_at
		 FROM purchase_order WHERE id = $1`,
		id,
	).Scan(&po.ID, &po.OrderNumber, &po.SupplierID, &po.BranchID, &po.UserID, &po.OrderDate,
		&expectedDate, &po.Status, &po.Subtotal, &po.DiscountTotal, &po.TaxTotal, &po.Total,
		&po.Notes, &po.CreatedAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	if expectedDate.Valid {
		po.ExpectedDate = &expectedDate.Time
	}
	if updatedAt.Valid {
		po.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		po.DeletedAt = &deletedAt.Time
	}
	return &po, nil
}

func (r *repository) GetPurchaseOrderByNumber(ctx context.Context, orderNumber string) (*PurchaseOrder, error) {
	var po PurchaseOrder
	var expectedDate, updatedAt, deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, order_number, supplier_id, branch_id, user_id, order_date, expected_date,
			status, subtotal, discount_total, tax_total, total, notes, created_at, updated_at, deleted_at
		 FROM purchase_order WHERE order_number = $1`,
		orderNumber,
	).Scan(&po.ID, &po.OrderNumber, &po.SupplierID, &po.BranchID, &po.UserID, &po.OrderDate,
		&expectedDate, &po.Status, &po.Subtotal, &po.DiscountTotal, &po.TaxTotal, &po.Total,
		&po.Notes, &po.CreatedAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	if expectedDate.Valid {
		po.ExpectedDate = &expectedDate.Time
	}
	if updatedAt.Valid {
		po.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		po.DeletedAt = &deletedAt.Time
	}
	return &po, nil
}

func (r *repository) UpdatePurchaseOrder(ctx context.Context, po *PurchaseOrder) error {
	now := time.Now()
	po.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE purchase_order SET status=$1, subtotal=$2, discount_total=$3, tax_total=$4, total=$5,
			notes=$6, updated_at=$7 WHERE id=$8`,
		po.Status, po.Subtotal, po.DiscountTotal, po.TaxTotal, po.Total, po.Notes, now, po.ID)
	return err
}

func (r *repository) ListPurchaseOrders(ctx context.Context, supplierID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]PurchaseOrder, int64, error) {
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

	baseQuery := "FROM purchase_order WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	if supplierID != nil {
		baseQuery += fmt.Sprintf(" AND supplier_id = $%d", argIndex)
		args = append(args, *supplierID)
		argIndex++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if startDate != nil {
		baseQuery += fmt.Sprintf(" AND order_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		baseQuery += fmt.Sprintf(" AND order_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT id, order_number, supplier_id, branch_id, user_id, order_date, expected_date, status, subtotal, discount_total, tax_total, total, notes, created_at, updated_at, deleted_at %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []PurchaseOrder
	for rows.Next() {
		var po PurchaseOrder
		var expectedDate, updatedAt, deletedAt sql.NullTime
		if err := rows.Scan(&po.ID, &po.OrderNumber, &po.SupplierID, &po.BranchID, &po.UserID, &po.OrderDate,
			&expectedDate, &po.Status, &po.Subtotal, &po.DiscountTotal, &po.TaxTotal, &po.Total,
			&po.Notes, &po.CreatedAt, &updatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		if expectedDate.Valid {
			po.ExpectedDate = &expectedDate.Time
		}
		if updatedAt.Valid {
			po.UpdatedAt = &updatedAt.Time
		}
		if deletedAt.Valid {
			po.DeletedAt = &deletedAt.Time
		}
		orders = append(orders, po)
	}
	return orders, total, nil
}

// ─── Purchase Order Items ──────────────────────────────────────────────────────

func (r *repository) CreatePurchaseOrderItem(ctx context.Context, item *PurchaseOrderItem) error {
	item.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO purchase_order_item (purchase_order_id, product_id, quantity, received_quantity, unit_cost, discount_amount, tax_rate, line_total, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		item.PurchaseOrderID, item.ProductID, item.Quantity, item.ReceivedQuantity, item.UnitCost,
		item.DiscountAmount, item.TaxRate, item.LineTotal, item.CreatedAt).Scan(&item.ID)
}

func (r *repository) GetPurchaseOrderItems(ctx context.Context, orderID int64) ([]PurchaseOrderItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, purchase_order_id, product_id, quantity, received_quantity, unit_cost, discount_amount, tax_rate, line_total, created_at
		 FROM purchase_order_item WHERE purchase_order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PurchaseOrderItem
	for rows.Next() {
		var item PurchaseOrderItem
		if err := rows.Scan(&item.ID, &item.PurchaseOrderID, &item.ProductID, &item.Quantity, &item.ReceivedQuantity,
			&item.UnitCost, &item.DiscountAmount, &item.TaxRate, &item.LineTotal, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *repository) UpdatePurchaseOrderItem(ctx context.Context, item *PurchaseOrderItem) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE purchase_order_item SET received_quantity = $1 WHERE id = $2`,
		item.ReceivedQuantity, item.ID)
	return err
}

// ─── Purchase ──────────────────────────────────────────────────────────────────

func (r *repository) CreatePurchase(ctx context.Context, p *Purchase) (int64, error) {
	p.CreatedAt = time.Now()
	p.PurchaseDate = time.Now()
	p.PendingAmount = p.Total - p.PaidAmount

	if p.PaidAmount >= p.Total {
		p.Status = StatusCompleted
	} else {
		p.Status = StatusPartial
	}

	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO purchase (purchase_number, purchase_order_id, supplier_id, branch_id, user_id, purchase_date,
			status, subtotal, discount_total, tax_total, total, paid_amount, pending_amount, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING id`,
		p.PurchaseNumber, p.PurchaseOrderID, p.SupplierID, p.BranchID, p.UserID, p.PurchaseDate,
		p.Status, p.Subtotal, p.DiscountTotal, p.TaxTotal, p.Total, p.PaidAmount, p.PendingAmount,
		p.Notes, p.CreatedAt).Scan(&id)
	return id, err
}

func (r *repository) GetPurchaseByID(ctx context.Context, id int64) (*Purchase, error) {
	var p Purchase
	var purchaseOrderID sql.NullInt64
	var updatedAt, deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, purchase_number, purchase_order_id, supplier_id, branch_id, user_id, purchase_date,
			status, subtotal, discount_total, tax_total, total, paid_amount, pending_amount, notes, created_at, updated_at, deleted_at
		 FROM purchase WHERE id = $1`, id,
	).Scan(&p.ID, &p.PurchaseNumber, &purchaseOrderID, &p.SupplierID, &p.BranchID, &p.UserID, &p.PurchaseDate,
		&p.Status, &p.Subtotal, &p.DiscountTotal, &p.TaxTotal, &p.Total, &p.PaidAmount, &p.PendingAmount,
		&p.Notes, &p.CreatedAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	if purchaseOrderID.Valid {
		p.PurchaseOrderID = &purchaseOrderID.Int64
	}
	if updatedAt.Valid {
		p.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		p.DeletedAt = &deletedAt.Time
	}
	return &p, nil
}

func (r *repository) GetPurchaseByNumber(ctx context.Context, purchaseNumber string) (*Purchase, error) {
	var p Purchase
	var purchaseOrderID sql.NullInt64
	var updatedAt, deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, purchase_number, purchase_order_id, supplier_id, branch_id, user_id, purchase_date,
			status, subtotal, discount_total, tax_total, total, paid_amount, pending_amount, notes, created_at, updated_at, deleted_at
		 FROM purchase WHERE purchase_number = $1`, purchaseNumber,
	).Scan(&p.ID, &p.PurchaseNumber, &purchaseOrderID, &p.SupplierID, &p.BranchID, &p.UserID, &p.PurchaseDate,
		&p.Status, &p.Subtotal, &p.DiscountTotal, &p.TaxTotal, &p.Total, &p.PaidAmount, &p.PendingAmount,
		&p.Notes, &p.CreatedAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	if purchaseOrderID.Valid {
		p.PurchaseOrderID = &purchaseOrderID.Int64
	}
	if updatedAt.Valid {
		p.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		p.DeletedAt = &deletedAt.Time
	}
	return &p, nil
}

func (r *repository) UpdatePurchase(ctx context.Context, p *Purchase) error {
	now := time.Now()
	p.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE purchase SET status=$1, paid_amount=$2, pending_amount=$3, notes=$4, updated_at=$5 WHERE id=$6`,
		p.Status, p.PaidAmount, p.PendingAmount, p.Notes, now, p.ID)
	return err
}

func (r *repository) ListPurchases(ctx context.Context, supplierID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Purchase, int64, error) {
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

	baseQuery := "FROM purchase WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	if supplierID != nil {
		baseQuery += fmt.Sprintf(" AND supplier_id = $%d", argIndex)
		args = append(args, *supplierID)
		argIndex++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if startDate != nil {
		baseQuery += fmt.Sprintf(" AND purchase_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		baseQuery += fmt.Sprintf(" AND purchase_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT id, purchase_number, purchase_order_id, supplier_id, branch_id, user_id, purchase_date, status, subtotal, discount_total, tax_total, total, paid_amount, pending_amount, notes, created_at, updated_at, deleted_at %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var purchases []Purchase
	for rows.Next() {
		var p Purchase
		var purchaseOrderID sql.NullInt64
		var updatedAt, deletedAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.PurchaseNumber, &purchaseOrderID, &p.SupplierID, &p.BranchID, &p.UserID,
			&p.PurchaseDate, &p.Status, &p.Subtotal, &p.DiscountTotal, &p.TaxTotal, &p.Total,
			&p.PaidAmount, &p.PendingAmount, &p.Notes, &p.CreatedAt, &updatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		if purchaseOrderID.Valid {
			p.PurchaseOrderID = &purchaseOrderID.Int64
		}
		if updatedAt.Valid {
			p.UpdatedAt = &updatedAt.Time
		}
		if deletedAt.Valid {
			p.DeletedAt = &deletedAt.Time
		}
		purchases = append(purchases, p)
	}
	return purchases, total, nil
}

// ─── Purchase Items ────────────────────────────────────────────────────────────

func (r *repository) CreatePurchaseItem(ctx context.Context, item *PurchaseItem) error {
	item.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO purchase_item (purchase_id, product_id, quantity, unit_cost, discount_amount, tax_rate, line_total, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		item.PurchaseID, item.ProductID, item.Quantity, item.UnitCost, item.DiscountAmount, item.TaxRate,
		item.LineTotal, item.CreatedAt).Scan(&item.ID)
}

func (r *repository) GetPurchaseItems(ctx context.Context, purchaseID int64) ([]PurchaseItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, purchase_id, product_id, quantity, unit_cost, discount_amount, tax_rate, line_total, created_at
		 FROM purchase_item WHERE purchase_id = $1`, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PurchaseItem
	for rows.Next() {
		var item PurchaseItem
		if err := rows.Scan(&item.ID, &item.PurchaseID, &item.ProductID, &item.Quantity, &item.UnitCost,
			&item.DiscountAmount, &item.TaxRate, &item.LineTotal, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// ─── Purchase Payments ────────────────────────────────────────────────────────

func (r *repository) CreatePurchasePayment(ctx context.Context, payment *PurchasePayment) error {
	payment.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO purchase_payment (purchase_id, payment_method, amount, reference_number, notes, user_id, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		payment.PurchaseID, payment.PaymentMethod, payment.Amount, payment.ReferenceNumber,
		payment.Notes, payment.UserID, payment.CreatedAt).Scan(&payment.ID)
}

func (r *repository) GetPurchasePayments(ctx context.Context, purchaseID int64) ([]PurchasePayment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, purchase_id, payment_method, amount, reference_number, notes, user_id, created_at
		 FROM purchase_payment WHERE purchase_id = $1 ORDER BY created_at`, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []PurchasePayment
	for rows.Next() {
		var pp PurchasePayment
		var refNum sql.NullString
		if err := rows.Scan(&pp.ID, &pp.PurchaseID, &pp.PaymentMethod, &pp.Amount, &refNum, &pp.Notes, &pp.UserID, &pp.CreatedAt); err != nil {
			return nil, err
		}
		if refNum.Valid {
			pp.ReferenceNumber = &refNum.String
		}
		payments = append(payments, pp)
	}
	return payments, nil
}

func (r *repository) GetSupplierSummary(ctx context.Context, supplierID int64) (*SupplierSummary, error) {
	var summary SupplierSummary
	summary.SupplierID = supplierID

	err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(tp.name, ''), COUNT(p.id), COALESCE(SUM(p.total),0), COALESCE(SUM(p.paid_amount),0), COALESCE(SUM(p.pending_amount),0)
		 FROM third_parties tp
		 LEFT JOIN purchase p ON p.supplier_id = tp.id AND p.deleted_at IS NULL
		 WHERE tp.id = $1
		 GROUP BY tp.id, tp.name`,
		supplierID,
	).Scan(&summary.SupplierName, &summary.TotalPurchases, &summary.TotalAmount, &summary.PaidAmount, &summary.PendingAmount)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
