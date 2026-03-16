package users

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestListRolesByMinLevel(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{q: db}

	// Test: roles with level >= 1 (ADMIN level)
	mock.ExpectQuery("SELECT id, name, description, level FROM public.roles").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "level"}).
			AddRow(2, "ADMIN", "Administrator", 1).
			AddRow(3, "SUPERVISOR", "Supervisor", 2).
			AddRow(4, "USER", "Standard user", 3))

	roles, err := repo.ListRolesByMinLevel(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, roles, 3)
	assert.Equal(t, "ADMIN", roles[0].Name)
	assert.Equal(t, 1, roles[0].Level)
}

func TestListRolesByMinLevel_SuperAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{q: db}

	// Test: roles with level >= 0 (SUPERADMIN level - all roles)
	mock.ExpectQuery("SELECT id, name, description, level FROM public.roles").
		WithArgs(0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "level"}).
			AddRow(1, "SUPERADMIN", "Super admin", 0).
			AddRow(2, "ADMIN", "Administrator", 1).
			AddRow(3, "SUPERVISOR", "Supervisor", 2))

	roles, err := repo.ListRolesByMinLevel(context.Background(), 0)
	assert.NoError(t, err)
	assert.Len(t, roles, 3)
	assert.Equal(t, "SUPERADMIN", roles[0].Name)
	assert.Equal(t, 0, roles[0].Level)
}

func TestListRolesByMinLevel_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{q: db}

	// Test: no roles found
	mock.ExpectQuery("SELECT id, name, description, level FROM public.roles").
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "level"}))

	roles, err := repo.ListRolesByMinLevel(context.Background(), 5)
	assert.NoError(t, err)
	assert.Len(t, roles, 0)
}

func TestService_AssignRoles_ValidatesLevel(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	svc := NewService(db, nil)
	ctx := context.Background()

	// Mock: check role level for role ID 1 (SUPERADMIN level 0)
	mock.ExpectQuery("SELECT level FROM public.roles").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"level"}).AddRow(0))

	// User has roleLevel 1 (ADMIN), trying to assign SUPERADMIN (level 0) should fail
	err = svc.AssignRoles(ctx, 1, []int64{1}, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no tiene permiso")
}

func TestService_AssignRoles_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	svc := NewService(db, nil)
	ctx := context.Background()

	// Mock: check role level for role ID 2 (ADMIN level 1)
	mock.ExpectQuery("SELECT level FROM public.roles").
		WithArgs(int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"level"}).AddRow(1))

	// Begin transaction
	mock.ExpectBegin()

	// Delete existing roles
	mock.ExpectExec("DELETE FROM public.user_roles").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Insert new role
	mock.ExpectExec("INSERT INTO public.user_roles").
		WithArgs(int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Commit
	mock.ExpectCommit()

	// User has roleLevel 1 (ADMIN), assigning ADMIN (level 1) should succeed
	err = svc.AssignRoles(ctx, 1, []int64{2}, 1)
	assert.NoError(t, err)
}

func TestService_AssignRoles_InvalidRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	svc := NewService(db, nil)
	ctx := context.Background()

	// Mock: role not found
	mock.ExpectQuery("SELECT level FROM public.roles").
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	err = svc.AssignRoles(ctx, 1, []int64{999}, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no encontrado")
}
