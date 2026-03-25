package invoices

import (
	"database/sql"
	"strconv"
	"time"

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

// GenerateInvoiceFromSale - POST /invoices/generate
func (h *Handler) GenerateInvoiceFromSale(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	userID := c.GetInt64("user_id")

	if enterpriseID == 0 || userID == 0 {
		response.BadRequest(c, "enterprise_id or user_id not found")
		return
	}

	var req struct {
		SalesOrderID int64 `json:"sales_order_id" binding:"required"`
		PrefixID     int64 `json:"prefix_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	inv, err := h.svc.GenerateInvoiceFromSale(c.Request.Context(), req.SalesOrderID, req.PrefixID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, inv)
}

// CreateInvoice - POST /invoices
func (h *Handler) CreateInvoice(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	userID := c.GetInt64("user_id")

	if enterpriseID == 0 || userID == 0 {
		response.BadRequest(c, "enterprise_id or user_id not found")
		return
	}

	var req struct {
		PrefixID      int64    `json:"prefix_id"`
		InvoiceType   string   `json:"invoice_type"`
		CustomerID    int64    `json:"customer_id" binding:"required"`
		BranchID      int64    `json:"branch_id" binding:"required"`
		SalesOrderID  *int64   `json:"sales_order_id"`
		DueDate       *string  `json:"due_date"`
		Subtotal      float64  `json:"subtotal" binding:"required"`
		DiscountTotal float64  `json:"discount_total"`
		TaxExempt     float64  `json:"tax_exempt"`
		TaxableAmount float64  `json:"taxable_amount"`
		Iva19         float64  `json:"iva_19"`
		Iva5          float64  `json:"iva_5"`
		Reteica       float64  `json:"reteica"`
		Retefuente    float64  `json:"retefuente"`
		Total         float64  `json:"total" binding:"required"`
		PaymentMethod string   `json:"payment_method"`
		Notes         string   `json:"notes"`
		Items         []struct {
			ProductID   int64   `json:"product_id" binding:"required"`
			ProductName string  `json:"product_name" binding:"required"`
			ProductSKU  string  `json:"product_sku"`
			Quantity    float64 `json:"quantity" binding:"required"`
			UnitPrice   float64 `json:"unit_price" binding:"required"`
			TaxRate     float64 `json:"tax_rate"`
			LineTotal   float64 `json:"line_total" binding:"required"`
		} `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get prefix to generate invoice number
	var invoiceNumber string
	if req.PrefixID > 0 {
		prefix, err := h.svc.GetInvoicePrefix(c.Request.Context(), req.BranchID, "")
		if err == nil {
			prefix.CurrentNumber++
			invoiceNumber = prefix.Prefix + "-" + strconv.FormatInt(prefix.CurrentNumber, 10)
		}
	}

	inv := &Invoice{
		InvoiceNumber: invoiceNumber,
		PrefixID:      req.PrefixID,
		InvoiceType:   req.InvoiceType,
		CustomerID:    req.CustomerID,
		BranchID:      req.BranchID,
		UserID:        userID,
		EnterpriseID: enterpriseID,
		SalesOrderID:  req.SalesOrderID,
		Subtotal:      req.Subtotal,
		DiscountTotal: req.DiscountTotal,
		TaxExempt:     req.TaxExempt,
		TaxableAmount: req.TaxableAmount,
		Iva19:         req.Iva19,
		Iva5:          req.Iva5,
		Reteica:       req.Reteica,
		Retefuente:    req.Retefuente,
		Total:         req.Total,
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
		InvoiceDate:   time.Now(),
		Status:        InvoiceStatusDraft,
	}

	if req.DueDate != nil {
		dueDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			inv.DueDate = &dueDate
		}
	}

	if err := h.svc.GenerateInvoice(c.Request.Context(), inv); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Create invoice items
	for _, itemReq := range req.Items {
		item := &InvoiceItem{
			InvoiceID:   inv.ID,
			ProductID:   itemReq.ProductID,
			ProductName: itemReq.ProductName,
			ProductSKU:  itemReq.ProductSKU,
			Quantity:    itemReq.Quantity,
			UnitPrice:   itemReq.UnitPrice,
			TaxRate:     itemReq.TaxRate,
			LineTotal:   itemReq.LineTotal,
		}
		// Would need to add CreateInvoiceItem to repo
		_ = item
	}

	response.Created(c, inv)
}

// GetInvoice - GET /invoices/:id
func (h *Handler) GetInvoice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	inv, err := h.svc.GetInvoice(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Factura no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, inv)
}

// GetInvoiceByNumber - GET /invoices/number/:invoiceNumber
func (h *Handler) GetInvoiceByNumber(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	invoiceNumber := c.Param("invoiceNumber")

	inv, err := h.svc.GetInvoiceByNumber(c.Request.Context(), invoiceNumber, enterpriseID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Factura no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, inv)
}

// ListInvoices - GET /invoices
func (h *Handler) ListInvoices(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var branchID, customerID *int64
	if bid := c.Query("branch_id"); bid != "" {
		if id, err := strconv.ParseInt(bid, 10, 64); err == nil {
			branchID = &id
		}
	}
	if cid := c.Query("customer_id"); cid != "" {
		if id, err := strconv.ParseInt(cid, 10, 64); err == nil {
			customerID = &id
		}
	}

	var startDate, endDate *time.Time
	if sd := c.Query("start_date"); sd != "" {
		if t, err := time.Parse("2006-01-02", sd); err == nil {
			startDate = &t
		}
	}
	if ed := c.Query("end_date"); ed != "" {
		if t, err := time.Parse("2006-01-02", ed); err == nil {
			endDate = &t
		}
	}

	filters := InvoiceFilters{
		Page:         page,
		Limit:        limit,
		Status:       c.Query("status"),
		InvoiceType:  c.Query("type"),
		BranchID:     branchID,
		CustomerID:   customerID,
		StartDate:    startDate,
		EndDate:      endDate,
	}

	invoices, err := h.svc.GetInvoices(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, invoices)
}

// IssueInvoice - POST /invoices/:id/issue
func (h *Handler) IssueInvoice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	if err := h.svc.IssueInvoice(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Factura emitidas"})
}

// CancelInvoice - POST /invoices/:id/cancel
func (h *Handler) CancelInvoice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.CancelInvoice(c.Request.Context(), id, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Factura cancelada"})
}

// GetInvoiceLogs - GET /invoices/:id/logs
func (h *Handler) GetInvoiceLogs(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	logs, err := h.svc.GetInvoiceLogs(c.Request.Context(), id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, logs)
}

// Invoice Prefix Handlers

// CreateInvoicePrefix - POST /invoice-prefixes
func (h *Handler) CreateInvoicePrefix(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	var req struct {
		BranchID         int64    `json:"branch_id" binding:"required"`
		Prefix           string   `json:"prefix" binding:"required"`
		ResolutionNumber string   `json:"resolution_number"`
		ResolutionDate   *string  `json:"resolution_date"`
		ValidFrom        *string  `json:"valid_from"`
		ValidUntil       *string  `json:"valid_until"`
		Description      string   `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	prefix := &InvoicePrefix{
		BranchID:        req.BranchID,
		EnterpriseID:    enterpriseID,
		Prefix:          req.Prefix,
		ResolutionNumber: req.ResolutionNumber,
		Description:     req.Description,
	}

	if req.ResolutionDate != nil {
		if rd, err := time.Parse("2006-01-02", *req.ResolutionDate); err == nil {
			prefix.ResolutionDate = &rd
		}
	}
	if req.ValidFrom != nil {
		if vf, err := time.Parse("2006-01-02", *req.ValidFrom); err == nil {
			prefix.ValidFrom = &vf
		}
	}
	if req.ValidUntil != nil {
		if vu, err := time.Parse("2006-01-02", *req.ValidUntil); err == nil {
			prefix.ValidUntil = &vu
		}
	}

	if err := h.svc.CreateInvoicePrefix(c.Request.Context(), prefix); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, prefix)
}

// ListInvoicePrefixes - GET /invoice-prefixes
func (h *Handler) ListInvoicePrefixes(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	prefixes, err := h.svc.GetInvoicePrefixes(c.Request.Context(), enterpriseID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, prefixes)
}
