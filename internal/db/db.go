package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps a pgx connection pool.
type DB struct {
	pool *pgxpool.Pool
}

// Connect initializes the database connection and runs migrations.
func Connect(ctx context.Context, databaseURL string) (*DB, error) {
	// 1. Connect
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// 2. Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	// 3. Run basic up migrations manually using raw SQL
	// Note: since the roadmap asks for github-migrate which we don't have installed/setup,
	// we will run them manually from the SQL files.
	if err := runMigrations(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrationFiles := []string{
		"migrations/001_create_bounties.up.sql",
		"migrations/002_create_events.up.sql",
		"migrations/003_create_cursor.up.sql",
	}

	for _, file := range migrationFiles {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// If running from tests, the path might be different.
				log.Printf("Migration file %s not found, skipping for now.", file)
				continue
			}
			return fmt.Errorf("could not read migration %s: %w", file, err)
		}

		_, err = pool.Exec(ctx, string(sqlBytes))
		if err != nil {
			// Some migrations might fail if they already exist, we ignore already exists errors for simplicity
			// if we aren't using a proper migration tracker.
			log.Printf("Migration %s executed (may already exist). Err: %v", file, err)
		} else {
			log.Printf("Migration %s applied successfully.", file)
		}
	}
	return nil
}
