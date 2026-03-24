package cash

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrDrawerNotFound   = errors.New("caja no encontrada")
	ErrShiftNotFound    = errors.New("turno no encontrado")
	ErrShiftAlreadyOpen = errors.New("ya existe un turno abierto para este usuario")
	ErrNoActiveShift    = errors.New("no hay turno activo")
	ErrInvalidAmount    = errors.New("monto inválido")
	ErrUnauthorized     = errors.New("no autorizado para cerrar turno")
)

type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

// ─── HU-CASH-001: Configure Cash Drawer ─────────────────────────────────────

func (s *service) ConfigureDrawer(ctx context.Context, branchID int64, name string, minFloat float64) (*CashDrawer, error) {
	// Check if drawer already exists
	existing, err := s.repo.GetDrawerByBranch(ctx, branchID)
	if err == nil && existing != nil {
		// Update existing drawer
		existing.Name = name
		existing.MinFloat = minFloat
		if err := s.repo.UpdateDrawer(ctx, existing); err != nil {
			return nil, fmt.Errorf("actualizando caja: %w", err)
		}
		return existing, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("verificando caja: %w", err)
	}

	// Create new drawer
	drawer := &CashDrawer{
		BranchID: branchID,
		Name:     name,
		MinFloat: minFloat,
	}
	if err := s.repo.CreateDrawer(ctx, drawer); err != nil {
		return nil, fmt.Errorf("creando caja: %w", err)
	}
	return drawer, nil
}

// ─── HU-CASH-002: Open Cash Shift ───────────────────────────────────────────

func (s *service) OpenShift(ctx context.Context, userID, branchID, cashDrawerID int64, openingAmount float64, notes string) (*CashShift, error) {
	if openingAmount < 0 {
		return nil, ErrInvalidAmount
	}

	// Check if user already has an active shift
	activeShift, err := s.repo.GetActiveShiftByUser(ctx, userID)
	if err == nil && activeShift != nil {
		return nil, ErrShiftAlreadyOpen
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("verificando turno activo: %w", err)
	}

	// Create new shift
	shift := &CashShift{
		CashDrawerID:  cashDrawerID,
		UserID:        userID,
		BranchID:      branchID,
		OpeningAmount: openingAmount,
		OpeningNotes:  notes,
	}

	if err := s.repo.CreateShift(ctx, shift); err != nil {
		return nil, fmt.Errorf("creando turno: %w", err)
	}

	return shift, nil
}

// ─── HU-CASH-003: Close Cash Shift ───────────────────────────────────────────

func (s *service) CloseShift(ctx context.Context, shiftID int64, userID int64, closingAmount float64, notes string) (*CashShift, error) {
	// Get shift
	shift, err := s.repo.GetShiftByID(ctx, shiftID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrShiftNotFound
		}
		return nil, fmt.Errorf("obteniendo turno: %w", err)
	}

	if shift.Status != StatusOpen {
		return nil, fmt.Errorf("el turno no está abierto")
	}

	// Get movements to calculate expected amount
	totalIn, totalOut, err := s.repo.GetMovementsSummary(ctx, shiftID)
	if err != nil {
		return nil, fmt.Errorf("obteniendo movimientos: %w", err)
	}

	// Calculate expected: opening + total IN - total OUT
	expectedAmount := shift.OpeningAmount + totalIn - totalOut

	// Calculate difference
	diff := closingAmount - expectedAmount

	now := time.Now()
	shift.ClosingAmount = &closingAmount
	shift.ExpectedAmount = &expectedAmount
	shift.Difference = &diff
	shift.ClosingNotes = notes
	shift.Status = StatusClosed
	shift.ClosedAt = &now
	shift.ClosedBy = &userID

	if err := s.repo.UpdateShift(ctx, shift); err != nil {
		return nil, fmt.Errorf("cerrando turno: %w", err)
	}

	return shift, nil
}

// ─── HU-CASH-004: Record Cash Movement ───────────────────────────────────────

func (s *service) RecordMovement(ctx context.Context, userID int64, req RecordMovementRequest) (*CashMovement, error) {
	// Validate shift is active
	shift, err := s.repo.GetShiftByID(ctx, req.ShiftID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrShiftNotFound
		}
		return nil, fmt.Errorf("verificando turno: %w", err)
	}

	if shift.Status != StatusOpen {
		return nil, fmt.Errorf("el turno no está abierto")
	}

	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Validate movement type
	if req.MovementType != MovementIn && req.MovementType != MovementOut {
		return nil, fmt.Errorf("tipo de movimiento inválido: use 'IN' o 'OUT'")
	}

	movement := &CashMovement{
		ShiftID:       req.ShiftID,
		MovementType:  req.MovementType,
		Reason:        req.Reason,
		Amount:        req.Amount,
		ReferenceID:   req.ReferenceID,
		ReferenceType: req.ReferenceType,
		Notes:         req.Notes,
		UserID:        userID,
	}

	if err := s.repo.CreateMovement(ctx, movement); err != nil {
		return nil, fmt.Errorf("registrando movimiento: %w", err)
	}

	return movement, nil
}

// ─── HU-CASH-005: Perform Cash Reconciliation ────────────────────────────────

func (s *service) ReconcileShift(ctx context.Context, shiftID int64, expectedAmount float64, notes string) (*CashShift, error) {
	shift, err := s.repo.GetShiftByID(ctx, shiftID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrShiftNotFound
		}
		return nil, fmt.Errorf("obteniendo turno: %w", err)
	}

	if shift.Status != StatusClosed {
		return nil, fmt.Errorf("el turno debe estar cerrado para conciliar")
	}

	// Update expected amount and set to audited
	shift.ExpectedAmount = &expectedAmount
	diff := shift.ClosingAmount
	if diff != nil {
		diffVal := *diff - expectedAmount
		shift.Difference = &diffVal
	}
	shift.Status = StatusAudited
	shift.ClosingNotes = notes

	if err := s.repo.UpdateShift(ctx, shift); err != nil {
		return nil, fmt.Errorf("conciliando turno: %w", err)
	}

	return shift, nil
}

// ─── HU-CASH-006: View Shift Summary ──────────────────────────────────────────

func (s *service) GetShiftSummary(ctx context.Context, shiftID int64) (*ShiftSummaryResponse, error) {
	shift, err := s.repo.GetShiftByID(ctx, shiftID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrShiftNotFound
		}
		return nil, fmt.Errorf("obteniendo turno: %w", err)
	}

	// Get movements
	movements, err := s.repo.GetMovementsByShift(ctx, shiftID)
	if err != nil {
		return nil, fmt.Errorf("obteniendo movimientos: %w", err)
	}

	// Calculate totals
	totalIn, totalOut := 0.0, 0.0
	for _, m := range movements {
		if m.MovementType == MovementIn {
			totalIn += m.Amount
		} else {
			totalOut += m.Amount
		}
	}

	// Calculate expected cash
	expectedCash := shift.OpeningAmount + totalIn - totalOut

	return &ShiftSummaryResponse{
		Shift:        *shift,
		Movements:    movements,
		TotalIn:      totalIn,
		TotalOut:     totalOut,
		ExpectedCash: expectedCash,
	}, nil
}

// ─── Additional Methods ─────────────────────────────────────────────────────

func (s *service) GetActiveShift(ctx context.Context, userID int64) (*CashShift, error) {
	return s.repo.GetActiveShiftByUser(ctx, userID)
}

func (s *service) GetDrawerByBranch(ctx context.Context, branchID int64) (*CashDrawer, error) {
	return s.repo.GetDrawerByBranch(ctx, branchID)
}

func (s *service) ListShifts(ctx context.Context, branchID int64, startDate, endDate *time.Time, status string, page, limit int) ([]CashShift, int64, error) {
	return s.repo.ListShifts(ctx, branchID, startDate, endDate, status, page, limit)
}
