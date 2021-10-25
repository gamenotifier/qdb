package db

import (
	"sync"
)

// Wrapper struct around a Rows instance
// Calls to Rows methods will wrap errors with the corresponding query.
type wrappedRows struct {
	query *Query
	rows  Rows

	// the last error encountered, wrapped
	lastErr *QueryError
	mu      sync.Mutex
}

func (w *wrappedRows) Close() error {
	return w.rows.Close()
}

func (w *wrappedRows) Scan(dest ...interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.lastErr == nil {
		if err := w.rows.Scan(dest...); err != nil {
			w.lastErr = newQueryError(err, w.query)
			return w.lastErr
		} else {
			// Scan returned no errors
			return nil
		}
	} else {
		// We already cached an error
		return w.lastErr
	}
}

func (w *wrappedRows) Next() bool {
	return w.rows.Next()
}

func (w *wrappedRows) Err() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.lastErr == nil {
		if err := w.rows.Err(); err != nil {
			w.lastErr = newQueryError(err, w.query)
			return w.lastErr
		} else {
			return nil
		}
	} else {
		// We already cached an error
		return w.lastErr
	}
}

func (w *wrappedRows) Columns() ([]string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.lastErr == nil {
		if cols, err := w.rows.Columns(); err != nil {
			w.lastErr = newQueryError(err, w.query)
			return nil, w.lastErr
		} else {
			return cols, nil
		}
	} else {
		// We already cached an error
		return nil, w.lastErr
	}
}

// wrapRows combines a query, a Rows instance, and a possible error into a wrappedRows instance
func wrapRows(query *Query, res Rows, err error) (Rows, error) {
	if err != nil {
		return nil, newQueryError(err, query)
	} else {
		return &wrappedRows{
			query: query,
			rows:  res,
		}, nil
	}
}

// Wrapper struct around a Row instance
// Calls to Row methods will wrap errors with the corresponding query.
type wrappedRow struct {
	query *Query
	row   Row

	// the last error encountered, wrapped
	lastErr *QueryError
	mu      sync.Mutex
}

func (w *wrappedRow) Scan(dest ...interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.lastErr == nil {
		if err := w.row.Scan(dest...); err != nil {
			// Scan returned an error, make sure to wrap it
			w.lastErr = newQueryError(err, w.query)
			return w.lastErr
		} else {
			// Scan returned no errors
			return nil
		}
	} else {
		// We already cached an error
		return w.lastErr
	}
}

// wrapRow combines a query and a Row instance into a wrappedRow instance
func wrapRow(query *Query, res Row) Row {
	return &wrappedRow{
		query: query,
		row:   res,
	}
}

// Wrapper struct around a Result instance
// Calls to Result methods will wrap errors with the corresponding query.
type wrappedResult struct {
	result Result
	query  *Query
}

func (w *wrappedResult) LastInsertId() (int64, error) {
	n, err := w.result.LastInsertId()
	if err != nil {
		return 0, newQueryError(err, w.query)
	} else {
		return n, nil
	}
}

func (w *wrappedResult) RowsAffected() (int64, error) {
	n, err := w.result.RowsAffected()
	if err != nil {
		return 0, newQueryError(err, w.query)
	} else {
		return n, nil
	}
}

// wrapRows combines a query, a Result instance, and a possible error into a wrappedResult instance
func wrapResult(query *Query, res Result, err error) (Result, error) {
	if err != nil {
		return nil, newQueryError(err, query)
	} else {
		return &wrappedResult{
			result: res,
			query:  query,
		}, nil
	}
}
