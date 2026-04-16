package brands

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// pqError simula un error de PostgreSQL
type pqError struct {
	code    string
	message string
}

func (e *pqError) Error() string {
	return e.message
}

func TestRepository_Create_WithActive(t *testing.T) {
	// Test: Verifica que se guarde el campo active correctamente
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	brand := &Brand{
		Name:         "Nueva Marca",
		Description:  "Descripción de marca",
		Active:       true,
		EnterpriseID: 1,
	}

	// Expect INSERT con el campo active - args: name, description, active, enterprise_id
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, createdAt)
	mock.ExpectQuery(`INSERT INTO "test_tenant".brand`).
		WithArgs("Nueva Marca", "Descripción de marca", true, int64(1)).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), "test_tenant", brand)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), brand.ID)
	assert.Equal(t, createdAt, brand.CreatedAt)
	assert.True(t, brand.Active)
}

func TestRepository_Create_WithInactive(t *testing.T) {
	// Test: Verifica que se guarde active=false correctamente
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	brand := &Brand{
		Name:         "Marca Inactiva",
		Description:  "Descripción",
		Active:       false,
		EnterpriseID: 1,
	}

	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(2, createdAt)
	mock.ExpectQuery(`INSERT INTO "test_tenant".brand`).
		WithArgs("Marca Inactiva", "Descripción", false, int64(1)).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), "test_tenant", brand)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), brand.ID)
	assert.False(t, brand.Active)
}

func TestRepository_Create_Success(t *testing.T) {
	// Test: Inserción exitosa de marca
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	brand := &Brand{
		Name:         "Marca Ejemplo",
		Description:  "Descripción de ejemplo",
		Active:       true,
		EnterpriseID: 1,
	}

	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, createdAt)
	mock.ExpectQuery(`INSERT INTO "test_tenant".brand`).
		WithArgs("Marca Ejemplo", "Descripción de ejemplo", true, int64(1)).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), "test_tenant", brand)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), brand.ID)
	assert.Equal(t, createdAt, brand.CreatedAt)
}

func TestRepository_GetByID_Success(t *testing.T) {
	// Test: Obtención por ID exitosa
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "active", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Marca Test", "Descripción test", true, 1, createdAt, nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE id`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	brand, err := repo.GetByID(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), brand.ID)
	assert.Equal(t, "Marca Test", brand.Name)
	assert.Equal(t, "Descripción test", brand.Description)
	assert.True(t, brand.Active)
	assert.Equal(t, createdAt, brand.CreatedAt)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	// Test: No encontrado
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE id`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	brand, err := repo.GetByID(context.Background(), "test_tenant", 999)

	assert.Nil(t, brand)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestRepository_List_FilterActive(t *testing.T) {
	// Test: Verifica que solo liste marcas activas (active = true)
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// El repositorio filtra por active = true
	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "active", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Marca Activa 1", "Desc 1", true, 1, createdAt, nil, nil).
		AddRow(2, "Marca Activa 2", "Desc 2", true, 1, createdAt, nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE enterprise_id`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	brands, err := repo.List(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Len(t, brands, 2)
	for _, b := range brands {
		assert.True(t, b.Active, "Solo debe listar marcas activas")
	}
}

func TestRepository_List_Success(t *testing.T) {
	// Test: Listado exitoso
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "active", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Marca A", "Desc A", true, 1, createdAt, nil, nil).
		AddRow(2, "Marca B", "Desc B", true, 1, createdAt, nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE enterprise_id`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	brands, err := repo.List(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Len(t, brands, 2)
	assert.Equal(t, "Marca A", brands[0].Name)
	assert.Equal(t, "Marca B", brands[1].Name)
}

func TestRepository_List_Empty(t *testing.T) {
	// Test: Lista vacía
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "active", "enterprise_id", "created_at", "updated_at", "deleted_at",
	})

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE enterprise_id`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	brands, err := repo.List(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Empty(t, brands)
}

func TestRepository_Update_Active(t *testing.T) {
	// Test: Verifica que se actualice el campo active correctamente
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	brand := &Brand{
		ID:          1,
		Name:        "Marca Actualizada",
		Description: "Nueva descripción",
		Active:      false,
	}

	// UPDATE debe incluir el campo active - args: name, description, active, id
	mock.ExpectExec(`UPDATE "test_tenant".brand SET`).
		WithArgs("Marca Actualizada", "Nueva descripción", false, int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), "test_tenant", brand)

	assert.NoError(t, err)
}

func TestRepository_Update_Success(t *testing.T) {
	// Test: Actualización exitosa
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	brand := &Brand{
		ID:          1,
		Name:        "Marca Modificada",
		Description: "Descripción modificada",
		Active:      true,
	}

	mock.ExpectExec(`UPDATE "test_tenant".brand SET`).
		WithArgs("Marca Modificada", "Descripción modificada", true, int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), "test_tenant", brand)

	assert.NoError(t, err)
}

func TestRepository_Delete_Success(t *testing.T) {
	// Test: Eliminación lógica exitosa
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectExec(`UPDATE "test_tenant".brand SET deleted_at`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
}

func TestRepository_Page_Success(t *testing.T) {
	// Test: Paginación exitosa
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// COUNT query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "test_tenant".brand`).
		WithArgs(int64(1)).
		WillReturnRows(countRows)

	// SELECT query con paginación - args: enterprise_id (1), limit (10), offset (0)
	selectRows := sqlmock.NewRows([]string{
		"id", "name", "description", "active", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Marca 1", "Desc 1", true, 1, createdAt, nil, nil).
		AddRow(2, "Marca 2", "Desc 2", true, 1, createdAt, nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE enterprise_id`).
		WithArgs(int64(1), int64(10), int64(0)).
		WillReturnRows(selectRows)

	result, err := repo.Page(context.Background(), "test_tenant", 1, 1, 10, "", "id", "asc", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, int64(1), result.Page)
	assert.Equal(t, int64(10), result.Limit)
	assert.Equal(t, int64(1), result.TotalPages)

	items := result.Items.([]Brand)
	assert.Len(t, items, 2)
	assert.Equal(t, "Marca 1", items[0].Name)
	assert.True(t, items[0].Active)
}

func TestRepository_Page_WithSearch(t *testing.T) {
	// Test: Paginación con búsqueda
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// COUNT query con búsqueda - args: enterprise_id (1), search (%)
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "test_tenant".brand`).
		WithArgs(int64(1), "%marca%").
		WillReturnRows(countRows)

	// SELECT query con búsqueda
	selectRows := sqlmock.NewRows([]string{
		"id", "name", "description", "active", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Marca Test", "Desc", true, 1, createdAt, nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE enterprise_id`).
		WithArgs(int64(1), "%marca%", int64(10), int64(0)).
		WillReturnRows(selectRows)

	result, err := repo.Page(context.Background(), "test_tenant", 1, 1, 10, "marca", "name", "asc", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
	assert.Contains(t, result.Items.([]Brand)[0].Name, "Marca")
}

func TestRepository_Page_Empty(t *testing.T) {
	// Test: Paginación vacía
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	// COUNT query vacío
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "test_tenant".brand`).
		WithArgs(int64(1)).
		WillReturnRows(countRows)

	// SELECT query vacío
	selectRows := sqlmock.NewRows([]string{
		"id", "name", "description", "active", "enterprise_id", "created_at", "updated_at", "deleted_at",
	})

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE enterprise_id`).
		WithArgs(int64(1), int64(10), int64(0)).
		WillReturnRows(selectRows)

	result, err := repo.Page(context.Background(), "test_tenant", 1, 1, 10, "", "id", "asc", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.Total)
	assert.Equal(t, int64(0), result.TotalPages)
	assert.Empty(t, result.Items)
}

func TestRepository_Page_InvalidSort(t *testing.T) {
	// Test: Paginación con ordenamiento inválido - usa valores por defecto
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// COUNT query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "test_tenant".brand`).
		WithArgs(int64(1)).
		WillReturnRows(countRows)

	// SELECT query - debe usar "id" y "asc" por defecto
	selectRows := sqlmock.NewRows([]string{
		"id", "name", "description", "active", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Marca", "Desc", true, 1, createdAt, nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE enterprise_id`).
		WithArgs(int64(1), int64(10), int64(0)).
		WillReturnRows(selectRows)

	result, err := repo.Page(context.Background(), "test_tenant", 1, 1, 10, "", "invalid_col", "invalid_order", nil)

	assert.NoError(t, err)
	// Verifica que la consulta se ejecutó (usa id y asc por defecto)
	assert.NotNil(t, result.Items)
}

func TestRepository_Page_VerifyActiveField(t *testing.T) {
	// Test: Verifica que el campo active se devuelva correctamente en paginación
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// COUNT query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "test_tenant".brand`).
		WithArgs(int64(1)).
		WillReturnRows(countRows)

	// SELECT con diferentes valores de active
	selectRows := sqlmock.NewRows([]string{
		"id", "name", "description", "active", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Activa", "Desc", true, 1, createdAt, nil, nil).
		AddRow(2, "Inactiva", "Desc", false, 1, createdAt, nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".brand WHERE enterprise_id`).
		WithArgs(int64(1), int64(10), int64(0)).
		WillReturnRows(selectRows)

	result, err := repo.Page(context.Background(), "test_tenant", 1, 1, 10, "", "id", "asc", nil)

	assert.NoError(t, err)
	items := result.Items.([]Brand)
	assert.Len(t, items, 2)
	assert.True(t, items[0].Active)
	assert.False(t, items[1].Active)
}
