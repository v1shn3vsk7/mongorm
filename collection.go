package mongorm

import (
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/v1shn3vsk7/mongorm/options"
)

type Collection struct {
	*mongo.Collection
}

func (db *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	mongoOpts := make([]*mongo_options.CollectionOptions, 0, len(opts))
	for _, opt := range opts {
		mongoOpts = append(mongoOpts, opt.ToMongo())
	}

	return &Collection{
		db.Database.Collection(name, mongoOpts...),
	}
}
