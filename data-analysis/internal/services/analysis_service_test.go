package services

import (
	"context"
	stderrors "errors"
	"testing"

	apperrors "dataanalysis/internal/errors"
	"dataanalysis/internal/models"
	"dataanalysis/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	findResult    []models.Beer
	findTopResult *models.Beer
	findErr       error
	findTopErr    error
}

func (m *mockRepo) Find(_ context.Context, _ repository.BeerFilter) ([]models.Beer, error) {
	return m.findResult, m.findErr
}

func (m *mockRepo) FindTop(_ context.Context, _ repository.SortField, _ bool) (*models.Beer, error) {
	return m.findTopResult, m.findTopErr
}

func beerFixtures() []models.Beer {
	return []models.Beer{
		{Name: "Heineken", Brand: "Heineken", Alcohol: 5.0, Year: 1873},
		{Name: "Duvel", Brand: "Duvel", Alcohol: 8.5, Year: 1923},
		{Name: "Hoegaarden", Brand: "InBev", Alcohol: 4.9, Year: 1966},
		{Name: "Chimay", Brand: "Chimay", Alcohol: 9.0, Year: 1862},
	}
}

func repoNotFoundErr() error {
	return &repository.RepoError{Category: repository.CategoryNotFound}
}

func repoTimeoutErr() error {
	return &repository.RepoError{Category: repository.CategoryTimeout}
}

func repoNetworkErr() error {
	return &repository.RepoError{Category: repository.CategoryNetwork}
}

func TestGetStrongest_Success(t *testing.T) {
	beer := &models.Beer{Name: "Chimay", Alcohol: 9.0}
	svc := NewAnalysisService(&mockRepo{findTopResult: beer})

	result, err := svc.GetStrongest(context.Background())

	require.NoError(t, err)
	assert.Equal(t, beer, result)
}

func TestGetStrongest_EmptyCollection(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findTopErr: repoNotFoundErr()})

	result, err := svc.GetStrongest(context.Background())

	assert.Nil(t, result)
	assert.True(t, stderrors.Is(err, ErrNoResults))
}

func TestGetStrongest_Timeout(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findTopErr: repoTimeoutErr()})

	_, err := svc.GetStrongest(context.Background())

	var appErr apperrors.AppError
	assert.True(t, stderrors.As(err, &appErr))
	assert.Equal(t, 500, appErr.StatusCode())
}

func TestGetStrongest_NetworkError(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findTopErr: repoNetworkErr()})

	_, err := svc.GetStrongest(context.Background())

	var appErr apperrors.AppError
	assert.True(t, stderrors.As(err, &appErr))
	assert.Equal(t, 500, appErr.StatusCode())
}

func TestGetWeakest_Success(t *testing.T) {
	beer := &models.Beer{Name: "Hoegaarden", Alcohol: 4.9}
	svc := NewAnalysisService(&mockRepo{findTopResult: beer})

	result, err := svc.GetWeakest(context.Background())

	require.NoError(t, err)
	assert.Equal(t, beer, result)
}

func TestGetWeakest_EmptyCollection(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findTopErr: repoNotFoundErr()})

	result, err := svc.GetWeakest(context.Background())

	assert.Nil(t, result)
	assert.True(t, stderrors.Is(err, ErrNoResults))
}

func TestGetWeakest_Timeout(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findTopErr: repoTimeoutErr()})

	_, err := svc.GetWeakest(context.Background())

	var appErr apperrors.AppError
	assert.True(t, stderrors.As(err, &appErr))
	assert.Equal(t, 500, appErr.StatusCode())
}

func TestGetOldest_Success(t *testing.T) {
	beer := &models.Beer{Name: "Chimay", Year: 1862}
	svc := NewAnalysisService(&mockRepo{findTopResult: beer})

	result, err := svc.GetOldest(context.Background())

	require.NoError(t, err)
	assert.Equal(t, beer, result)
}

func TestGetOldest_EmptyCollection(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findTopErr: repoNotFoundErr()})

	result, err := svc.GetOldest(context.Background())

	assert.Nil(t, result)
	assert.True(t, stderrors.Is(err, ErrNoResults))
}

func TestGetOldest_NetworkError(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findTopErr: repoNetworkErr()})

	_, err := svc.GetOldest(context.Background())

	var appErr apperrors.AppError
	assert.True(t, stderrors.As(err, &appErr))
	assert.Equal(t, 500, appErr.StatusCode())
}

func TestGetNewest_Success(t *testing.T) {
	beer := &models.Beer{Name: "Hoegaarden", Year: 1966}
	svc := NewAnalysisService(&mockRepo{findTopResult: beer})

	result, err := svc.GetNewest(context.Background())

	require.NoError(t, err)
	assert.Equal(t, beer, result)
}

func TestGetNewest_EmptyCollection(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findTopErr: repoNotFoundErr()})

	result, err := svc.GetNewest(context.Background())

	assert.Nil(t, result)
	assert.True(t, stderrors.Is(err, ErrNoResults))
}

func TestGetNewest_Timeout(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findTopErr: repoTimeoutErr()})

	_, err := svc.GetNewest(context.Background())

	var appErr apperrors.AppError
	assert.True(t, stderrors.As(err, &appErr))
	assert.Equal(t, 500, appErr.StatusCode())
}

func TestGetStats_Success(t *testing.T) {
	beers := beerFixtures()
	svc := NewAnalysisService(&mockRepo{findResult: beers})

	result, err := svc.GetStats(context.Background())

	require.NoError(t, err)
	assert.Equal(t, int64(4), result.Total)
	assert.Equal(t, int64(4), result.TotalBrands)
	assert.InDelta(t, 6.85, result.AvgAlcohol, 0.01)
	assert.Equal(t, 4.9, result.MinAlcohol)
	assert.Equal(t, 9.0, result.MaxAlcohol)
	assert.Equal(t, 1862, result.OldestYear)
	assert.Equal(t, 1966, result.NewestYear)
}

func TestGetStats_EmptyCollection(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findResult: []models.Beer{}})

	result, err := svc.GetStats(context.Background())

	assert.Nil(t, result)
	assert.True(t, stderrors.Is(err, ErrNoResults))
}

func TestGetStats_RepositoryError(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findErr: repoTimeoutErr()})

	_, err := svc.GetStats(context.Background())

	var appErr apperrors.AppError
	assert.True(t, stderrors.As(err, &appErr))
	assert.Equal(t, 500, appErr.StatusCode())
}

func TestGetStats_NetworkError(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findErr: repoNetworkErr()})

	_, err := svc.GetStats(context.Background())

	var appErr apperrors.AppError
	assert.True(t, stderrors.As(err, &appErr))
	assert.Equal(t, 500, appErr.StatusCode())
}

func TestGetStats_SingleBeer(t *testing.T) {
	beers := []models.Beer{{Name: "Duvel", Brand: "Duvel", Alcohol: 8.5, Year: 1923}}
	svc := NewAnalysisService(&mockRepo{findResult: beers})

	result, err := svc.GetStats(context.Background())

	require.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
	assert.Equal(t, int64(1), result.TotalBrands)
	assert.Equal(t, 8.5, result.AvgAlcohol)
	assert.Equal(t, 8.5, result.MinAlcohol)
	assert.Equal(t, 8.5, result.MaxAlcohol)
	assert.Equal(t, 1923, result.OldestYear)
	assert.Equal(t, 1923, result.NewestYear)
}

func TestGetStatsByBrand_Success(t *testing.T) {
	beers := []models.Beer{
		{Brand: "Chimay", Alcohol: 7.0, Year: 1862},
		{Brand: "Chimay", Alcohol: 8.0, Year: 1900},
		{Brand: "Duvel", Alcohol: 8.5, Year: 1923},
	}
	svc := NewAnalysisService(&mockRepo{findResult: beers})

	result, err := svc.GetStatsByBrand(context.Background())

	require.NoError(t, err)
	assert.Len(t, result, 2)

	brandMap := make(map[string]models.BrandStats)
	for _, s := range result {
		brandMap[s.Brand] = s
	}

	chimay := brandMap["Chimay"]
	assert.Equal(t, int64(2), chimay.Count)
	assert.Equal(t, 7.0, chimay.MinAlcohol)
	assert.Equal(t, 8.0, chimay.MaxAlcohol)
	assert.InDelta(t, 7.5, chimay.AvgAlcohol, 0.01)

	duvel := brandMap["Duvel"]
	assert.Equal(t, int64(1), duvel.Count)
	assert.Equal(t, 8.5, duvel.MinAlcohol)
	assert.Equal(t, 8.5, duvel.MaxAlcohol)
}

func TestGetStatsByBrand_EmptyCollection(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findResult: []models.Beer{}})

	result, err := svc.GetStatsByBrand(context.Background())

	assert.Nil(t, result)
	assert.True(t, stderrors.Is(err, ErrNoResults))
}

func TestGetStatsByBrand_RepositoryError(t *testing.T) {
	svc := NewAnalysisService(&mockRepo{findErr: repoTimeoutErr()})

	_, err := svc.GetStatsByBrand(context.Background())

	var appErr apperrors.AppError
	assert.True(t, stderrors.As(err, &appErr))
	assert.Equal(t, 500, appErr.StatusCode())
}

func TestGetStatsByBrand_SingleBrand(t *testing.T) {
	beers := []models.Beer{
		{Brand: "Duvel", Alcohol: 8.5, Year: 1923},
		{Brand: "Duvel", Alcohol: 6.0, Year: 1950},
	}
	svc := NewAnalysisService(&mockRepo{findResult: beers})

	result, err := svc.GetStatsByBrand(context.Background())

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Duvel", result[0].Brand)
	assert.Equal(t, int64(2), result[0].Count)
	assert.Equal(t, 6.0, result[0].MinAlcohol)
	assert.Equal(t, 8.5, result[0].MaxAlcohol)
	assert.InDelta(t, 7.25, result[0].AvgAlcohol, 0.01)
}
