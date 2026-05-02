// Package services implements business logic for beer analysis.
package services

import (
	"context"
	"fmt"

	"dataanalysis/internal/errors"
	"dataanalysis/internal/models"
	"dataanalysis/internal/repository"
)

// AnalysisService implements AnalysisServiceInterface.
type AnalysisService struct {
	repo repository.AnalysisRepositoryInterface
}

// NewAnalysisService creates a new AnalysisService.
func NewAnalysisService(repo repository.AnalysisRepositoryInterface) *AnalysisService {
	return &AnalysisService{repo: repo}
}

// GetByID retrieves a beer by its identifier.
func (s *AnalysisService) GetByID(ctx context.Context, id string) (*models.Beer, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id cannot be empty", ErrInvalidBeerQuery)
	}

	beer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, s.mapError(err)
	}
	return beer, nil
}

// FindBeers retrieves beers matching the query.
func (s *AnalysisService) FindBeers(ctx context.Context, q BeerQuery) ([]models.Beer, error) {
	normalized, err := q.Validate()
	if err != nil {
		return nil, err
	}

	beers, err := s.repo.Find(ctx, s.translateFilter(normalized.Filter))
	if err != nil {
		return nil, s.mapError(err)
	}
	if len(beers) == 0 {
		return nil, fmt.Errorf("%w", ErrNoResults)
	}
	return beers, nil
}

// GetStats computes global statistics from all beers.
func (s *AnalysisService) GetStats(ctx context.Context) (*models.GeneralStats, error) {
	beers, err := s.repo.Find(ctx, repository.BeerFilter{})
	if err != nil {
		return nil, s.mapError(err)
	}
	if len(beers) == 0 {
		return nil, fmt.Errorf("%w", ErrNoResults)
	}
	return computeStats(beers), nil
}

// GetStatsByBrand computes statistics grouped by brand.
func (s *AnalysisService) GetStatsByBrand(ctx context.Context) ([]models.BrandStats, error) {
	beers, err := s.repo.Find(ctx, repository.BeerFilter{})
	if err != nil {
		return nil, s.mapError(err)
	}
	if len(beers) == 0 {
		return nil, fmt.Errorf("%w", ErrNoResults)
	}
	return computeStatsByBrand(beers), nil
}

// computeStats calculates aggregate statistics from beers.
func computeStats(beers []models.Beer) *models.GeneralStats {
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

	return stats
}

// computeStatsByBrand calculates statistics per brand.
func computeStatsByBrand(beers []models.Beer) []models.BrandStats {
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

	return result
}

// translateFilter maps service filter to repository filter.
func (s *AnalysisService) translateFilter(f BeerFilter) repository.BeerFilter {
	return repository.BeerFilter{
		Brand:      f.Brand,
		Name:       f.Name,
		NameLike:   f.NameLike,
		MinAlcohol: f.MinAlcohol,
		MaxAlcohol: f.MaxAlcohol,
		FromYear:   f.FromYear,
		ToYear:     f.ToYear,
	}
}

// mapError converts repository errors to service-level errors.
func (s *AnalysisService) mapError(err error) error {
	if err == nil {
		return nil
	}

	switch err {
	case repository.ErrNoBeerFound:
		return fmt.Errorf("%w", ErrBeerNotFound)
	case repository.ErrEmptyCollection:
		return fmt.Errorf("%w", ErrNoResults)
	case repository.ErrInvalidID:
		return fmt.Errorf("%w: invalid id format", ErrInvalidBeerQuery)
	}

	return errors.Internal(err)
}
