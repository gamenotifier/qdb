package db

import "fmt"

type Query struct {
	name  string
	query string
	args  []interface{}
}

func NewQuery(name string, query string, args ...interface{}) *Query {
	return &Query{
		name:  name,
		query: query,
		args:  args,
	}
}

func (q *Query) Name() string {
	return q.name
}

func (q *Query) String() string {
	return fmt.Sprintf("query %q: %q { %v }", q.name, q.query, q.args)
}
