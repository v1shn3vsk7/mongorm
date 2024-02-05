package options

import mongo_options "go.mongodb.org/mongo-driver/mongo/options"

type CollectionOptions struct {
	*mongo_options.CollectionOptions
}

func Collection() *CollectionOptions {
	return &CollectionOptions{}
}

func (c *CollectionOptions) ToMongo() *mongo_options.CollectionOptions {
	return c.CollectionOptions
}
