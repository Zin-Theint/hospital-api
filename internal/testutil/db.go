package testutil

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Pool *pgxpool.Pool
	DSN  string
	term func() error
}

func NewTestDB(ctx context.Context) (*TestDB, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		Env:          map[string]string{"POSTGRES_PASSWORD": "pw", "POSTGRES_USER": "api", "POSTGRES_DB": "hospital"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		AutoRemove:   true,
	}
	ct, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, _ := ct.Host(ctx)
	p, _ := ct.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://api:pw@%s:%s/hospital?sslmode=disable", host, p.Port())

	// pgxpool
	cfg, _ := pgxpool.ParseConfig(dsn)
	cfg.MaxConnLifetime = time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		_ = ct.Terminate(ctx)
		return nil, err
	}

	migDir, _ := filepath.Abs("../../migrations")
	if err := goose.UpContext(ctx, pool, migDir); err != nil {
		pool.Close()
		_ = ct.Terminate(ctx)
		return nil, err
	}

	return &TestDB{
		Pool: pool,
		DSN:  dsn,
		term: func() error {
			pool.Close()
			return ct.Terminate(ctx)
		},
	}, nil
}

func (t *TestDB) Terminate(ctx context.Context) { _ = t.term() }
