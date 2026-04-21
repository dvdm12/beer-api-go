// Package db provides MongoDB abstractions.
package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoCollectionAdapter wraps a MongoDB collection.
type MongoCollectionAdapter struct {
	collection *mongo.Collection
}

// NewMongoCollectionAdapter creates a new collection adapter.
func NewMongoCollectionAdapter(c *mongo.Collection) *MongoCollectionAdapter {
	return &MongoCollectionAdapter{collection: c}
}

// InsertOne inserts a document into the collection.
func (m *MongoCollectionAdapter) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	result, err := m.collection.InsertOne(ctx, document)
	return result, err
}

// CountDocuments returns the number of documents matching a filter.
func (m *MongoCollectionAdapter) CountDocuments(ctx context.Context, filter bson.M) (int64, error) {
	return m.collection.CountDocuments(ctx, filter)
}