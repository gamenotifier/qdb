package db

import (
	"context"
	"database/sql"
	"sync"
)

type Tx interface {
	DB
	Commit() error
	Rollback() error
}

type txImpl struct {
	db     DB
	tx     *sql.Tx
	done   bool
	doneMu sync.Mutex
}

func (t *txImpl) Query(ctx context.Context, query *Query) (Rows, error) {
	res, err := t.tx.QueryContext(ctx, query.query, query.args...)
	return wrapRows(query, res, err)
}

func (t *txImpl) QueryRow(ctx context.Context, query *Query) Row {
	res := t.tx.QueryRowContext(ctx, query.query, query.args...)
	return wrapRow(query, res)
}

func (t *txImpl) Exec(ctx context.Context, query *Query) (Result, error) {
	res, err := t.tx.ExecContext(ctx, query.query, query.args...)
	return wrapResult(query, res, err)
}

func (t *txImpl) Commit() error {
	t.doneMu.Lock()
	defer t.doneMu.Unlock()
	if t.done {
		return ErrTxDone
	}
	t.done = true
	return t.tx.Commit()
}

func (t *txImpl) Rollback() error {
	t.doneMu.Lock()
	defer t.doneMu.Unlock()
	if t.done {
		return ErrTxDone
	}
	t.done = true
	return t.tx.Rollback()
}

func (t *txImpl) Ping(ctx context.Context) error {
	// Use underlying database to ping
	return t.db.Ping(ctx)
}

func (t *txImpl) Begin(context.Context, *sql.TxOptions) (Tx, error) {
	return nil, ErrAlreadyInTx
}
