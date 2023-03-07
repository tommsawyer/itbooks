package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/matryer/is"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
var db DB

// DB is postgres database.
type DB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Connect connects to postgres.
func Connect(ctx context.Context, postgresURI string) error {
	conn, err := pgx.Connect(ctx, postgresURI)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %v", err)
	}
	db = conn
	return nil
}

type transactionKey struct{}

// RunInTransaction runs provided callback in transaction.
// Commits if callback returns nil-error, rollback otherwise.
func RunInTransaction(ctx context.Context, cb func(ctx context.Context) error) error {
	db := getDB(ctx)
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, transactionKey{}, tx)

	err = cb(ctx)
	if err != nil {
		return tx.Rollback(ctx)
	}

	return tx.Commit(ctx)
}

// RunAndRollback run transaction and rollbacks everything.
// Used in tests only.
func RunAndRollback(t *testing.T) (context.Context, *is.I, func()) {
	ctx := context.Background()
	db := getDB(ctx)
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	ctx = context.WithValue(ctx, transactionKey{}, tx)

	return ctx, is.New(t), func() {
		if err := tx.Rollback(ctx); err != nil {
			t.Errorf("cannot rollback transcation: %v", err)
		}
	}
}

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
