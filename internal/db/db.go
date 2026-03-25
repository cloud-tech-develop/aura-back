package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
	DSN    string
	Driver string
}

func New(driver, dsn string) (*DB, error) {
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir db (%s): %w", driver, err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("conectar db (%s): %w", driver, err)
	}
	return &DB{DB: conn, DSN: dsn, Driver: driver}, nil
}

// Querier interface matches the one used in modules
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Wrap returns a Querier that automatically adapts queries for the current driver
func (db *DB) Wrap(q Querier) Querier {
	return &wrappedQuerier{q: q, driver: db.Driver}
}

type wrappedQuerier struct {
	q      Querier
	driver string
}

func (w *wrappedQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return w.q.QueryContext(ctx, adaptQuery(query, w.driver), args...)
}

func (w *wrappedQuerier) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return w.q.QueryRowContext(ctx, adaptQuery(query, w.driver), args...)
}

func (w *wrappedQuerier) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return w.q.ExecContext(ctx, adaptQuery(query, w.driver), args...)
}

func adaptQuery(query string, driver string) string {
	if driver != "sqlite" {
		return query
	}
	newQuery := query
	// Replace ILIKE with LIKE if SQLite (basic support)
	newQuery = strings.ReplaceAll(newQuery, "ILIKE", "LIKE")
	// Replace NOW() with datetime('now') if SQLite
	newQuery = strings.ReplaceAll(newQuery, "NOW()", "datetime('now')")

	for i := 100; i >= 1; i-- {
		placeholder := fmt.Sprintf("$%d", i)
		if !strings.Contains(newQuery, placeholder) {
			continue
		}
		newQuery = strings.ReplaceAll(newQuery, placeholder, "?")
	}
	return newQuery
}
