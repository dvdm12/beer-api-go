// Package services implements beer analysis business logic.
package services

import (
	"context"
	stderrors "errors"
	"fmt"

	apperrors "dataanalysis/internal/errors"
	"dataanalysis/internal/models"
	"dataanalysis/internal/repository"
)

// AnalysisService provides beer analysis operations.
type AnalysisService struct {
	repo repository.AnalysisRepositoryInterface
}

// NewAnalysisService returns a new AnalysisService.
func NewAnalysisService(repo repository.AnalysisRepositoryInterface) *AnalysisService {
	return &AnalysisService{repo: repo}
}

// GetStrongest returns the beer with the highest alcohol content.
func (s *AnalysisService) GetStrongest(ctx context.Context) (*models.Beer, error) {
	beer, err := s.repo.FindTop(ctx, repository.SortByAlcohol, true)
	if err != nil {
		return nil, s.mapError(err)
	}
	return beer, nil
}

// GetWeakest returns the beer with the lowest alcohol content.
func (s *AnalysisService) GetWeakest(ctx context.Context) (*models.Beer, error) {
	beer, err := s.repo.FindTop(ctx, repository.SortByAlcohol, false)
	if err != nil {
		return nil, s.mapError(err)
	}
	return beer, nil
}

// GetOldest returns the beer with the earliest year.
func (s *AnalysisService) GetOldest(ctx context.Context) (*models.Beer, error) {
	beer, err := s.repo.FindTop(ctx, repository.SortByYear, false)
	if err != nil {
		return nil, s.mapError(err)
	}
	return beer, nil
}

// GetNewest returns the beer with the latest year.
func (s *AnalysisService) GetNewest(ctx context.Context) (*models.Beer, error) {
	beer, err := s.repo.FindTop(ctx, repository.SortByYear, true)
	if err != nil {
		return nil, s.mapError(err)
	}
	return beer, nil
}

// GetStats returns global beer statistics.
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

// GetStatsByBrand returns statistics grouped by brand.
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

// computeStats calculates global statistics.
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

// computeStatsByBrand calculates statistics by brand.
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

// mapError converts repository errors.
func (s *AnalysisService) mapError(err error) error {
	if err == nil {
		return nil
	}

	var repoErr *repository.RepoError
	if stderrors.As(err, &repoErr) {
		if repoErr.Category == repository.CategoryNotFound {
			return fmt.Errorf("%w", ErrNoResults)
		}
		return apperrors.Internal(err)
	}

	return apperrors.Internal(err)
}
