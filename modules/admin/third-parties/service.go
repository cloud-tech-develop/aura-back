package thirdparties

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type repository struct {
	db db.Querier
}

func NewRepository(db db.Querier) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, tenantSlug string, tp *ThirdParty) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s".third_parties (
			user_id, first_name, last_name, document_number, document_type,
			personal_email, commercial_name, address, phone, additional_email,
			tax_responsibility, is_client, is_provider, is_employee,
			municipality_id, municipality
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query,
		tp.UserID, tp.FirstName, tp.LastName, tp.DocumentNumber, tp.DocumentType,
		tp.PersonalEmail, tp.CommercialName, tp.Address, tp.Phone, tp.AdditionalEmail,
		tp.TaxResponsibility, tp.IsClient, tp.IsProvider, tp.IsEmployee,
		tp.MunicipalityID, tp.Municipality,
	).Scan(&tp.ID, &tp.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create third party: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*ThirdParty, error) {
	tp := &ThirdParty{}
	query := fmt.Sprintf(`
		SELECT id, user_id, first_name, last_name, document_number, document_type,
			personal_email, commercial_name, address, phone, additional_email,
			tax_responsibility, is_client, is_provider, is_employee,
			municipality_id, municipality, created_at, deleted_at
		FROM "%s".third_parties WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tp.ID, &tp.UserID, &tp.FirstName, &tp.LastName, &tp.DocumentNumber, &tp.DocumentType,
		&tp.PersonalEmail, &tp.CommercialName, &tp.Address, &tp.Phone, &tp.AdditionalEmail,
		&tp.TaxResponsibility, &tp.IsClient, &tp.IsProvider, &tp.IsEmployee,
		&tp.MunicipalityID, &tp.Municipality, &tp.CreatedAt, &tp.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get third party: %w", err)
	}
	return tp, nil
}

func (r *repository) GetByDocument(ctx context.Context, tenantSlug string, docNumber string) (*ThirdParty, error) {
	tp := &ThirdParty{}
	query := fmt.Sprintf(`
		SELECT id, user_id, first_name, last_name, document_number, document_type,
			personal_email, commercial_name, address, phone, additional_email,
			tax_responsibility, is_client, is_provider, is_employee,
			municipality_id, municipality, created_at, deleted_at
		FROM "%s".third_parties WHERE document_number = $1 AND deleted_at IS NULL`, tenantSlug)

	err := r.db.QueryRowContext(ctx, query, docNumber).Scan(
		&tp.ID, &tp.UserID, &tp.FirstName, &tp.LastName, &tp.DocumentNumber, &tp.DocumentType,
		&tp.PersonalEmail, &tp.CommercialName, &tp.Address, &tp.Phone, &tp.AdditionalEmail,
		&tp.TaxResponsibility, &tp.IsClient, &tp.IsProvider, &tp.IsEmployee,
		&tp.MunicipalityID, &tp.Municipality, &tp.CreatedAt, &tp.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get third party by document: %w", err)
	}
	return tp, nil
}

func (r *repository) Update(ctx context.Context, tenantSlug string, id int64, tp *ThirdParty) error {
	query := fmt.Sprintf(`
		UPDATE "%s".third_parties SET
			user_id = $1, first_name = $2, last_name = $3, document_type = $4,
			personal_email = $5, commercial_name = $6, address = $7, phone = $8,
			additional_email = $9, tax_responsibility = $10, is_client = $11,
			is_provider = $12, is_employee = $13, municipality_id = $14, municipality = $15
		WHERE id = $16 AND deleted_at IS NULL`, tenantSlug)

	result, err := r.db.ExecContext(ctx, query,
		tp.UserID, tp.FirstName, tp.LastName, tp.DocumentType,
		tp.PersonalEmail, tp.CommercialName, tp.Address, tp.Phone,
		tp.AdditionalEmail, tp.TaxResponsibility, tp.IsClient,
		tp.IsProvider, tp.IsEmployee, tp.MunicipalityID, tp.Municipality,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update third party: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	query := fmt.Sprintf(`UPDATE "%s".third_parties SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, tenantSlug)
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete third party: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ThirdPartyFilters) ([]ThirdParty, error) {
	query := fmt.Sprintf(`
		SELECT id, user_id, first_name, last_name, document_number, document_type,
			personal_email, commercial_name, address, phone, additional_email,
			tax_responsibility, is_client, is_provider, is_employee,
			municipality_id, municipality, created_at, deleted_at
		FROM "%s".third_parties WHERE deleted_at IS NULL`, tenantSlug)

	args := []interface{}{}
	argPos := 1

	// Filter by type
	if filters.Type != "" {
		switch filters.Type {
		case "client":
			query += fmt.Sprintf(" AND is_client = TRUE")
		case "provider":
			query += fmt.Sprintf(" AND is_provider = TRUE")
		case "employee":
			query += fmt.Sprintf(" AND is_employee = TRUE")
		}
	}

	// Search by name or document
	if filters.Search != "" {
		search := "%" + strings.ToLower(filters.Search) + "%"
		query += fmt.Sprintf(" AND (LOWER(first_name) LIKE $%d OR LOWER(last_name) LIKE $%d OR LOWER(commercial_name) LIKE $%d OR LOWER(document_number) LIKE $%d)", argPos, argPos, argPos, argPos)
		args = append(args, search)
		argPos++
	}

	query += " ORDER BY created_at DESC"

	// Pagination
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
		return nil, fmt.Errorf("failed to list third parties: %w", err)
	}
	defer rows.Close()

	var thirdParties []ThirdParty
	for rows.Next() {
		var tp ThirdParty
		if err := rows.Scan(
			&tp.ID, &tp.UserID, &tp.FirstName, &tp.LastName, &tp.DocumentNumber, &tp.DocumentType,
			&tp.PersonalEmail, &tp.CommercialName, &tp.Address, &tp.Phone, &tp.AdditionalEmail,
			&tp.TaxResponsibility, &tp.IsClient, &tp.IsProvider, &tp.IsEmployee,
			&tp.MunicipalityID, &tp.Municipality, &tp.CreatedAt, &tp.DeletedAt,
		); err != nil {
			return nil, err
		}
		thirdParties = append(thirdParties, tp)
	}
	return thirdParties, nil
}

func (r *repository) Count(ctx context.Context, tenantSlug string, enterpriseID int64, filters ThirdPartyFilters) (int, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM "%s".third_parties WHERE deleted_at IS NULL`, tenantSlug)

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count third parties: %w", err)
	}
	return count, nil
}

// Service implementation
type service struct {
	repo Repository
}

func NewService(db db.Querier) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) Create(ctx context.Context, tenantSlug string, tp *ThirdParty) error {
	// Validate document type
	validDocTypes := map[string]bool{
		DocumentTypeCC:       true,
		DocumentTypeCE:       true,
		DocumentTypeNIT:      true,
		DocumentTypePASSPORT: true,
		DocumentTypeRUT:      true,
	}
	if !validDocTypes[tp.DocumentType] {
		return fmt.Errorf("invalid document type: %s", tp.DocumentType)
	}

	// Validate tax responsibility
	validTaxResp := map[string]bool{
		TaxRespResponsible:    true,
		TaxRespNotResponsible: true,
	}
	if !validTaxResp[tp.TaxResponsibility] {
		return fmt.Errorf("invalid tax responsibility: %s", tp.TaxResponsibility)
	}

	// Check for duplicate document number
	existing, err := s.repo.GetByDocument(ctx, tenantSlug, tp.DocumentNumber)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existing != nil {
		return fmt.Errorf("document number already exists: %s", tp.DocumentNumber)
	}

	return s.repo.Create(ctx, tenantSlug, tp)
}

func (s *service) GetByID(ctx context.Context, tenantSlug string, id int64) (*ThirdParty, error) {
	return s.repo.GetByID(ctx, tenantSlug, id)
}

func (s *service) GetByDocument(ctx context.Context, tenantSlug string, docNumber string) (*ThirdParty, error) {
	return s.repo.GetByDocument(ctx, tenantSlug, docNumber)
}

func (s *service) Update(ctx context.Context, tenantSlug string, id int64, tp *ThirdParty) error {
	// Validate document type
	validDocTypes := map[string]bool{
		DocumentTypeCC:       true,
		DocumentTypeCE:       true,
		DocumentTypeNIT:      true,
		DocumentTypePASSPORT: true,
		DocumentTypeRUT:      true,
	}
	if !validDocTypes[tp.DocumentType] {
		return fmt.Errorf("invalid document type: %s", tp.DocumentType)
	}

	// Validate tax responsibility
	validTaxResp := map[string]bool{
		TaxRespResponsible:    true,
		TaxRespNotResponsible: true,
	}
	if !validTaxResp[tp.TaxResponsibility] {
		return fmt.Errorf("invalid tax responsibility: %s", tp.TaxResponsibility)
	}

	return s.repo.Update(ctx, tenantSlug, id, tp)
}

func (s *service) Delete(ctx context.Context, tenantSlug string, id int64) error {
	return s.repo.Delete(ctx, tenantSlug, id)
}

func (s *service) List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ThirdPartyFilters) ([]ThirdParty, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	return s.repo.List(ctx, tenantSlug, enterpriseID, filters)
}
