package mockdb

import (
	"context"
	"database/sql"
	db "github.com/gamenotifier/qdb"
)

type mockTx struct {
	db *MockDB
}

func newMockTx(source *MockDB) *mockTx {
	return &mockTx{
		db: source,
	}
}

func (m *mockTx) Query(ctx context.Context, query *db.Query) (db.Rows, error) {
	return m.db.Query(ctx, query)
}

func (m *mockTx) QueryRow(ctx context.Context, query *db.Query) db.Row {
	return m.db.QueryRow(ctx, query)
}

func (m *mockTx) Exec(ctx context.Context, query *db.Query) (db.Result, error) {
	return m.db.Exec(ctx, query)
}

func (m *mockTx) Begin(ctx context.Context, opts *sql.TxOptions) (db.Tx, error) {
	return nil, db.ErrAlreadyInTx
}

// Ping is a NOP for mockTx
func (m *mockTx) Ping(ctx context.Context) error {
	return nil
}

// Commit is a NOP for mockTx
func (m *mockTx) Commit() error {
	return nil
}

// Rollback is a NOP for mockTx
func (m *mockTx) Rollback() error {
	return nil
}
