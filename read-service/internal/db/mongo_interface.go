package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SingleResultInterface abstracts a MongoDB single result.
// It allows mocking the Decode method in tests.
type SingleResultInterface interface {
	// Decode decodes the result into the provided value.
	Decode(v interface{}) error
}

// MongoCollectionInterface abstracts MongoDB collection operations.
// It enables mocking database interactions in tests.
type MongoCollectionInterface interface {
	// FindOne retrieves a single document matching the filter.
	// It returns a SingleResultInterface for easier testing.
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResultInterface

	// Find retrieves multiple documents matching the filter.
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}
