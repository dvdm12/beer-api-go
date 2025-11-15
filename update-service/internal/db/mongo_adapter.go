package db

import (
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollection struct {
    Collection *mongo.Collection
}

func NewMongoCollection(c *mongo.Collection) *MongoCollection {
    return &MongoCollection{Collection: c}
}

func (m *MongoCollection) UpdateOne(
    ctx context.Context,
    filter interface{},
    update interface{},
    opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
    return m.Collection.UpdateOne(ctx, filter, update, opts...)
}
