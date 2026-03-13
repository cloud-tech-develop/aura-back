package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Migrator struct {
	db *sql.DB
}

func New(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// MigrateSchema aplica todas las migraciones pendientes a un esquema específico
func (m *Migrator) MigrateSchema(schemaName string) error {
	// 1. Crear el esquema si no existe
	if _, err := m.db.Exec(fmt.Sprintf(
		"CREATE SCHEMA IF NOT EXISTS %q", schemaName,
	)); err != nil {
		return fmt.Errorf("crear esquema: %w", err)
	}

	// 2. Crear tabla de control de migraciones en el esquema
	if _, err := m.db.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %q.schema_migrations (
            version     BIGINT PRIMARY KEY,
            name        TEXT NOT NULL,
            applied_at  TIMESTAMPTZ DEFAULT NOW()
        )`, schemaName,
	)); err != nil {
		return fmt.Errorf("crear tabla de migraciones: %w", err)
	}

	// 3. Leer versiones ya aplicadas
	applied, err := m.appliedVersions(schemaName)
	if err != nil {
		return err
	}

	// 4. Leer archivos de migración
	files, err := m.migrationFiles()
	if err != nil {
		return err
	}

	// 5. Aplicar las pendientes en una transacción por cada una
	for _, f := range files {
		if applied[f.version] {
			continue
		}
		if err := m.applyMigration(schemaName, f); err != nil {
			return fmt.Errorf("aplicando migración %s: %w", f.name, err)
		}
		fmt.Printf("[%s] ✓ migración aplicada: %s\n", schemaName, f.name)
	}

	return nil
}

type migrationFile struct {
	version int64
	name    string
	path    string
}

func (m *Migrator) migrationFiles() ([]migrationFile, error) {
	entries, err := filepath.Glob("internal/migration/migrations/public/*.sql")
	if err != nil {
		return nil, err
	}

	var files []migrationFile
	for _, path := range entries {
		base := filepath.Base(path)
		parts := strings.SplitN(strings.TrimSuffix(base, ".sql"), "_", 2)
		if len(parts) < 2 {
			continue
		}
		var version int64
		fmt.Sscanf(parts[0], "%d", &version)
		files = append(files, migrationFile{
			version: version,
			name:    base,
			path:    path,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].version < files[j].version
	})
	return files, nil
}

func (m *Migrator) appliedVersions(schema string) (map[int64]bool, error) {
	rows, err := m.db.Query(fmt.Sprintf(
		`SELECT version FROM %q.schema_migrations`, schema,
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := map[int64]bool{}
	for rows.Next() {
		var v int64
		rows.Scan(&v)
		applied[v] = true
	}
	return applied, nil
}

func (m *Migrator) applyMigration(schema string, f migrationFile) error {
	content, err := os.ReadFile(f.path)
	if err != nil {
		return err
	}

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Establecer search_path para la transacción
	if _, err := tx.Exec(fmt.Sprintf("SET LOCAL search_path TO %q, public", schema)); err != nil {
		return err
	}

	// Ejecutar el SQL de la migración
	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	// Registrar la versión aplicada
	if _, err := tx.Exec(fmt.Sprintf(
		`INSERT INTO %q.schema_migrations (version, name) VALUES ($1, $2)`,
		schema,
	), f.version, f.name); err != nil {
		return err
	}

	return tx.Commit()
}
