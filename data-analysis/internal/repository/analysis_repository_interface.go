package repository

import (
	"context"
	"dataanalysis/internal/models"
	"errors"
	"fmt"
)

// SortField defines allowed fields for sorting operations.
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

// BeerFilter defines optional query parameters.
// Nil fields are ignored during query construction.
type BeerFilter struct {
	Brand      *string
	Name       *string
	NameLike   *string
	MinAlcohol *float64
	MaxAlcohol *float64
	FromYear   *int
	ToYear     *int
}

// ErrInvalidFilter indicates an invalid filter configuration.
var ErrInvalidFilter = errors.New("invalid filter")

// Validate ensures the filter is logically consistent.
func (f BeerFilter) Validate() error {
	if f.MinAlcohol != nil && *f.MinAlcohol < 0 {
		return fmt.Errorf("%w: MinAlcohol cannot be negative", ErrInvalidFilter)
	}
	if f.MaxAlcohol != nil && *f.MaxAlcohol > 100 {
		return fmt.Errorf("%w: MaxAlcohol cannot exceed 100", ErrInvalidFilter)
	}
	if f.Brand != nil && *f.Brand == "" {
		return fmt.Errorf("%w: Brand cannot be empty", ErrInvalidFilter)
	}
	if f.MinAlcohol != nil && f.MaxAlcohol != nil && *f.MinAlcohol > *f.MaxAlcohol {
		return fmt.Errorf("%w: MinAlcohol > MaxAlcohol", ErrInvalidFilter)
	}
	if f.FromYear != nil && f.ToYear != nil && *f.FromYear > *f.ToYear {
		return fmt.Errorf("%w: FromYear > ToYear", ErrInvalidFilter)
	}

	if f.Name != nil && f.NameLike != nil {
		return fmt.Errorf("%w: Name and NameLike are mutually exclusive", ErrInvalidFilter)
	}

	return nil
}

// AnalysisRepositoryInterface defines read operations for beer analysis.
type AnalysisRepositoryInterface interface {
	// GetByID returns a beer by its ID.
	GetByID(ctx context.Context, id string) (*models.Beer, error)

	// Find returns beers matching the filter.
	// Always returns a non-nil slice.
	Find(ctx context.Context, filter BeerFilter) ([]models.Beer, error)

	// FindOne returns the first match for the filter.
	// Result is non-deterministic unless filter is unique.
	FindOne(ctx context.Context, filter BeerFilter) (*models.Beer, error)

	// FindTop returns the top beer sorted by field.
	// desc=true applies descending order.
	FindTop(ctx context.Context, field SortField, desc bool) (*models.Beer, error)
}
