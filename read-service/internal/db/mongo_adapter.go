package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollectionAdapter struct {
	Collection *mongo.Collection
}

func NewMongoCollectionAdapter(collection *mongo.Collection) *MongoCollectionAdapter {
	return &MongoCollectionAdapter{
		Collection: collection,
	}
}

func (m *MongoCollectionAdapter) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.Collection.FindOne(ctx, filter, opts...)
}

func (m *MongoCollectionAdapter) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return m.Collection.Find(ctx, filter, opts...)
}
