// Package repository provides data access logic for the create service.
package repository

import (
	"context"
	"createservice/internal/db"
	"createservice/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateRepository handles persistence operations for beers.
type CreateRepository struct {
	collection db.MongoCollectionInterface
}

// NewCreateRepository creates a new repository instance.
func NewCreateRepository(collection db.MongoCollectionInterface) *CreateRepository {
	return &CreateRepository{
		collection: collection,
	}
}

// CreateBeer inserts a new beer into the database.
func (r *CreateRepository) CreateBeer(beer models.Beer) (*struct{}, error) {
	// Create context with timeout for DB operation.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert document into collection.
	_, err := r.collection.InsertOne(ctx, beer)
	return nil, err
}

// ExistsBeer checks if a beer with the given name exists.
func (r *CreateRepository) ExistsBeer(name string) (bool, error) {
	// Create context with timeout for DB operation.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Count documents matching the filter.
	count, err := r.collection.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}