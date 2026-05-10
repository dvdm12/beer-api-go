// Package services implements business logic for beer analysis.
package services

import (
	"context"
	stderrors "errors"

	"dataanalysis/internal/models"
)

// ErrNoResults indicates the collection is empty.
var ErrNoResults = stderrors.New("no results found")

// BeerRankingService defines ranking operations.
type BeerRankingService interface {
	GetStrongest(ctx context.Context) (*models.Beer, error)
	GetWeakest(ctx context.Context) (*models.Beer, error)
	GetOldest(ctx context.Context) (*models.Beer, error)
	GetNewest(ctx context.Context) (*models.Beer, error)
}

// BeerStatsService defines aggregation operations.
type BeerStatsService interface {
	GetStats(ctx context.Context) (*models.GeneralStats, error)
	GetStatsByBrand(ctx context.Context) ([]models.BrandStats, error)
}

// AnalysisServiceInterface composes ranking and stats services.
type AnalysisServiceInterface interface {
	BeerRankingService
	BeerStatsService
}
