package payments

import (
	"database/sql"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(database *db.DB) *Handler {
	q := database.Wrap(database.DB)
	return &Handler{svc: NewService(q)}
}

// ProcessPayment - POST /payments
func (h *Handler) ProcessPayment(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	userID := c.GetInt64("user_id")

	if enterpriseID == 0 || userID == 0 {
		response.BadRequest(c, "enterprise_id or user_id not found")
		return
	}

	var req struct {
		ReferenceID      int64   `json:"reference_id" binding:"required"`
		ReferenceType    string  `json:"reference_type" binding:"required"`
		PaymentMethod    string  `json:"payment_method" binding:"required"`
		Amount           float64 `json:"amount" binding:"required,gt=0"`
		ReferenceNumber  string  `json:"reference_number"`
		CardType         string  `json:"card_type"`
		CardLastDigits   string  `json:"card_last_digits"`
		BankName         string  `json:"bank_name"`
		AuthorizationCode string `json:"authorization_code"`
		ChangeAmount     float64 `json:"change_amount"`
		CashDrawerID     *int64  `json:"cash_drawer_id"`
		BranchID         int64   `json:"branch_id" binding:"required"`
		Notes            string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	payment := &Payment{
		ReferenceID:       req.ReferenceID,
		ReferenceType:     req.ReferenceType,
		PaymentMethod:     req.PaymentMethod,
		Amount:            req.Amount,
		ReferenceNumber:   req.ReferenceNumber,
		CardType:          req.CardType,
		CardLastDigits:    req.CardLastDigits,
		BankName:          req.BankName,
		AuthorizationCode: req.AuthorizationCode,
		ChangeAmount:      req.ChangeAmount,
		CashDrawerID:      req.CashDrawerID,
		UserID:            userID,
		BranchID:          req.BranchID,
		EnterpriseID:      enterpriseID,
		Notes:             req.Notes,
	}

	if err := h.svc.ProcessPayment(c.Request.Context(), payment); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, payment)
}

// ProcessMultiplePayments - POST /payments/batch
func (h *Handler) ProcessMultiplePayments(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	userID := c.GetInt64("user_id")

	if enterpriseID == 0 || userID == 0 {
		response.BadRequest(c, "enterprise_id or user_id not found")
		return
	}

	var req struct {
		ReferenceID   int64    `json:"reference_id" binding:"required"`
		ReferenceType string   `json:"reference_type" binding:"required"`
		BranchID       int64    `json:"branch_id" binding:"required"`
		CashDrawerID   *int64   `json:"cash_drawer_id"`
		Payments      []struct {
			PaymentMethod    string  `json:"payment_method" binding:"required"`
			Amount           float64 `json:"amount" binding:"required,gt=0"`
			ReferenceNumber  string  `json:"reference_number"`
			CardType         string  `json:"card_type"`
			CardLastDigits   string  `json:"card_last_digits"`
			BankName         string  `json:"bank_name"`
			AuthorizationCode string `json:"authorization_code"`
		} `json:"payments" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var payments []Payment
	for _, p := range req.Payments {
		payments = append(payments, Payment{
			ReferenceID:       req.ReferenceID,
			ReferenceType:      req.ReferenceType,
			PaymentMethod:      p.PaymentMethod,
			Amount:             p.Amount,
			ReferenceNumber:    p.ReferenceNumber,
			CardType:           p.CardType,
			CardLastDigits:     p.CardLastDigits,
			BankName:           p.BankName,
			AuthorizationCode:  p.AuthorizationCode,
			CashDrawerID:       req.CashDrawerID,
			BranchID:           req.BranchID,
			EnterpriseID:       enterpriseID,
			UserID:             userID,
		})
	}

	if err := h.svc.ProcessMultiplePayments(c.Request.Context(), payments); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, gin.H{"message": "payments processed", "count": len(payments)})
}

// ListPayments - GET /payments
func (h *Handler) ListPayments(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var refID *int64
	if rid := c.Query("reference_id"); rid != "" {
		if id, err := strconv.ParseInt(rid, 10, 64); err == nil {
			refID = &id
		}
	}

	filters := PaymentFilters{
		Page:            page,
		Limit:           limit,
		PaymentMethod:   c.Query("method"),
		Status:          c.Query("status"),
		ReferenceID:     refID,
	}

	list, err := h.svc.ListPayments(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// GetPayment - GET /payments/:id
func (h *Handler) GetPayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid payment ID")
		return
	}

	payment, err := h.svc.GetPayment(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "payment not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, payment)
}

// GetPaymentsByOrder - GET /payments/reference/:referenceType/:referenceId
func (h *Handler) GetPaymentsByOrder(c *gin.Context) {
	refType := c.Param("referenceType")
	refID, err := strconv.ParseInt(c.Param("referenceId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reference ID")
		return
	}

	payments, err := h.svc.GetPaymentsByOrder(c.Request.Context(), refType, refID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, payments)
}

// CancelPayment - POST /payments/:id/cancel
func (h *Handler) CancelPayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid payment ID")
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.BadRequest(c, "user_id not found")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.CancelPayment(c.Request.Context(), id, userID, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "payment cancelled"})
}

// Cash Drawer Handlers

// OpenCashDrawer - POST /cash-drawers
func (h *Handler) OpenCashDrawer(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	userID := c.GetInt64("user_id")

	if enterpriseID == 0 || userID == 0 {
		response.BadRequest(c, "enterprise_id or user_id not found")
		return
	}

	var req struct {
		BranchID       int64   `json:"branch_id" binding:"required"`
		OpeningBalance float64 `json:"opening_balance" binding:"required"`
		Notes          string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	drawer := &CashDrawer{
		UserID:         userID,
		BranchID:       req.BranchID,
		EnterpriseID:   enterpriseID,
		OpeningBalance: req.OpeningBalance,
		Notes:          req.Notes,
	}

	if err := h.svc.OpenCashDrawer(c.Request.Context(), drawer); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, drawer)
}

// GetCashDrawer - GET /cash-drawers/:id
func (h *Handler) GetCashDrawer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid drawer ID")
		return
	}

	drawer, err := h.svc.GetCashDrawer(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "cash drawer not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, drawer)
}

// GetOpenDrawer - GET /cash-drawers/open
func (h *Handler) GetOpenDrawer(c *gin.Context) {
	userID := c.GetInt64("user_id")
	branchID, err := strconv.ParseInt(c.Query("branch_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid branch_id")
		return
	}
	if userID == 0 {
		response.BadRequest(c, "user_id not found")
		return
	}

	drawer, err := h.svc.GetOpenDrawer(c.Request.Context(), userID, branchID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "no open cash drawer")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, drawer)
}

// CloseCashDrawer - POST /cash-drawers/:id/close
func (h *Handler) CloseCashDrawer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid drawer ID")
		return
	}

	var req struct {
		ClosingBalance float64 `json:"closing_balance" binding:"required"`
		Notes           string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.CloseCashDrawer(c.Request.Context(), id, req.ClosingBalance, req.Notes); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "cash drawer closed"})
}

// ListCashDrawers - GET /cash-drawers
func (h *Handler) ListCashDrawers(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	var userID *int64
	if uid := c.Query("user_id"); uid != "" {
		if id, err := strconv.ParseInt(uid, 10, 64); err == nil {
			userID = &id
		}
	}

	status := c.Query("status")

	list, err := h.svc.ListCashDrawers(c.Request.Context(), enterpriseID, userID, status)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, list)
}

// AddCashIn - POST /cash-drawers/:id/cash-in
func (h *Handler) AddCashIn(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid drawer ID")
		return
	}

	var req struct {
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AddCashIn(c.Request.Context(), id, req.Amount, req.Description); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "cash added"})
}

// AddCashOut - POST /cash-drawers/:id/cash-out
func (h *Handler) AddCashOut(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid drawer ID")
		return
	}

	var req struct {
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AddCashOut(c.Request.Context(), id, req.Amount, req.Description); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "cash removed"})
}
