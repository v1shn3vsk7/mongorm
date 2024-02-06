package examples

import (
	"context"

	"github.com/v1shn3vsk7/mongorm"
)

func Where() {
	ctx := context.Background()
	client, _ := mongorm.New(ctx, nil)

	usersCL := client.Database("admin").Collection("users")

	usersCL.Query().
		Where("user_id", mongorm.EQ, "00000").
		And().
		Where("user_name", mongorm.NE, "dolbayeb").
		FindOne(ctx)

}
