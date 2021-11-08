package mockdb

import (
	db "github.com/gamenotifier/qdb"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMockDB_Query(t *testing.T) {
	mdb := New()
	q1 := db.NewQuery("query_1", "SELECT name, age FROM people")
	mdb.AddQueryHook("query_1", func() [][]interface{} {
		return [][]interface{}{
			{"Abby", 20, true},
			{"Braden", 43, true},
			{"Caitlyn", 37, false},
		}
	})

	res, err := mdb.Query(nil, q1)
	assert.Nil(t, err)

	type person struct {
		name      string
		age       int
		likesCake bool
	}

	expected := []person{
		{"Abby", 20, true},
		{"Braden", 43, true},
		{"Caitlyn", 37, false},
	}

	var people []person
	for res.Next() {
		var newPerson person
		err := res.Scan(&newPerson.name, &newPerson.age, &newPerson.likesCake)
		assert.Nil(t, err)
		people = append(people, newPerson)
	}

	assert.Nil(t, res.Close())
	assert.Equal(t, expected, people)
}

func TestMockDB_QueryRow(t *testing.T) {
	janFirst := time.Unix(1609477200, 0).In(time.UTC)

	mdb := New()
	q1 := db.NewQuery("query_1", "SELECT username, email, verified FROM users WHERE id=$1", 123)
	mdb.AddQueryRowHook("query_1", func() []interface{} {
		return []interface{}{
			"john.smith.32",
			"john.smith@example.com",
			false,
			janFirst,
		}
	})

	type user struct {
		username          string
		email             string
		verified          bool
		creationTimestamp string // test converting different types
	}

	expected := user{
		"john.smith.32",
		"john.smith@example.com",
		false,
		janFirst.Format(time.RFC3339Nano),
	}

	var u user
	res := mdb.QueryRow(nil, q1)
	err := res.Scan(&u.username, &u.email, &u.verified, &u.creationTimestamp)
	assert.Nil(t, err)
	assert.Equal(t, expected, u)
}

func TestMockDB_Exec(t *testing.T) {
	mdb := New()
	q1 := db.NewQuery("query_1", "UPDATE rooms SET open=$1 WHERE building=$2", false, "Chemistry Building")
	mdb.AddQueryResultHook("query_1", func() (int64, int64) {
		return 0, 24
	})

	q2 := db.NewQuery("query_2", "INSERT INTO rooms(name, building, open) VALUES($1, $2, $3)", "Lecture Hall 123", "Chemistry Building", true)
	mdb.AddQueryResultHook("query_2", func() (int64, int64) {
		return 25, 1
	})

	// Run query 1
	res, err := mdb.Exec(nil, q1)
	assert.Nil(t, err)

	_, err = res.LastInsertId()
	assert.Nil(t, err)
	rowsAffected, err := res.RowsAffected()
	assert.Nil(t, err)
	assert.Equal(t, rowsAffected, int64(24))

	// Run query 2
	res, err = mdb.Exec(nil, q2)
	assert.Nil(t, err)

	lastInsertID, err := res.LastInsertId()
	assert.Nil(t, err)
	rowsAffected, err = res.RowsAffected()
	assert.Nil(t, err)
	assert.Equal(t, lastInsertID, int64(25))
	assert.Equal(t, rowsAffected, int64(1))
}

func TestMockDB_Query_Fail(t *testing.T) {
	mdb := New()
	q1 := db.NewQuery("query_1", "SELECT happy FROM people")
	_, err := mdb.Query(nil, q1)
	assert.Equal(t, ErrNoQueryHook, err)

	// note that this is not AddQueryHook
	mdb.AddQueryRowHook("query_1", func() []interface{} {
		return []interface{}{
			true,
		}
	})
	_, err = mdb.Query(nil, q1)
	assert.Equal(t, ErrQueryValuesNil, err)

	// Now actually add the hook
	mdb.AddQueryHook("query_1", func() [][]interface{} {
		return [][]interface{}{
			{true},
		}
	})
	_, err = mdb.Query(nil, q1)
	assert.Nil(t, err)
}

func TestMockDB_QueryRow_Fail(t *testing.T) {
	var useless bool

	mdb := New()
	q1 := db.NewQuery("query_1", "SELECT happy FROM people WHERE id=$1", 123)
	res := mdb.QueryRow(nil, q1)
	err := res.Scan(&useless)
	assert.Equal(t, ErrNoQueryHook, err)

	// note that this is not AddQueryRowHook
	mdb.AddQueryHook("query_1", func() [][]interface{} {
		return [][]interface{}{
			{true},
		}
	})
	res = mdb.QueryRow(nil, q1)
	err = res.Scan(&useless)
	assert.Equal(t, ErrQueryRowValuesNil, err)

	// Now actually add the hook
	mdb.AddQueryRowHook("query_1", func() []interface{} {
		return []interface{}{
			true,
		}
	})

	res = mdb.QueryRow(nil, q1)
	err = res.Scan(&useless)
	assert.Nil(t, err)
}

func TestMockDB_Exec_Fail(t *testing.T) {
	mdb := New()
	q1 := db.NewQuery("query_1", "UPDATE people SET happy=$1", true)
	_, err := mdb.Exec(nil, q1)
	assert.Equal(t, ErrNoQueryHook, err)

	// note that this is not AddQueryResultHook
	mdb.AddQueryHook("query_1", func() [][]interface{} {
		return [][]interface{}{
			{},
		}
	})
	_, err = mdb.Exec(nil, q1)
	assert.Equal(t, ErrResultValuesNil, err)

	// Now actually add the hook
	mdb.AddQueryResultHook("query_1", func() (int64, int64) {
		return 0, 0
	})
	_, err = mdb.Exec(nil, q1)
	assert.Nil(t, err)
}