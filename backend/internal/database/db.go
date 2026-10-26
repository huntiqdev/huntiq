package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/config"
)

// ensureSchema creates any tables that must exist for the app to function
// but may not have been created by the migration runner yet.
// This is idempotent — safe to run on every startup.
func ensureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS admin_permissions (
			admin_id    INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			permission  VARCHAR(50) NOT NULL,
			enabled     BOOLEAN     NOT NULL DEFAULT FALSE,
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (admin_id, permission)
		);
		CREATE INDEX IF NOT EXISTS idx_admin_permissions_admin ON admin_permissions(admin_id);

		ALTER TABLE admin_notifications
			ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL DEFAULT 'info',
			ADD COLUMN IF NOT EXISTS read_by INT[] NOT NULL DEFAULT '{}';
	`)
	return err
}

func newPool(cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}

	if v := cfg.DBMaxConns; v > 0 {
		poolCfg.MaxConns = int32(v)
	}
	if v := cfg.DBMinConns; v > 0 {
		poolCfg.MinConns = int32(v)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pool, nil
}

// NewPool creates a pool and runs ensureSchema. Use for the main app server.
func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := newPool(cfg)
	if err != nil {
		return nil, err
	}
	if err := ensureSchema(context.Background(), pool); err != nil {
		return nil, fmt.Errorf("ensureSchema: %w", err)
	}
	return pool, nil
}

// NewPoolNoSchema creates a pool without running ensureSchema. Use for the migrate command.
func NewPoolNoSchema(cfg *config.Config) (*pgxpool.Pool, error) {
	return newPool(cfg)
}
