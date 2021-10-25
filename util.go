package db

import (
	"database/sql"
	"io"
)

type Result sql.Result

type Row interface {
	Scan(dest ...interface{}) error
}

type Rows interface {
	io.Closer
	Scan(dest ...interface{}) error
	Next() bool
	Err() error
	Columns() ([]string, error)
}

// ExpectRowsAffected returns whether RowsAffected() succeeded, and an error if the count was not as expected.
func ExpectRowsAffected(res Result, countExpected int64, errIfNot error) (bool, error) {
	count, err := res.RowsAffected()
	if err != nil {
		// Something happened, but perhaps the operation went through.
		// Should we care? Probably not, but return false so the caller knows.
		return false, nil
	} else if count != countExpected {
		return true, errIfNot
	}
	return true, nil
}
