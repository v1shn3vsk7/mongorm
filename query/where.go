package query

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/v1shn3vsk7/mongorm/internal/operators"
)

func (q *Query) Where(key string, cond int, value interface{}) *Query {
	q.mx.Lock()
	defer q.mx.Unlock()

	elem := bson.E{Key: key}

	switch cond {
	case EQ:
		elem.Value = value
	case NE:
		elem.Value = bson.D{{Key: operators.NE, Value: value}}

		// TODO: implement other operators
	}

	q.query = append(q.query, elem)

	return q
}
