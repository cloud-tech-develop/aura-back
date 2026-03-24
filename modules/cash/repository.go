package cash

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

// ─── Cash Drawer ─────────────────────────────────────────────────────────────

func (r *repository) CreateDrawer(ctx context.Context, d *CashDrawer) error {
	d.CreatedAt = time.Now()
	d.IsActive = true
	if d.Name == "" {
		d.Name = "MAIN"
	}

	query := `
		INSERT INTO cash_drawer (branch_id, name, is_active, min_float, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		d.BranchID, d.Name, d.IsActive, d.MinFloat, d.CreatedAt,
	).Scan(&d.ID)
}

func (r *repository) GetDrawerByBranch(ctx context.Context, branchID int64) (*CashDrawer, error) {
	var d CashDrawer
	var updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, branch_id, name, is_active, min_float, created_at, updated_at 
		 FROM cash_drawer WHERE branch_id = $1`,
		branchID,
	).Scan(&d.ID, &d.BranchID, &d.Name, &d.IsActive, &d.MinFloat, &d.CreatedAt, &updatedAt)

	if err != nil {
		return nil, err
	}
	if updatedAt.Valid {
		d.UpdatedAt = &updatedAt.Time
	}
	return &d, nil
}

func (r *repository) UpdateDrawer(ctx context.Context, d *CashDrawer) error {
	now := time.Now()
	d.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE cash_drawer SET name = $1, is_active = $2, min_float = $3, updated_at = $4 WHERE id = $5`,
		d.Name, d.IsActive, d.MinFloat, d.UpdatedAt, d.ID,
	)
	return err
}

// ─── Cash Shift ──────────────────────────────────────────────────────────────

func (r *repository) CreateShift(ctx context.Context, s *CashShift) error {
	s.CreatedAt = time.Now()
	s.OpenedAt = time.Now()
	s.Status = StatusOpen

	query := `
		INSERT INTO cash_shift (cash_drawer_id, user_id, branch_id, opening_amount, opening_notes, status, opened_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		s.CashDrawerID, s.UserID, s.BranchID, s.OpeningAmount, s.OpeningNotes, s.Status, s.OpenedAt, s.CreatedAt,
	).Scan(&s.ID)
}

func (r *repository) GetShiftByID(ctx context.Context, id int64) (*CashShift, error) {
	var s CashShift
	var closedAt, closedBy, authorizedBy sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, cash_drawer_id, user_id, branch_id, opening_amount, closing_amount, expected_amount, 
		        difference, opening_notes, closing_notes, status, opened_at, closed_at, closed_by, authorized_by, created_at
		 FROM cash_shift WHERE id = $1`,
		id,
	).Scan(&s.ID, &s.CashDrawerID, &s.UserID, &s.BranchID, &s.OpeningAmount, &s.ClosingAmount,
		&s.ExpectedAmount, &s.Difference, &s.OpeningNotes, &s.ClosingNotes, &s.Status, &s.OpenedAt,
		&closedAt, &closedBy, &authorizedBy, &s.CreatedAt)

	if err != nil {
		return nil, err
	}
	if closedAt.Valid {
		t := time.Unix(closedAt.Int64, 0)
		s.ClosedAt = &t
	}
	if closedBy.Valid {
		s.ClosedBy = &closedBy.Int64
	}
	if authorizedBy.Valid {
		s.AuthorizedBy = &authorizedBy.Int64
	}
	return &s, nil
}

func (r *repository) GetActiveShiftByUser(ctx context.Context, userID int64) (*CashShift, error) {
	var s CashShift
	var closedAt, closedBy, authorizedBy sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, cash_drawer_id, user_id, branch_id, opening_amount, closing_amount, expected_amount, 
		        difference, opening_notes, closing_notes, status, opened_at, closed_at, closed_by, authorized_by, created_at
		 FROM cash_shift WHERE user_id = $1 AND status = 'OPEN' ORDER BY opened_at DESC LIMIT 1`,
		userID,
	).Scan(&s.ID, &s.CashDrawerID, &s.UserID, &s.BranchID, &s.OpeningAmount, &s.ClosingAmount,
		&s.ExpectedAmount, &s.Difference, &s.OpeningNotes, &s.ClosingNotes, &s.Status, &s.OpenedAt,
		&closedAt, &closedBy, &authorizedBy, &s.CreatedAt)

	if err != nil {
		return nil, err
	}
	if closedAt.Valid {
		t := time.Unix(closedAt.Int64, 0)
		s.ClosedAt = &t
	}
	if closedBy.Valid {
		s.ClosedBy = &closedBy.Int64
	}
	if authorizedBy.Valid {
		s.AuthorizedBy = &authorizedBy.Int64
	}
	return &s, nil
}

func (r *repository) UpdateShift(ctx context.Context, s *CashShift) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE cash_shift SET closing_amount = $1, expected_amount = $2, difference = $3, 
		        closing_notes = $4, status = $5, closed_at = $6, closed_by = $7, authorized_by = $8
		 WHERE id = $9`,
		s.ClosingAmount, s.ExpectedAmount, s.Difference, s.ClosingNotes, s.Status, s.ClosedAt, s.ClosedBy, s.AuthorizedBy, s.ID,
	)
	return err
}

func (r *repository) ListShifts(ctx context.Context, branchID int64, startDate, endDate *time.Time, status string, page, limit int) ([]CashShift, int64, error) {
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

	baseQuery := "FROM cash_shift WHERE branch_id = $1"
	args := []interface{}{branchID}
	argIndex := 2

	if startDate != nil {
		baseQuery += fmt.Sprintf(" AND opened_at >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		baseQuery += fmt.Sprintf(" AND opened_at <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated data
	dataQuery := fmt.Sprintf("SELECT id, cash_drawer_id, user_id, branch_id, opening_amount, closing_amount, expected_amount, difference, opening_notes, closing_notes, status, opened_at, closed_at, closed_by, authorized_by, created_at %s ORDER BY opened_at DESC LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var shifts []CashShift
	for rows.Next() {
		var s CashShift
		var closedAt, closedBy, authorizedBy sql.NullInt64

		if err := rows.Scan(&s.ID, &s.CashDrawerID, &s.UserID, &s.BranchID, &s.OpeningAmount, &s.ClosingAmount,
			&s.ExpectedAmount, &s.Difference, &s.OpeningNotes, &s.ClosingNotes, &s.Status, &s.OpenedAt,
			&closedAt, &closedBy, &authorizedBy, &s.CreatedAt); err != nil {
			return nil, 0, err
		}
		if closedAt.Valid {
			t := time.Unix(closedAt.Int64, 0)
			s.ClosedAt = &t
		}
		if closedBy.Valid {
			s.ClosedBy = &closedBy.Int64
		}
		if authorizedBy.Valid {
			s.AuthorizedBy = &authorizedBy.Int64
		}
		shifts = append(shifts, s)
	}

	return shifts, total, nil
}

// ─── Cash Movement ───────────────────────────────────────────────────────────

func (r *repository) CreateMovement(ctx context.Context, m *CashMovement) error {
	m.CreatedAt = time.Now()

	query := `
		INSERT INTO cash_movement (shift_id, movement_type, reason, amount, reference_id, reference_type, notes, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		m.ShiftID, m.MovementType, m.Reason, m.Amount, m.ReferenceID, m.ReferenceType, m.Notes, m.UserID, m.CreatedAt,
	).Scan(&m.ID)
}

func (r *repository) GetMovementsByShift(ctx context.Context, shiftID int64) ([]CashMovement, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, shift_id, movement_type, reason, amount, reference_id, reference_type, notes, user_id, created_at
		 FROM cash_movement WHERE shift_id = $1 ORDER BY created_at ASC`,
		shiftID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []CashMovement
	for rows.Next() {
		var m CashMovement
		var referenceID sql.NullInt64
		var referenceType sql.NullString

		if err := rows.Scan(&m.ID, &m.ShiftID, &m.MovementType, &m.Reason, &m.Amount, &referenceID, &referenceType, &m.Notes, &m.UserID, &m.CreatedAt); err != nil {
			return nil, err
		}
		if referenceID.Valid {
			m.ReferenceID = &referenceID.Int64
		}
		if referenceType.Valid {
			m.ReferenceType = &referenceType.String
		}
		movements = append(movements, m)
	}

	return movements, nil
}

func (r *repository) GetMovementsSummary(ctx context.Context, shiftID int64) (totalIn, totalOut float64, err error) {
	err = r.db.QueryRowContext(ctx,
		`SELECT 
			COALESCE(SUM(CASE WHEN movement_type = 'IN' THEN amount ELSE 0 END), 0) as total_in,
			COALESCE(SUM(CASE WHEN movement_type = 'OUT' THEN amount ELSE 0 END), 0) as total_out
		 FROM cash_movement WHERE shift_id = $1`,
		shiftID,
	).Scan(&totalIn, &totalOut)
	return
}
