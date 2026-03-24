package shrinkage

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

// ─── ShrinkageReason ──────────────────────────────────────────────────────────

func (r *repository) CreateReason(ctx context.Context, reason *ShrinkageReason) error {
	reason.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO shrinkage_reason (code, name, description, requires_authorization, authorization_threshold, is_active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		reason.Code, reason.Name, reason.Description, reason.RequiresAuthorization, reason.AuthorizationThreshold, reason.IsActive, reason.CreatedAt,
	).Scan(&reason.ID)
}

func (r *repository) GetReasonByID(ctx context.Context, id int64) (*ShrinkageReason, error) {
	var reason ShrinkageReason
	var threshold sql.NullFloat64
	err := r.db.QueryRowContext(ctx,
		`SELECT id, code, name, description, requires_authorization, authorization_threshold, is_active, created_at
		 FROM shrinkage_reason WHERE id = $1`, id,
	).Scan(&reason.ID, &reason.Code, &reason.Name, &reason.Description, &reason.RequiresAuthorization, &threshold, &reason.IsActive, &reason.CreatedAt)
	if err != nil {
		return nil, err
	}
	if threshold.Valid {
		reason.AuthorizationThreshold = &threshold.Float64
	}
	return &reason, nil
}

func (r *repository) GetReasonByCode(ctx context.Context, code string) (*ShrinkageReason, error) {
	var reason ShrinkageReason
	var threshold sql.NullFloat64
	err := r.db.QueryRowContext(ctx,
		`SELECT id, code, name, description, requires_authorization, authorization_threshold, is_active, created_at
		 FROM shrinkage_reason WHERE code = $1`, code,
	).Scan(&reason.ID, &reason.Code, &reason.Name, &reason.Description, &reason.RequiresAuthorization, &threshold, &reason.IsActive, &reason.CreatedAt)
	if err != nil {
		return nil, err
	}
	if threshold.Valid {
		reason.AuthorizationThreshold = &threshold.Float64
	}
	return &reason, nil
}

func (r *repository) ListReasons(ctx context.Context, activeOnly bool) ([]ShrinkageReason, error) {
	query := `SELECT id, code, name, description, requires_authorization, authorization_threshold, is_active, created_at FROM shrinkage_reason`
	if activeOnly {
		query += " WHERE is_active = TRUE"
	}
	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reasons []ShrinkageReason
	for rows.Next() {
		var reason ShrinkageReason
		var threshold sql.NullFloat64
		if err := rows.Scan(&reason.ID, &reason.Code, &reason.Name, &reason.Description, &reason.RequiresAuthorization, &threshold, &reason.IsActive, &reason.CreatedAt); err != nil {
			return nil, err
		}
		if threshold.Valid {
			reason.AuthorizationThreshold = &threshold.Float64
		}
		reasons = append(reasons, reason)
	}
	return reasons, nil
}

func (r *repository) UpdateReason(ctx context.Context, reason *ShrinkageReason) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE shrinkage_reason SET name=$1, description=$2, requires_authorization=$3, authorization_threshold=$4, is_active=$5 WHERE id=$6`,
		reason.Name, reason.Description, reason.RequiresAuthorization, reason.AuthorizationThreshold, reason.IsActive, reason.ID)
	return err
}

func (r *repository) IsReasonUsed(ctx context.Context, reasonID int64) (bool, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM shrinkage WHERE reason_id = $1`, reasonID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ─── Shrinkage ────────────────────────────────────────────────────────────────

func (r *repository) CreateShrinkage(ctx context.Context, s *Shrinkage) (int64, error) {
	s.CreatedAt = time.Now()
	if s.ShrinkageDate.IsZero() {
		s.ShrinkageDate = time.Now()
	}

	return s.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO shrinkage (shrinkage_number, branch_id, user_id, reason_id, shrinkage_date, total_value, status, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		s.ShrinkageNumber, s.BranchID, s.UserID, s.ReasonID, s.ShrinkageDate, s.TotalValue, s.Status, s.Notes, s.CreatedAt,
	).Scan(&s.ID)
}

func (r *repository) GetShrinkageByID(ctx context.Context, id int64) (*Shrinkage, error) {
	var s Shrinkage
	var authBy, cancelledBy sql.NullInt64
	var authAt, cancelledAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, shrinkage_number, branch_id, user_id, reason_id, shrinkage_date, total_value, status, notes,
			authorized_by, authorized_at, cancellation_reason, cancelled_by, cancelled_at, created_at, updated_at
		 FROM shrinkage WHERE id = $1`, id,
	).Scan(&s.ID, &s.ShrinkageNumber, &s.BranchID, &s.UserID, &s.ReasonID, &s.ShrinkageDate, &s.TotalValue, &s.Status, &s.Notes,
		&authBy, &authAt, &s.CancellationReason, &cancelledBy, &cancelledAt, &s.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if authBy.Valid {
		s.AuthorizedBy = &authBy.Int64
	}
	if authAt.Valid {
		s.AuthorizedAt = &authAt.Time
	}
	if cancelledBy.Valid {
		s.CancelledBy = &cancelledBy.Int64
	}
	if cancelledAt.Valid {
		s.CancelledAt = &cancelledAt.Time
	}
	if updatedAt.Valid {
		s.UpdatedAt = &updatedAt.Time
	}
	return &s, nil
}

func (r *repository) GetShrinkageByNumber(ctx context.Context, number string) (*Shrinkage, error) {
	var s Shrinkage
	var authBy, cancelledBy sql.NullInt64
	var authAt, cancelledAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT id, shrinkage_number, branch_id, user_id, reason_id, shrinkage_date, total_value, status, notes,
			authorized_by, authorized_at, cancellation_reason, cancelled_by, cancelled_at, created_at, updated_at
		 FROM shrinkage WHERE shrinkage_number = $1`, number,
	).Scan(&s.ID, &s.ShrinkageNumber, &s.BranchID, &s.UserID, &s.ReasonID, &s.ShrinkageDate, &s.TotalValue, &s.Status, &s.Notes,
		&authBy, &authAt, &s.CancellationReason, &cancelledBy, &cancelledAt, &s.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if authBy.Valid {
		s.AuthorizedBy = &authBy.Int64
	}
	if authAt.Valid {
		s.AuthorizedAt = &authAt.Time
	}
	if cancelledBy.Valid {
		s.CancelledBy = &cancelledBy.Int64
	}
	if cancelledAt.Valid {
		s.CancelledAt = &cancelledAt.Time
	}
	if updatedAt.Valid {
		s.UpdatedAt = &updatedAt.Time
	}
	return &s, nil
}

func (r *repository) UpdateShrinkage(ctx context.Context, s *Shrinkage) error {
	now := time.Now()
	s.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE shrinkage SET status=$1, notes=$2, authorized_by=$3, authorized_at=$4, 
		 cancellation_reason=$5, cancelled_by=$6, cancelled_at=$7, updated_at=$8 WHERE id=$9`,
		s.Status, s.Notes, s.AuthorizedBy, s.AuthorizedAt, s.CancellationReason, s.CancelledBy, s.CancelledAt, now, s.ID)
	return err
}

func (r *repository) ListShrinkages(ctx context.Context, branchID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Shrinkage, int64, error) {
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

	baseQuery := "FROM shrinkage WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if branchID != nil {
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argIndex)
		args = append(args, *branchID)
		argIndex++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if startDate != nil {
		baseQuery += fmt.Sprintf(" AND shrinkage_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		baseQuery += fmt.Sprintf(" AND shrinkage_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT id, shrinkage_number, branch_id, user_id, reason_id, shrinkage_date, total_value, status, notes, authorized_by, authorized_at, cancellation_reason, cancelled_by, cancelled_at, created_at, updated_at %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var shrinkages []Shrinkage
	for rows.Next() {
		var s Shrinkage
		var authBy, cancelledBy sql.NullInt64
		var authAt, cancelledAt, updatedAt sql.NullTime
		if err := rows.Scan(&s.ID, &s.ShrinkageNumber, &s.BranchID, &s.UserID, &s.ReasonID, &s.ShrinkageDate, &s.TotalValue, &s.Status, &s.Notes,
			&authBy, &authAt, &s.CancellationReason, &cancelledBy, &cancelledAt, &s.CreatedAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		if authBy.Valid {
			s.AuthorizedBy = &authBy.Int64
		}
		if authAt.Valid {
			s.AuthorizedAt = &authAt.Time
		}
		if cancelledBy.Valid {
			s.CancelledBy = &cancelledBy.Int64
		}
		if cancelledAt.Valid {
			s.CancelledAt = &cancelledAt.Time
		}
		if updatedAt.Valid {
			s.UpdatedAt = &updatedAt.Time
		}
		shrinkages = append(shrinkages, s)
	}
	return shrinkages, total, nil
}

// ─── ShrinkageItem ────────────────────────────────────────────────────────────

func (r *repository) CreateShrinkageItem(ctx context.Context, item *ShrinkageItem) error {
	item.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO shrinkage_item (shrinkage_id, product_id, batch_number, serial_number, quantity, unit_cost, total_value, reason_detail, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		item.ShrinkageID, item.ProductID, item.BatchNumber, item.SerialNumber, item.Quantity, item.UnitCost, item.TotalValue, item.ReasonDetail, item.CreatedAt,
	).Scan(&item.ID)
}

func (r *repository) GetShrinkageItems(ctx context.Context, shrinkageID int64) ([]ShrinkageItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, shrinkage_id, product_id, batch_number, serial_number, quantity, unit_cost, total_value, reason_detail, created_at
		 FROM shrinkage_item WHERE shrinkage_id = $1`, shrinkageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ShrinkageItem
	for rows.Next() {
		var item ShrinkageItem
		var batch, serial sql.NullString
		if err := rows.Scan(&item.ID, &item.ShrinkageID, &item.ProductID, &batch, &serial, &item.Quantity, &item.UnitCost, &item.TotalValue, &item.ReasonDetail, &item.CreatedAt); err != nil {
			return nil, err
		}
		if batch.Valid {
			item.BatchNumber = &batch.String
		}
		if serial.Valid {
			item.SerialNumber = &serial.String
		}
		items = append(items, item)
	}
	return items, nil
}

// ─── Reporting ────────────────────────────────────────────────────────────────

func (r *repository) GetShrinkageReport(ctx context.Context, branchID *int64, startDate, endDate *time.Time) ([]ShrinkageReportItem, error) {
	query := `
		SELECT sr.id as reason_id, sr.name as reason_name, 
			COUNT(s.id) as total_count, 
			COALESCE(SUM(si.quantity), 0) as total_quantity,
			COALESCE(SUM(s.total_value), 0) as total_value
		FROM shrinkage_reason sr
		LEFT JOIN shrinkage s ON s.reason_id = sr.id AND s.status = 'APPROVED'
		LEFT JOIN shrinkage_item si ON si.shrinkage_id = s.id
		WHERE sr.is_active = TRUE`
	args := []interface{}{}
	argIndex := 1

	if branchID != nil {
		query += fmt.Sprintf(" AND s.branch_id = $%d", argIndex)
		args = append(args, *branchID)
		argIndex++
	}
	if startDate != nil {
		query += fmt.Sprintf(" AND s.shrinkage_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		query += fmt.Sprintf(" AND s.shrinkage_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	query += " GROUP BY sr.id, sr.name ORDER BY total_value DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report []ShrinkageReportItem
	for rows.Next() {
		var item ShrinkageReportItem
		if err := rows.Scan(&item.ReasonID, &item.ReasonName, &item.TotalCount, &item.TotalQuantity, &item.TotalValue); err != nil {
			return nil, err
		}
		report = append(report, item)
	}
	return report, nil
}
