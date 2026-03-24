package invoices

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Invoice operations

func (r *repository) CreateInvoice(ctx context.Context, inv *Invoice) error {
	query := `
		INSERT INTO invoice (invoice_number, prefix_id, invoice_type, reference_id, reference_type, sales_order_id, customer_id, branch_id, user_id, enterprise_id, invoice_date, due_date, subtotal, discount_total, tax_exempt, taxable_amount, iva_19, iva_5, reteica, retefuente, reteica_rate, retefuente_rate, total, payment_method, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		inv.InvoiceNumber, inv.PrefixID, inv.InvoiceType, inv.ReferenceID, inv.ReferenceType,
		inv.SalesOrderID, inv.CustomerID, inv.BranchID, inv.UserID, inv.EnterpriseID,
		inv.InvoiceDate, inv.DueDate, inv.Subtotal, inv.DiscountTotal, inv.TaxExempt, inv.TaxableAmount,
		inv.Iva19, inv.Iva5, inv.Reteica, inv.Retefuente, inv.ReteicaRate, inv.RetefuenteRate,
		inv.Total, inv.PaymentMethod, inv.Status, inv.Notes,
	).Scan(&inv.ID, &inv.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}
	return nil
}

func (r *repository) GetInvoiceByID(ctx context.Context, id int64) (*Invoice, error) {
	inv := &Invoice{}
	query := `
		SELECT id, invoice_number, prefix_id, invoice_type, reference_id, reference_type, sales_order_id, customer_id, branch_id, user_id, enterprise_id, invoice_date, due_date, subtotal, discount_total, tax_exempt, taxable_amount, iva_19, iva_5, reteica, retefuente, reteica_rate, retefuente_rate, total, payment_method, status, notes, cancelled_at, cancelled_by, cancellation_reason, credit_note_id, created_at, updated_at, deleted_at
		FROM invoice WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inv.ID, &inv.InvoiceNumber, &inv.PrefixID, &inv.InvoiceType, &inv.ReferenceID, &inv.ReferenceType,
		&inv.SalesOrderID, &inv.CustomerID, &inv.BranchID, &inv.UserID, &inv.EnterpriseID,
		&inv.InvoiceDate, &inv.DueDate, &inv.Subtotal, &inv.DiscountTotal, &inv.TaxExempt, &inv.TaxableAmount,
		&inv.Iva19, &inv.Iva5, &inv.Reteica, &inv.Retefuente, &inv.ReteicaRate, &inv.RetefuenteRate,
		&inv.Total, &inv.PaymentMethod, &inv.Status, &inv.Notes, &inv.CancelledAt, &inv.CancelledBy,
		&inv.CancellationReason, &inv.CreditNoteID, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}
	return inv, nil
}

func (r *repository) GetInvoiceByNumber(ctx context.Context, invoiceNumber string, enterpriseID int64) (*Invoice, error) {
	inv := &Invoice{}
	query := `
		SELECT id, invoice_number, prefix_id, invoice_type, reference_id, reference_type, sales_order_id, customer_id, branch_id, user_id, enterprise_id, invoice_date, due_date, subtotal, discount_total, tax_exempt, taxable_amount, iva_19, iva_5, reteica, retefuente, reteica_rate, retefuente_rate, total, payment_method, status, notes, cancelled_at, cancelled_by, cancellation_reason, credit_note_id, created_at, updated_at, deleted_at
		FROM invoice WHERE invoice_number = $1 AND enterprise_id = $2 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, invoiceNumber, enterpriseID).Scan(
		&inv.ID, &inv.InvoiceNumber, &inv.PrefixID, &inv.InvoiceType, &inv.ReferenceID, &inv.ReferenceType,
		&inv.SalesOrderID, &inv.CustomerID, &inv.BranchID, &inv.UserID, &inv.EnterpriseID,
		&inv.InvoiceDate, &inv.DueDate, &inv.Subtotal, &inv.DiscountTotal, &inv.TaxExempt, &inv.TaxableAmount,
		&inv.Iva19, &inv.Iva5, &inv.Reteica, &inv.Retefuente, &inv.ReteicaRate, &inv.RetefuenteRate,
		&inv.Total, &inv.PaymentMethod, &inv.Status, &inv.Notes, &inv.CancelledAt, &inv.CancelledBy,
		&inv.CancellationReason, &inv.CreditNoteID, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get invoice by number: %w", err)
	}
	return inv, nil
}

func (r *repository) UpdateInvoice(ctx context.Context, inv *Invoice) error {
	query := `
		UPDATE invoice SET status = $1, notes = $2, updated_at = NOW()
		WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, inv.Status, inv.Notes, inv.ID)
	if err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}
	return nil
}

func (r *repository) UpdateInvoiceStatus(ctx context.Context, invoiceID int64, status string) error {
	query := `UPDATE invoice SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, invoiceID)
	if err != nil {
		return fmt.Errorf("failed to update invoice status: %w", err)
	}
	return nil
}

func (r *repository) CancelInvoice(ctx context.Context, invoiceID int64, cancelledBy int64, reason string) error {
	query := `UPDATE invoice SET status = 'CANCELLED', cancelled_at = NOW(), cancelled_by = $1, cancellation_reason = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, cancelledBy, reason, invoiceID)
	if err != nil {
		return fmt.Errorf("failed to cancel invoice: %w", err)
	}
	return nil
}

func (r *repository) GetInvoices(ctx context.Context, enterpriseID int64, filters InvoiceFilters) ([]Invoice, error) {
	query := `
		SELECT id, invoice_number, prefix_id, invoice_type, reference_id, reference_type, sales_order_id, customer_id, branch_id, user_id, enterprise_id, invoice_date, due_date, subtotal, discount_total, tax_exempt, taxable_amount, iva_19, iva_5, reteica, retefuente, reteica_rate, retefuente_rate, total, payment_method, status, notes, cancelled_at, cancelled_by, cancellation_reason, credit_note_id, created_at, updated_at, deleted_at
		FROM invoice WHERE enterprise_id = $1 AND deleted_at IS NULL`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, filters.Status)
		argPos++
	}

	if filters.InvoiceType != "" {
		query += fmt.Sprintf(" AND invoice_type = $%d", argPos)
		args = append(args, filters.InvoiceType)
		argPos++
	}

	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	if filters.CustomerID != nil {
		query += fmt.Sprintf(" AND customer_id = $%d", argPos)
		args = append(args, *filters.CustomerID)
		argPos++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND invoice_date >= $%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND invoice_date <= $%d", argPos)
		args = append(args, *filters.EndDate)
		argPos++
	}

	query += " ORDER BY invoice_date DESC, id DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filters.Limit)
		argPos++

		offset := (filters.Page - 1) * filters.Limit
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoices: %w", err)
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		if err := rows.Scan(
			&inv.ID, &inv.InvoiceNumber, &inv.PrefixID, &inv.InvoiceType, &inv.ReferenceID, &inv.ReferenceType,
			&inv.SalesOrderID, &inv.CustomerID, &inv.BranchID, &inv.UserID, &inv.EnterpriseID,
			&inv.InvoiceDate, &inv.DueDate, &inv.Subtotal, &inv.DiscountTotal, &inv.TaxExempt, &inv.TaxableAmount,
			&inv.Iva19, &inv.Iva5, &inv.Reteica, &inv.Retefuente, &inv.ReteicaRate, &inv.RetefuenteRate,
			&inv.Total, &inv.PaymentMethod, &inv.Status, &inv.Notes, &inv.CancelledAt, &inv.CancelledBy,
			&inv.CancellationReason, &inv.CreditNoteID, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (r *repository) CountInvoices(ctx context.Context, enterpriseID int64, filters InvoiceFilters) (int, error) {
	query := `SELECT COUNT(*) FROM invoice WHERE enterprise_id = $1 AND deleted_at IS NULL`
	var count int
	err := r.db.QueryRowContext(ctx, query, enterpriseID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count invoices: %w", err)
	}
	return count, nil
}

// Invoice Item operations

func (r *repository) CreateInvoiceItem(ctx context.Context, item *InvoiceItem) error {
	query := `
		INSERT INTO invoice_item (invoice_id, product_id, product_name, product_sku, quantity, unit_price, discount_amount, tax_rate, tax_amount, line_total)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		item.InvoiceID, item.ProductID, item.ProductName, item.ProductSKU,
		item.Quantity, item.UnitPrice, item.DiscountAmount, item.TaxRate, item.TaxAmount, item.LineTotal,
	).Scan(&item.ID, &item.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create invoice item: %w", err)
	}
	return nil
}

func (r *repository) GetInvoiceItems(ctx context.Context, invoiceID int64) ([]InvoiceItem, error) {
	query := `
		SELECT id, invoice_id, product_id, product_name, product_sku, quantity, unit_price, discount_amount, tax_rate, tax_amount, line_total, created_at
		FROM invoice_item WHERE invoice_id = $1
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}
	defer rows.Close()

	var items []InvoiceItem
	for rows.Next() {
		var item InvoiceItem
		if err := rows.Scan(
			&item.ID, &item.InvoiceID, &item.ProductID, &item.ProductName, &item.ProductSKU,
			&item.Quantity, &item.UnitPrice, &item.DiscountAmount, &item.TaxRate, &item.TaxAmount,
			&item.LineTotal, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// Invoice Prefix operations

func (r *repository) GetInvoicePrefix(ctx context.Context, branchID int64, prefix string) (*InvoicePrefix, error) {
	p := &InvoicePrefix{}
	query := `
		SELECT id, prefix, branch_id, enterprise_id, current_number, resolution_number, resolution_date, valid_from, valid_until, description, is_active, created_at, updated_at
		FROM invoice_prefix WHERE branch_id = $1 AND prefix = $2 AND is_active = TRUE`

	err := r.db.QueryRowContext(ctx, query, branchID, prefix).Scan(
		&p.ID, &p.Prefix, &p.BranchID, &p.EnterpriseID, &p.CurrentNumber,
		&p.ResolutionNumber, &p.ResolutionDate, &p.ValidFrom, &p.ValidUntil,
		&p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get invoice prefix: %w", err)
	}
	return p, nil
}

func (r *repository) GetInvoicePrefixByID(ctx context.Context, id int64) (*InvoicePrefix, error) {
	p := &InvoicePrefix{}
	query := `
		SELECT id, prefix, branch_id, enterprise_id, current_number, resolution_number, resolution_date, valid_from, valid_until, description, is_active, created_at, updated_at
		FROM invoice_prefix WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Prefix, &p.BranchID, &p.EnterpriseID, &p.CurrentNumber,
		&p.ResolutionNumber, &p.ResolutionDate, &p.ValidFrom, &p.ValidUntil,
		&p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get invoice prefix: %w", err)
	}
	return p, nil
}

func (r *repository) CreateInvoicePrefix(ctx context.Context, p *InvoicePrefix) error {
	query := `
		INSERT INTO invoice_prefix (prefix, branch_id, enterprise_id, current_number, resolution_number, resolution_date, valid_from, valid_until, description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		p.Prefix, p.BranchID, p.EnterpriseID, p.CurrentNumber,
		p.ResolutionNumber, p.ResolutionDate, p.ValidFrom, p.ValidUntil, p.Description, p.IsActive,
	).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create invoice prefix: %w", err)
	}
	return nil
}

func (r *repository) UpdateInvoicePrefix(ctx context.Context, p *InvoicePrefix) error {
	query := `UPDATE invoice_prefix SET current_number = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, p.CurrentNumber, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update invoice prefix: %w", err)
	}
	return nil
}

func (r *repository) ListInvoicePrefixes(ctx context.Context, enterpriseID int64) ([]InvoicePrefix, error) {
	query := `
		SELECT id, prefix, branch_id, enterprise_id, current_number, resolution_number, resolution_date, valid_from, valid_until, is_active, created_at, updated_at
		FROM invoice_prefix WHERE enterprise_id = $1 AND is_active = TRUE
		ORDER BY prefix`

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoice prefixes: %w", err)
	}
	defer rows.Close()

	var prefixes []InvoicePrefix
	for rows.Next() {
		var p InvoicePrefix
		if err := rows.Scan(
			&p.ID, &p.Prefix, &p.BranchID, &p.EnterpriseID, &p.CurrentNumber,
			&p.ResolutionNumber, &p.ResolutionDate, &p.ValidFrom, &p.ValidUntil,
			&p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		prefixes = append(prefixes, p)
	}
	return prefixes, nil
}

// Invoice Log operations

func (r *repository) CreateInvoiceLog(ctx context.Context, log *InvoiceLog) error {
	query := `
		INSERT INTO invoice_log (invoice_id, action, user_id, details)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, log.InvoiceID, log.Action, log.UserID, log.Details).Scan(&log.ID, &log.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create invoice log: %w", err)
	}
	return nil
}

func (r *repository) GetInvoiceLogs(ctx context.Context, invoiceID int64) ([]InvoiceLog, error) {
	query := `
		SELECT id, invoice_id, action, user_id, details, created_at
		FROM invoice_log WHERE invoice_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice logs: %w", err)
	}
	defer rows.Close()

	var logs []InvoiceLog
	for rows.Next() {
		var log InvoiceLog
		if err := rows.Scan(&log.ID, &log.InvoiceID, &log.Action, &log.UserID, &log.Details, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// Service implementation
type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) GenerateInvoiceFromSale(ctx context.Context, salesOrderID int64, prefixID int64) (*Invoice, error) {
	// Get prefix info
	prefix, err := s.repo.GetInvoicePrefixByID(ctx, prefixID)
	if err != nil {
		return nil, fmt.Errorf("invoice prefix not found: %w", err)
	}

	// Generate invoice number
	prefix.CurrentNumber++
	invoiceNumber := fmt.Sprintf("%s-%06d", prefix.Prefix, prefix.CurrentNumber)

	// Create invoice (status DRAFT initially)
	inv := &Invoice{
		InvoiceNumber: invoiceNumber,
		PrefixID:      prefixID,
		InvoiceType:   InvoiceTypeSale,
		SalesOrderID:  &salesOrderID,
		Status:        InvoiceStatusDraft,
		InvoiceDate:   time.Now(),
	}

	if err := s.repo.CreateInvoice(ctx, inv); err != nil {
		return nil, err
	}

	// Update prefix sequence
	if err := s.repo.UpdateInvoicePrefix(ctx, prefix); err != nil {
		return nil, err
	}

	return inv, nil
}

func (s *service) GenerateInvoice(ctx context.Context, inv *Invoice) error {
	if inv.Status == "" {
		inv.Status = InvoiceStatusDraft
	}
	if inv.InvoiceType == "" {
		inv.InvoiceType = InvoiceTypeSale
	}
	if inv.InvoiceDate.IsZero() {
		inv.InvoiceDate = time.Now()
	}
	return s.repo.CreateInvoice(ctx, inv)
}

func (s *service) GetInvoice(ctx context.Context, id int64) (*Invoice, error) {
	inv, err := s.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetInvoiceItems(ctx, id)
	if err != nil {
		return nil, err
	}

	inv.Items = items
	return inv, nil
}

func (s *service) GetInvoiceByNumber(ctx context.Context, invoiceNumber string, enterpriseID int64) (*Invoice, error) {
	inv, err := s.repo.GetInvoiceByNumber(ctx, invoiceNumber, enterpriseID)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetInvoiceItems(ctx, inv.ID)
	if err != nil {
		return nil, err
	}

	inv.Items = items
	return inv, nil
}

func (s *service) GetInvoices(ctx context.Context, enterpriseID int64, filters InvoiceFilters) ([]Invoice, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	return s.repo.GetInvoices(ctx, enterpriseID, filters)
}

func (s *service) IssueInvoice(ctx context.Context, invoiceID int64) error {
	inv, err := s.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	if inv.Status != InvoiceStatusDraft {
		return fmt.Errorf("invoice can only be issued from DRAFT status")
	}

	if err := s.repo.UpdateInvoiceStatus(ctx, invoiceID, InvoiceStatusIssued); err != nil {
		return err
	}

	// Log the action
	s.repo.CreateInvoiceLog(ctx, &InvoiceLog{
		InvoiceID: invoiceID,
		Action:    "ISSUED",
		UserID:    inv.UserID,
		Details:   "Invoice issued",
	})

	return nil
}

func (s *service) CancelInvoice(ctx context.Context, invoiceID int64, reason string) error {
	inv, err := s.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	if inv.Status == InvoiceStatusCancelled {
		return fmt.Errorf("invoice already cancelled")
	}

	if err := s.repo.CancelInvoice(ctx, invoiceID, inv.UserID, reason); err != nil {
		return err
	}

	// Log the action
	s.repo.CreateInvoiceLog(ctx, &InvoiceLog{
		InvoiceID: invoiceID,
		Action:    "CANCELLED",
		UserID:    inv.UserID,
		Details:   "Invoice cancelled: " + reason,
	})

	return nil
}

func (s *service) CreateInvoicePrefix(ctx context.Context, prefix *InvoicePrefix) error {
	prefix.CurrentNumber = 0
	prefix.IsActive = true
	return s.repo.CreateInvoicePrefix(ctx, prefix)
}

func (s *service) GetInvoicePrefixes(ctx context.Context, enterpriseID int64) ([]InvoicePrefix, error) {
	return s.repo.ListInvoicePrefixes(ctx, enterpriseID)
}

func (s *service) GetInvoicePrefix(ctx context.Context, branchID int64, prefix string) (*InvoicePrefix, error) {
	return s.repo.GetInvoicePrefix(ctx, branchID, prefix)
}

func (s *service) GetInvoiceLogs(ctx context.Context, invoiceID int64) ([]InvoiceLog, error) {
	return s.repo.GetInvoiceLogs(ctx, invoiceID)
}

// Helper to format invoice number
func formatInvoiceNumber(prefix string, sequence int64) string {
	return fmt.Sprintf("%s-%06d", prefix, sequence)
}

// Helper to parse invoice number
func parseInvoiceNumber(invoiceNumber string) (prefix string, sequence int64, err error) {
	_, err = fmt.Sscanf(invoiceNumber, "%s-%d", &prefix, &sequence)
	if err != nil {
		return "", 0, err
	}
	return prefix, sequence, nil
}

// Need strconv for parsing
var _ = strconv.ParseInt
