package products

import (
	"database/sql"
	"strconv"

	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	categorySvc CategoryService
	brandSvc    BrandService
	productSvc  ProductService
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		categorySvc: NewCategoryService(db),
		brandSvc:    NewBrandService(db),
		productSvc:  NewProductService(db),
	}
}

// Category Handlers
func (h *Handler) CreateCategory(c *gin.Context) {
	empresaID := c.GetInt64("empresa_id")
	if empresaID == 0 {
		response.BadRequest(c, "empresa_id not found")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ParentID    *int64 `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category := &Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		EmpresaID:   empresaID,
	}

	if err := h.categorySvc.Create(c.Request.Context(), category); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, category)
}

func (h *Handler) ListCategories(c *gin.Context) {
	empresaID := c.GetInt64("empresa_id")
	if empresaID == 0 {
		response.BadRequest(c, "empresa_id not found")
		return
	}

	categories, err := h.categorySvc.List(c.Request.Context(), empresaID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, categories)
}

func (h *Handler) GetCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	category, err := h.categorySvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Categoría no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, category)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ParentID    *int64 `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category := &Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	if err := h.categorySvc.Update(c.Request.Context(), id, category); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Categoría no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, category)
}

// Brand Handlers
func (h *Handler) CreateBrand(c *gin.Context) {
	empresaID := c.GetInt64("empresa_id")
	if empresaID == 0 {
		response.BadRequest(c, "empresa_id not found")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	brand := &Brand{
		Name:        req.Name,
		Description: req.Description,
		EmpresaID:   empresaID,
	}

	if err := h.brandSvc.Create(c.Request.Context(), brand); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, brand)
}

func (h *Handler) ListBrands(c *gin.Context) {
	empresaID := c.GetInt64("empresa_id")
	if empresaID == 0 {
		response.BadRequest(c, "empresa_id not found")
		return
	}

	brands, err := h.brandSvc.List(c.Request.Context(), empresaID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, brands)
}

func (h *Handler) GetBrand(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	brand, err := h.brandSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Marca no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, brand)
}

func (h *Handler) UpdateBrand(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	brand := &Brand{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.brandSvc.Update(c.Request.Context(), id, brand); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Marca no encontrada")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, brand)
}

// Product Handlers
func (h *Handler) CreateProduct(c *gin.Context) {
	empresaID := c.GetInt64("empresa_id")
	if empresaID == 0 {
		response.BadRequest(c, "empresa_id not found")
		return
	}

	var req struct {
		SKU         string  `json:"sku" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		CategoryID  *int64  `json:"category_id"`
		BrandID     *int64  `json:"brand_id"`
		CostPrice   float64 `json:"cost_price" binding:"required"`
		SalePrice   float64 `json:"sale_price" binding:"required"`
		TaxRate     float64 `json:"tax_rate"`
		MinStock    int     `json:"min_stock"`
		ImageURL    string  `json:"image_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	product := &Product{
		SKU:        req.SKU,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		BrandID:    req.BrandID,
		CostPrice:  req.CostPrice,
		SalePrice:  req.SalePrice,
		TaxRate:    req.TaxRate,
		MinStock:   req.MinStock,
		ImageURL:   req.ImageURL,
		EmpresaID:  empresaID,
	}

	if err := h.productSvc.Create(c.Request.Context(), product); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, product)
}

func (h *Handler) ListProducts(c *gin.Context) {
	empresaID := c.GetInt64("empresa_id")
	if empresaID == 0 {
		response.BadRequest(c, "empresa_id not found")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	var categoryID, brandID *int64
	if catIDStr := c.Query("category_id"); catIDStr != "" {
		if id, err := strconv.ParseInt(catIDStr, 10, 64); err == nil {
			categoryID = &id
		}
	}
	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		if id, err := strconv.ParseInt(brandIDStr, 10, 64); err == nil {
			brandID = &id
		}
	}

	filters := ListFilters{
		Page:       page,
		Limit:      limit,
		Search:     search,
		CategoryID: categoryID,
		BrandID:    brandID,
	}

	products, err := h.productSvc.List(c.Request.Context(), empresaID, filters)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, products)
}

func (h *Handler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	product, err := h.productSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Producto no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, product)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		SKU         string  `json:"sku"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		CategoryID  *int64  `json:"category_id"`
		BrandID     *int64  `json:"brand_id"`
		CostPrice   float64 `json:"cost_price"`
		SalePrice   float64 `json:"sale_price"`
		TaxRate     float64 `json:"tax_rate"`
		MinStock    int     `json:"min_stock"`
		ImageURL    string  `json:"image_url"`
		Status      string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	product := &Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		BrandID:     req.BrandID,
		CostPrice:   req.CostPrice,
		SalePrice:   req.SalePrice,
		TaxRate:     req.TaxRate,
		MinStock:    req.MinStock,
		ImageURL:    req.ImageURL,
		Status:      req.Status,
	}

	if err := h.productSvc.Update(c.Request.Context(), id, product); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Producto no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, product)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	if err := h.productSvc.Delete(c.Request.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Producto no encontrado")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	c.Status(204)
}
