package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/tenant"
)

type service struct {
	db       *sql.DB
	repo     Repository
	eventBus events.EventBus
}

func NewService(db *sql.DB, eventBus events.EventBus) Service {
	return &service{
		db:       db,
		repo:     NewRepository(db),
		eventBus: eventBus,
	}
}

func (s *service) Create(ctx context.Context, tenantSlug string, user *User, password string) error {
	// 1. Validate input
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}

	// 2. Hash password
	hash, err := tenant.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = hash
	user.Active = true // Default active

	// 3. Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create repository instance with transaction
	repoTx := NewRepositoryWithQuerier(tx)

	// Check email uniqueness within the transaction
	existingUser, err := repoTx.GetByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("el email %s ya está registrado", user.Email)
	}

	// 4. Insert user into public.users
	if err := repoTx.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// 5. Insert third party into tenant.third_parties
	// Construct query with tenant schema
	thirdPartyQuery := fmt.Sprintf(`
		INSERT INTO %q.third_parties (user_id, first_name, last_name, document_number, document_type, personal_email, tax_responsibility, is_employee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, tenantSlug)

	// Use default values if not provided
	firstName := user.FirstName
	if firstName == "" {
		firstName = user.Name // Fallback to user name if first name not provided
	}

	_, err = tx.ExecContext(ctx, thirdPartyQuery,
		user.ID,
		firstName,
		user.LastName,
		user.DocumentNumber,
		user.DocumentType,
		user.PersonalEmail,
		user.TaxResponsibility,
		user.IsEmployee,
	)
	if err != nil {
		return fmt.Errorf("failed to create third party: %w", err)
	}

	// 6. Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 7. Publish event
	s.publish(NewCreatedEvent(user))

	return nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListByEnterprise(ctx context.Context, enterpriseID int64, page, limit int) ([]User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.ListByEnterprise(ctx, enterpriseID, limit, offset)
}

func (s *service) Update(ctx context.Context, user *User) error {
	// Check if email is already taken by another user
	existingUser, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil && existingUser.ID != user.ID {
		return fmt.Errorf("el email %s ya está registrado", user.Email)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	s.publish(NewUpdatedEvent(user))
	return nil
}

func (s *service) UpdateStatus(ctx context.Context, id int64, active bool) error {
	return s.repo.UpdateStatus(ctx, id, active)
}

func (s *service) AssignRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	// Start transaction for role assignment
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing roles for this user
	_, err = tx.ExecContext(ctx, `DELETE FROM public.user_roles WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to clear existing roles: %w", err)
	}

	// Insert new roles
	for _, roleID := range roleIDs {
		_, err = tx.ExecContext(ctx, `INSERT INTO public.user_roles (user_id, role_id) VALUES ($1, $2)`, userID, roleID)
		if err != nil {
			return fmt.Errorf("failed to assign role %d: %w", roleID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit role assignment: %w", err)
	}

	return nil
}

func (s *service) publish(event events.Event) {
	if s.eventBus == nil {
		return
	}
	if err := s.eventBus.Publish(event); err != nil {
		fmt.Printf("[users.Service] warn: publish failed: %v\n", err)
	}
}

// Event constructors
func NewCreatedEvent(user *User) events.Event {
	return events.NewBaseEvent(EventCreated, user)
}

func NewUpdatedEvent(user *User) events.Event {
	return events.NewBaseEvent(EventUpdated, user)
}
