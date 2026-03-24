package commissions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrRuleNotFound       = errors.New("regla de comisión no encontrada")
	ErrCommissionNotFound = errors.New("comisión no encontrada")
	ErrInvalidType        = errors.New("tipo de comisión inválido")
	ErrAlreadySettled     = errors.New("comisión ya liquidada")
	ErrAlreadyCancelled   = errors.New("comisión ya cancelada")
)

type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

// ─── HU-COMM-001: Configure Commission Rules ─────────────────────────────────

func (s *service) CreateRule(ctx context.Context, req CreateRuleRequest) (*CommissionRule, error) {
	if req.CommissionType != TypePercentageSale && req.CommissionType != TypePercentageMargin && req.CommissionType != TypeFixedAmount {
		return nil, ErrInvalidType
	}

	rule := &CommissionRule{
		Name:           req.Name,
		CommissionType: req.CommissionType,
		EmployeeID:     req.EmployeeID,
		ProductID:      req.ProductID,
		CategoryID:     req.CategoryID,
		Value:          req.Value,
		MinSaleAmount:  req.MinSaleAmount,
	}

	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			rule.StartDate = &t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			rule.EndDate = &t
		}
	}

	id, err := s.repo.CreateRule(ctx, rule)
	if err != nil {
		return nil, fmt.Errorf("creando regla: %w", err)
	}
	rule.ID = id

	return rule, nil
}

func (s *service) UpdateRule(ctx context.Context, id int64, req UpdateRuleRequest) (*CommissionRule, error) {
	rule, err := s.repo.GetRuleByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRuleNotFound
		}
		return nil, fmt.Errorf("obteniendo regla: %w", err)
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	if req.Value > 0 {
		rule.Value = req.Value
	}
	rule.MinSaleAmount = req.MinSaleAmount
	rule.IsActive = req.IsActive

	if err := s.repo.UpdateRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("actualizando regla: %w", err)
	}

	return rule, nil
}

func (s *service) DeleteRule(ctx context.Context, id int64) error {
	if err := s.repo.DeleteRule(ctx, id); err != nil {
		return fmt.Errorf("eliminando regla: %w", err)
	}
	return nil
}

func (s *service) ListRules(ctx context.Context, activeOnly bool) ([]CommissionRule, error) {
	return s.repo.ListRules(ctx, activeOnly)
}

// ─── HU-COMM-002: Calculate Commissions on Sale ──────────────────────────────

func (s *service) CalculateCommissions(ctx context.Context, req CalculateCommissionRequest) ([]Commission, error) {
	// Get applicable rules
	rules, err := s.repo.GetApplicableRules(ctx, &req.EmployeeID, req.ProductID, req.CategoryID, req.SaleAmount)
	if err != nil {
		return nil, fmt.Errorf("obteniendo reglas aplicables: %w", err)
	}

	var commissions []Commission
	for _, rule := range rules {
		comm := &Commission{
			SalesOrderID:   req.SalesOrderID,
			EmployeeID:     req.EmployeeID,
			BranchID:       req.BranchID,
			RuleID:         &rule.ID,
			SaleAmount:     req.SaleAmount,
			CommissionType: rule.CommissionType,
			CommissionRate: rule.Value,
		}

		// Calculate commission amount based on type
		switch rule.CommissionType {
		case TypePercentageSale:
			comm.CommissionAmount = req.SaleAmount * rule.Value / 100
		case TypePercentageMargin:
			// Default margin if not provided
			margin := 30.0
			comm.ProfitMargin = &margin
			comm.CommissionAmount = (req.SaleAmount * margin / 100) * rule.Value / 100
		case TypeFixedAmount:
			comm.CommissionAmount = rule.Value
		}

		id, err := s.repo.CreateCommission(ctx, comm)
		if err != nil {
			return nil, fmt.Errorf("creando comisión: %w", err)
		}
		comm.ID = id
		commissions = append(commissions, *comm)
	}

	return commissions, nil
}

// ─── HU-COMM-003: View Commission History ─────────────────────────────────────

func (s *service) ListCommissions(ctx context.Context, employeeID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Commission, int64, error) {
	return s.repo.ListCommissions(ctx, employeeID, status, startDate, endDate, page, limit)
}

// ─── HU-COMM-004: Settle Commissions ─────────────────────────────────────────

func (s *service) SettleCommissions(ctx context.Context, settledBy int64, req SettleCommissionsRequest) error {
	if err := s.repo.SettleCommissions(ctx, req.CommissionIDs, settledBy, req.SettlementPeriod, req.Notes); err != nil {
		return fmt.Errorf("liquidando comisiones: %w", err)
	}
	return nil
}

// ─── HU-COMM-005: Commission Reports ─────────────────────────────────────────

func (s *service) GetCommissionReport(ctx context.Context, filter CommissionReportFilter) ([]CommissionSummary, error) {
	return s.repo.GetCommissionSummary(ctx, filter.EmployeeID, filter.StartDate, filter.EndDate)
}

func (s *service) GetCommissionTotals(ctx context.Context, employeeID *int64, startDate, endDate *time.Time) (map[string]float64, error) {
	totalSales, totalCommissions, pendingAmount, settledAmount, err := s.repo.GetCommissionTotals(ctx, employeeID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"total_sales":       totalSales,
		"total_commissions": totalCommissions,
		"pending_amount":    pendingAmount,
		"settled_amount":    settledAmount,
	}, nil
}

// ─── Additional ──────────────────────────────────────────────────────────────

func (s *service) GetCommissionByID(ctx context.Context, id int64) (*Commission, error) {
	return s.repo.GetCommissionByID(ctx, id)
}
