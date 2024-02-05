package mongorm

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	*mongo.Database
}
