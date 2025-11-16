package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCollectionInterface interface {
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
}
