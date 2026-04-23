package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoSingleResultAdapter wraps a mongo.SingleResult.
// It implements SingleResultInterface for testability.
type MongoSingleResultAdapter struct {
	sr *mongo.SingleResult
}

// Decode delegates decoding to the underlying MongoDB result.
func (m *MongoSingleResultAdapter) Decode(v interface{}) error {
	return m.sr.Decode(v)
}

// MongoCollectionAdapter wraps a mongo.Collection.
// It implements MongoCollectionInterface for easier mocking.
type MongoCollectionAdapter struct {
	Collection *mongo.Collection
}

// NewMongoCollectionAdapter creates a new MongoCollectionAdapter.
func NewMongoCollectionAdapter(collection *mongo.Collection) *MongoCollectionAdapter {
	return &MongoCollectionAdapter{Collection: collection}
}

// FindOne executes a query and wraps the result in a SingleResultInterface.
func (m *MongoCollectionAdapter) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) SingleResultInterface {
	result := m.Collection.FindOne(ctx, filter, opts...)
	return &MongoSingleResultAdapter{sr: result}
}

// Find executes a query and returns a MongoDB cursor.
func (m *MongoCollectionAdapter) Find(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (*mongo.Cursor, error) {
	return m.Collection.Find(ctx, filter, opts...)
}
