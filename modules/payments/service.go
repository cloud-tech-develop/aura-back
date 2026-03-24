package payments

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

// Payment operations

func (r *repository) CreatePayment(ctx context.Context, p *Payment) error {
	query := `
		INSERT INTO payment (payment_type, reference_id, reference_type, payment_method, amount, reference_number, bank_name, card_type, card_last_digits, authorization_code, change_amount, cash_drawer_id, branch_id, enterprise_id, user_id, notes, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		p.PaymentType, p.ReferenceID, p.ReferenceType, p.PaymentMethod, p.Amount, p.ReferenceNumber,
		p.BankName, p.CardType, p.CardLastDigits, p.AuthorizationCode, p.ChangeAmount, p.CashDrawerID,
		p.BranchID, p.EnterpriseID, p.UserID, p.Notes, p.Status,
	).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

func (r *repository) GetPaymentByID(ctx context.Context, id int64) (*Payment, error) {
	p := &Payment{}
	query := `
		SELECT id, payment_type, reference_id, reference_type, payment_method, amount, reference_number, bank_name, card_type, card_last_digits, authorization_code, change_amount, cash_drawer_id, branch_id, enterprise_id, user_id, notes, status, cancelled_at, cancelled_by, cancellation_reason, created_at, updated_at
		FROM payment WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.PaymentType, &p.ReferenceID, &p.ReferenceType, &p.PaymentMethod, &p.Amount,
		&p.ReferenceNumber, &p.BankName, &p.CardType, &p.CardLastDigits, &p.AuthorizationCode,
		&p.ChangeAmount, &p.CashDrawerID, &p.BranchID, &p.EnterpriseID, &p.UserID, &p.Notes,
		&p.Status, &p.CancelledAt, &p.CancelledBy, &p.CancellationReason, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return p, nil
}

func (r *repository) GetPaymentsByReference(ctx context.Context, referenceType string, referenceID int64) ([]Payment, error) {
	query := `
		SELECT id, payment_type, reference_id, reference_type, payment_method, amount, reference_number, bank_name, card_type, card_last_digits, authorization_code, change_amount, cash_drawer_id, branch_id, enterprise_id, user_id, notes, status, cancelled_at, cancelled_by, cancellation_reason, created_at, updated_at
		FROM payment WHERE reference_type = $1 AND reference_id = $2 AND status != 'CANCELLED'
		ORDER BY created_at`

	rows, err := r.db.QueryContext(ctx, query, referenceType, referenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by reference: %w", err)
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(
			&p.ID, &p.PaymentType, &p.ReferenceID, &p.ReferenceType, &p.PaymentMethod, &p.Amount,
			&p.ReferenceNumber, &p.BankName, &p.CardType, &p.CardLastDigits, &p.AuthorizationCode,
			&p.ChangeAmount, &p.CashDrawerID, &p.BranchID, &p.EnterpriseID, &p.UserID, &p.Notes,
			&p.Status, &p.CancelledAt, &p.CancelledBy, &p.CancellationReason, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *repository) UpdatePaymentStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE payment SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	return nil
}

func (r *repository) CancelPayment(ctx context.Context, id int64, cancelledBy int64, reason string) error {
	query := `UPDATE payment SET status = 'CANCELLED', cancelled_at = NOW(), cancelled_by = $1, cancellation_reason = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, cancelledBy, reason, id)
	if err != nil {
		return fmt.Errorf("failed to cancel payment: %w", err)
	}
	return nil
}

func (r *repository) ListPayments(ctx context.Context, enterpriseID int64, filters PaymentFilters) ([]Payment, error) {
	query := `
		SELECT id, payment_type, reference_id, reference_type, payment_method, amount, reference_number, bank_name, card_type, card_last_digits, authorization_code, change_amount, cash_drawer_id, branch_id, enterprise_id, user_id, notes, status, cancelled_at, cancelled_by, cancellation_reason, created_at, updated_at
		FROM payment WHERE enterprise_id = $1`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.PaymentMethod != "" {
		query += fmt.Sprintf(" AND payment_method = $%d", argPos)
		args = append(args, filters.PaymentMethod)
		argPos++
	}

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, filters.Status)
		argPos++
	}

	if filters.ReferenceID != nil {
		query += fmt.Sprintf(" AND reference_id = $%d", argPos)
		args = append(args, *filters.ReferenceID)
		argPos++
	}

	query += " ORDER BY created_at DESC"

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
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(
			&p.ID, &p.PaymentType, &p.ReferenceID, &p.ReferenceType, &p.PaymentMethod, &p.Amount,
			&p.ReferenceNumber, &p.BankName, &p.CardType, &p.CardLastDigits, &p.AuthorizationCode,
			&p.ChangeAmount, &p.CashDrawerID, &p.BranchID, &p.EnterpriseID, &p.UserID, &p.Notes,
			&p.Status, &p.CancelledAt, &p.CancelledBy, &p.CancellationReason, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *repository) CreateTransaction(ctx context.Context, tx *PaymentTransaction) error {
	query := `
		INSERT INTO payment_transaction (payment_id, transaction_type, amount, previous_balance, new_balance, processor_reference, processor_response)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		tx.PaymentID, tx.TransactionType, tx.Amount, tx.PreviousBalance, tx.NewBalance,
		tx.ProcessorReference, tx.ProcessorResponse,
	).Scan(&tx.ID, &tx.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *repository) GetTransactionsByPayment(ctx context.Context, paymentID int64) ([]PaymentTransaction, error) {
	query := `
		SELECT id, payment_id, transaction_type, amount, previous_balance, new_balance, processor_reference, processor_response, created_at
		FROM payment_transaction WHERE payment_id = $1
		ORDER BY created_at`

	rows, err := r.db.QueryContext(ctx, query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var txs []PaymentTransaction
	for rows.Next() {
		var tx PaymentTransaction
		if err := rows.Scan(
			&tx.ID, &tx.PaymentID, &tx.TransactionType, &tx.Amount, &tx.PreviousBalance,
			&tx.NewBalance, &tx.ProcessorReference, &tx.ProcessorResponse, &tx.CreatedAt,
		); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

// CashDrawer operations

func (r *repository) CreateCashDrawer(ctx context.Context, drawer *CashDrawer) error {
	query := `
		INSERT INTO cash_drawer (user_id, branch_id, enterprise_id, opening_balance, cash_in, cash_out, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, opened_at`

	err := r.db.QueryRowContext(ctx, query,
		drawer.UserID, drawer.BranchID, drawer.EnterpriseID, drawer.OpeningBalance,
		drawer.CashIn, drawer.CashOut, drawer.Status, drawer.Notes,
	).Scan(&drawer.ID, &drawer.OpenedAt)
	if err != nil {
		return fmt.Errorf("failed to create cash drawer: %w", err)
	}
	return nil
}

func (r *repository) GetCashDrawerByID(ctx context.Context, id int64) (*CashDrawer, error) {
	d := &CashDrawer{}
	query := `
		SELECT id, user_id, branch_id, enterprise_id, opening_balance, closing_balance, cash_in, cash_out, status, opened_at, closed_at, notes
		FROM cash_drawer WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.UserID, &d.BranchID, &d.EnterpriseID, &d.OpeningBalance, &d.ClosingBalance,
		&d.CashIn, &d.CashOut, &d.Status, &d.OpenedAt, &d.ClosedAt, &d.Notes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get cash drawer: %w", err)
	}
	return d, nil
}

func (r *repository) GetOpenCashDrawer(ctx context.Context, userID, branchID int64) (*CashDrawer, error) {
	d := &CashDrawer{}
	query := `
		SELECT id, user_id, branch_id, enterprise_id, opening_balance, closing_balance, cash_in, cash_out, status, opened_at, closed_at, notes
		FROM cash_drawer WHERE user_id = $1 AND branch_id = $2 AND status = 'OPEN'`

	err := r.db.QueryRowContext(ctx, query, userID, branchID).Scan(
		&d.ID, &d.UserID, &d.BranchID, &d.EnterpriseID, &d.OpeningBalance, &d.ClosingBalance,
		&d.CashIn, &d.CashOut, &d.Status, &d.OpenedAt, &d.ClosedAt, &d.Notes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get open cash drawer: %w", err)
	}
	return d, nil
}

func (r *repository) UpdateCashDrawer(ctx context.Context, drawer *CashDrawer) error {
	query := `
		UPDATE cash_drawer SET closing_balance = $1, cash_in = $2, cash_out = $3, status = $4, closed_at = $5, notes = $6
		WHERE id = $7`

	_, err := r.db.ExecContext(ctx, query,
		drawer.ClosingBalance, drawer.CashIn, drawer.CashOut, drawer.Status, drawer.ClosedAt, drawer.Notes, drawer.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update cash drawer: %w", err)
	}
	return nil
}

func (r *repository) ListCashDrawers(ctx context.Context, enterpriseID int64, userID *int64, status string) ([]CashDrawer, error) {
	query := `
		SELECT id, user_id, branch_id, enterprise_id, opening_balance, closing_balance, cash_in, cash_out, status, opened_at, closed_at, notes
		FROM cash_drawer WHERE enterprise_id = $1`

	args := []interface{}{enterpriseID}
	argPos := 2

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *userID)
		argPos++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, status)
		argPos++
	}

	query += " ORDER BY opened_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list cash drawers: %w", err)
	}
	defer rows.Close()

	var drawers []CashDrawer
	for rows.Next() {
		var d CashDrawer
		if err := rows.Scan(
			&d.ID, &d.UserID, &d.BranchID, &d.EnterpriseID, &d.OpeningBalance, &d.ClosingBalance,
			&d.CashIn, &d.CashOut, &d.Status, &d.OpenedAt, &d.ClosedAt, &d.Notes,
		); err != nil {
			return nil, err
		}
		drawers = append(drawers, d)
	}
	return drawers, nil
}

func (r *repository) CreateCashMovement(ctx context.Context, mov *CashMovement) error {
	query := `
		INSERT INTO cash_movement (cash_drawer_id, movement_type, amount, description, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		mov.CashDrawerID, mov.MovementType, mov.Amount, mov.Description, mov.UserID,
	).Scan(&mov.ID, &mov.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create cash movement: %w", err)
	}
	return nil
}

func (r *repository) GetCashMovements(ctx context.Context, drawerID int64) ([]CashMovement, error) {
	query := `
		SELECT id, cash_drawer_id, movement_type, amount, description, user_id, created_at
		FROM cash_movement WHERE cash_drawer_id = $1
		ORDER BY created_at`

	rows, err := r.db.QueryContext(ctx, query, drawerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cash movements: %w", err)
	}
	defer rows.Close()

	var movements []CashMovement
	for rows.Next() {
		var mov CashMovement
		if err := rows.Scan(
			&mov.ID, &mov.CashDrawerID, &mov.MovementType, &mov.Amount, &mov.Description, &mov.UserID, &mov.CreatedAt,
		); err != nil {
			return nil, err
		}
		movements = append(movements, mov)
	}
	return movements, nil
}

// Service implementation
type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) ProcessPayment(ctx context.Context, p *Payment) error {
	if p.Status == "" {
		p.Status = PaymentStatusCompleted
	}
	if p.PaymentType == "" {
		p.PaymentType = PaymentTypeSale
	}
	return s.repo.CreatePayment(ctx, p)
}

func (s *service) ProcessMultiplePayments(ctx context.Context, payments []Payment) error {
	for i := range payments {
		if payments[i].Status == "" {
			payments[i].Status = PaymentStatusCompleted
		}
		if payments[i].PaymentType == "" {
			payments[i].PaymentType = PaymentTypeSale
		}
		if err := s.repo.CreatePayment(ctx, &payments[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) CalculateChange(ctx context.Context, amount float64, paymentMethod string) float64 {
	if paymentMethod != MethodCash {
		return 0
	}
	return 0
}

func (s *service) GetPayment(ctx context.Context, id int64) (*Payment, error) {
	return s.repo.GetPaymentByID(ctx, id)
}

func (s *service) GetPaymentsByOrder(ctx context.Context, referenceType string, referenceID int64) ([]Payment, error) {
	return s.repo.GetPaymentsByReference(ctx, referenceType, referenceID)
}

func (s *service) CancelPayment(ctx context.Context, id int64, cancelledBy int64, reason string) error {
	return s.repo.CancelPayment(ctx, id, cancelledBy, reason)
}

func (s *service) ListPayments(ctx context.Context, enterpriseID int64, filters PaymentFilters) ([]Payment, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	return s.repo.ListPayments(ctx, enterpriseID, filters)
}

func (s *service) OpenCashDrawer(ctx context.Context, drawer *CashDrawer) error {
	existing, err := s.repo.GetOpenCashDrawer(ctx, drawer.UserID, drawer.BranchID)
	if err == nil && existing != nil {
		return fmt.Errorf("cash drawer already open for this user and branch")
	}

	drawer.Status = DrawerStatusOpen
	drawer.CashIn = 0
	drawer.CashOut = 0
	return s.repo.CreateCashDrawer(ctx, drawer)
}

func (s *service) CloseCashDrawer(ctx context.Context, drawerID int64, closingBalance float64, notes string) error {
	drawer, err := s.repo.GetCashDrawerByID(ctx, drawerID)
	if err != nil {
		return err
	}
	if drawer.Status != DrawerStatusOpen {
		return fmt.Errorf("cash drawer is not open")
	}

	now := time.Now()
	drawer.ClosingBalance = &closingBalance
	drawer.Status = DrawerStatusClosed
	drawer.ClosedAt = &now
	drawer.Notes = notes

	return s.repo.UpdateCashDrawer(ctx, drawer)
}

func (s *service) GetCashDrawer(ctx context.Context, drawerID int64) (*CashDrawer, error) {
	return s.repo.GetCashDrawerByID(ctx, drawerID)
}

func (s *service) GetOpenDrawer(ctx context.Context, userID, branchID int64) (*CashDrawer, error) {
	return s.repo.GetOpenCashDrawer(ctx, userID, branchID)
}

func (s *service) AddCashIn(ctx context.Context, drawerID int64, amount float64, description string) error {
	drawer, err := s.repo.GetCashDrawerByID(ctx, drawerID)
	if err != nil {
		return err
	}
	if drawer.Status != DrawerStatusOpen {
		return fmt.Errorf("cash drawer is not open")
	}

	mov := &CashMovement{
		CashDrawerID: drawerID,
		MovementType: CashMovementIn,
		Amount:       amount,
		Description:  description,
		UserID:       drawer.UserID,
	}

	if err := s.repo.CreateCashMovement(ctx, mov); err != nil {
		return err
	}

	drawer.CashIn += amount
	return s.repo.UpdateCashDrawer(ctx, drawer)
}

func (s *service) AddCashOut(ctx context.Context, drawerID int64, amount float64, description string) error {
	drawer, err := s.repo.GetCashDrawerByID(ctx, drawerID)
	if err != nil {
		return err
	}
	if drawer.Status != DrawerStatusOpen {
		return fmt.Errorf("cash drawer is not open")
	}

	currentBalance := drawer.OpeningBalance + drawer.CashIn - drawer.CashOut
	if currentBalance < amount {
		return fmt.Errorf("insufficient cash in drawer")
	}

	mov := &CashMovement{
		CashDrawerID: drawerID,
		MovementType: CashMovementOut,
		Amount:       amount,
		Description:  description,
		UserID:       drawer.UserID,
	}

	if err := s.repo.CreateCashMovement(ctx, mov); err != nil {
		return err
	}

	drawer.CashOut += amount
	return s.repo.UpdateCashDrawer(ctx, drawer)
}

func (s *service) ListCashDrawers(ctx context.Context, enterpriseID int64, userID *int64, status string) ([]CashDrawer, error) {
	return s.repo.ListCashDrawers(ctx, enterpriseID, userID, status)
}
