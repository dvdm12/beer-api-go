// Package repository handles data persistence.
package repository

import (
	"context"
	"deleteservice/internal/db"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sentinel errors for domain cases.
var ErrBeerNotFound = errors.New("beer not found")
var ErrInvalidID    = errors.New("invalid beer ID format")

// DeleteRepository interacts with MongoDB.
type DeleteRepository struct {
	collection db.MongoCollectionInterface
}

// NewDeleteRepository creates a repository instance.
func NewDeleteRepository(collection db.MongoCollectionInterface) *DeleteRepository {
	return &DeleteRepository{collection: collection}
}

// DeleteBeer removes a beer by ID.
func (r *DeleteRepository) DeleteBeer(id string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Parse ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	// Execute delete
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	// Check if document was deleted
	if result.DeletedCount == 0 {
		return ErrBeerNotFound
	}

	return nil
}