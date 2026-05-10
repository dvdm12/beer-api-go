package repository

import (
	"context"
	"dataanalysis/internal/db"
	"dataanalysis/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = CollectionBeers

// AnalysisRepository provides read access to the beers collection.
type AnalysisRepository struct {
	collection db.MongoCollectionInterface
	logger     Logger
}

// NewAnalysisRepository creates a new AnalysisRepository.
func NewAnalysisRepository(collection db.MongoCollectionInterface, logger Logger) *AnalysisRepository {
	if logger == nil {
		logger = defaultLogger
	}
	return &AnalysisRepository{collection: collection, logger: logger}
}

// buildQuery translates a BeerFilter into a MongoDB query.
func buildQuery(filter BeerFilter) bson.M {
	query := bson.M{}

	if filter.Brand != nil {
		query["brand"] = *filter.Brand
	}
	if filter.Name != nil {
		query["name"] = *filter.Name
	}
	if filter.NameLike != nil {
		query["name"] = primitive.Regex{Pattern: *filter.NameLike, Options: "i"}
	}

	alcohol := bson.M{}
	if filter.MinAlcohol != nil {
		alcohol["$gte"] = *filter.MinAlcohol
	}
	if filter.MaxAlcohol != nil {
		alcohol["$lte"] = *filter.MaxAlcohol
	}
	if len(alcohol) > 0 {
		query["alcohol"] = alcohol
	}

	year := bson.M{}
	if filter.FromYear != nil {
		year["$gte"] = *filter.FromYear
	}
	if filter.ToYear != nil {
		year["$lte"] = *filter.ToYear
	}
	if len(year) > 0 {
		query["year"] = year
	}

	return query
}

// Find returns beers matching the filter. Always returns a non-nil slice.
func (r *AnalysisRepository) Find(ctx context.Context, filter BeerFilter) ([]models.Beer, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, buildQuery(filter))
	if err != nil {
		return nil, MapMongoError(err, OpFind, collectionName, r.logger)
	}
	defer cursor.Close(ctx)

	beers := []models.Beer{}
	if err = cursor.All(ctx, &beers); err != nil {
		return nil, MapMongoError(err, OpFind, collectionName, r.logger)
	}
	return beers, nil
}

// FindTop returns the top-ranked beer by field.
// desc=true for descending, desc=false for ascending.
func (r *AnalysisRepository) FindTop(ctx context.Context, field SortField, desc bool) (*models.Beer, error) {
	if !field.IsValid() {
		return nil, ErrInvalidFilter
	}

	order := 1
	if desc {
		order = -1
	}

	opts := options.FindOne().SetSort(bson.D{{Key: string(field), Value: order}})

	var beer models.Beer
	err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&beer)
	if err != nil {
		return nil, MapMongoError(err, OpFindTop, collectionName, r.logger)
	}
	return &beer, nil
}
