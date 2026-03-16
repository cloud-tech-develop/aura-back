package users

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
)

// querier is an internal interface to support both *sql.DB and *sql.Tx.
type querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type repository struct {
	q querier
}

func NewRepository(db *sql.DB) Repository {
	return &repository{q: db}
}

// NewRepositoryWithQuerier creates a repository with a specific querier (e.g., a transaction).
func NewRepositoryWithQuerier(q querier) Repository {
	return &repository{q: q}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO public.users (enterprise_id, name, email, password_hash, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.q.QueryRowContext(ctx, query,
		user.EnterpriseID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("repository: failed to create user: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	var user User
	query := `
		SELECT id, enterprise_id, name, email, password_hash, active, created_at, updated_at, deleted_at
		FROM public.users WHERE id = $1 AND deleted_at IS NULL`

	err := r.q.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.EnterpriseID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Active, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("repository: failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email vo.Email) (*User, error) {
	var user User
	query := `
		SELECT id, enterprise_id, name, email, password_hash, active, created_at, updated_at, deleted_at
		FROM public.users WHERE email = $1 AND deleted_at IS NULL`

	err := r.q.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.EnterpriseID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Active, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("repository: failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *repository) ListByEnterprise(ctx context.Context, enterpriseID int64, limit, offset int) ([]User, error) {
	query := `
		SELECT id, enterprise_id, name, email, password_hash, active, created_at, updated_at, deleted_at
		FROM public.users 
		WHERE enterprise_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.q.QueryContext(ctx, query, enterpriseID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID, &user.EnterpriseID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Active, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("repository: failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: rows error: %w", err)
	}

	return users, nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE public.users 
		SET name = $1, email = $2, password_hash = $3, active = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL`

	user.UpdatedAt = time.Now()
	_, err := r.q.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Active,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("repository: failed to update user: %w", err)
	}
	return nil
}

func (r *repository) UpdateStatus(ctx context.Context, id int64, active bool) error {
	query := `
		UPDATE public.users 
		SET active = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL`

	_, err := r.q.ExecContext(ctx, query, active, time.Now(), id)
	if err != nil {
		return fmt.Errorf("repository: failed to update user status: %w", err)
	}
	return nil
}
