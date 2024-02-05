package examples

import (
	"context"

	"github.com/v1shn3vsk7/mongorm"
	"github.com/v1shn3vsk7/mongorm/query"
)

func Where() {
	ctx := context.Background()
	client, _ := mongorm.New(ctx, nil)

	usersCL := client.Database("admin").Collection("users")

	query := query.New()
	query.
		Where("user_id", mongorm.EQ, "00000").
		And().
		Where("user_name", mongorm.NE, "dolbayeb")

	_ = usersCL.FindOne(ctx, query.Bson(), nil).Decode(struct{}{})
}
