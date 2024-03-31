package query

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Query struct {
	query bson.D

	operators map[int]*empty
}

type empty struct{}

// New - deprecated, use high level func
func New() *Query {
	return &Query{
		query:     make(bson.D, 0),
		operators: make(map[int]*empty),
	}
}

func (q *Query) Bson() bson.D {
	return q.query
}
