package db

import (
	"context"
	"database/sql"
)

type DB interface {
	Query(ctx context.Context, query *Query) (Rows, error)
	QueryRow(ctx context.Context, query *Query) Row
	Exec(ctx context.Context, query *Query) (Result, error)
	Begin(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Ping(ctx context.Context) error
}

type dbImpl struct {
	db *sql.DB
}

func (d *dbImpl) Query(ctx context.Context, query *Query) (Rows, error) {
	res, err := d.db.QueryContext(ctx, query.query, query.args...)
	return wrapRows(query, res, err)
}

func (d *dbImpl) QueryRow(ctx context.Context, query *Query) Row {
	res := d.db.QueryRowContext(ctx, query.query, query.args...)
	return wrapRow(query, res)
}

func (d *dbImpl) Exec(ctx context.Context, query *Query) (Result, error) {
	res, err := d.db.ExecContext(ctx, query.query, query.args...)
	return wrapResult(query, res, err)
}

func (d *dbImpl) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *dbImpl) Begin(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := d.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &txImpl{db: d, tx: tx, done: false}, nil
}

func New(driver string, uri string) (DB, error) {
	db, err := sql.Open(driver, uri)
	if err != nil {
		return nil, err
	}

	return &dbImpl{
		db: db,
	}, nil
}
