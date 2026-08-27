package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DB {
	return &DB{Pool: pool}
}

func (db *DB) ensureUser(ctx context.Context, userID string) error {
	_, err := db.Pool.Exec(ctx, `INSERT INTO users (id) VALUES ($1) ON CONFLICT DO NOTHING`, userID)
	return err
}
