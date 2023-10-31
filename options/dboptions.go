package options

import mongo_options "go.mongodb.org/mongo-driver/mongo/options"

type DatabaseOptions struct {
	*mongo_options.DatabaseOptions
}

func (d *DatabaseOptions) ToMongoOptions() *mongo_options.DatabaseOptions {
	return d.DatabaseOptions
}
