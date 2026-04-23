package repository

import (
	"context"
	"time"

	"updateservice/internal/db"
	"updateservice/internal/errors"
	"updateservice/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateRepository handles database operations for updating beers.
type UpdateRepository struct {
	collection db.MongoCollectionInterface
}

// NewUpdateRepository creates a new UpdateRepository instance.
func NewUpdateRepository(collection db.MongoCollectionInterface) *UpdateRepository {
	return &UpdateRepository{collection: collection}
}

// UpdateBeer updates a beer document using the provided context.
func (r *UpdateRepository) UpdateBeer(ctx context.Context, id string, beer models.Beer) error {
	// Add a timeout to the inherited context.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Parse the string ID into a MongoDB ObjectID.
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequest("Invalid beer ID format")
	}

	// Build the query filter and update payload.
	filter := bson.M{"_id": oid}
	update := buildBeerUpdate(beer)

	// Execute the update operation.
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Internal(err)
	}

	if result.MatchedCount == 0 {
		return errors.NotFound("Beer")
	}

	return nil
}

// buildBeerUpdate constructs the BSON update payload from the domain entity.
func buildBeerUpdate(beer models.Beer) bson.M {
	return bson.M{
		"$set": bson.M{
			"name":    beer.Name,
			"brand":   beer.Brand,
			"alcohol": beer.Alcohol,
			"year":    beer.Year,
		},
	}
}
