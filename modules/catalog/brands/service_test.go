package brands

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── Mocks ─────────────────────────────────────────────────────────────────────

// MockRepository implementa el interfaz Repository para tests
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, tenantSlug string, b *Brand) error {
	args := m.Called(ctx, tenantSlug, b)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Brand, error) {
	args := m.Called(ctx, tenantSlug, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Brand), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]BrandList, error) {
	args := m.Called(ctx, tenantSlug, enterpriseID)
	return args.Get(0).([]BrandList), args.Error(1)
}

func (m *MockRepository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	args := m.Called(ctx, tenantSlug, enterpriseID, page, limit, search, sort, order, params)
	return args.Get(0).(domain.PageResult), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, tenantSlug string, b *Brand) error {
	args := m.Called(ctx, tenantSlug, b)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	args := m.Called(ctx, tenantSlug, id)
	return args.Error(0)
}

// ─── Tests del Service ─────────────────────────────────────────────────────────

func TestService_Create_ValidInput(t *testing.T) {
	// Test: Crear marca con datos válidos
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	brand := &Brand{
		Name:         "Marca Ejemplo",
		Description:  "Descripción de marca",
		Active:       true,
		EnterpriseID: 1,
	}

	// Expect repository Create to be called
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.AnythingOfType("*brands.Brand")).Return(nil).Once()

	err := svc.Create(context.Background(), "test_tenant", brand)

	assert.NoError(t, err)
	assert.Equal(t, "Marca Ejemplo", brand.Name)
	assert.Equal(t, "Descripción de marca", brand.Description)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_MissingName(t *testing.T) {
	// Test: Error cuando name está vacío
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	brand := &Brand{
		Name:         "",
		Description:  "Descripción",
		Active:       true,
		EnterpriseID: 1,
	}

	err := svc.Create(context.Background(), "test_tenant", brand)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestService_Create_DefaultActive(t *testing.T) {
	// Test: Verifica que active sea true por defecto al crear
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	brand := &Brand{
		Name:         "Marca Test",
		Description:  "Descripción test",
		Active:       true,
		EnterpriseID: 1,
	}

	mockRepo.On("Create", mock.Anything, "test_tenant", mock.MatchedBy(func(b *Brand) bool {
		return b.Active == true
	})).Return(nil).Once()

	err := svc.Create(context.Background(), "test_tenant", brand)

	assert.NoError(t, err)
	assert.True(t, brand.Active)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_CustomActive(t *testing.T) {
	// Test: Verifica que active sea false cuando se pasa false
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	brand := &Brand{
		Name:         "Marca Inactiva",
		Description:  "Marca desactivada intencionalmente",
		Active:       false,
		EnterpriseID: 1,
	}

	mockRepo.On("Create", mock.Anything, "test_tenant", mock.MatchedBy(func(b *Brand) bool {
		return b.Active == false
	})).Return(nil).Once()

	err := svc.Create(context.Background(), "test_tenant", brand)

	assert.NoError(t, err)
	assert.False(t, brand.Active)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_Success(t *testing.T) {
	// Test: Obtener marca por ID exitosamente
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	now := time.Now()
	expectedBrand := &Brand{
		ID:           1,
		Name:         "Marca Test",
		Description:  "Descripción test",
		Active:       true,
		EnterpriseID: 1,
		CreatedAt:    now,
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(expectedBrand, nil).Once()

	brand, err := svc.GetByID(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), brand.ID)
	assert.Equal(t, "Marca Test", brand.Name)
	assert.True(t, brand.Active)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_NotFound(t *testing.T) {
	// Test: Error cuando la marca no existe
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(999)).Return(nil, sql.ErrNoRows).Once()

	brand, err := svc.GetByID(context.Background(), "test_tenant", 999)

	assert.Nil(t, brand)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	mockRepo.AssertExpectations(t)
}

func TestService_List_Empty(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	mockRepo.On("List", mock.Anything, "test_tenant", int64(1)).Return([]BrandList{}, nil).Once()

	brands, err := svc.List(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Empty(t, brands)
	mockRepo.AssertExpectations(t)
}

func TestService_List_OnlyActive(t *testing.T) {
	// Test: Verifica que List solo retorna marcas activas
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	brands := []BrandList{
		{ID: 1, Name: "Marca Activa"},
	}

	mockRepo.On("List", mock.Anything, "test_tenant", int64(1)).Return(brands, nil).Once()

	result, err := svc.List(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	for _, b := range result {
		assert.Equal(t, "Marca Activa", b.Name)
	}
	mockRepo.AssertExpectations(t)
}

func TestService_Update_Valid(t *testing.T) {
	// Test: Actualización exitosa
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	existingBrand := &Brand{
		ID:           1,
		Name:         "Marca Original",
		Description:  "Descripción original",
		Active:       true,
		EnterpriseID: 1,
		CreatedAt:    time.Now(),
	}

	updatedBrand := &Brand{
		Name:        "Marca Actualizada",
		Description: "Descripción actualizada",
		Active:      false,
	}

	// First GetByID to get existing
	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingBrand, nil).Once()
	// Then Update
	mockRepo.On("Update", mock.Anything, "test_tenant", mock.AnythingOfType("*brands.Brand")).Return(nil).Once()

	err := svc.Update(context.Background(), "test_tenant", 1, updatedBrand)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Update_Active(t *testing.T) {
	// Test: Verifica que se actualice el campo active correctamente
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	existingBrand := &Brand{
		ID:           1,
		Name:         "Marca Activa",
		Description:  "Era activa",
		Active:       true,
		EnterpriseID: 1,
		CreatedAt:    time.Now(),
	}

	// Actualizar para desactivar la marca
	updatedBrand := &Brand{
		Name:        "Marca Activa",
		Description: "Ahora inactive",
		Active:      false,
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingBrand, nil).Once()
	mockRepo.On("Update", mock.Anything, "test_tenant", mock.MatchedBy(func(b *Brand) bool {
		return b.Active == false
	})).Return(nil).Once()

	err := svc.Update(context.Background(), "test_tenant", 1, updatedBrand)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Update_NotFound(t *testing.T) {
	// Test: Error al actualizar cuando no existe
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	updatedBrand := &Brand{
		Name: "NoExiste",
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(999)).Return(nil, sql.ErrNoRows).Once()

	err := svc.Update(context.Background(), "test_tenant", 999, updatedBrand)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_Valid(t *testing.T) {
	// Test: Eliminación lógica exitosa
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	mockRepo.On("Delete", mock.Anything, "test_tenant", int64(1)).Return(nil).Once()

	err := svc.Delete(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Page_DefaultPagination(t *testing.T) {
	// Test: Paginación con valores por defecto
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	now := time.Now()
	expectedResult := domain.PageResult{
		Items:      []Brand{{ID: 1, Name: "Marca Test", Description: "Desc", Active: true, EnterpriseID: 1, CreatedAt: now}},
		Total:      1,
		Page:       1,
		Limit:      10,
		TotalPages: 1,
	}

	// Page defaults to 1 and limit to 10 when values < 1 are passed
	mockRepo.On("Page", mock.Anything, "test_tenant", int64(1), int64(1), int64(10), "", "id", "asc", mock.Anything).Return(expectedResult, nil).Once()

	result, err := svc.Page(context.Background(), "test_tenant", 1, 0, 0, "", "id", "asc", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Page)
	assert.Equal(t, int64(10), result.Limit)
	mockRepo.AssertExpectations(t)
}

func TestService_Page_CustomPagination(t *testing.T) {
	// Test: Paginación con valores personalizados
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	expectedResult := domain.PageResult{
		Items:      []Brand{},
		Total:      0,
		Page:       2,
		Limit:      20,
		TotalPages: 0,
	}

	mockRepo.On("Page", mock.Anything, "test_tenant", int64(1), int64(2), int64(20), "search", "name", "desc", mock.Anything).Return(expectedResult, nil).Once()

	result, err := svc.Page(context.Background(), "test_tenant", 1, 2, 20, "search", "name", "desc", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.Page)
	assert.Equal(t, int64(20), result.Limit)
	mockRepo.AssertExpectations(t)
}
