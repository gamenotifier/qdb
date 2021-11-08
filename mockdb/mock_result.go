package mockdb

type mockResult struct {
	lastInsertID int64
	rowsAffected int64
}

func newMockResult(lastInsertID, rowsAffected int64) *mockResult {
	return &mockResult{
		lastInsertID: lastInsertID,
		rowsAffected: rowsAffected,
	}
}

func (m *mockResult) LastInsertId() (int64, error) {
	return m.lastInsertID, nil
}

func (m *mockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}
