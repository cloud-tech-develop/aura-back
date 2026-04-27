package products

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository implements the Repository interface for testing
// Provides mock implementations of all repository methods using testify
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, tenantSlug string, p *Product) error {
	args := m.Called(ctx, tenantSlug, p)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error) {
	args := m.Called(ctx, tenantSlug, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func (m *MockRepository) GetBySKU(ctx context.Context, tenantSlug string, sku string, enterpriseID int64) (*Product, error) {
	args := m.Called(ctx, tenantSlug, sku, enterpriseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func (m *MockRepository) GetByBarcode(ctx context.Context, tenantSlug string, barcode string, enterpriseID int64) (*Product, error) {
	args := m.Called(ctx, tenantSlug, barcode, enterpriseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func (m *MockRepository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	args := m.Called(ctx, tenantSlug, enterpriseID, page, limit, search, sort, order, params)
	return args.Get(0).(domain.PageResult), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, tenantSlug string, enterpriseID int64, filters ListFilters) ([]Product, error) {
	args := m.Called(ctx, tenantSlug, enterpriseID, filters)
	return args.Get(0).([]Product), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, tenantSlug string, p *Product) error {
	args := m.Called(ctx, tenantSlug, p)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	args := m.Called(ctx, tenantSlug, id)
	return args.Error(0)
}

// ─── Service Tests ───────────────────────────────────────────────────────────

// TestService_Create_ValidInput tests successful product creation with valid input
func TestService_Create_ValidInput(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		SKU:          "sk-u1",
		Name:         "Producto test",
		Barcode:      "1123255241",
		UnitID:       6,
		ProductType:  "ESTANDAR",
		CostPrice:    17000,
		SalePrice:    18558,
		EnterpriseID: 1,
	}

	// Expect repository GetBySKU (to check uniqueness) to return ErrNoRows (no existing product)
	mockRepo.On("GetBySKU", mock.Anything, "test_tenant", "sk-u1", int64(1)).Return(nil, sql.ErrNoRows).Once()
	// Expect repository GetByBarcode to return ErrNoRows
	mockRepo.On("GetByBarcode", mock.Anything, "test_tenant", "1123255241", int64(1)).Return(nil, sql.ErrNoRows).Once()
	// Expect repository Create to be called
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).Return(nil).Once()

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.NoError(t, err)
	assert.Equal(t, "sk-u1", product.SKU)
	assert.Equal(t, "Producto test", product.Name)
	assert.Equal(t, "ESTANDAR", product.ProductType)
	assert.True(t, product.Active) // Should be set to true by default
	mockRepo.AssertExpectations(t)
}

// TestService_Create_MissingSKU tests error when SKU is missing
func TestService_Create_MissingSKU(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		Name:         "Producto test",
		UnitID:       6,
		EnterpriseID: 1,
	}

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sku is required")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestService_Create_MissingName tests error when name is missing
func TestService_Create_MissingName(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		SKU:          "sk-u1",
		UnitID:       6,
		EnterpriseID: 1,
	}

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestService_Create_MissingUnitID tests error when unit measure is missing
func TestService_Create_MissingUnitID(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		SKU:          "sk-u1",
		Name:         "Producto test",
		EnterpriseID: 1,
	}

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unit_id is required")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestService_Create_InvalidProductType tests error when product type is invalid
func TestService_Create_InvalidProductType(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		SKU:          "sk-u1",
		Name:         "Producto test",
		ProductType:  "INVALIDO",
		UnitID:       6,
		EnterpriseID: 1,
	}

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid product type")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestService_Create_DefaultProductType tests that ESTANDAR is set as default when not specified
func TestService_Create_DefaultProductType(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		SKU:          "sk-u1",
		Name:         "Producto test",
		UnitID:       6,
		EnterpriseID: 1,
	}

	// Expect GetBySKU to return ErrNoRows (product doesn't exist)
	mockRepo.On("GetBySKU", mock.Anything, "test_tenant", "sk-u1", int64(1)).Return(nil, sql.ErrNoRows).Once()
	// Expect repository Create
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.MatchedBy(func(p *Product) bool {
		return p.ProductType == "ESTANDAR"
	})).Return(nil).Once()

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.NoError(t, err)
	assert.Equal(t, "ESTANDAR", product.ProductType)
	mockRepo.AssertExpectations(t)
}

// TestService_Create_DuplicateSKU tests error when SKU already exists
func TestService_Create_DuplicateSKU(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	existingProduct := &Product{
		ID:           1,
		SKU:          "sk-u1",
		EnterpriseID: 1,
	}

	product := &Product{
		SKU:          "sk-u1",
		Name:         "Producto duplicado",
		UnitID:       6,
		EnterpriseID: 1,
	}

	// Expect GetBySKU to return existing product
	mockRepo.On("GetBySKU", mock.Anything, "test_tenant", "sk-u1", int64(1)).Return(existingProduct, nil).Once()

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestService_Create_ValidBarcode tests barcode validation when provided
func TestService_Create_ValidBarcode(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		SKU:          "sk-u1",
		Name:         "Producto test",
		Barcode:      "1123255241",
		UnitID:       6,
		EnterpriseID: 1,
	}

	mockRepo.On("GetBySKU", mock.Anything, "test_tenant", "sk-u1", int64(1)).Return(nil, sql.ErrNoRows).Once()
	mockRepo.On("GetByBarcode", mock.Anything, "test_tenant", "1123255241", int64(1)).Return(nil, sql.ErrNoRows).Once()
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).Return(nil).Once()

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.NoError(t, err)
	assert.Equal(t, "1123255241", product.Barcode)
	mockRepo.AssertExpectations(t)
}

// TestService_Create_DuplicateBarcode tests error when barcode already exists
func TestService_Create_DuplicateBarcode(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	existingProduct := &Product{
		ID:           1,
		Barcode:      "1123255241",
		EnterpriseID: 1,
	}

	product := &Product{
		SKU:          "sk-u2",
		Name:         "Producto con barcode duplicado",
		Barcode:      "1123255241",
		UnitID:       6,
		EnterpriseID: 1,
	}

	// SKU check passes (no duplicate)
	mockRepo.On("GetBySKU", mock.Anything, "test_tenant", "sk-u2", int64(1)).Return(nil, sql.ErrNoRows).Once()
	// Barcode check returns existing product
	mockRepo.On("GetByBarcode", mock.Anything, "test_tenant", "1123255241", int64(1)).Return(existingProduct, nil).Once()

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestService_Create_ValidPricing tests validation of non-negative prices
func TestService_Create_ValidPricing(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		SKU:          "sk-u1",
		Name:         "Producto test",
		CostPrice:    17000,
		SalePrice:    18558,
		UnitID:       6,
		EnterpriseID: 1,
	}

	mockRepo.On("GetBySKU", mock.Anything, "test_tenant", "sk-u1", int64(1)).Return(nil, sql.ErrNoRows).Once()
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).Return(nil).Once()

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.NoError(t, err)
}

// TestService_Create_NegativeCostPrice tests error when cost price is negative
func TestService_Create_NegativeCostPrice(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		SKU:          "sk-u1",
		Name:         "Producto test",
		CostPrice:    -100,
		SalePrice:    100,
		UnitID:       6,
		EnterpriseID: 1,
	}

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cost_price cannot be negative")
}

// TestService_Create_NegativeSalePrice tests error when sale price is negative
func TestService_Create_NegativeSalePrice(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	product := &Product{
		SKU:          "sk-u1",
		Name:         "Producto test",
		CostPrice:    100,
		SalePrice:    -50,
		UnitID:       6,
		EnterpriseID: 1,
	}

	err := svc.Create(context.Background(), "test_tenant", product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sale_price cannot be negative")
}

// TestService_GetByID_Success tests successful product retrieval by ID
func TestService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	now := vo.DateTime(time.Now())
	expectedProduct := &Product{
		ID:           1,
		SKU:          "sk-u1",
		Name:         "Producto test",
		ProductType:  "ESTANDAR",
		Active:       true,
		EnterpriseID: 1,
		CreatedAt:    now,
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(expectedProduct, nil).Once()

	product, err := svc.GetByID(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), product.ID)
	assert.Equal(t, "sk-u1", product.SKU)
	assert.Equal(t, "Producto test", product.Name)
	assert.True(t, product.Active)
	mockRepo.AssertExpectations(t)
}

// TestService_GetByID_NotFound tests error when product doesn't exist
func TestService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(999)).Return(nil, sql.ErrNoRows).Once()

	product, err := svc.GetByID(context.Background(), "test_tenant", 999)

	assert.Nil(t, product)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	mockRepo.AssertExpectations(t)
}

// TestService_Update_NotFound tests error when updating non-existent product
func TestService_Update_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	updatedProduct := &Product{
		Name: "Producto actualizado",
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(999)).Return(nil, sql.ErrNoRows).Once()

	err := svc.Update(context.Background(), "test_tenant", 999, updatedProduct)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
}

// TestService_Update_InvalidProductType tests error with invalid product type
func TestService_Update_InvalidProductType(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	existingProduct := &Product{
		ID:           1,
		SKU:          "sk-u1",
		Name:         "Producto existente",
		ProductType:  "ESTANDAR",
		EnterpriseID: 1,
	}

	updatedProduct := &Product{
		Name:        "Producto actualizado",
		ProductType: "TIPO_INVALIDO",
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingProduct, nil).Once()

	err := svc.Update(context.Background(), "test_tenant", 1, updatedProduct)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid product type")
}

// TestService_Delete_Valid tests successful soft delete
func TestService_Delete_Valid(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	// First GetByID to verify product exists
	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(&Product{ID: 1}, nil).Once()
	// Then Delete
	mockRepo.On("Delete", mock.Anything, "test_tenant", int64(1)).Return(nil).Once()

	err := svc.Delete(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestService_Delete_NotFound tests error when deleting non-existent product
func TestService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(999)).Return(nil, sql.ErrNoRows).Once()

	err := svc.Delete(context.Background(), "test_tenant", 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
}

// TestService_Page_DefaultPagination tests default pagination values
func TestService_Page_DefaultPagination(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	expectedResult := domain.PageResult{
		Items:      []Product{{ID: 1, SKU: "sk-u1", Name: "Producto 1"}},
		Total:      1,
		Page:       1,
		Limit:      10,
		TotalPages: 1,
	}

	mockRepo.On("Page", mock.Anything, "test_tenant", int64(1), int64(1), int64(10), "", "id", "asc", mock.Anything).Return(expectedResult, nil).Once()

	result, err := svc.Page(context.Background(), "test_tenant", 1, 0, 0, "", "id", "asc", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Page)
	assert.Equal(t, int64(10), result.Limit)
	mockRepo.AssertExpectations(t)
}

// TestService_Page_Custom tests custom pagination
func TestService_Page_Custom(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	expectedResult := domain.PageResult{
		Items:      []Product{},
		Total:      0,
		Page:       3,
		Limit:      25,
		TotalPages: 0,
	}

	mockRepo.On("Page", mock.Anything, "test_tenant", int64(1), int64(3), int64(25), "searchterm", "name", "desc", mock.Anything).Return(expectedResult, nil).Once()

	result, err := svc.Page(context.Background(), "test_tenant", 1, 3, 25, "searchterm", "name", "desc", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), result.Page)
	assert.Equal(t, int64(25), result.Limit)
	mockRepo.AssertExpectations(t)
}

// TestService_List_EmptyList tests empty list result
func TestService_List_EmptyList(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	mockRepo.On("List", mock.Anything, "test_tenant", int64(1), mock.Anything).Return([]Product{}, nil).Once()

	result, err := svc.List(context.Background(), "test_tenant", 1, ListFilters{})

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

// ─── Domain Tests ─────────────────────────────────────────────────────────

// TestIsValidProductType_ValidTypes tests valid product type values
func TestIsValidProductType_ValidTypes(t *testing.T) {
	tests := []struct {
		name        string
		productType string
		want        bool
	}{
		{"ESTANDAR is valid", "ESTANDAR", true},
		{"SERVICIO is valid", "SERVICIO", true},
		{"COMBO is valid", "COMBO", true},
		{"RECETA is valid", "RECETA", true},
		{"LOWERCASE ESTANDAR", "estandar", false},
		{"empty string", "", false},
		{"invalid type", "INVALIDO", false},
		{"different case", "Estandar", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidProductType(tt.productType)
			assert.Equal(t, tt.want, got)
		})
	}
}
