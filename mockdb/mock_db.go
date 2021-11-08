package mockdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gamenotifier/qdb"
)

var (
	ErrNoQueryHook       = fmt.Errorf("mockdb: no query hook found")
	ErrQueryValuesNil    = fmt.Errorf("mockdb: QueryValuesFunc is nil")
	ErrQueryRowValuesNil = fmt.Errorf("mockdb: QueryRowValuesFunc is nil")
	ErrResultValuesNil   = fmt.Errorf("mockdb: ResultValuesFunc is nil")
)

// QueryValuesFunc returns a slice of many column values
type QueryValuesFunc func() [][]interface{}

// QueryRowValuesFunc returns a single slice of column values
type QueryRowValuesFunc func() []interface{}

// ResultValuesFunc returns the last insert id and the number of rows affected
type ResultValuesFunc func() (int64, int64)

type HookCounter interface {
	Name() string
	Triggered() int
}

type queryHook struct {
	name      string
	triggered int
	query     QueryValuesFunc
	queryRow  QueryRowValuesFunc
	result    ResultValuesFunc
}

func (q *queryHook) Name() string {
	return q.name
}

func (q *queryHook) Triggered() int {
	return q.triggered
}

func newQueryHook(name string, query QueryValuesFunc, queryRow QueryRowValuesFunc, result ResultValuesFunc) *queryHook {
	return &queryHook{
		name:      name,
		triggered: 0,
		query:     query,
		queryRow:  queryRow,
		result:    result,
	}
}

type MockDB struct {
	hooks map[string]*queryHook
}

func New() *MockDB {
	return &MockDB{
		hooks: make(map[string]*queryHook),
	}
}

func (m *MockDB) AddHook(queryName string, query QueryValuesFunc, queryRow QueryRowValuesFunc, result ResultValuesFunc) HookCounter {
	hook := newQueryHook(queryName, query, queryRow, result)
	m.hooks[queryName] = hook
	return hook
}

func (m *MockDB) AddQueryHook(queryName string, queryFunc QueryValuesFunc) HookCounter {
	return m.AddHook(queryName, queryFunc, nil, nil)
}


func (m *MockDB) AddQueryRowHook(queryName string, queryRowFunc QueryRowValuesFunc) HookCounter {
	return m.AddHook(queryName, nil, queryRowFunc, nil)
}

func (m *MockDB) AddQueryResultHook(queryName string, resultFunc ResultValuesFunc) HookCounter {
	return m.AddHook(queryName, nil, nil, resultFunc)
}

func (m *MockDB) RemoveHook(queryName string) {
	delete(m.hooks, queryName)
}

func (m *MockDB) ClearHooks() {
	m.hooks = make(map[string]*queryHook)
}

func (m *MockDB) Query(_ context.Context, query *db.Query) (db.Rows, error) {
	hook, ok := m.hooks[query.Name()]
	if !ok {
		return nil, ErrNoQueryHook
	} else if hook.query == nil {
		return nil, ErrQueryValuesNil
	}

	hook.triggered++
	return newMockRows(hook.query()), nil
}

func (m *MockDB) QueryRow(_ context.Context, query *db.Query) db.Row {
	hook, ok := m.hooks[query.Name()]
	if !ok {
		return newMockRowError(ErrNoQueryHook)
	} else if hook.queryRow == nil {
		return newMockRowError(ErrQueryRowValuesNil)
	}

	hook.triggered++
	return newMockRow(hook.queryRow())
}

func (m *MockDB) Exec(_ context.Context, query *db.Query) (db.Result, error) {
	hook, ok := m.hooks[query.Name()]
	if !ok {
		return nil, ErrNoQueryHook
	} else if hook.result == nil {
		return nil, ErrResultValuesNil
	}

	hook.triggered++
	return newMockResult(hook.result()), nil
}

func (m *MockDB) Begin(context.Context, *sql.TxOptions) (db.Tx, error) {
	return newMockTx(m), nil
}

// Ping is a NOP for MockDb
func (m *MockDB) Ping(context.Context) error {
	return nil
}
