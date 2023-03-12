package postgres

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/tommsawyer/itbooks/postgres/postgrestest"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	flag.Parse()

	uri, stopPostgres, err := postgrestest.RunContainer(ctx, testing.Verbose())
	if err != nil {
		log.Fatalf("cannot run testing postgres: %v", err)
	}

	if err := Connect(ctx, uri); err != nil {
		log.Fatalf("cannot connect to testing postgres: %v", err)
	}

	code := m.Run()
	stopPostgres()
	os.Exit(code)
}

func testTransaction(t *testing.T) (context.Context, *is.I, func()) {
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
