package examples

import (
	"context"

	"github.com/v1shn3vsk7/mongorm"
)

func Where() {
	ctx := context.Background()
	client, _ := mongorm.New(ctx, nil)

	usersCL := client.Database("<database>").Collection("<collection>")

	var userDTO struct {
		UserID   string `bson:"user_id"`
		Username string `bson:"user_name"`
	}

	query := usersCL.
		Query().
		Where("<key>", mongorm.EQ, "<value>").
		And().
		Where("<key>", mongorm.NE, "<value>")

	err := usersCL.FindOne(ctx, query.Bson(), nil).Decode(&userDTO)
	if err != nil {
		// handle error
	}
}
