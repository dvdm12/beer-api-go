package repository

import (
	"context"
	"time"

	"deleteservice/internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteRepository struct {
	collection db.MongoCollectionInterface
}

func NewDeleteRepository(collection db.MongoCollectionInterface) *DeleteRepository {
	return &DeleteRepository{
		collection: collection,
	}
}

func (r *DeleteRepository) DeleteBeer(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return primitive.ErrInvalidHex
	}

	return nil
}
