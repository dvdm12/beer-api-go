package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoCollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}) (interface{}, error)
	CountDocuments(ctx context.Context, filter bson.M) (int64, error)
}