package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type txKey struct{}

// injectTx injects transaction to context
func InjectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

var (
	ErrDuplicate = errors.New("login already in use")
)

type Storage struct {
	DB *sql.DB
}

// NewStorage opens a connection to the database and applies migrations.
func NewStorage(addr string) (*Storage, error) {
	db, err := goose.OpenDBWithDriver("pgx", addr)
	if err != nil {
		return nil, fmt.Errorf("goose: failed to open DB: %w", err)
	}

	err = goose.Up(db, "./migrations")
	if err != nil {
		return nil, fmt.Errorf("goose: failed to migrate: %w", err)
	}

	return &Storage{
		DB: db,
	}, nil
}

// CloseDB closes the database connection.
func (s *Storage) Close() error {
	return s.DB.Close()
}
