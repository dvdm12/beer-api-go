package repository

import (
	"context"
	"dataanalysis/internal/models"
	"errors"
	"fmt"
)

// SortField defines allowed sort fields.
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

// ErrInvalidFilter indicates a structurally incoherent filter.
var ErrInvalidFilter = errors.New("invalid filter")

// Validate checks structural coherence of the filter.
// Semantic rules (bounds, required fields) are enforced at the service layer.
func (f BeerFilter) Validate() error {
	if f.MinAlcohol != nil && f.MaxAlcohol != nil && *f.MinAlcohol > *f.MaxAlcohol {
		return fmt.Errorf("%w: MinAlcohol cannot exceed MaxAlcohol", ErrInvalidFilter)
	}
	if f.FromYear != nil && f.ToYear != nil && *f.FromYear > *f.ToYear {
		return fmt.Errorf("%w: FromYear cannot exceed ToYear", ErrInvalidFilter)
	}
	if f.Name != nil && f.NameLike != nil {
		return fmt.Errorf("%w: Name and NameLike are mutually exclusive", ErrInvalidFilter)
	}
	return nil
}

// AnalysisRepositoryInterface defines raw data access operations.
type AnalysisRepositoryInterface interface {
	// Find returns beers matching the filter. Always returns a non-nil slice.
	Find(ctx context.Context, filter BeerFilter) ([]models.Beer, error)

	// FindTop returns the top-ranked beer by field.
	// desc=true for descending, desc=false for ascending.
	FindTop(ctx context.Context, field SortField, desc bool) (*models.Beer, error)
}
