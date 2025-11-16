package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCollectionAdapter struct {
	Collection *mongo.Collection
}

func NewMongoCollectionAdapter(collection *mongo.Collection) *MongoCollectionAdapter {
	return &MongoCollectionAdapter{
		Collection: collection,
	}
}

func (m *MongoCollectionAdapter) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return m.Collection.DeleteOne(ctx, filter)
}
