package db

import (
	"database/sql"
	"fmt"
)

var (
	ErrNoRows      = sql.ErrNoRows
	ErrTxDone      = sql.ErrTxDone
	ErrAlreadyInTx = fmt.Errorf("tx: already in transaction")
)

// QueryError describes an error that occurred while processing the result
// of an SQL query, whether during query time or scan time.
type QueryError struct {
	err   error
	query *Query
}

func newQueryError(err error, query *Query) *QueryError {
	return &QueryError{
		err:   err,
		query: query,
	}
}

// Error implements the error interface
func (qErr *QueryError) Error() string {
	return fmt.Sprintf("query %q: %s", qErr.query.name, qErr.err.Error())
}

// Unwrap implements the Wrapper interface
func (qErr *QueryError) Unwrap() error {
	return qErr.err
}

// String implements the Stringer interface. Simply returns Error().
func (qErr *QueryError) String() string {
	return qErr.Error()
}

// QueryBody exposes the wrapped query body
func (qErr *QueryError) QueryBody() string {
	return qErr.query.query
}

// QueryArgs exposes the wrapped query arguments
func (qErr *QueryError) QueryArgs() []interface{} {
	return qErr.query.args
}
