package examples

import (
	"context"
	"log"

	"github.com/v1shn3vsk7/mongorm"
	"github.com/v1shn3vsk7/mongorm/options"
)

func init_example(ctx context.Context, mngDsn string) *mongorm.Client {
	opts := options.Client().ApplyURI(mngDsn)
	client, err := mongorm.New(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
