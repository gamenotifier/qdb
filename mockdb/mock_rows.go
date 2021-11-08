package mockdb

import (
	"fmt"
	db "github.com/gamenotifier/qdb"
)

type mockRows struct {
	values  [][]interface{}
	currRow int
	lastErr error
	closed  bool
}

func newMockRows(values [][]interface{}) *mockRows {
	return &mockRows{
		values:  values,
		currRow: 0,
		lastErr: nil,
		closed:  false,
	}
}

// Close is a NOP for mockRows
func (m *mockRows) Close() error {
	m.closed = true
	return nil
}

func (m *mockRows) Scan(dest ...interface{}) (err error) {
	if m.lastErr != nil {
		return m.lastErr
	}

	// set lastErr automagically
	defer func() {
		if err != nil {
			m.lastErr = err
		}
	}()

	// check for some error conditions
	if !m.Next() {
		return fmt.Errorf("mockdb: no Rows available")
	} else if m.closed {
		return fmt.Errorf("mockdb: Rows are closed")
	} else if m.values == nil || len(m.values) == 0 || m.values[m.currRow] == nil {
		return db.ErrNoRows
	} else if len(dest) != len(m.values[m.currRow]) {
		return fmt.Errorf("mockdb: Expected %d arguments in Scan, not %d", len(m.values[m.currRow]), len(dest))
	}

	for i, sv := range m.values[m.currRow] {
		err := convertAssign(dest[i], sv)
		if err != nil {
			return fmt.Errorf("mockdb: Scan error on column index %d: %w", i, err)
		}
	}

	m.currRow++
	return nil
}

func (m *mockRows) Next() bool {
	return len(m.values) > m.currRow
}

func (m *mockRows) Err() error {
	return m.lastErr
}

func (m *mockRows) Columns() ([]string, error) {
	emptyCols := make([]string, len(m.values))
	return emptyCols, nil
}
