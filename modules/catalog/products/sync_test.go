package products

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventBus implementa events.EventBus para testing
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(event events.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventName string, handler events.EventHandler) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventName string, handler events.EventHandler) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

func (m *MockEventBus) Start() error {
	return nil
}

func (m *MockEventBus) Stop() error {
	return nil
}

// ─── SyncProductPayload Tests ──────────────────────────────────────────────────

// TestProduct_ToSyncPayload_AllFields verifica que todos los campos del producto
// se mapean correctamente al payload de sincronización
func TestProduct_ToSyncPayload_AllFields(t *testing.T) {
	// Crear un producto con todos los campos
	catID := int64(10)
	brandID := int64(5)
	price3 := 25000.0

	product := &Product{
		ID:                 1,
		SKU:                "SKU-001",
		Barcode:            "123456789",
		Name:               "Producto Test",
		Description:        "Descripción del producto",
		CategoryID:         &catID,
		BrandID:            &brandID,
		UnitID:             1,
		ProductType:        "STANDARD",
		Active:             true,
		VisibleInPOS:       true,
		CostPrice:          10000.0,
		SalePrice:          15000.0,
		Price2:             14000.0,
		Price3:             &price3,
		IVAPercentage:      19.0,
		ConsumptionTax:     0.0,
		CurrentStock:       50,
		MinStock:           10,
		MaxStock:           100,
		ManagesInventory:   true,
		ManagesBatches:     false,
		ManagesSerial:      false,
		AllowNegativeStock: false,
		ImageURL:          "http://example.com/image.jpg",
		EnterpriseID:      1,
	}

	// Convertir a payload de sincronización
	payload := product.ToSyncPayload("create", "online")

	// Verificar que todos los campos se mapean correctamente
	assert.Equal(t, "", payload.TenantSlug, "TenantSlug debe estar vacío initially")
	assert.Equal(t, int64(1), payload.ProductID)
	assert.Equal(t, "SKU-001", payload.SKU)
	assert.Equal(t, "123456789", payload.Barcode)
	assert.Equal(t, "Producto Test", payload.Name)
	assert.Equal(t, "Descripción del producto", payload.Description)
	assert.Equal(t, &catID, payload.CategoryID)
	assert.Equal(t, &brandID, payload.BrandID)
	assert.Equal(t, int64(1), payload.UnitID)
	assert.Equal(t, "STANDARD", payload.ProductType)
	assert.True(t, payload.Active)
	assert.True(t, payload.VisibleInPOS)
	assert.Equal(t, 10000.0, payload.CostPrice)
	assert.Equal(t, 15000.0, payload.SalePrice)
	assert.Equal(t, 14000.0, payload.Price2)
	assert.Equal(t, &price3, payload.Price3)
	assert.Equal(t, 19.0, payload.IVAPercentage)
	assert.Equal(t, 0.0, payload.ConsumptionTax)
	assert.Equal(t, 50, payload.CurrentStock)
	assert.Equal(t, 10, payload.MinStock)
	assert.Equal(t, 100, payload.MaxStock)
	assert.True(t, payload.ManagesInventory)
	assert.False(t, payload.ManagesBatches)
	assert.False(t, payload.ManagesSerial)
	assert.False(t, payload.AllowNegativeStock)
	assert.Equal(t, "http://example.com/image.jpg", payload.ImageURL)
	assert.Equal(t, int64(1), payload.EnterpriseID)
	assert.Equal(t, "create", payload.Action)
	assert.Equal(t, "online", payload.Source)
}

// TestProduct_ToSyncPayloadWithTenant verifica que el tenant slug se establece
// correctamente cuando se usa ToSyncPayloadWithTenant
func TestProduct_ToSyncPayloadWithTenant(t *testing.T) {
	product := &Product{
		ID:           1,
		SKU:          "SKU-001",
		Name:         "Producto Test",
		EnterpriseID: 1,
	}

	payload := product.ToSyncPayloadWithTenant("test_tenant", "create", "offline")

	assert.Equal(t, "test_tenant", payload.TenantSlug)
	assert.Equal(t, "create", payload.Action)
	assert.Equal(t, "offline", payload.Source)
}

// TestProduct_ToSyncPayload_Timestamp se verifica que el timestamp se establece
// al momento de la conversión
func TestProduct_ToSyncPayload_Timestamp(t *testing.T) {
	before := time.Now()
	product := &Product{
		ID:           1,
		SKU:          "SKU-001",
		Name:         "Producto Test",
		EnterpriseID: 1,
	}
	after := time.Now()

	payload := product.ToSyncPayload("create", "online")

	// El timestamp debe estar entre before y after
	assert.True(t, payload.Timestamp.After(before) || payload.Timestamp.Equal(before))
	assert.True(t, payload.Timestamp.Before(after) || payload.Timestamp.Equal(after))
}

// TestSyncEventCreation verifica la creación de eventos de sincronización
func TestSyncEventCreation(t *testing.T) {
	product := &Product{
		ID:           1,
		SKU:          "SKU-001",
		Name:         "Producto Test",
		EnterpriseID: 1,
	}

	t.Run("NewSyncCreatedEvent", func(t *testing.T) {
		event := NewSyncCreatedEvent("test_tenant", product)
		assert.Equal(t, EventProductOfflineCreated, event.GetName())
		payload := event.GetPayload().(SyncProductPayload)
		assert.Equal(t, "test_tenant", payload.TenantSlug)
		assert.Equal(t, "create", payload.Action)
		assert.Equal(t, "online", payload.Source)
	})

	t.Run("NewSyncUpdatedEvent", func(t *testing.T) {
		event := NewSyncUpdatedEvent("test_tenant", product)
		assert.Equal(t, EventProductOfflineUpdated, event.GetName())
		payload := event.GetPayload().(SyncProductPayload)
		assert.Equal(t, "test_tenant", payload.TenantSlug)
		assert.Equal(t, "update", payload.Action)
	})

	t.Run("NewSyncDeletedEvent", func(t *testing.T) {
		event := NewSyncDeletedEvent("test_tenant", product)
		assert.Equal(t, EventProductOfflineDeleted, event.GetName())
		payload := event.GetPayload().(SyncProductPayload)
		assert.Equal(t, "delete", payload.Action)
	})

	t.Run("NewSyncCreatedEventFromOffline", func(t *testing.T) {
		event := NewSyncCreatedEventFromOffline("test_tenant", product)
		assert.Equal(t, EventProductOnlineCreated, event.GetName())
		payload := event.GetPayload().(SyncProductPayload)
		assert.Equal(t, "offline", payload.Source)
	})

	t.Run("NewSyncUpdatedEventFromOffline", func(t *testing.T) {
		event := NewSyncUpdatedEventFromOffline("test_tenant", product)
		assert.Equal(t, EventProductOnlineUpdated, event.GetName())
	})

	t.Run("NewSyncDeletedEventFromOffline", func(t *testing.T) {
		event := NewSyncDeletedEventFromOffline("test_tenant", product)
		assert.Equal(t, EventProductOnlineDeleted, event.GetName())
	})
}

// ─── SyncHandler Tests ─────────────────────────────────────────────────────────

// TestSyncHandler_HandleCreate verifica la creación de productos desde sync
func TestSyncHandler_HandleCreate(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// Preparar payload como mapa (como viene de RabbitMQ)
	payload := map[string]interface{}{
		"tenant_slug":        "test_tenant",
		"product_id":        float64(1),
		"sku":               "SKU-001",
		"barcode":           "123456789",
		"name":              "Producto Test",
		"description":       "Descripción",
		"product_type":      "STANDARD",
		"active":            true,
		"visible_in_pos":    true,
		"cost_price":        10000.0,
		"sale_price":        15000.0,
		"price_2":           14000.0,
		"iva_percentage":    19.0,
		"consumption_tax":   0.0,
		"manages_inventory": true,
		"manages_batches":   false,
		"manages_serial":    false,
		"allow_negative_stock": false,
		"image_url":        "http://example.com/img.jpg",
		"enterprise_id":    float64(1),
		"action":           "create",
	}

	// Expect repository Create to be called
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.MatchedBy(func(p *Product) bool {
		return p.SKU == "SKU-001" && p.Name == "Producto Test"
	})).Return(nil).Once()

	err := handler.handleCreate(context.Background(), "test_tenant", payload)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestSyncHandler_HandleCreate_WithOptionalFields verifica la creación con campos opcionales
func TestSyncHandler_HandleCreate_WithOptionalFields(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// El handler espera valores directos (float64, no *float64)
	// Para campos que pueden ser nil, no los incluimos en el payload
	payload := map[string]interface{}{
		"tenant_slug":        "test_tenant",
		"product_id":        float64(1),
		"sku":               "SKU-001",
		"barcode":           "123456789",
		"name":              "Producto Test",
		"description":       "Descripción",
		"product_type":      "STANDARD",
		"active":            true,
		"visible_in_pos":    true,
		"cost_price":        10000.0,
		"sale_price":        15000.0,
		"price_2":           14000.0,
		"price_3":           25000.0, // valor directo para price_3
		"iva_percentage":    19.0,
		"consumption_tax":   0.0,
		"manages_inventory": true,
		"manages_batches":   false,
		"manages_serial":    false,
		"allow_negative_stock": false,
		"image_url":        "http://example.com/img.jpg",
		"enterprise_id":    float64(1),
		"category_id":      float64(10),  // valor directo (handler convierte a int64)
		"brand_id":         float64(5),    // valor directo
		"unit_measure_id":  float64(1),
		"min_stock":        float64(5),    // valor directo (handler convierte a int)
		"max_stock":        float64(50),   // valor directo
		"current_stock":    float64(25),   // valor directo
		"action":           "create",
	}

	mockRepo.On("Create", mock.Anything, "test_tenant", mock.MatchedBy(func(p *Product) bool {
		return p.CategoryID != nil && *p.CategoryID == 10 &&
			p.BrandID != nil && *p.BrandID == 5 &&
			p.Price3 != nil && *p.Price3 == 25000.0 &&
			p.MinStock == 5 && p.MaxStock == 50 && p.CurrentStock == 25
	})).Return(nil).Once()

	err := handler.handleCreate(context.Background(), "test_tenant", payload)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestSyncHandler_HandleUpdate verifica la actualización de productos desde sync
func TestSyncHandler_HandleUpdate(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	payload := map[string]interface{}{
		"tenant_slug":        "test_tenant",
		"product_id":        float64(1),
		"sku":               "SKU-001-UPDATED",
		"barcode":           "123456789",
		"name":              "Producto Actualizado",
		"description":       "Nueva descripción",
		"product_type":      "STANDARD",
		"active":            true,
		"visible_in_pos":    true,
		"cost_price":        12000.0,
		"sale_price":        18000.0,
		"price_2":           16000.0,
		"iva_percentage":    19.0,
		"consumption_tax":   0.0,
		"manages_inventory": true,
		"manages_batches":   false,
		"manages_serial":    false,
		"allow_negative_stock": false,
		"image_url":        "http://example.com/new.jpg",
		"enterprise_id":    float64(1),
		"action":           "update",
	}

	mockRepo.On("Update", mock.Anything, "test_tenant", mock.MatchedBy(func(p *Product) bool {
		return p.Name == "Producto Actualizado" && p.SKU == "SKU-001-UPDATED"
	})).Return(nil).Once()

	err := handler.handleUpdate(context.Background(), "test_tenant", payload)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestSyncHandler_HandleDelete verifica la eliminación de productos desde sync
func TestSyncHandler_HandleDelete(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	payload := map[string]interface{}{
		"tenant_slug": "test_tenant",
		"product_id": float64(1),
		"action":     "delete",
	}

	mockRepo.On("Delete", mock.Anything, "test_tenant", int64(1)).Return(nil).Once()

	err := handler.handleDelete(context.Background(), "test_tenant", payload)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestSyncHandler_HandleProductSync_InvalidPayload verifica error con payload inválido
func TestSyncHandler_HandleProductSync_InvalidPayload(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// Payload inválido (no es map[string]interface{})
	event := events.NewBaseEvent("test.event", "invalid payload")

	err := handler.HandleProductSync(context.Background(), event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sync payload type")
}

// TestSyncHandler_HandleProductSync_MissingTenantSlug verifica error sin tenant_slug
func TestSyncHandler_HandleProductSync_MissingTenantSlug(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	payload := map[string]interface{}{
		"product_id": float64(1),
		"action":     "create",
	}
	event := events.NewBaseEvent("test.event", payload)

	err := handler.HandleProductSync(context.Background(), event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant_slug is required")
}

// TestSyncHandler_HandleProductSync_MissingAction verifica error sin acción
func TestSyncHandler_HandleProductSync_MissingAction(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	payload := map[string]interface{}{
		"tenant_slug": "test_tenant",
		"product_id": float64(1),
	}
	event := events.NewBaseEvent("test.event", payload)

	err := handler.HandleProductSync(context.Background(), event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "action is required")
}

// TestSyncHandler_HandleProductSync_UnknownAction verifica comportamiento con acción desconocida
func TestSyncHandler_HandleProductSync_UnknownAction(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	payload := map[string]interface{}{
		"tenant_slug": "test_tenant",
		"product_id": float64(1),
		"action":     "unknown_action",
	}
	event := events.NewBaseEvent("test.event", payload)

	err := handler.HandleProductSync(context.Background(), event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown sync action")
}

// ─── ProductSyncFromOffline Tests ─────────────────────────────────────────────

// TestSyncHandler_ProductSyncFromOffline_Create procesa productos nuevos desde offline
func TestSyncHandler_ProductSyncFromOffline_Create(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// El producto no existe, debe crearse
	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(nil, sql.ErrNoRows).Once()
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).Return(nil).Once()

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-OFFLINE-1",
			Name:        "Producto Offline 1",
			EnterpriseID: 1,
			Action:      "create",
			Timestamp:   time.Now(),
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	assert.Equal(t, "success", results[0].Status)
	assert.Equal(t, "create", results[0].Action)
	assert.Contains(t, results[0].Message, "created successfully")
	mockRepo.AssertExpectations(t)
}

// TestSyncHandler_ProductSyncFromOffline_Update procesa actualizaciones desde offline
func TestSyncHandler_ProductSyncFromOffline_Update(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// El producto existe pero con timestamp nil (significa que no hay update previo, offline gana)
	existingProduct := &Product{
		ID:        1,
		SKU:       "SKU-OFFLINE-1",
		Name:      "Producto old",
		UpdatedAt: nil, // nil significa que no hay update previo, offline gana
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingProduct, nil).Once()
	mockRepo.On("Update", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).Return(nil).Once()

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-OFFLINE-1",
			Name:        "Producto Offline Actualizado",
			EnterpriseID: 1,
			Action:      "update",
			Timestamp:   time.Now(),
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	assert.Equal(t, "success", results[0].Status)
	assert.Contains(t, results[0].Message, "updated successfully")
	mockRepo.AssertExpectations(t)
}

// TestSyncHandler_ProductSyncFromOffline_Delete procesa eliminaciones desde offline
func TestSyncHandler_ProductSyncFromOffline_Delete(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// El producto existe pero tiene timestamp nil (sync inicial)
	existingProduct := &Product{
		ID:        1,
		SKU:       "SKU-OFFLINE-1",
		UpdatedAt: nil,
	}

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingProduct, nil).Once()
	mockRepo.On("Delete", mock.Anything, "test_tenant", int64(1)).Return(nil).Once()

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-OFFLINE-1",
			Action:      "delete",
			Timestamp:   time.Now(),
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	assert.Equal(t, "success", results[0].Status)
	assert.Contains(t, results[0].Message, "deleted successfully")
	mockRepo.AssertExpectations(t)
}

// TestSyncHandler_ProductSyncFromOffline_MultipleProducts procesa múltiples productos
func TestSyncHandler_ProductSyncFromOffline_MultipleProducts(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// Producto 1: nuevo (crear)
	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(nil, sql.ErrNoRows).Once()
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).Return(nil).Once()

	// Producto 2: existe (actualizar) - timestamp nil = offline gana
	existingProduct := &Product{
		ID:        2,
		UpdatedAt: nil,
	}
	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(2)).Return(existingProduct, nil).Once()
	mockRepo.On("Update", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).Return(nil).Once()

	// Producto 3: no existe en online (crear como nuevo)
	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(3)).Return(nil, sql.ErrNoRows).Once()
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).Return(nil).Once()

	products := []SyncProductPayload{
		{ProductID: 1, SKU: "SKU-1", Action: "create", Timestamp: time.Now()},
		{ProductID: 2, SKU: "SKU-2", Action: "update", Timestamp: time.Now()},
		{ProductID: 3, SKU: "SKU-3", Action: "create", Timestamp: time.Now()},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 3)
	for _, r := range results {
		assert.Equal(t, "success", r.Status)
	}
}

// TestSyncHandler_ProductSyncFromOffline_ErrorOnCreate verifica manejo de errores en creación
func TestSyncHandler_ProductSyncFromOffline_ErrorOnCreate(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(nil, sql.ErrNoRows).Once()
	mockRepo.On("Create", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).
		Return(assert.AnError).Once()

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-ERROR",
			Name:        "Producto Error",
			EnterpriseID: 1,
			Action:      "create",
			Timestamp:   time.Now(),
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	assert.Equal(t, "error", results[0].Status)
	assert.Equal(t, assert.AnError.Error(), results[0].Message)
}

// ─── Conflict Resolution Tests ─────────────────────────────────────────────────

// TestSyncHandler_Conflict_OnlineWins_Create verifica que online gana cuando tiene
// versión más reciente en operación de creación
func TestSyncHandler_Conflict_OnlineWins_Create(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// El producto ya existe y online tiene timestamp más reciente
	onlineTime := time.Now()
	onlineUpdatedAt := vo.DateTime(onlineTime)

	existingProduct := &Product{
		ID:        1,
		SKU:       "SKU-ONLINE",
		Name:      "Producto Online",
		UpdatedAt: &onlineUpdatedAt,
	}

	// Offline viene con timestamp anterior
	offlineTimestamp := onlineTime.Add(-time.Hour)

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingProduct, nil).Once()
	// No debe llamar a Update porque online gana

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-OFFLINE",
			Name:        "Producto Offline",
			EnterpriseID: 1,
			Action:      "create",
			Timestamp:   offlineTimestamp,
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	// El conflicto se resuelve a favor de online, resolveConflict retorna nil
	// así que el resultado es success (aunque no hace update)
	assert.Equal(t, "success", results[0].Status)
	assert.Contains(t, results[0].Message, "created successfully")
	mockRepo.AssertNotCalled(t, "Update")
}

// TestSyncHandler_Conflict_OnlineWins_Update verifica que online gana en actualización
func TestSyncHandler_Conflict_OnlineWins_Update(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// Online tiene versión más reciente
	onlineTime := time.Now()
	onlineUpdatedAt := vo.DateTime(onlineTime)

	existingProduct := &Product{
		ID:        1,
		SKU:       "SKU-001",
		Name:      "Producto Online Actualizado",
		UpdatedAt: &onlineUpdatedAt,
	}

	// Offline viene con timestamp anterior
	offlineTimestamp := onlineTime.Add(-time.Hour)

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingProduct, nil).Once()

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-001",
			Name:        "Producto Offline Viejos",
			EnterpriseID: 1,
			Action:      "update",
			Timestamp:   offlineTimestamp,
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	assert.Equal(t, "success", results[0].Status)
	// No debe llamar a Update porque online gana el conflicto
	mockRepo.AssertNotCalled(t, "Update")
}

// TestSyncHandler_Conflict_OfflineWins verifica que offline gana cuando tiene
// versión más reciente
func TestSyncHandler_Conflict_OfflineWins(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// El producto existe pero con timestamp antiguo (nil = sin update previo)
	existingProduct := &Product{
		ID:        1,
		SKU:       "SKU-001",
		Name:      "Producto Viejo",
		UpdatedAt: nil,
	}

	// Offline viene con timestamp reciente
	offlineTimestamp := time.Now()

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingProduct, nil).Once()
	mockRepo.On("Update", mock.Anything, "test_tenant", mock.AnythingOfType("*products.Product")).Return(nil).Once()

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-001",
			Name:        "Producto Offline Reciente",
			EnterpriseID: 1,
			Action:      "update",
			Timestamp:   offlineTimestamp,
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	assert.Equal(t, "success", results[0].Status)
	assert.Contains(t, results[0].Message, "updated successfully")
	mockRepo.AssertExpectations(t)
}

// TestSyncHandler_Conflict_Delete_OnlineWins verifica que online gana en conflicto de eliminación
func TestSyncHandler_Conflict_Delete_OnlineWins(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// Online tiene versión más reciente
	onlineTime := time.Now()
	onlineUpdatedAt := vo.DateTime(onlineTime)

	existingProduct := &Product{
		ID:        1,
		SKU:       "SKU-001",
		UpdatedAt: &onlineUpdatedAt,
	}

	// Offline intenta eliminar pero con timestamp anterior
	offlineTimestamp := onlineTime.Add(-time.Hour)

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingProduct, nil).Once()
	// No debe llamar a Delete porque online tiene versión más reciente

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-001",
			Action:      "delete",
			Timestamp:   offlineTimestamp,
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	// El handler retorna error cuando online tiene versión más reciente
	assert.Equal(t, "error", results[0].Status)
	assert.Contains(t, results[0].Message, "conflict")
	mockRepo.AssertNotCalled(t, "Delete")
}

// TestSyncHandler_Conflict_Delete_OfflineWins verifica que offline puede eliminar
// cuando tiene versión más reciente
func TestSyncHandler_Conflict_Delete_OfflineWins(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// Producto existente con timestamp antiguo (nil)
	existingProduct := &Product{
		ID:        1,
		SKU:       "SKU-001",
		UpdatedAt: nil,
	}

	// Offline viene con timestamp reciente
	offlineTimestamp := time.Now()

	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(existingProduct, nil).Once()
	mockRepo.On("Delete", mock.Anything, "test_tenant", int64(1)).Return(nil).Once()

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-001",
			Action:      "delete",
			Timestamp:   offlineTimestamp,
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	assert.Equal(t, "success", results[0].Status)
	assert.Contains(t, results[0].Message, "deleted successfully")
	mockRepo.AssertExpectations(t)
}

// TestSyncHandler_Conflict_Delete_ProductNotFound verifica que no hay conflicto
// cuando el producto ya fue eliminado en online
func TestSyncHandler_Conflict_Delete_ProductNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// El producto no existe en online (ya fue eliminado)
	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(nil, sql.ErrNoRows).Once()
	// No debe llamar a Delete

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-001",
			Action:      "delete",
			Timestamp:   time.Now(),
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	// El handler retorna nil (sin error) cuando el producto no existe
	assert.Equal(t, "success", results[0].Status)
	assert.Contains(t, results[0].Message, "deleted successfully")
	mockRepo.AssertNotCalled(t, "Delete")
}

// TestSyncHandler_Conflict_UnknownAction verifica comportamiento con acción desconocida
func TestSyncHandler_Conflict_UnknownAction(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-001",
			Action:      "unknown_action",
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	assert.Equal(t, "error", results[0].Status)
	assert.Contains(t, results[0].Message, "Unknown action")
}

// TestSyncHandler_Conflict_GetByIDError verifica manejo de errores al obtener producto
func TestSyncHandler_Conflict_GetByIDError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockEventBus := new(MockEventBus)
	handler := NewSyncHandler(mockRepo, mockEventBus)

	// Error al obtener producto
	mockRepo.On("GetByID", mock.Anything, "test_tenant", int64(1)).Return(nil, assert.AnError).Once()

	products := []SyncProductPayload{
		{
			ProductID:   1,
			SKU:         "SKU-001",
			Action:      "update",
			Timestamp:   time.Now(),
		},
	}

	results := handler.ProductSyncFromOffline(context.Background(), "test_tenant", products)

	assert.Len(t, results, 1)
	assert.Equal(t, "error", results[0].Status)
}

// ─── Helper Functions Tests ───────────────────────────────────────────────────

// TestIsOnlineNewer verifica la lógica de comparación de timestamps
func TestIsOnlineNewer(t *testing.T) {
	now := time.Now()
	hourAgo := now.Add(-time.Hour)
	hourLater := now.Add(time.Hour)

	tests := []struct {
		name             string
		onlineTimestamp  *vo.DateTime
		offlineTimestamp time.Time
		expected         bool
	}{
		{
			name:             "nil online timestamp returns false",
			onlineTimestamp:  nil,
			offlineTimestamp: now,
			expected:         false,
		},
		{
			name:             "online timestamp vacío es más antiguo que offline reciente",
			onlineTimestamp:  &vo.DateTime{},
			offlineTimestamp: now,
			expected:         false,
		},
		{
			name:             "online con timestamp reciente es más reciente que offline antiguo",
			onlineTimestamp:  func() *vo.DateTime { t := vo.DateTime(now); return &t }(),
			offlineTimestamp: hourAgo,
			expected:         true,
		},
		{
			name:             "offline es más reciente que online",
			onlineTimestamp:  func() *vo.DateTime { t := vo.DateTime(now); return &t }(),
			offlineTimestamp: hourLater,
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOnlineNewer(tt.onlineTimestamp, tt.offlineTimestamp)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSyncResult_Structure verifica la estructura de SyncResult
func TestSyncResult_Structure(t *testing.T) {
	result := SyncResult{
		ProductID: 1,
		SKU:       "SKU-001",
		Action:    "created",
		Status:   "success",
		Message:  "Product created successfully",
	}

	assert.Equal(t, int64(1), result.ProductID)
	assert.Equal(t, "SKU-001", result.SKU)
	assert.Equal(t, "created", result.Action)
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, "Product created successfully", result.Message)
}

// TestSyncConstants verifica que las constantes de eventos están definidas correctamente
func TestSyncConstants(t *testing.T) {
	assert.Equal(t, "product.offline.created", EventProductOfflineCreated)
	assert.Equal(t, "product.offline.updated", EventProductOfflineUpdated)
	assert.Equal(t, "product.offline.deleted", EventProductOfflineDeleted)
	assert.Equal(t, "product.online.created", EventProductOnlineCreated)
	assert.Equal(t, "product.online.updated", EventProductOnlineUpdated)
	assert.Equal(t, "product.online.deleted", EventProductOnlineDeleted)
}