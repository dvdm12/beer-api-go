package db

import (
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollectionAdapter struct {
    collection *mongo.Collection
}

func NewMongoCollectionAdapter(c *mongo.Collection) *MongoCollectionAdapter {
    return &MongoCollectionAdapter{collection: c}
}

func (m *MongoCollectionAdapter) UpdateOne(
    ctx context.Context,
    filter interface{},
    update interface{},
    opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
    return m.collection.UpdateOne(ctx, filter, update, opts...)
}
