package postgrestest

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tommsawyer/itbooks/postgres/migrations"
)

// RunContainer runs testing postgres in docker.
func RunContainer(ctx context.Context, verbose bool) (uri string, cleanup func(), err error) {
	if !verbose {
		testcontainers.Logger = log.New(io.Discard, "", log.LstdFlags)
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       "test",
			"POSTGRES_PORT":     "5432",
		},
		WaitingFor: wait.ForExposedPort(),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	endpoint, err := postgresContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("failed to get postgres endpoint: %v", err)
	}
	postgresURI := "postgres://postgres:secret@" + endpoint + "/test?sslmode=disable"
	if err := applyMigrations(postgresURI); err != nil {
		log.Fatalf("failed to apply postgres migrations: %v", err)
	}

	return postgresURI, func() {
		if err := postgresContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate postgres: %v", err)
		}
	}, nil
}

func applyMigrations(postgresURI string) error {
	d, err := iofs.New(migrations.Migrations, ".")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, postgresURI)
	if err != nil {
		return err
	}
	if testing.Verbose() {
		m.Log = &migrateLogger{log.New(os.Stdout, "", log.LstdFlags)}
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		return sourceErr
	}

	return dbErr
}

type migrateLogger struct {
	*log.Logger
}

func (*migrateLogger) Verbose() bool { return false }
