package transfers

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

// ─── Transfer ─────────────────────────────────────────────────────────────────

func (r *repository) CreateTransfer(ctx context.Context, t *Transfer) (int64, error) {
	t.CreatedAt = time.Now()
	t.RequestedDate = time.Now()
	t.Status = StatusPending

	return t.ID, r.db.QueryRowContext(ctx,
		`INSERT INTO transfer (transfer_number, origin_branch_id, destination_branch_id, user_id, 
			status, requested_date, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		t.TransferNumber, t.OriginBranchID, t.DestinationBranchID, t.UserID,
		t.Status, t.RequestedDate, t.Notes, t.CreatedAt,
	).Scan(&t.ID)
}

func (r *repository) GetTransferByID(ctx context.Context, id int64) (*Transfer, error) {
	var t Transfer
	var shippedDate, receivedDate, cancelledAt, updatedAt sql.NullTime
	var shippedBy, receivedBy, cancelledBy sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, transfer_number, origin_branch_id, destination_branch_id, user_id, status,
			requested_date, shipped_date, received_date, notes, shipped_by, received_by,
			cancellation_reason, cancelled_by, cancelled_at, created_at, updated_at
		 FROM transfer WHERE id = $1`, id,
	).Scan(&t.ID, &t.TransferNumber, &t.OriginBranchID, &t.DestinationBranchID, &t.UserID, &t.Status,
		&t.RequestedDate, &shippedDate, &receivedDate, &t.Notes, &shippedBy, &receivedBy,
		&t.CancellationReason, &cancelledBy, &cancelledAt, &t.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if shippedDate.Valid {
		t.ShippedDate = &shippedDate.Time
	}
	if receivedDate.Valid {
		t.ReceivedDate = &receivedDate.Time
	}
	if shippedBy.Valid {
		t.ShippedBy = &shippedBy.Int64
	}
	if receivedBy.Valid {
		t.ReceivedBy = &receivedBy.Int64
	}
	if cancelledBy.Valid {
		t.CancelledBy = &cancelledBy.Int64
	}
	if cancelledAt.Valid {
		t.CancelledAt = &cancelledAt.Time
	}
	if updatedAt.Valid {
		t.UpdatedAt = &updatedAt.Time
	}
	return &t, nil
}

func (r *repository) GetTransferByNumber(ctx context.Context, number string) (*Transfer, error) {
	var t Transfer
	var shippedDate, receivedDate, cancelledAt, updatedAt sql.NullTime
	var shippedBy, receivedBy, cancelledBy sql.NullInt64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, transfer_number, origin_branch_id, destination_branch_id, user_id, status,
			requested_date, shipped_date, received_date, notes, shipped_by, received_by,
			cancellation_reason, cancelled_by, cancelled_at, created_at, updated_at
		 FROM transfer WHERE transfer_number = $1`, number,
	).Scan(&t.ID, &t.TransferNumber, &t.OriginBranchID, &t.DestinationBranchID, &t.UserID, &t.Status,
		&t.RequestedDate, &shippedDate, &receivedDate, &t.Notes, &shippedBy, &receivedBy,
		&t.CancellationReason, &cancelledBy, &cancelledAt, &t.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if shippedDate.Valid {
		t.ShippedDate = &shippedDate.Time
	}
	if receivedDate.Valid {
		t.ReceivedDate = &receivedDate.Time
	}
	if shippedBy.Valid {
		t.ShippedBy = &shippedBy.Int64
	}
	if receivedBy.Valid {
		t.ReceivedBy = &receivedBy.Int64
	}
	if cancelledBy.Valid {
		t.CancelledBy = &cancelledBy.Int64
	}
	if cancelledAt.Valid {
		t.CancelledAt = &cancelledAt.Time
	}
	if updatedAt.Valid {
		t.UpdatedAt = &updatedAt.Time
	}
	return &t, nil
}

func (r *repository) UpdateTransfer(ctx context.Context, t *Transfer) error {
	now := time.Now()
	t.UpdatedAt = &now
	_, err := r.db.ExecContext(ctx,
		`UPDATE transfer SET status=$1, notes=$2, shipped_date=$3, received_date=$4,
			shipped_by=$5, received_by=$6, cancellation_reason=$7, cancelled_by=$8, 
			cancelled_at=$9, updated_at=$10 WHERE id=$11`,
		t.Status, t.Notes, t.ShippedDate, t.ReceivedDate, t.ShippedBy, t.ReceivedBy,
		t.CancellationReason, t.CancelledBy, t.CancelledAt, now, t.ID)
	return err
}

func (r *repository) ListTransfers(ctx context.Context, originBranchID, destBranchID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Transfer, int64, error) {
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

	baseQuery := "FROM transfer WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if originBranchID != nil {
		baseQuery += fmt.Sprintf(" AND origin_branch_id = $%d", argIndex)
		args = append(args, *originBranchID)
		argIndex++
	}
	if destBranchID != nil {
		baseQuery += fmt.Sprintf(" AND destination_branch_id = $%d", argIndex)
		args = append(args, *destBranchID)
		argIndex++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if startDate != nil {
		baseQuery += fmt.Sprintf(" AND requested_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		baseQuery += fmt.Sprintf(" AND requested_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf("SELECT id, transfer_number, origin_branch_id, destination_branch_id, user_id, status, requested_date, shipped_date, received_date, notes, shipped_by, received_by, cancellation_reason, cancelled_by, cancelled_at, created_at, updated_at %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", baseQuery, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transfers []Transfer
	for rows.Next() {
		var t Transfer
		var shippedDate, receivedDate, cancelledAt, updatedAt sql.NullTime
		var shippedBy, receivedBy, cancelledBy sql.NullInt64
		if err := rows.Scan(&t.ID, &t.TransferNumber, &t.OriginBranchID, &t.DestinationBranchID, &t.UserID, &t.Status,
			&t.RequestedDate, &shippedDate, &receivedDate, &t.Notes, &shippedBy, &receivedBy,
			&t.CancellationReason, &cancelledBy, &cancelledAt, &t.CreatedAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		if shippedDate.Valid {
			t.ShippedDate = &shippedDate.Time
		}
		if receivedDate.Valid {
			t.ReceivedDate = &receivedDate.Time
		}
		if shippedBy.Valid {
			t.ShippedBy = &shippedBy.Int64
		}
		if receivedBy.Valid {
			t.ReceivedBy = &receivedBy.Int64
		}
		if cancelledBy.Valid {
			t.CancelledBy = &cancelledBy.Int64
		}
		if cancelledAt.Valid {
			t.CancelledAt = &cancelledAt.Time
		}
		if updatedAt.Valid {
			t.UpdatedAt = &updatedAt.Time
		}
		transfers = append(transfers, t)
	}
	return transfers, total, nil
}

// ─── TransferItem ──────────────────────────────────────────────────────────────

func (r *repository) CreateTransferItem(ctx context.Context, item *TransferItem) error {
	item.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx,
		`INSERT INTO transfer_item (transfer_id, product_id, requested_quantity, shipped_quantity, received_quantity, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		item.TransferID, item.ProductID, item.RequestedQuantity, item.ShippedQuantity, item.ReceivedQuantity, item.CreatedAt,
	).Scan(&item.ID)
}

func (r *repository) GetTransferItems(ctx context.Context, transferID int64) ([]TransferItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, transfer_id, product_id, requested_quantity, shipped_quantity, received_quantity, created_at
		 FROM transfer_item WHERE transfer_id = $1`, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TransferItem
	for rows.Next() {
		var item TransferItem
		if err := rows.Scan(&item.ID, &item.TransferID, &item.ProductID, &item.RequestedQuantity,
			&item.ShippedQuantity, &item.ReceivedQuantity, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *repository) UpdateTransferItem(ctx context.Context, item *TransferItem) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE transfer_item SET shipped_quantity=$1, received_quantity=$2 WHERE id=$3`,
		item.ShippedQuantity, item.ReceivedQuantity, item.ID)
	return err
}
