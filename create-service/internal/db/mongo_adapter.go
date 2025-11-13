package db

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
)

type MongoCollectionAdapter struct {
    collection *mongo.Collection
}

func NewMongoCollectionAdapter(c *mongo.Collection) *MongoCollectionAdapter {
    return &MongoCollectionAdapter{collection: c}
}

func (m *MongoCollectionAdapter) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
    result, err := m.collection.InsertOne(ctx, document)
    return result, err
}
