package units

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create_Success(t *testing.T) {
	// Test: Inserción exitosa de unidad
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	unit := &Unit{
		Name:          "Kilogramo",
		Abbreviation:  "kg",
		Active:        true,
		AllowDecimals: true,
		EnterpriseID:  1,
	}

	// Expect INSERT with RETURNING - QueryRowContext executes single query with INSERT+RETURNING
	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, "2024-01-01T00:00:00Z")
	mock.ExpectQuery(`INSERT INTO "test_tenant"`).
		WithArgs("Kilogramo", "kg", true, true, int64(1)).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), "test_tenant", unit)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), unit.ID)
	assert.Equal(t, "2024-01-01T00:00:00Z", unit.CreatedAt)
}

func TestRepository_Create_DuplicateName(t *testing.T) {
	// Test: Error de nombre duplicado
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	unit := &Unit{
		Name:          "Kilogramo",
		Abbreviation:  "kg",
		Active:        true,
		AllowDecimals: true,
		EnterpriseID:  1,
	}

	// Simular error de PostgreSQL por restricción única
	mock.ExpectQuery(`INSERT INTO "test_tenant"`).
		WithArgs("Kilogramo", "kg", true, true, int64(1)).
		WillReturnError(&pqError{code: "23505", message: "duplicate key value violates unique constraint"})

	err = repo.Create(context.Background(), "test_tenant", unit)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
}

// pqError simula un error de PostgreSQL
type pqError struct {
	code    string
	message string
}

func (e *pqError) Error() string {
	return e.message
}

func TestRepository_GetByID_Success(t *testing.T) {
	// Test: Obtención por ID exitosa
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "name", "abbreviation", "active", "allow_decimals", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Kilogramo", "kg", true, true, 1, "2024-01-01T00:00:00Z", nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".unit`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	unit, err := repo.GetByID(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), unit.ID)
	assert.Equal(t, "Kilogramo", unit.Name)
	assert.Equal(t, "kg", unit.Abbreviation)
	assert.True(t, unit.Active)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	// Test: No encontrado
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".unit`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	unit, err := repo.GetByID(context.Background(), "test_tenant", 999)

	assert.Nil(t, unit)
	assert.Error(t, err)
}

func TestRepository_List_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Kilogramo").
		AddRow(2, "Litro")

	mock.ExpectQuery(`SELECT id, name FROM "test_tenant".unit WHERE enterprise_id`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	units, err := repo.List(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Len(t, units, 2)
	assert.Equal(t, "Kilogramo", units[0].Name)
	assert.Equal(t, "Litro", units[1].Name)
}

func TestRepository_List_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "name"})

	mock.ExpectQuery(`SELECT id, name FROM "test_tenant".unit WHERE enterprise_id`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	units, err := repo.List(context.Background(), "test_tenant", 1)

	assert.NoError(t, err)
	assert.Empty(t, units)
}

func TestRepository_Update_Success(t *testing.T) {
	// Test: Actualización exitosa
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	unit := &Unit{
		ID:            1,
		Name:          "Kilogramo Updated",
		Abbreviation:  "kgg",
		Active:        false,
		AllowDecimals: false,
	}

	mock.ExpectExec(`UPDATE "test_tenant".unit SET`).
		WithArgs("Kilogramo Updated", "kgg", false, false, int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), "test_tenant", unit)

	assert.NoError(t, err)
}

func TestRepository_Delete_Success(t *testing.T) {
	// Test: Eliminación lógica exitosa
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectExec(`UPDATE "test_tenant".unit SET deleted_at`).
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

	// COUNT query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "test_tenant".unit`).
		WithArgs(int64(1)).
		WillReturnRows(countRows)

	// SELECT query con paginación - args: enterprise_id (1), limit (10), offset (0)
	selectRows := sqlmock.NewRows([]string{
		"id", "name", "abbreviation", "active", "allow_decimals", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Kilogramo", "kg", true, true, 1, "2024-01-01T00:00:00Z", nil, nil).
		AddRow(2, "Litro", "L", true, false, 1, "2024-01-01T00:00:00Z", nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".unit WHERE enterprise_id`).
		WithArgs(int64(1), int64(10), int64(0)).
		WillReturnRows(selectRows)

	result, err := repo.Page(context.Background(), "test_tenant", 1, 1, 10, "", "id", "asc", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, int64(1), result.Page)
	assert.Equal(t, int64(10), result.Limit)
	assert.Equal(t, int64(1), result.TotalPages)

	items := result.Items.([]Unit)
	assert.Len(t, items, 2)
	assert.Equal(t, "Kilogramo", items[0].Name)
}

func TestRepository_Page_WithSearch(t *testing.T) {
	// Test: Paginación con búsqueda
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	// COUNT query con búsqueda - args: enterprise_id (1), search (kg%)
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "test_tenant".unit`).
		WithArgs(int64(1), "%kg%").
		WillReturnRows(countRows)

	// SELECT query con búsqueda - args: enterprise_id (1), search (kg%), limit (10), offset (0)
	selectRows := sqlmock.NewRows([]string{
		"id", "name", "abbreviation", "active", "allow_decimals", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "Kilogramo", "kg", true, true, 1, "2024-01-01T00:00:00Z", nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".unit WHERE enterprise_id`).
		WithArgs(int64(1), "%kg%", int64(10), int64(0)).
		WillReturnRows(selectRows)

	result, err := repo.Page(context.Background(), "test_tenant", 1, 1, 10, "kg", "name", "asc", nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
	assert.Contains(t, result.Items.([]Unit)[0].Name, "Kilogramo")
}

func TestRepository_Page_Empty(t *testing.T) {
	// Test: Paginación vacía
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	// COUNT query vacío
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "test_tenant".unit`).
		WithArgs(int64(1)).
		WillReturnRows(countRows)

	// SELECT query vacío
	selectRows := sqlmock.NewRows([]string{
		"id", "name", "abbreviation", "active", "allow_decimals", "enterprise_id", "created_at", "updated_at", "deleted_at",
	})

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".unit WHERE enterprise_id`).
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

	// COUNT query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "test_tenant".unit`).
		WithArgs(int64(1)).
		WillReturnRows(countRows)

	// SELECT query - debe usar "id" y "asc" por defecto
	selectRows := sqlmock.NewRows([]string{
		"id", "name", "abbreviation", "active", "allow_decimals", "enterprise_id", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "kg", "kg", true, true, 1, "2024-01-01T00:00:00Z", nil, nil)

	mock.ExpectQuery(`SELECT .* FROM "test_tenant".unit WHERE enterprise_id`).
		WithArgs(int64(1), int64(10), int64(0)).
		WillReturnRows(selectRows)

	result, err := repo.Page(context.Background(), "test_tenant", 1, 1, 10, "", "invalid_col", "invalid_order", nil)

	assert.NoError(t, err)
	// Verifica que la consulta se ejecutó (usa id y asc por defecto)
	assert.NotNil(t, result.Items)
}
