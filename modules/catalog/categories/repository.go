package categories

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

func (r *repository) Create(ctx context.Context, c *Category) error {
	query := `
		INSERT INTO category (name, description, parent_id, enterprise_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, c.Name, c.Description, c.ParentID, c.EnterpriseID).
		Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Category, error) {
	c := &Category{}
	query := `
		SELECT id, name, description, parent_id, enterprise_id, created_at, updated_at, deleted_at
		FROM category WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.ParentID, &c.EnterpriseID,
		&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return c, nil
}

func (r *repository) List(ctx context.Context, enterpriseID int64) ([]Category, error) {
	query := `
		SELECT id, name, description, parent_id, enterprise_id, created_at, updated_at, deleted_at
		FROM category WHERE enterprise_id = $1 AND deleted_at IS NULL
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var list []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.EnterpriseID,
			&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *repository) Update(ctx context.Context, c *Category) error {
	query := `
		UPDATE category SET name = $1, description = $2, parent_id = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, c.Name, c.Description, c.ParentID, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE category SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}
