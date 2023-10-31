package options

import mongo_options "go.mongodb.org/mongo-driver/mongo/options"

type ClientOptions struct {
	opts *mongo_options.ClientOptions
}

func Client() *ClientOptions {
	return &ClientOptions{
		opts: mongo_options.Client(),
	}
}

func (c *ClientOptions) ApplyURI(uri string) *ClientOptions {
	return &ClientOptions{
		opts: c.opts.ApplyURI(uri),
	}
}

func (c *ClientOptions) MongoOptions() *mongo_options.ClientOptions {
	return c.opts
}
