package tests

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PGContainer struct {
	pgc *postgres.PostgresContainer
	db  *sql.DB
}

func NewPGContainer() (*PGContainer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbName := "rssagg"
	dbUser := "root"
	dbPassword := "secret"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	connString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("cannot get the connection string: %w", err)
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the db: %w", err)
	}

	if err = database.MigrateUp(db, "../sqlc/schema"); err != nil {
		return nil, fmt.Errorf("cannot migrate up: %w", err)
	}

	return &PGContainer{
		pgc: postgresContainer,
		db:  db,
	}, nil
}

func (p *PGContainer) DB() *sql.DB {
	return p.db
}

func (p *PGContainer) Terminate(ctx context.Context) error {
	if err := p.pgc.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate container: %w", err)
	}

	return nil
}
