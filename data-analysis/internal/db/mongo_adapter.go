package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCollectionAdapter wraps a MongoDB collection.
type MongoCollectionAdapter struct {
	collection *mongo.Collection
}

// NewMongoCollectionAdapter creates a new collection adapter.
func NewMongoCollectionAdapter(c *mongo.Collection) *MongoCollectionAdapter {
	return &MongoCollectionAdapter{collection: c}
}

// FindOne returns the first document matching the filter.
func (m *MongoCollectionAdapter) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.collection.FindOne(ctx, filter, opts...)
}

// Find returns a cursor over documents matching the filter.
func (m *MongoCollectionAdapter) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return m.collection.Find(ctx, filter, opts...)
}
