package units

import (
	"context"
	"database/sql"
	"testing"

	"github.com/cloud-tech-develop/aura-back/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── Mocks ───────────────────────────────────────────────────────────────────────

// MockRepository implementa el interfaz Repository para tests
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, tenantSlug string, u *Unit) error {
	args := m.Called(ctx, tenantSlug, u)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Unit, error) {
	args := m.Called(ctx, tenantSlug, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Unit), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, tenantSlug string, enterpriseID int64) ([]domain.ListId, error) {
	args := m.Called(ctx, tenantSlug, enterpriseID)
	return args.Get(0).([]domain.ListId), args.Error(1)
}

func (m *MockRepository) Page(ctx context.Context, tenantSlug string, enterpriseID int64, page int64, limit int64, search string, sort string, order string, params map[string]any) (domain.PageResult, error) {
	args := m.Called(ctx, tenantSlug, enterpriseID, page, limit, search, sort, order, params)
	return args.Get(0).(domain.PageResult), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, tenantSlug string, u *Unit) error {
	args := m.Called(ctx, tenantSlug, u)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, tenantSlug string, id int64) error {
	args := m.Called(ctx, tenantSlug, id)
	return args.Error(0)
}

// ─── Tests del Service ─────────────────────────────────────────────────────────

func TestService_Create_ValidInput(t *testing.T) {
	// Test: Crear unidad con datos válidos
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	unit := &Unit{
		Name:          "Kilogramo",
		Abbreviation:  "kg",
		Active:        true,
		AllowDecimals: true,
		EnterpriseID:  1,
	}

	// Expect repository Create to be called
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.AnythingOfType("*units.Unit")).Return(nil).Once()

	err := svc.Create(context.Background(), "test_tenant", unit)

	assert.NoError(t, err)
	assert.Equal(t, "Kilogramo", unit.Name)
	assert.Equal(t, "kg", unit.Abbreviation)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_MissingName(t *testing.T) {
	// Test: Error cuando name está vacío
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	unit := &Unit{
		Name:          "",
		Abbreviation:  "kg",
		Active:        true,
		AllowDecimals: true,
		EnterpriseID:  1,
	}

	err := svc.Create(context.Background(), "test_tenant", unit)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestService_Create_MissingAbbreviation(t *testing.T) {
	// Test: Error cuando abbreviation está vacío
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	unit := &Unit{
		Name:          "Kilogramo",
		Abbreviation:  "",
		Active:        true,
		AllowDecimals: true,
		EnterpriseID:  1,
	}

	err := svc.Create(context.Background(), "test_tenant", unit)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "abbreviation is required")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestService_GetByID_Success(t *testing.T) {
	// Test: Obtener unidad por ID exitosamente
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	expectedUnit := &Unit{
		ID:            1,
		Name:          "Kilogramo",
		Abbreviation:  "kg",
		Active:        true,
		AllowDecimals: true,
		EnterpriseID:  1,
		CreatedAt:     "2024-01-01T00:00:00Z",
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(expectedUnit, nil).Once()

	unit, err := svc.GetByID(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), unit.ID)
	assert.Equal(t, "Kilogramo", unit.Name)
	assert.Equal(t, "kg", unit.Abbreviation)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_NotFound(t *testing.T) {
	// Test: Error cuando la unidad no existe
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(999)).Return(nil, sql.ErrNoRows).Once()

	unit, err := svc.GetByID(context.Background(), "test_tenant", 999)

	assert.Nil(t, unit)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	mockRepo.AssertExpectations(t)
}

func TestService_List_Empty(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	mockRepo.On("List", mock.Anything, "test_tenant", int64(1)).Return([]domain.ListId{}, nil).Once()

	units, err := svc.List(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Empty(t, units)
	mockRepo.AssertExpectations(t)
}

func TestService_List_WithData(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	units := []domain.ListId{
		{Id: 1, Name: "Kilogramo"},
		{Id: 2, Name: "Litro"},
	}

	mockRepo.On("List", mock.Anything, "test_tenant", int64(1)).Return(units, nil).Once()

	result, err := svc.List(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Kilogramo", result[0].Name)
	assert.Equal(t, "Litro", result[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestService_Update_Valid(t *testing.T) {
	// Test: Actualización exitosa
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	existingUnit := &Unit{
		ID:            1,
		Name:          "Kilogramo",
		Abbreviation:  "kg",
		Active:        true,
		AllowDecimals: true,
		EnterpriseID:  1,
		CreatedAt:     "2024-01-01T00:00:00Z",
	}

	updatedUnit := &Unit{
		Name:          "Kilogramo Actualizado",
		Abbreviation:  "kgg",
		Active:        false,
		AllowDecimals: false,
	}

	// First GetByID to get existing
	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingUnit, nil).Once()
	// Then Update
	mockRepo.On("Update", mock.Anything, "test_tenant", mock.AnythingOfType("*units.Unit")).Return(nil).Once()

	err := svc.Update(context.Background(), "test_tenant", 1, updatedUnit)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Update_NotFound(t *testing.T) {
	// Test: Error al actualizar cuando no existe
	mockRepo := new(MockRepository)
	svc := &service{repo: mockRepo}

	updatedUnit := &Unit{
		Name: "NoExiste",
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(999)).Return(nil, sql.ErrNoRows).Once()

	err := svc.Update(context.Background(), "test_tenant", 999, updatedUnit)

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

	expectedResult := domain.PageResult{
		Items:      []Unit{{ID: 1, Name: "kg", Abbreviation: "kg", Active: true, EnterpriseID: 1}},
		Total:      1,
		Page:       1,
		Limit:      10,
		TotalPages: 1,
	}

	// Page defaults to 1 and limit to 10 when values < 1 are passed
	// Use mock.Anything for the params map since service may convert nil to empty map
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
		Items:      []Unit{},
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
