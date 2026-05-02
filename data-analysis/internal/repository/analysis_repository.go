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

// AnalysisRepository provides read access to beers.
type AnalysisRepository struct {
	collection db.MongoCollectionInterface
	logger     Logger
}

// NewAnalysisRepository creates a new repository instance.
func NewAnalysisRepository(collection db.MongoCollectionInterface, logger Logger) *AnalysisRepository {
	if logger == nil {
		logger = defaultLogger
	}
	return &AnalysisRepository{collection: collection, logger: logger}
}

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

// GetByID retrieves a single beer by its unique identifier.
func (r *AnalysisRepository) GetByID(ctx context.Context, id string) (*models.Beer, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	var beer models.Beer
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&beer)
	if err != nil {
		return nil, MapMongoError(err, OpFindOne, collectionName, r.logger)
	}
	return &beer, nil
}

// Find retrieves all beers matching the given filter.
// Always returns a non-nil slice.
func (r *AnalysisRepository) Find(ctx context.Context, filter BeerFilter) ([]models.Beer, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, buildQuery(filter))
	if err != nil {
		return nil, MapMongoError(err, OpFind, collectionName, r.logger)
	}
	defer cursor.Close(ctx)

	beers := []models.Beer{} // non-nil empty slice by contract
	if err = cursor.All(ctx, &beers); err != nil {
		return nil, MapMongoError(err, OpFind, collectionName, r.logger)
	}
	return beers, nil
}

// FindOne retrieves the first beer matching the given filter.
func (r *AnalysisRepository) FindOne(ctx context.Context, filter BeerFilter) (*models.Beer, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	var beer models.Beer
	err := r.collection.FindOne(ctx, buildQuery(filter)).Decode(&beer)
	if err != nil {
		return nil, MapMongoError(err, OpFindOne, collectionName, r.logger)
	}
	return &beer, nil
}

// FindTop retrieves the beer ranked first by the given SortField.
func (r *AnalysisRepository) FindTop(ctx context.Context, field SortField, desc bool) (*models.Beer, error) {
	if !field.IsValid() {
		return nil, ErrInvalidID
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

// GetAll returns all beers.
func (r *AnalysisRepository) GetAll(ctx context.Context) ([]models.Beer, error) {
	return r.Find(ctx, BeerFilter{})
}

// GetByBrand returns beers by brand.
func (r *AnalysisRepository) GetByBrand(ctx context.Context, brand string) ([]models.Beer, error) {
	return r.Find(ctx, BeerFilter{Brand: &brand})
}

// GetByAlcoholMin returns beers with alcohol >= min.
func (r *AnalysisRepository) GetByAlcoholMin(ctx context.Context, min float64) ([]models.Beer, error) {
	return r.Find(ctx, BeerFilter{MinAlcohol: &min})
}

// GetByAlcoholMax returns beers with alcohol <= max.
func (r *AnalysisRepository) GetByAlcoholMax(ctx context.Context, max float64) ([]models.Beer, error) {
	return r.Find(ctx, BeerFilter{MaxAlcohol: &max})
}

// GetByAlcoholRange returns beers within alcohol range.
func (r *AnalysisRepository) GetByAlcoholRange(ctx context.Context, min, max float64) ([]models.Beer, error) {
	return r.Find(ctx, BeerFilter{MinAlcohol: &min, MaxAlcohol: &max})
}

// GetByYearRange returns beers within year range.
func (r *AnalysisRepository) GetByYearRange(ctx context.Context, from, to int) ([]models.Beer, error) {
	return r.Find(ctx, BeerFilter{FromYear: &from, ToYear: &to})
}

// GetByName returns a beer by exact name.
func (r *AnalysisRepository) GetByName(ctx context.Context, name string) (*models.Beer, error) {
	return r.FindOne(ctx, BeerFilter{Name: &name})
}

// GetByNameLike returns beers matching name pattern.
func (r *AnalysisRepository) GetByNameLike(ctx context.Context, name string) ([]models.Beer, error) {
	return r.Find(ctx, BeerFilter{NameLike: &name})
}

// GetStrongest returns the beer with highest alcohol.
func (r *AnalysisRepository) GetStrongest(ctx context.Context) (*models.Beer, error) {
	return r.FindTop(ctx, SortByAlcohol, true)
}

// GetOldest returns the oldest beer.
func (r *AnalysisRepository) GetOldest(ctx context.Context) (*models.Beer, error) {
	return r.FindTop(ctx, SortByYear, false)
}

// GetNewest returns the newest beer.
func (r *AnalysisRepository) GetNewest(ctx context.Context) (*models.Beer, error) {
	return r.FindTop(ctx, SortByYear, true)
}

// GetStats returns general statistics.
func (r *AnalysisRepository) GetStats(ctx context.Context) (*models.GeneralStats, error) {
	beers, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(beers) == 0 {
		return nil, ErrEmptyCollection
	}

	stats := &models.GeneralStats{Total: int64(len(beers))}
	brandSet := make(map[string]struct{})
	minAlcohol := beers[0].Alcohol
	maxAlcohol := beers[0].Alcohol
	minYear := beers[0].Year
	maxYear := beers[0].Year
	totalAlcohol := 0.0

	for _, b := range beers {
		brandSet[b.Brand] = struct{}{}
		totalAlcohol += b.Alcohol
		if b.Alcohol < minAlcohol {
			minAlcohol = b.Alcohol
		}
		if b.Alcohol > maxAlcohol {
			maxAlcohol = b.Alcohol
		}
		if b.Year < minYear {
			minYear = b.Year
		}
		if b.Year > maxYear {
			maxYear = b.Year
		}
	}

	stats.AvgAlcohol = totalAlcohol / float64(len(beers))
	stats.MinAlcohol = minAlcohol
	stats.MaxAlcohol = maxAlcohol
	stats.OldestYear = minYear
	stats.NewestYear = maxYear
	stats.TotalBrands = int64(len(brandSet))

	return stats, nil
}

// GetStatsByBrand returns stats grouped by brand.
func (r *AnalysisRepository) GetStatsByBrand(ctx context.Context) ([]models.BrandStats, error) {
	beers, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(beers) == 0 {
		return nil, ErrEmptyCollection
	}

	brandMap := make(map[string]*models.BrandStats)

	for _, b := range beers {
		if _, ok := brandMap[b.Brand]; !ok {
			brandMap[b.Brand] = &models.BrandStats{
				Brand:      b.Brand,
				MinAlcohol: b.Alcohol,
				MaxAlcohol: b.Alcohol,
			}
		}
		s := brandMap[b.Brand]
		s.Count++
		s.TotalAlcohol += b.Alcohol
		if b.Alcohol < s.MinAlcohol {
			s.MinAlcohol = b.Alcohol
		}
		if b.Alcohol > s.MaxAlcohol {
			s.MaxAlcohol = b.Alcohol
		}
	}

	result := make([]models.BrandStats, 0, len(brandMap))
	for _, s := range brandMap {
		s.AvgAlcohol = s.TotalAlcohol / float64(s.Count)
		result = append(result, *s)
	}

	return result, nil
}
