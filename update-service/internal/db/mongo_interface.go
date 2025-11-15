package db

import (
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollectionInterface interface {
    UpdateOne(
        ctx context.Context,
        filter interface{},
        update interface{},
        opts ...*options.UpdateOptions,
    ) (*mongo.UpdateResult, error)
}
