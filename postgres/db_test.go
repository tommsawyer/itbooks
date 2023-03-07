package postgres

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

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
