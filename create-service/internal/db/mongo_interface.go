package db

import "context"

type MongoCollectionInterface interface {
    InsertOne(ctx context.Context, document interface{}) (interface{}, error)
}
