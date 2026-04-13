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

// WithSchema returns a Querier that sets the search_path before queries
func (db *DB) WithSchema(q Querier, schema string) Querier {
	return &schemaQuerier{db: db.DB, q: q, driver: db.Driver, schema: schema}
}

type schemaQuerier struct {
	db     *sql.DB
	q      Querier
	driver string
	schema string
}

func (s *schemaQuerier) execWithSchema(ctx context.Context, fn func() error) error {
	if s.schema == "" {
		return fn()
	}
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("obtener conexión: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, fmt.Sprintf("SET search_path TO %q", s.schema)); err != nil {
		return fmt.Errorf("set search_path: %w", err)
	}
	return fn()
}

func (s *schemaQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error

	if s.schema == "" {
		rows, err = s.q.QueryContext(ctx, adaptQuery(query, s.driver), args...)
	} else {
		err = s.execWithSchema(ctx, func() error {
			rows, err = s.q.QueryContext(ctx, adaptQuery(query, s.driver), args...)
			return err
		})
	}
	return rows, err
}

func (s *schemaQuerier) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if s.schema == "" {
		return s.q.QueryRowContext(ctx, adaptQuery(query, s.driver), args...)
	}
	// For schema queries, we need to use a connection with the schema set
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return &sql.Row{}
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, fmt.Sprintf("SET search_path TO %q", s.schema)); err != nil {
		return &sql.Row{}
	}

	return conn.QueryRowContext(ctx, adaptQuery(query, s.driver), args...)
}

func (s *schemaQuerier) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error

	if s.schema == "" {
		result, err = s.q.ExecContext(ctx, adaptQuery(query, s.driver), args...)
	} else {
		err = s.execWithSchema(ctx, func() error {
			result, err = s.q.ExecContext(ctx, adaptQuery(query, s.driver), args...)
			return err
		})
	}
	return result, err
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
