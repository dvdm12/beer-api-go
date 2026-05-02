// Package services defines domain-level models and interfaces for beer analysis.
// It encapsulates filtering, sorting, pagination, and service contracts.
package services

import (
	"context"
	stderrors "errors"
	"fmt"

	"dataanalysis/internal/models"
)

const (
	// MaxAlcoholPct is the upper bound for alcohol percentage.
	MaxAlcoholPct = 100.0

	// MinAlcoholPct is the lower bound for alcohol percentage.
	MinAlcoholPct = 0.0

	// DefaultLimit is the default number of results returned when no limit is specified.
	DefaultLimit = 50

	// MaxLimit is the maximum allowed number of results per query.
	MaxLimit = 100
)

// SortField defines allowed fields for sorting.
// The zero value is invalid; use predefined constants.
type SortField string

const (
	SortByAlcohol SortField = "alcohol"
	SortByYear    SortField = "year"
	SortByName    SortField = "name"
)

// IsValid reports whether the SortField is supported.
func (s SortField) IsValid() bool {
	switch s {
	case SortByAlcohol, SortByYear, SortByName:
		return true
	default:
		return false
	}
}

// SortOrder defines the direction of sorting.
// The zero value is invalid; use SortAsc or SortDesc.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// IsValid reports whether the SortOrder is supported.
func (o SortOrder) IsValid() bool {
	return o == SortAsc || o == SortDesc
}

// Sort represents a single sorting criterion.
type Sort struct {
	Field SortField
	Order SortOrder
}

// Validate checks whether the Sort is valid.
func (s Sort) Validate() error {
	if !s.Field.IsValid() {
		return fmt.Errorf("%w: invalid sort field '%s'", ErrInvalidBeerQuery, s.Field)
	}
	if !s.Order.IsValid() {
		return fmt.Errorf("%w: invalid sort order '%s'", ErrInvalidBeerQuery, s.Order)
	}
	return nil
}

// Pagination defines offset-based pagination parameters.
type Pagination struct {
	Limit  int
	Offset int
}

// Validate checks whether pagination parameters are valid.
func (p Pagination) Validate() error {
	if p.Limit < 0 {
		return fmt.Errorf("%w: limit must be non-negative", ErrInvalidBeerQuery)
	}
	if p.Offset < 0 {
		return fmt.Errorf("%w: offset must be non-negative", ErrInvalidBeerQuery)
	}
	return nil
}

// Normalize returns a new Pagination with defaults and bounds applied.
// It does not mutate the receiver.
func (p Pagination) Normalize() Pagination {
	n := p

	if n.Limit <= 0 {
		n.Limit = DefaultLimit
	}
	if n.Limit > MaxLimit {
		n.Limit = MaxLimit
	}
	if n.Offset < 0 {
		n.Offset = 0
	}

	return n
}

// ErrInvalidBeerQuery indicates invalid query configuration.
// It is intended for internal validation and should be wrapped before returning.
var ErrInvalidBeerQuery = stderrors.New("invalid beer query")

// ErrBeerNotFound indicates no beer matched the given identifier.
var ErrBeerNotFound = stderrors.New("beer not found")

// ErrNoResults indicates a query returned no results.
var ErrNoResults = stderrors.New("no results found")

// BeerFilter defines optional filtering parameters.
// Nil fields are ignored.
type BeerFilter struct {
	Brand      *string
	Name       *string
	NameLike   *string
	MinAlcohol *float64
	MaxAlcohol *float64
	FromYear   *int
	ToYear     *int
}

// validate performs internal validation of filter constraints.
func (f BeerFilter) validate() error {
	if f.Brand != nil && *f.Brand == "" {
		return fmt.Errorf("%w: brand cannot be empty", ErrInvalidBeerQuery)
	}
	if f.Name != nil && f.NameLike != nil {
		return fmt.Errorf("%w: name and nameLike are mutually exclusive", ErrInvalidBeerQuery)
	}
	if f.MinAlcohol != nil && *f.MinAlcohol < MinAlcoholPct {
		return fmt.Errorf("%w: minAlcohol cannot be negative", ErrInvalidBeerQuery)
	}
	if f.MaxAlcohol != nil && *f.MaxAlcohol > MaxAlcoholPct {
		return fmt.Errorf("%w: maxAlcohol cannot exceed %.0f", ErrInvalidBeerQuery, MaxAlcoholPct)
	}
	if f.MinAlcohol != nil && f.MaxAlcohol != nil && *f.MinAlcohol > *f.MaxAlcohol {
		return fmt.Errorf("%w: minAlcohol cannot exceed maxAlcohol", ErrInvalidBeerQuery)
	}
	if f.FromYear != nil && f.ToYear != nil && *f.FromYear > *f.ToYear {
		return fmt.Errorf("%w: fromYear cannot exceed toYear", ErrInvalidBeerQuery)
	}
	return nil
}

// BeerQuery aggregates filtering, sorting, and pagination.
type BeerQuery struct {
	Filter     BeerFilter
	Sort       []Sort
	Pagination Pagination
}

// Validate validates all query components and returns a normalized copy.
// It does not mutate the original query.
func (q BeerQuery) Validate() (BeerQuery, error) {
	if err := q.Filter.validate(); err != nil {
		return BeerQuery{}, err
	}

	for _, s := range q.Sort {
		if err := s.Validate(); err != nil {
			return BeerQuery{}, err
		}
	}

	if err := q.Pagination.Validate(); err != nil {
		return BeerQuery{}, err
	}

	normalized := q
	normalized.Pagination = q.Pagination.Normalize()

	return normalized, nil
}

// BeerQueryService defines read operations for beers.
type BeerQueryService interface {
	// GetByID retrieves a beer by its identifier.
	// The id must be a 24-character hexadecimal string.
	GetByID(ctx context.Context, id string) (*models.Beer, error)

	// FindBeers retrieves beers matching the provided query.
	// The query must be validated before execution.
	// Returns a non-nil slice on success.
	FindBeers(ctx context.Context, q BeerQuery) ([]models.Beer, error)
}

// BeerStatsService defines aggregation operations.
type BeerStatsService interface {
	// GetStats returns global statistics for all beers.
	GetStats(ctx context.Context) (*models.GeneralStats, error)

	// GetStatsByBrand returns statistics grouped by brand.
	GetStatsByBrand(ctx context.Context) ([]models.BrandStats, error)
}

// AnalysisServiceInterface composes all service operations.
type AnalysisServiceInterface interface {
	BeerQueryService
	BeerStatsService
}
