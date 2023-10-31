package mongorm

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/v1shn3vsk7/mongorm/options"
)

type Client struct {
	*mongo.Client
}

func New(ctx context.Context, opts ...*options.ClientOptions) (*Client, error) {
	mongoOpts := make([]*mongo_options.ClientOptions, 0, len(opts))
	for _, opt := range opts {
		mongoOpts = append(mongoOpts, opt.MongoOptions())
	}

	client, err := mongo.Connect(ctx, mongoOpts...)
	if err != nil {
		return nil, fmt.Errorf("mongorm: err connect to mongo client: %v", err)
	}

	return &Client{
		client,
	}, nil
}

func (c *Client) Database(name string, opts ...*options.DatabaseOptions) *Database {
	mongoOpts := make([]*mongo_options.DatabaseOptions, 0, len(opts))
	for _, opt := range opts {
		mongoOpts = append(mongoOpts, opt.ToMongoOptions())
	}

	return &Database{
		c.Client.Database(name, mongoOpts...),
	}
}
