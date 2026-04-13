package brands

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type querier = db.Querier

type repository struct {
	db querier
}

func NewRepository(db querier) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, b *Brand) error {
	query := `
		INSERT INTO brand (name, description, enterprise_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, b.Name, b.Description, b.EnterpriseID).
		Scan(&b.ID, &b.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create brand: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Brand, error) {
	b := &Brand{}
	query := `
		SELECT id, name, description, enterprise_id, created_at, updated_at, deleted_at
		FROM brand WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.Name, &b.Description, &b.EnterpriseID,
		&b.CreatedAt, &b.UpdatedAt, &b.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get brand: %w", err)
	}
	return b, nil
}

func (r *repository) List(ctx context.Context, enterpriseID int64) ([]Brand, error) {
	query := `
		SELECT id, name, description, enterprise_id, created_at, updated_at, deleted_at
		FROM brand WHERE enterprise_id = $1 AND deleted_at IS NULL
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list brands: %w", err)
	}
	defer rows.Close()

	var list []Brand
	for rows.Next() {
		var b Brand
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.EnterpriseID,
			&b.CreatedAt, &b.UpdatedAt, &b.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *repository) Update(ctx context.Context, b *Brand) error {
	query := `
		UPDATE brand SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, b.Name, b.Description, b.ID)
	if err != nil {
		return fmt.Errorf("failed to update brand: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE brand SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}
	return nil
}
