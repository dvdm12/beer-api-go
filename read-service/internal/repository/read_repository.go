// Package repository handles data persistence for read operations.
package repository

import (
	"context"
	"errors"
	"readservice/internal/db"
	"readservice/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Sentinel errors representing domain-specific cases.
var ErrBeerNotFound = errors.New("beer not found")
var ErrInvalidID = errors.New("invalid beer ID format")

// ReadRepository provides read access to the data store.
type ReadRepository struct {
	collection db.MongoCollectionInterface
}

// NewReadRepository creates a new ReadRepository instance.
func NewReadRepository(collection db.MongoCollectionInterface) *ReadRepository {
	return &ReadRepository{collection: collection}
}

// GetBeerByID retrieves a beer by its ID.
//
// Workflow:
//  1. Validates the ID format.
//  2. Queries the database.
//  3. Maps database errors to domain errors.
//
// Returns:
//   - A pointer to Beer if found.
//   - ErrInvalidID if the ID format is invalid.
//   - ErrBeerNotFound if no document exists.
//   - A raw error for unexpected database failures.
func (r *ReadRepository) GetBeerByID(id string) (*models.Beer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Validate ID format
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	// Query database
	var beer models.Beer
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&beer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrBeerNotFound
		}
		return nil, err
	}

	return &beer, nil
}

// GetAllBeers retrieves all beers from the collection.
//
// Returns:
//   - A slice of Beer (can be empty if no records exist).
//   - An error if the query fails.
func (r *ReadRepository) GetAllBeers() ([]models.Beer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var beers []models.Beer
	if err = cursor.All(ctx, &beers); err != nil {
		return nil, err
	}

	return beers, nil
}
