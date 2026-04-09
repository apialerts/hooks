package db

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations
var migrationsFS embed.FS

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database url: %w", err)
	}
	config.MaxConns = 5
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return pool, nil
}

func RunMigrations(databaseURL string) error {
	goose.SetBaseFS(migrationsFS)

	db, err := goose.OpenDBWithDriver("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("goose open: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}
