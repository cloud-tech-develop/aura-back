package reports

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		svc: NewService(db),
	}
}

func (h *Handler) GetSalesSummary(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	filters := h.parseReportFilters(c)

	report, err := h.svc.GetSalesSummary(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, report)
}

func (h *Handler) GetProductSales(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	filters := h.parseReportFilters(c)

	reports, err := h.svc.GetProductSales(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, reports)
}

func (h *Handler) GetPaymentMethodBreakdown(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	filters := h.parseReportFilters(c)

	reports, err := h.svc.GetPaymentMethodBreakdown(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, reports)
}

func (h *Handler) GetDailySales(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	filters := h.parseReportFilters(c)

	reports, err := h.svc.GetDailySales(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, reports)
}

func (h *Handler) GetTopCustomers(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	filters := h.parseReportFilters(c)

	customers, err := h.svc.GetTopCustomers(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, customers)
}

func (h *Handler) GetSalesByPeriod(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		response.BadRequest(c, "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.BadRequest(c, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.BadRequest(c, "invalid end_date format")
		return
	}

	if endDate.After(time.Now()) {
		response.BadRequest(c, "end_date cannot be in the future")
		return
	}

	daysDiff := int(endDate.Sub(startDate).Hours() / 24)
	if daysDiff > 365 {
		response.BadRequest(c, "date range cannot exceed 365 days")
		return
	}

	filters := ReportFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if id, err := strconv.ParseInt(branchIDStr, 10, 64); err == nil {
			filters.BranchID = &id
		}
	}

	if groupBy := c.Query("group_by"); groupBy != "" {
		filters.GroupBy = groupBy
	}

	report, err := h.svc.GetSalesByPeriod(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, report)
}

func (h *Handler) GetSalesByProduct(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		response.BadRequest(c, "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.BadRequest(c, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.BadRequest(c, "invalid end_date format")
		return
	}

	filters := ReportFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if id, err := strconv.ParseInt(branchIDStr, 10, 64); err == nil {
			filters.BranchID = &id
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if id, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			filters.CategoryID = &id
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	report, err := h.svc.GetSalesByProduct(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, report)
}

func (h *Handler) GetSalesByEmployee(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		response.BadRequest(c, "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.BadRequest(c, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.BadRequest(c, "invalid end_date format")
		return
	}

	filters := ReportFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if id, err := strconv.ParseInt(branchIDStr, 10, 64); err == nil {
			filters.BranchID = &id
		}
	}

	report, err := h.svc.GetSalesByEmployee(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, report)
}

func (h *Handler) GetInventoryStatus(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	filters := ReportFilters{}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if id, err := strconv.ParseInt(branchIDStr, 10, 64); err == nil {
			filters.BranchID = &id
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if id, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			filters.CategoryID = &id
		}
	}

	stockFilter := c.Query("stock_filter")

	report, err := h.svc.GetInventoryStatus(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	if stockFilter == "low" {
		var filteredProducts []InventoryProduct
		for _, p := range report.Products {
			if p.Status == "LOW_STOCK" {
				filteredProducts = append(filteredProducts, p)
			}
		}
		report.Products = filteredProducts
	} else if stockFilter == "out" {
		var filteredProducts []InventoryProduct
		for _, p := range report.Products {
			if p.Status == "OUT_OF_STOCK" {
				filteredProducts = append(filteredProducts, p)
			}
		}
		report.Products = filteredProducts
	} else if stockFilter == "normal" {
		var filteredProducts []InventoryProduct
		for _, p := range report.Products {
			if p.Status == "NORMAL" {
				filteredProducts = append(filteredProducts, p)
			}
		}
		report.Products = filteredProducts
	}

	response.OK(c, report)
}

func (h *Handler) GetMovementHistory(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		response.BadRequest(c, "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.BadRequest(c, "invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.BadRequest(c, "invalid end_date format")
		return
	}

	filters := ReportFilters{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	if productIDStr := c.Query("product_id"); productIDStr != "" {
		if id, err := strconv.ParseInt(productIDStr, 10, 64); err == nil {
			filters.ProductID = &id
		}
	}

	if movementType := c.Query("movement_type"); movementType != "" {
		filters.GroupBy = movementType
	}

	report, err := h.svc.GetMovementHistory(c.Request.Context(), enterpriseID, filters)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, report)
}

func (h *Handler) ExportToPDF(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	reportType := c.Param("reportType")
	if reportType == "" {
		response.BadRequest(c, "reportType is required")
		return
	}

	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	data, err := h.svc.ExportToPDF(c.Request.Context(), enterpriseID, reportType, req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=\""+reportType+"-report.pdf\"")
	c.Data(200, "application/pdf", data)
}

func (h *Handler) ExportToExcel(c *gin.Context) {
	enterpriseID := c.GetInt64("enterprise_id")
	if enterpriseID == 0 {
		response.BadRequest(c, "enterprise_id not found")
		return
	}

	reportType := c.Param("reportType")
	if reportType == "" {
		response.BadRequest(c, "reportType is required")
		return
	}

	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	data, err := h.svc.ExportToExcel(c.Request.Context(), enterpriseID, reportType, req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\""+reportType+"-report.xlsx\"")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *Handler) parseReportFilters(c *gin.Context) ReportFilters {
	filters := ReportFilters{}

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters.StartDate = &t
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filters.EndDate = &t
		}
	}

	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		if id, err := strconv.ParseInt(branchIDStr, 10, 64); err == nil {
			filters.BranchID = &id
		}
	}

	if productIDStr := c.Query("product_id"); productIDStr != "" {
		if id, err := strconv.ParseInt(productIDStr, 10, 64); err == nil {
			filters.ProductID = &id
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if id, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			filters.CategoryID = &id
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	return filters
}
