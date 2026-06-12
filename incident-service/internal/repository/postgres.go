package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

// NewPool creates a pgx connection pool and verifies the database
// is reachable. The caller is responsible for calling pool.Close().
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errs.Wrapf(errs.KindInternal, "repository.NewPool", err, "create pool")
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, errs.Wrapf(errs.KindUnavailable, "repository.NewPool", err, "ping database")
	}

	return pool, nil
}
