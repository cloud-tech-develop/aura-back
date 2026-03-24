package commissions

import (
	"context"
	"time"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	EventRuleCreated = "commission.rule.created"
	EventRuleUpdated = "commission.rule.updated"
	EventCalculated  = "commission.calculated"
	EventSettled     = "commission.settled"
	EventCancelled   = "commission.cancelled"
)

const (
	TypePercentageSale   = "PERCENTAGE_SALE"
	TypePercentageMargin = "PERCENTAGE_MARGIN"
	TypeFixedAmount      = "FIXED_AMOUNT"
)

const (
	StatusPending   = "PENDING"
	StatusSettled   = "SETTLED"
	StatusCancelled = "CANCELLED"
)

// ─── Entities ─────────────────────────────────────────────────────────────────

// CommissionRule represents a commission configuration rule
type CommissionRule struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	CommissionType string     `json:"commission_type"`
	EmployeeID     *int64     `json:"employee_id,omitempty"`
	ProductID      *int64     `json:"product_id,omitempty"`
	CategoryID     *int64     `json:"category_id,omitempty"`
	Value          float64    `json:"value"`
	MinSaleAmount  float64    `json:"min_sale_amount"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

// Commission represents a calculated commission record
type Commission struct {
	ID               int64      `json:"id"`
	SalesOrderID     int64      `json:"sales_order_id"`
	EmployeeID       int64      `json:"employee_id"`
	BranchID         int64      `json:"branch_id"`
	RuleID           *int64     `json:"rule_id,omitempty"`
	SaleAmount       float64    `json:"sale_amount"`
	ProfitMargin     *float64   `json:"profit_margin,omitempty"`
	CommissionType   string     `json:"commission_type"`
	CommissionRate   float64    `json:"commission_rate"`
	CommissionAmount float64    `json:"commission_amount"`
	Status           string     `json:"status"`
	SettledAt        *time.Time `json:"settled_at,omitempty"`
	SettledBy        *int64     `json:"settled_by,omitempty"`
	SettlementPeriod string     `json:"settlement_period"`
	Notes            string     `json:"notes"`
	CreatedAt        time.Time  `json:"created_at"`
}

// CommissionSummary for reporting
type CommissionSummary struct {
	EmployeeID       int64   `json:"employee_id"`
	EmployeeName     string  `json:"employee_name"`
	TotalSales       float64 `json:"total_sales"`
	TotalCommissions float64 `json:"total_commissions"`
	PendingAmount    float64 `json:"pending_amount"`
	SettledAmount    float64 `json:"settled_amount"`
	SalesCount       int64   `json:"sales_count"`
}

// ─── Request DTOs ──────────────────────────────────────────────────────────────

// CreateRuleRequest for HU-COMM-001
type CreateRuleRequest struct {
	Name           string  `json:"name" binding:"required"`
	CommissionType string  `json:"commission_type" binding:"required"`
	EmployeeID     *int64  `json:"employee_id"`
	ProductID      *int64  `json:"product_id"`
	CategoryID     *int64  `json:"category_id"`
	Value          float64 `json:"value" binding:"required,gt=0"`
	MinSaleAmount  float64 `json:"min_sale_amount"`
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
}

// UpdateRuleRequest
type UpdateRuleRequest struct {
	Name          string  `json:"name"`
	Value         float64 `json:"value"`
	MinSaleAmount float64 `json:"min_sale_amount"`
	IsActive      bool    `json:"is_active"`
}

// CalculateCommissionRequest for HU-COMM-002
type CalculateCommissionRequest struct {
	SalesOrderID int64   `json:"sales_order_id" binding:"required"`
	EmployeeID   int64   `json:"employee_id" binding:"required"`
	BranchID     int64   `json:"branch_id" binding:"required"`
	SaleAmount   float64 `json:"sale_amount" binding:"required"`
	ProductID    *int64  `json:"product_id"`
	CategoryID   *int64  `json:"category_id"`
}

// SettleCommissionsRequest for HU-COMM-004
type SettleCommissionsRequest struct {
	CommissionIDs    []int64 `json:"commission_ids" binding:"required,min=1"`
	SettlementPeriod string  `json:"settlement_period"`
	Notes            string  `json:"notes"`
}

// CommissionReportFilter for HU-COMM-005
type CommissionReportFilter struct {
	EmployeeID *int64
	BranchID   *int64
	Status     string
	StartDate  *time.Time
	EndDate    *time.Time
}

// ─── Repository Interface ─────────────────────────────────────────────────────

type Repository interface {
	// CommissionRule
	CreateRule(ctx context.Context, rule *CommissionRule) (int64, error)
	GetRuleByID(ctx context.Context, id int64) (*CommissionRule, error)
	UpdateRule(ctx context.Context, rule *CommissionRule) error
	DeleteRule(ctx context.Context, id int64) error
	ListRules(ctx context.Context, activeOnly bool) ([]CommissionRule, error)
	GetApplicableRules(ctx context.Context, employeeID, productID, categoryID *int64, saleAmount float64) ([]CommissionRule, error)

	// Commission
	CreateCommission(ctx context.Context, c *Commission) (int64, error)
	GetCommissionByID(ctx context.Context, id int64) (*Commission, error)
	ListCommissions(ctx context.Context, employeeID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Commission, int64, error)
	SettleCommissions(ctx context.Context, ids []int64, settledBy int64, period string, notes string) error
	CancelCommission(ctx context.Context, id int64, notes string) error

	// Reporting
	GetCommissionSummary(ctx context.Context, employeeID *int64, startDate, endDate *time.Time) ([]CommissionSummary, error)
	GetCommissionTotals(ctx context.Context, employeeID *int64, startDate, endDate *time.Time) (totalSales, totalCommissions, pendingAmount, settledAmount float64, err error)
}

// ─── Service Interface ────────────────────────────────────────────────────────

type Service interface {
	// HU-COMM-001: Configure Commission Rules
	CreateRule(ctx context.Context, req CreateRuleRequest) (*CommissionRule, error)
	UpdateRule(ctx context.Context, id int64, req UpdateRuleRequest) (*CommissionRule, error)
	DeleteRule(ctx context.Context, id int64) error
	ListRules(ctx context.Context, activeOnly bool) ([]CommissionRule, error)

	// HU-COMM-002: Calculate Commissions on Sale
	CalculateCommissions(ctx context.Context, req CalculateCommissionRequest) ([]Commission, error)

	// HU-COMM-003: View Commission History
	ListCommissions(ctx context.Context, employeeID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Commission, int64, error)

	// HU-COMM-004: Settle Commissions
	SettleCommissions(ctx context.Context, settledBy int64, req SettleCommissionsRequest) error

	// HU-COMM-005: Commission Reports
	GetCommissionReport(ctx context.Context, filter CommissionReportFilter) ([]CommissionSummary, error)
	GetCommissionTotals(ctx context.Context, employeeID *int64, startDate, endDate *time.Time) (map[string]float64, error)

	// Additional
	GetCommissionByID(ctx context.Context, id int64) (*Commission, error)
}
