package offline

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// ─── Enterprise Operations ─────────────────────────────────────────────────────

// UpsertEnterprise inserts an enterprise if it doesn't exist
func (r *repository) UpsertEnterprise(ctx context.Context, e *Enterprise) error {
	// Check if enterprise already exists
	existing, err := r.GetEnterpriseBySlug(ctx, e.Slug)
	if err == nil && existing != nil {
		// Update existing
		return r.updateEnterprise(ctx, e)
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Set defaults
	if e.Status == "" {
		e.Status = "ACTIVE"
	}
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	settingsJSON, err := json.Marshal(e.Settings)
	if err != nil || string(settingsJSON) == "null" {
		settingsJSON = []byte("{}")
	}

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO enterprises 
		 (id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at, deleted_at) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain,
		e.Email, e.Document, e.DV, e.Phone, e.MunicipalityID, e.Municipality,
		e.Status, settingsJSON, e.CreatedAt, e.UpdatedAt, nil,
	)
	return err
}

// updateEnterprise updates an existing enterprise
func (r *repository) updateEnterprise(ctx context.Context, e *Enterprise) error {
	e.UpdatedAt = time.Now()

	settingsJSON, err := json.Marshal(e.Settings)
	if err != nil || string(settingsJSON) == "null" {
		settingsJSON = []byte("{}")
	}

	_, err = r.db.ExecContext(ctx,
		`UPDATE enterprises SET 
		 tenant_id = ?, name = ?, commercial_name = ?, slug = ?, sub_domain = ?, email = ?, document = ?, dv = ?, phone = ?, municipality_id = ?, municipality = ?, status = ?, settings = ?, updated_at = ?
		 WHERE slug = ?`,
		e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain,
		e.Email, e.Document, e.DV, e.Phone, e.MunicipalityID, e.Municipality,
		e.Status, settingsJSON, e.UpdatedAt, e.Slug,
	)
	return err
}

// GetEnterpriseBySlug retrieves an enterprise by its slug
func (r *repository) GetEnterpriseBySlug(ctx context.Context, slug string) (*Enterprise, error) {
	var e Enterprise
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM enterprises WHERE slug = ? AND deleted_at IS NULL`,
		slug,
	).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
		&e.Email, &e.Document, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality,
		&e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(settingsJSON, &e.Settings)
	return &e, nil
}

// ListEnterprises returns all enterprises stored locally
func (r *repository) ListEnterprises(ctx context.Context) ([]Enterprise, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM enterprises WHERE deleted_at IS NULL ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Enterprise
	for rows.Next() {
		var e Enterprise
		var settingsJSON []byte
		if err := rows.Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
			&e.Email, &e.Document, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality,
			&e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(settingsJSON, &e.Settings)
		list = append(list, e)
	}

	return list, nil
}

// ─── Plan Operations ─────────────────────────────────────────────────────

// UpsertPlan inserts or updates a plan
func (r *repository) UpsertPlan(ctx context.Context, p *Plan) error {
	// Check if exists
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM plans WHERE id = ?)", p.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE plans SET enterprise_id = ?, max_users = ?, max_enterprises = ?, trial_until = ?, updated_at = ?
			 WHERE id = ?`,
			p.EnterpriseID, p.MaxUsers, p.MaxEnterprises, p.TrialUntil, time.Now(), p.ID,
		)
	} else {
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO plans (id, enterprise_id, max_users, max_enterprises, trial_until, created_at, updated_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			p.ID, p.EnterpriseID, p.MaxUsers, p.MaxEnterprises, p.TrialUntil, p.CreatedAt, p.UpdatedAt, nil,
		)
	}
	return err
}

// ListPlans returns all plans stored locally
func (r *repository) ListPlans(ctx context.Context) ([]Plan, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, enterprise_id, max_users, max_enterprises, trial_until, created_at, updated_at, deleted_at FROM plans")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.EnterpriseID, &p.MaxUsers, &p.MaxEnterprises, &p.TrialUntil, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// ─── User Operations ─────────────────────────────────────────────────────

// UpsertUser inserts or updates a user
func (r *repository) UpsertUser(ctx context.Context, u *User) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", u.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE users SET enterprise_id = ?, name = ?, email = ?, active = ?, updated_at = ?
			 WHERE id = ?`,
			u.EnterpriseID, u.Name, u.Email, u.Active, time.Now(), u.ID,
		)
	} else {
		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO users (id, enterprise_id, name, email, active, password_hash, created_at, updated_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			u.ID, u.EnterpriseID, u.Name, u.Email, u.Active, u.PasswordHash, u.CreatedAt, u.UpdatedAt, nil,
		)
	}
	return err
}

// ListUsers returns all users stored locally
func (r *repository) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, enterprise_id, name, email, active, created_at, updated_at, deleted_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.EnterpriseID, &u.Name, &u.Email, &u.Active, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}

// ─── UserRole Operations ─────────────────────────────────────────────────────

// UpsertUserRole inserts or updates a user role
func (r *repository) UpsertUserRole(ctx context.Context, ur *UserRole) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM user_roles WHERE id = ?)", ur.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE user_roles SET user_id = ?, role_id = ? WHERE id = ?`,
			ur.UserID, ur.RoleID, ur.ID,
		)
	} else {
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO user_roles (id, user_id, role_id) VALUES (?, ?, ?)`,
			ur.ID, ur.UserID, ur.RoleID,
		)
	}
	return err
}

// ListUserRoles returns all user roles stored locally
func (r *repository) ListUserRoles(ctx context.Context) ([]UserRole, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, role_id FROM user_roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []UserRole
	for rows.Next() {
		var ur UserRole
		if err := rows.Scan(&ur.ID, &ur.UserID, &ur.RoleID); err != nil {
			return nil, err
		}
		list = append(list, ur)
	}
	return list, nil
}
