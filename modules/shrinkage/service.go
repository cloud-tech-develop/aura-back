package shrinkage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrShrinkageNotFound   = errors.New("merma no encontrada")
	ErrReasonNotFound      = errors.New("razón de merma no encontrada")
	ErrReasonAlreadyExists = errors.New("razón de merma ya existe")
	ErrReasonInUse         = errors.New("razón de merma en uso, no se puede eliminar")
	ErrAlreadyAuthorized   = errors.New("merma ya fue autorizada")
	ErrAlreadyCancelled    = errors.New("merma ya fue cancelada")
	ErrNoAuthorization     = errors.New("merma no requiere autorización")
	ErrInvalidItems        = errors.New("items inválidos")
)

type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) generateShrinkageNumber(branchID int64) string {
	return fmt.Sprintf("SHR-%d-%d", branchID, time.Now().UnixNano())
}

// ─── HU-SHR-001: Register Shrinkage ──────────────────────────────────────────

func (s *service) RegisterShrinkage(ctx context.Context, userID int64, req RegisterShrinkageRequest) (*Shrinkage, error) {
	if len(req.Items) == 0 {
		return nil, ErrInvalidItems
	}

	// Verify reason exists
	reason, err := s.repo.GetReasonByID(ctx, req.ReasonID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrReasonNotFound
		}
		return nil, fmt.Errorf("verificando razón: %w", err)
	}

	shrinkDate := time.Now()
	if req.ShrinkageDate != "" {
		if t, err := time.Parse("2006-01-02", req.ShrinkageDate); err == nil {
			shrinkDate = t
		}
	}

	shrinkage := &Shrinkage{
		ShrinkageNumber: s.generateShrinkageNumber(req.BranchID),
		BranchID:        req.BranchID,
		UserID:          userID,
		ReasonID:        req.ReasonID,
		ShrinkageDate:   shrinkDate,
		Notes:           req.Notes,
	}

	// Calculate total value and create items
	var totalValue float64
	var totalQuantity float64
	for _, itemReq := range req.Items {
		totalValue += itemReq.Quantity * itemReq.UnitCost
		totalQuantity += itemReq.Quantity
	}

	shrinkage.TotalValue = totalValue

	// Check if authorization is required
	if reason.RequiresAuthorization {
		if reason.AuthorizationThreshold == nil || totalValue >= *reason.AuthorizationThreshold {
			shrinkage.Status = StatusPending
		} else {
			shrinkage.Status = StatusApproved
		}
	} else {
		shrinkage.Status = StatusApproved
	}

	shrinkageID, err := s.repo.CreateShrinkage(ctx, shrinkage)
	if err != nil {
		return nil, fmt.Errorf("creando merma: %w", err)
	}
	shrinkage.ID = shrinkageID

	// Create shrinkage items
	for _, itemReq := range req.Items {
		item := &ShrinkageItem{
			ShrinkageID:  shrinkageID,
			ProductID:    itemReq.ProductID,
			BatchNumber:  &itemReq.BatchNumber,
			SerialNumber: &itemReq.SerialNumber,
			Quantity:     itemReq.Quantity,
			UnitCost:     itemReq.UnitCost,
			TotalValue:   itemReq.Quantity * itemReq.UnitCost,
			ReasonDetail: itemReq.ReasonDetail,
		}
		if *item.BatchNumber == "" {
			item.BatchNumber = nil
		}
		if *item.SerialNumber == "" {
			item.SerialNumber = nil
		}
		if err := s.repo.CreateShrinkageItem(ctx, item); err != nil {
			return nil, fmt.Errorf("creando item de merma: %w", err)
		}
	}

	return shrinkage, nil
}

// ─── HU-SHR-002: Configure Shrinkage Reasons ──────────────────────────────────

func (s *service) CreateReason(ctx context.Context, req CreateReasonRequest) (*ShrinkageReason, error) {
	// Check if code already exists
	existing, err := s.repo.GetReasonByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrReasonAlreadyExists
	}

	reason := &ShrinkageReason{
		Code:                   req.Code,
		Name:                   req.Name,
		Description:            req.Description,
		RequiresAuthorization:  req.RequiresAuthorization,
		AuthorizationThreshold: req.AuthorizationThreshold,
		IsActive:               true,
	}

	if err := s.repo.CreateReason(ctx, reason); err != nil {
		return nil, fmt.Errorf("creando razón: %w", err)
	}
	return reason, nil
}

func (s *service) ListReasons(ctx context.Context, activeOnly bool) ([]ShrinkageReason, error) {
	return s.repo.ListReasons(ctx, activeOnly)
}

// ─── HU-SHR-003: Authorize High-Value Shrinkage ────────────────────────────────

func (s *service) AuthorizeShrinkage(ctx context.Context, shrinkageID int64, userID int64, approved bool, notes string) (*Shrinkage, error) {
	shrinkage, err := s.repo.GetShrinkageByID(ctx, shrinkageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrShrinkageNotFound
		}
		return nil, fmt.Errorf("obteniendo merma: %w", err)
	}

	if shrinkage.Status != StatusPending {
		return nil, ErrAlreadyAuthorized
	}

	// Check if reason requires authorization
	reason, err := s.repo.GetReasonByID(ctx, shrinkage.ReasonID)
	if err != nil {
		return nil, fmt.Errorf("obteniendo razón: %w", err)
	}
	if !reason.RequiresAuthorization {
		return nil, ErrNoAuthorization
	}

	now := time.Now()
	shrinkage.AuthorizedBy = &userID
	shrinkage.AuthorizedAt = &now
	if approved {
		shrinkage.Status = StatusApproved
	} else {
		shrinkage.Status = StatusRejected
	}
	shrinkage.Notes = notes

	if err := s.repo.UpdateShrinkage(ctx, shrinkage); err != nil {
		return nil, fmt.Errorf("autorizando merma: %w", err)
	}
	return shrinkage, nil
}

// ─── HU-SHR-004: View Shrinkage Report ─────────────────────────────────────────

func (s *service) GetShrinkageReport(ctx context.Context, branchID *int64, startDate, endDate *time.Time) ([]ShrinkageReportItem, error) {
	return s.repo.GetShrinkageReport(ctx, branchID, startDate, endDate)
}

func (s *service) ListShrinkages(ctx context.Context, branchID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Shrinkage, int64, error) {
	return s.repo.ListShrinkages(ctx, branchID, status, startDate, endDate, page, limit)
}

// ─── HU-SHR-005: Cancel Shrinkage ──────────────────────────────────────────────

func (s *service) CancelShrinkage(ctx context.Context, shrinkageID int64, userID int64, reason string) (*Shrinkage, error) {
	shrinkage, err := s.repo.GetShrinkageByID(ctx, shrinkageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrShrinkageNotFound
		}
		return nil, fmt.Errorf("obteniendo merma: %w", err)
	}

	if shrinkage.Status == StatusCancelled {
		return nil, ErrAlreadyCancelled
	}

	now := time.Now()
	shrinkage.Status = StatusCancelled
	shrinkage.CancellationReason = reason
	shrinkage.CancelledBy = &userID
	shrinkage.CancelledAt = &now

	if err := s.repo.UpdateShrinkage(ctx, shrinkage); err != nil {
		return nil, fmt.Errorf("cancelando merma: %w", err)
	}
	return shrinkage, nil
}

// ─── Additional ──────────────────────────────────────────────────────────────

func (s *service) GetShrinkageByID(ctx context.Context, id int64) (*Shrinkage, error) {
	return s.repo.GetShrinkageByID(ctx, id)
}
