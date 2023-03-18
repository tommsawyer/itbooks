package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	db   DB
)

// Fields is used for update queries.
type Fields = map[string]any

// DB is postgres database.
type DB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Connect connects to postgres.
func Connect(ctx context.Context, uri string) error {
	conn, err := pgx.Connect(ctx, uri)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	db = conn
	return nil
}

type transactionKey struct{}

// getDB returns transaction attached to context
// or global database connection.
func getDB(ctx context.Context) DB {
	tx, ok := ctx.Value(transactionKey{}).(pgx.Tx)
	if ok {
		return tx
	}

	if db == nil {
		panic("call postgres.Connect() before using queries")
	}

	return db
}
