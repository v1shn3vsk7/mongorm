package examples

import (
	"context"

	"github.com/v1shn3vsk7/mongorm"
	"github.com/v1shn3vsk7/mongorm/options"
)

func _init(ctx context.Context, mngDsn string) *mongorm.Collection {
	opts := options.Client()
	opts.ApplyURI(mngDsn)

	client, err := mongorm.New(ctx, opts)
	if err != nil {
		// handle error
	}

	return client.Database("admin").Collection("users")
}
