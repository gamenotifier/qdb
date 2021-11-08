package mockdb

import (
	"fmt"
	db "github.com/gamenotifier/qdb"
)

type mockRow struct {
	values  []interface{}
	lastErr error
}

func newMockRow(values []interface{}) *mockRow {
	return &mockRow{
		values:  values,
		lastErr: nil,
	}
}

func newMockRowError(err error) *mockRow {
	return &mockRow{
		values:  nil,
		lastErr: err,
	}
}

func (m *mockRow) Scan(dest ...interface{}) (err error) {
	if m.lastErr != nil {
		return m.lastErr
	}

	// set lastErr automagically
	defer func() {
		if err != nil {
			m.lastErr = err
		}
	}()

	if m.values == nil || len(m.values) == 0 {
		return db.ErrNoRows
	}

	if len(dest) != len(m.values) {
		return fmt.Errorf("mockdb: Expected %d arguments in Scan, not %d", len(m.values), len(dest))
	}

	for i, sv := range m.values {
		err := convertAssign(dest[i], sv)
		if err != nil {
			return fmt.Errorf("mockdb: Scan error on column index %d: %w", i, err)
		}
	}

	return nil
}
