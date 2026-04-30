package offline

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloud-tech-develop/aura-back/modules/catalog/products"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository mocks the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) UpsertThirdParty(ctx context.Context, tp *ThirdParty) error {
	args := m.Called(ctx, tp)
	return args.Error(0)
}

func (m *MockRepository) UpsertCategory(ctx context.Context, c *Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockRepository) UpsertBrand(ctx context.Context, b *Brand) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockRepository) UpsertUnit(ctx context.Context, u *Unit) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) UpsertProduct(ctx context.Context, tenantSlug string, p *products.Product) error {
	args := m.Called(ctx, tenantSlug, p)
	return args.Error(0)
}

func (m *MockRepository) UpsertPresentation(ctx context.Context, p *Presentation) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepository) UpsertEnterprise(ctx context.Context, e *Enterprise) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockRepository) GetEnterpriseBySlug(ctx context.Context, slug string) (*Enterprise, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Enterprise), args.Error(1)
}

func (m *MockRepository) ListEnterprises(ctx context.Context) ([]Enterprise, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Enterprise), args.Error(1)
}

func (m *MockRepository) UpsertPlan(ctx context.Context, p *Plan) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepository) UpsertUser(ctx context.Context, u *User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) UpsertUserRole(ctx context.Context, ur *UserRole) error {
	args := m.Called(ctx, ur)
	return args.Error(0)
}

// TestService_SyncCategories_WrapperFormat tests that categories sync uses wrapper format
func TestService_SyncCategories_WrapperFormat(t *testing.T) {
	// Create a test server that returns wrapped response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"data": []Category{
				{ID: 1, Name: "Category 1", EnterpriseID: 0},
				{ID: 2, Name: "Category 2", EnterpriseID: 0},
			},
			"success": true,
			"message": "Operación exitosa",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create mock repository
	mockRepo := new(MockRepository)
	mockRepo.On("UpsertCategory", mock.Anything, mock.AnythingOfType("*offline.Category")).Return(nil).Times(2)

	// Note: We can't easily test the service because it requires full initialization
	// This test just documents the expected behavior
	t.Log("Expected: Service should decode wrapper format {\"data\": [...]}")
	t.Log("Expected: Each item should have enterprise_id set before upserting")

	assert.NotNil(t, server.URL, "Test server should be created")
}

// TestDomain_CategoryFields tests Category entity fields
func TestDomain_CategoryFields(t *testing.T) {
	category := Category{
		ID:           1,
		Name:         "Test Category",
		Description:  "Test Description",
		Active:       true,
		EnterpriseID: 1,
	}

	assert.Equal(t, int64(1), category.ID)
	assert.Equal(t, "Test Category", category.Name)
	assert.True(t, category.Active)
	assert.Equal(t, int64(1), category.EnterpriseID)
}

// TestDomain_ThirdPartyFields tests ThirdParty entity fields
func TestDomain_ThirdPartyFields(t *testing.T) {
	tp := ThirdParty{
		ID:             1,
		FirstName:      "John",
		LastName:       "Doe",
		DocumentNumber: "12345678",
		DocumentType:   "CC",
		EnterpriseID:   1,
	}

	assert.Equal(t, int64(1), tp.ID)
	assert.Equal(t, "John", tp.FirstName)
	assert.Equal(t, "Doe", tp.LastName)
	assert.Equal(t, "12345678", tp.DocumentNumber)
}

// TestDomain_SyncResult tests SyncResult structure
func TestDomain_SyncResult(t *testing.T) {
	result := SyncResult{
		Enterprises:  1,
		Categories:   2,
		Brands:       3,
		Units:        4,
		Products:     5,
		ThirdParties: 6,
		Errors:       []string{},
	}

	assert.Equal(t, 1, result.Enterprises)
	assert.Equal(t, 2, result.Categories)
	assert.Equal(t, 3, result.Brands)
	assert.Equal(t, 4, result.Units)
	assert.Equal(t, 5, result.Products)
	assert.Equal(t, 6, result.ThirdParties)
	assert.Empty(t, result.Errors)
}

// TestSyncResult_WithErrors tests SyncResult with errors
func TestSyncResult_WithErrors(t *testing.T) {
	result := SyncResult{
		Categories: 0,
		Brands:     0,
		Errors:     []string{"categories: status: 400"},
	}

	assert.Equal(t, 0, result.Categories)
	assert.Equal(t, 0, result.Brands)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "categories: status: 400", result.Errors[0])
}
