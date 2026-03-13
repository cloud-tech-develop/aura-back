package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
	DSN string
}

func New(dsn string) (*DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir db: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("conectar db: %w", err)
	}
	return &DB{DB: conn, DSN: dsn}, nil
}
