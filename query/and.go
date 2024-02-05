package query

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/v1shn3vsk7/mongorm/internal/operators"
)

func (q *Query) And() *Query {
	q.mx.Lock()
	defer q.mx.Unlock()

	// make new query and copy previous
	newQuery := make(bson.D, 0, len(q.query)+1)
	for i := 1; i < len(q.query); i++ {
		newQuery[i] = q.query[i]
	}

	// TODO: try to use copy to copy previous query
	//copy(newQuery, q.query)

	// add $and operator to beginning of the query
	newQuery[0] = bson.E{Key: operators.AND, Value: newQuery[1:]}
	q.query = newQuery

	return q
}
