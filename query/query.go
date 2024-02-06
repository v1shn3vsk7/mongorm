package query

import (
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

type Query struct {
	query bson.D
	mx    sync.RWMutex // TODO: decide if mutex is needed
}

// New - deprecated, use high level func
func New() *Query {
	return &Query{
		query: make(bson.D, 0),
		mx:    sync.RWMutex{},
	}
}

func (q *Query) Bson() bson.D {
	q.mx.Lock()
	defer q.mx.Unlock()

	return q.query
}
