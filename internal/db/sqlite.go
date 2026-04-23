package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// AdaptQueryForSQLite transforms a PostgreSQL query to SQLite-compatible syntax
func AdaptQueryForSQLite(query string) string {
	newQuery := query

	// Replace ILIKE with LIKE
	newQuery = strings.ReplaceAll(newQuery, "ILIKE", "LIKE")

	// Replace NOW() with datetime('now')
	newQuery = strings.ReplaceAll(newQuery, "NOW()", "datetime('now')")

	// Replace BIGSERIAL with INTEGER AUTOINCREMENT
	newQuery = strings.ReplaceAll(newQuery, "BIGSERIAL", "INTEGER AUTOINCREMENT")

	// Replace SERIAL with INTEGER AUTOINCREMENT
	newQuery = strings.ReplaceAll(newQuery, "SERIAL", "INTEGER AUTOINCREMENT")

	// Replace TIMESTAMPTZ with TEXT
	newQuery = strings.ReplaceAll(newQuery, "TIMESTAMPTZ", "TEXT")

	// Replace BOOLEAN with INTEGER
	newQuery = strings.ReplaceAll(newQuery, "BOOLEAN", "INTEGER")

	// Replace RETURNING clause (SQLite doesn't support it well)
	newQuery = strings.ReplaceAll(newQuery, "RETURNING id", "")
	newQuery = strings.ReplaceAll(newQuery, "RETURNING id,", "")
	newQuery = strings.ReplaceAll(newQuery, "RETURNING ", "")

	// Replace public. prefix (schema) with nothing
	newQuery = strings.ReplaceAll(newQuery, "public.", "")

	// TODO: si esta en modo offline reemplazar tenant desde el token
	newQuery = strings.ReplaceAll(newQuery, "tenant.", "")

	// Replace ON CONFLICT (...) DO UPDATE SET (PostgreSQL upsert) with OR REPLACE
	newQuery = strings.ReplaceAll(newQuery, "ON CONFLICT (slug) DO UPDATE SET", "OR REPLACE")

	// Convert $1, $2 placeholders to ?
	for i := 100; i >= 1; i-- {
		placeholder := fmt.Sprintf("$%d", i)
		if !strings.Contains(newQuery, placeholder) {
			continue
		}
		newQuery = strings.ReplaceAll(newQuery, placeholder, "?")
	}

	return newQuery
}

// SQLiteAdapter wraps *sql.DB to automatically adapt queries for SQLite
type SQLiteAdapter struct {
	db *sql.DB
}

// NewSQLiteAdapter creates a new SQLite adapter
func NewSQLiteAdapter(db *sql.DB) *SQLiteAdapter {
	return &SQLiteAdapter{db: db}
}

// DB returns the underlying *sql.DB
func (s *SQLiteAdapter) DB() *sql.DB {
	return s.db
}

// QueryContext executes a query with automatic adaptation
func (s *SQLiteAdapter) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, AdaptQueryForSQLite(query), args...)
}

// QueryRowContext executes a row query with automatic adaptation
func (s *SQLiteAdapter) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, AdaptQueryForSQLite(query), args...)
}

// ExecContext executes a statement with automatic adaptation
func (s *SQLiteAdapter) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, AdaptQueryForSQLite(query), args...)
}

// PrepareContext prepares a statement with automatic adaptation
func (s *SQLiteAdapter) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.db.PrepareContext(ctx, AdaptQueryForSQLite(query))
}

// Begin starts a transaction
func (s *SQLiteAdapter) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

// BeginTx starts a transaction with context
func (s *SQLiteAdapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, opts)
}

// Close closes the database
func (s *SQLiteAdapter) Close() error {
	return s.db.Close()
}

// Ping checks the connection
func (s *SQLiteAdapter) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// SetMaxIdleConns sets the max idle connections
func (s *SQLiteAdapter) SetMaxIdleConns(n int) {
	s.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the max open connections
func (s *SQLiteAdapter) SetMaxOpenConns(n int) {
	s.db.SetMaxOpenConns(n)
}
