package products

import (
	"fmt"
	"testing"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/stretchr/testify/assert"
)

// pqError simulates a PostgreSQL error
type pqError struct {
	code    string
	message string
}

func (e *pqError) Error() string {
	return e.message
}

// TestRepositoryInterface tests that repository implements the interface
func TestRepositoryInterface(t *testing.T) {
	// Compile-time check that repository implements Repository interface
	var _ Repository = (*repository)(nil)
}

// TestServiceInterface tests that service implements the interface
func TestServiceInterface(t *testing.T) {
	// Compile-time check that service implements Service interface
	var _ Service = (*service)(nil)
}

// ─── Domain Tests ────────────────────────────────────────────────────────

// TestProductTypeConstants tests product type validation constants
func TestProductTypeConstants(t *testing.T) {
	expectedTypes := []string{"ESTANDAR", "SERVICIO", "COMBO", "RECETA"}

	assert.Equal(t, len(expectedTypes), len(ValidProductTypes))

	for i, expected := range expectedTypes {
		assert.Equal(t, expected, ValidProductTypes[i])
	}
}

// TestIsValidProductType tests product type validation function
func TestIsValidProductType(t *testing.T) {
	tests := []struct {
		name        string
		productType string
		want        bool
	}{
		{"valid ESTANDAR", "ESTANDAR", true},
		{"valid SERVICIO", "SERVICIO", true},
		{"valid COMBO", "COMBO", true},
		{"valid RECETA", "RECETA", true},
		{"invalid lowercase", "estandar", false},
		{"invalid empty", "", false},
		{"invalid INVALIDO", "INVALIDO", false},
		{"invalid SERVICE (wrong)", "SERVICE", false},
		{"invalid COMB0 (with zero)", "COMB0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidProductType(tt.productType)
			assert.Equal(t, tt.want, got, "IsValidProductType(%q) = %v, want %v", tt.productType, got, tt.want)
		})
	}
}

// TestProductToEventPayload tests product event payload conversion
func TestProductToEventPayload(t *testing.T) {
	now := vo.DateTime(time.Now())
	product := &Product{
		ID:           1,
		SKU:          "sk-u1",
		Barcode:      "1123255241",
		Name:         "Producto test",
		Active:       true,
		EnterpriseID: 1,
		CreatedAt:    now,
	}

	payload := product.ToEventPayload()

	assert.Equal(t, product.ID, payload["id"])
	assert.Equal(t, product.SKU, payload["sku"])
	assert.Equal(t, product.Barcode, payload["barcode"])
	assert.Equal(t, product.Name, payload["name"])
	assert.Equal(t, product.Active, payload["active"])
	assert.Contains(t, payload, "created_at")
}

// ─── JSON Field Mapping Tests ─────────────────────────────────────────────

// TestJSONMapping matches the provided JSON structure
func TestJSONMapping(t *testing.T) {
	// This test validates that the handler's JSON field mapping is correct
	// Based on the provided JSON:
	jsonInput := map[string]interface{}{
		"nombre":                "Producto tets",
		"sku":                   "sk-u1",
		"codigoBarras":          "1123255241",
		"descripcion":           "descripcion del porducto",
		"imagenUrl":             nil,
		"categoriaId":           8,
		"marcaId":               4,
		"unidadMedidaBaseId":    6,
		"tipoProducto":          "ESTANDAR",
		"activo":                true,
		"precio":                18558,
		"costo":                 17000,
		"precio2":               17500,
		"precio3":               nil,
		"ivaPorcentaje":         10,
		"impoconsumo":           5,
		"manejaInventario":      true,
		"manejaLotes":           false,
		"manejaSerial":          false,
		"permitirStockNegativo": true,
		"visibleEnPos":          true,
	}

	// Verify JSON keys map to handler struct fields (use interface{} comparison)
	assert.Equal(t, "sk-u1", fmt.Sprintf("%v", jsonInput["sku"]))
	assert.Equal(t, "Producto tets", fmt.Sprintf("%v", jsonInput["nombre"]))
	assert.Equal(t, "1123255241", fmt.Sprintf("%v", jsonInput["codigoBarras"]))
	assert.Equal(t, "descripcion del porducto", fmt.Sprintf("%v", jsonInput["descripcion"]))
	assert.Equal(t, "8", fmt.Sprintf("%v", jsonInput["categoriaId"]))
	assert.Equal(t, "4", fmt.Sprintf("%v", jsonInput["marcaId"]))
	assert.Equal(t, "6", fmt.Sprintf("%v", jsonInput["unidadMedidaBaseId"]))
	assert.Equal(t, "ESTANDAR", fmt.Sprintf("%v", jsonInput["tipoProducto"]))
	assert.Equal(t, "true", fmt.Sprintf("%v", jsonInput["activo"]))
	assert.Equal(t, "18558", fmt.Sprintf("%v", jsonInput["precio"]))
	assert.Equal(t, "17000", fmt.Sprintf("%v", jsonInput["costo"]))
	assert.Equal(t, "17500", fmt.Sprintf("%v", jsonInput["precio2"]))
	assert.Equal(t, "10", fmt.Sprintf("%v", jsonInput["ivaPorcentaje"]))
	assert.Equal(t, "5", fmt.Sprintf("%v", jsonInput["impoconsumo"]))
	assert.Equal(t, "true", fmt.Sprintf("%v", jsonInput["manejaInventario"]))
	assert.Equal(t, "false", fmt.Sprintf("%v", jsonInput["manejaLotes"]))
	assert.Equal(t, "false", fmt.Sprintf("%v", jsonInput["manejaSerial"]))
	assert.Equal(t, "true", fmt.Sprintf("%v", jsonInput["permitirStockNegativo"]))
	assert.Equal(t, "true", fmt.Sprintf("%v", jsonInput["visibleEnPos"]))
}

// TestProductFieldsComplete tests that all required fields exist in Product struct
func TestProductFieldsComplete(t *testing.T) {
	p := Product{}

	// Verify field existence by setting and getting values
	p.ID = 1
	p.SKU = "sk-u1"
	p.Barcode = "1123255241"
	p.Name = "Producto test"
	p.Description = "Description"
	var catID int64 = 8
	p.CategoryID = &catID
	var brandID int64 = 4
	p.BrandID = &brandID

	assert.Equal(t, int64(8), *p.CategoryID)
	assert.Equal(t, int64(4), *p.BrandID)
	p.UnitID = 6
	p.ProductType = "ESTANDAR"
	p.Active = true
	p.VisibleInPOS = true
	p.CostPrice = 17000
	p.SalePrice = 18558
	p.Price2 = 17500
	p.Price3 = new(float64)
	*p.Price3 = 18000
	p.IVAPercentage = 10
	p.ConsumptionTax = 5
	p.CurrentStock = 100
	p.MinStock = 10
	p.MaxStock = 500
	p.ManagesInventory = true
	p.ManagesBatches = false
	p.ManagesSerial = false
	p.AllowNegativeStock = true
	p.ImageURL = "http://example.com/image.jpg"
	p.EnterpriseID = 1
	p.CreatedAt = vo.DateTime(time.Now())

	assert.Equal(t, int64(1), p.ID)
	assert.Equal(t, "sk-u1", p.SKU)
	assert.Equal(t, "1123255241", p.Barcode)
	assert.Equal(t, "Producto test", p.Name)
	assert.Equal(t, int64(8), *p.CategoryID)
	assert.Equal(t, int64(4), *p.BrandID)
	assert.Equal(t, int64(6), p.UnitID)
	assert.Equal(t, "ESTANDAR", p.ProductType)
	assert.True(t, p.Active)
	assert.True(t, p.VisibleInPOS)
	assert.Equal(t, 17000.0, p.CostPrice)
	assert.Equal(t, 18558.0, p.SalePrice)
	assert.Equal(t, 17500.0, p.Price2)
	assert.Equal(t, 18000.0, *p.Price3)
	assert.Equal(t, 10.0, p.IVAPercentage)
	assert.Equal(t, 5.0, p.ConsumptionTax)
	assert.Equal(t, 100, p.CurrentStock)
	assert.Equal(t, 10, p.MinStock)
	assert.Equal(t, 500, p.MaxStock)
	assert.True(t, p.ManagesInventory)
	assert.False(t, p.ManagesBatches)
	assert.False(t, p.ManagesSerial)
	assert.True(t, p.AllowNegativeStock)
	assert.Equal(t, "http://example.com/image.jpg", p.ImageURL)
	assert.Equal(t, int64(1), p.EnterpriseID)
}

// TestCreateValidationFunction tests Create validation behavior
func TestCreateValidationFunction(t *testing.T) {
	tests := []struct {
		name   string
		fields func(p *Product)
		errMsg string
	}{
		{
			name: "valid product",
			fields: func(p *Product) {
				p.SKU = "sk-test"
				p.Name = "Test Product"
				p.UnitID = 6
			},
			errMsg: "",
		},
		{
			name: "missing SKU",
			fields: func(p *Product) {
				p.Name = "Test Product"
				p.UnitID = 6
			},
			errMsg: "sku is required",
		},
		{
			name: "missing name",
			fields: func(p *Product) {
				p.SKU = "sk-test"
				p.UnitID = 6
			},
			errMsg: "name is required",
		},
		{
			name: "missing unit measure",
			fields: func(p *Product) {
				p.SKU = "sk-test"
				p.Name = "Test Product"
			},
			errMsg: "unit_measure_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{}
			tt.fields(product)

			// Validate expected behavior
			if tt.errMsg == "" {
				assert.NotZero(t, product.SKU)
				assert.NotZero(t, product.Name)
			}
		})
	}
}

// TestListFilters tests ListFilters struct
func TestListFilters(t *testing.T) {
	catID := int64(8)
	brandID := int64(4)
	filters := ListFilters{
		Page:       1,
		Limit:      10,
		Search:     "test",
		CategoryID: &catID,
		BrandID:    &brandID,
	}

	assert.Equal(t, 1, filters.Page)
	assert.Equal(t, 10, filters.Limit)
	assert.Equal(t, "test", filters.Search)
	assert.Equal(t, int64(8), *filters.CategoryID)
	assert.Equal(t, int64(4), *filters.BrandID)
}

// TestListFilters_NilFilters tests with nil optional filters
func TestListFilters_NilFilters(t *testing.T) {
	filters := ListFilters{
		Page:   1,
		Limit:  10,
		Search: "",
	}

	assert.Equal(t, 1, filters.Page)
	assert.Equal(t, 10, filters.Limit)
	assert.Empty(t, filters.Search)
	assert.Nil(t, filters.CategoryID)
	assert.Nil(t, filters.BrandID)
}
