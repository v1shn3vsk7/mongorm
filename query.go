package mongorm

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/v1shn3vsk7/mongorm/internal/operators"
)

type Query struct {
	query      bson.D
	collection *Collection
}

func (c *Collection) Query() *Query {
	return &Query{
		collection: c,
		query:      make(bson.D, 0),
	}
}

func (q *Query) Where(key string, cond CondOperator, value interface{}) *Query {
	elem := bson.E{Key: key}

	// TODO: change checking for map querying
	switch cond {
	case EQ:
		elem.Value = value
	case NE:
		elem.Value = bson.D{{Key: operators.NE, Value: value}}
	case GT:
		elem.Value = bson.D{{Key: operators.GT, Value: value}}
	case GTE:
		elem.Value = bson.D{{Key: operators.GTE, Value: value}}
	case LT:
		elem.Value = bson.D{{Key: operators.LT, Value: value}}
	case LTE:
		elem.Value = bson.D{{Key: operators.LTE, Value: value}}
	}

	q.query = append(q.query, elem)

	return q
}

func (q *Query) And() *Query {
	// make new query and copy previous
	newQuery := make(bson.D, 0, len(q.query)+1)
	for i := 1; i < len(q.query); i++ {
		newQuery[i] = q.query[i]
	}

	// add $and operator to beginning of the query
	newQuery[0] = bson.E{Key: operators.AND, Value: newQuery[1:]}

	q.query = newQuery

	return q
}

func (q *Query) Or() *Query {
	newQuery := make(bson.D, 0, len(q.query)+1)
	for i := 1; i < len(q.query); i++ {
		newQuery[i] = q.query[i]
	}

	newQuery[0] = bson.E{Key: operators.OR, Value: newQuery[1:]}

	q.query = newQuery

	return q
}

func (q *Query) FindOne(ctx context.Context) {
	q.collection.FindOne(ctx, q.query)
}
