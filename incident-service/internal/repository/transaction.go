package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

// Querier is the subset of pgx methods the repositories rely on.
// It is satisfied by both *pgxpool.Pool and pgx.Tx, so a repository can
// run either against the connection pool directly or inside a transaction.
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// TxManager runs unit-of-work operations inside a single database transaction.
type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// WithTx begins a transaction and invokes fn with incident and event
// repositories bound to it. The transaction is committed if fn returns nil
// and rolled back otherwise (including on panic). Errors returned by fn are
// propagated unchanged so callers can branch on their Kind.
func (m *TxManager) WithTx(ctx context.Context, fn func(incident.Repository, event.Repository) error) error {
	const op = "repository.TxManager.WithTx"

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return errs.Wrapf(errs.KindUnavailable, op, err, "begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(NewIncidentRepository(tx), NewEventRepository(tx)); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return errs.Wrapf(errs.KindUnavailable, op, err, "commit transaction")
	}

	return nil
}
