package mongorm

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/v1shn3vsk7/mongorm/internal/action"
	"github.com/v1shn3vsk7/mongorm/internal/operators"
)

type Query struct {
	collection *Collection
	actions    []*action.Action
	query      bson.D
}

func (c *Collection) Query() *Query {
	return &Query{
		collection: c,
		actions:    make([]*action.Action, 0),
	}
}

func (q *Query) Where(key string, cond CondOperator, value interface{}) *Query {
	q.actions = append(q.actions, &action.Action{
		Type:     action.Where,
		Key:      key,
		Operator: cond.Uint8(),
		Value:    value,
	})

	return q

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

func (q *Query) where(key string, cond CondOperator, value interface{}) {
	elem := bson.E{Key: key}

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
}

func (q *Query) And() *Query {
	q.actions = append(q.actions, &action.Action{
		Type: action.And,
	})

	return q
}

func (q *Query) and() {
	// make new query and copy previous
	newQuery := make(bson.D, 0, len(q.query)+1)
	for i := 1; i < len(q.query); i++ {
		newQuery[i] = q.query[i]
	}

	// add $and operator to beginning of the query
	newQuery[0] = bson.E{Key: operators.AND, Value: newQuery[1:]}

	q.query = newQuery
}

func (q *Query) Or() *Query {
	q.actions = append(q.actions, &action.Action{
		Type: action.Or,
	})

	return q
}

func (q *Query) or() {
	newQuery := make(bson.D, 0, len(q.query)+1)
	for i := 1; i < len(q.query); i++ {
		newQuery[i] = q.query[i]
	}

	newQuery[0] = bson.E{Key: operators.OR, Value: newQuery[1:]}

	q.query = newQuery
}

func (q *Query) Bson() bson.D {
	if len(q.actions) == 0 {
		return q.query
	}

	for _, actionF := range q.actions {
		switch actionF.Type {
		case action.Where:
			q.where(actionF.Key, CondOperator(actionF.Operator), actionF.Value)
		case action.Or:
			q.or()
		case action.And:
			q.and()
		}
	}

	return q.query
}
