package services

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "updateservice/internal/errors"
	"updateservice/internal/models"

	"github.com/stretchr/testify/assert"
)

// mockUpdateRepo simulates the repository layer.
type mockUpdateRepo struct {
	called bool
	err    error
}

// UpdateBeer implements the updated UpdateRepositoryInterface.
func (m *mockUpdateRepo) UpdateBeer(ctx context.Context, id string, beer models.Beer) error {
	m.called = true
	return m.err
}

func TestUpdateService_UpdateBeer_Success(t *testing.T) {
	mock := &mockUpdateRepo{}
	service := NewUpdateService(mock)

	beer := models.Beer{Name: "TestBeer", Brand: "TestBrand", Alcohol: 4.8, Year: 2022}

	err := service.UpdateBeer(context.Background(), "507f1f77bcf86cd799439011", beer)

	assert.NoError(t, err)
	assert.True(t, mock.called)
}

func TestUpdateService_UpdateBeer_MissingData(t *testing.T) {
	mock := &mockUpdateRepo{}
	service := NewUpdateService(mock)

	// Missing Brand
	beer := models.Beer{Name: "TestBeer", Alcohol: 4.8, Year: 2022}

	err := service.UpdateBeer(context.Background(), "507f1f77bcf86cd799439011", beer)

	assert.Error(t, err)
	assert.False(t, mock.called, "Repository should not be called if validation fails")

	// Verify custom domain error code
	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeMissingData, appErr.Code())
}

func TestUpdateService_UpdateBeer_InvalidAlcohol(t *testing.T) {
	mock := &mockUpdateRepo{}
	service := NewUpdateService(mock)

	// Negative alcohol content
	beer := models.Beer{Name: "TestBeer", Brand: "TestBrand", Alcohol: -1.5, Year: 2022}

	err := service.UpdateBeer(context.Background(), "123", beer)

	assert.Error(t, err)
	assert.False(t, mock.called)

	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInvalidAlcohol, appErr.Code())
}

func TestUpdateService_UpdateBeer_InvalidYear(t *testing.T) {
	mock := &mockUpdateRepo{}
	service := NewUpdateService(mock)

	// Future year
	futureYear := time.Now().Year() + 1
	beer := models.Beer{Name: "TestBeer", Brand: "TestBrand", Alcohol: 4.5, Year: futureYear}

	err := service.UpdateBeer(context.Background(), "123", beer)

	assert.Error(t, err)
	assert.False(t, mock.called)

	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInvalidYear, appErr.Code())
}

func TestUpdateService_UpdateBeer_RepositoryError(t *testing.T) {
	// Simulate a database failure
	mock := &mockUpdateRepo{err: errors.New("database timeout")}
	service := NewUpdateService(mock)

	beer := models.Beer{Name: "TestBeer", Brand: "TestBrand", Alcohol: 4.5, Year: 2020}

	err := service.UpdateBeer(context.Background(), "123", beer)

	assert.Error(t, err)
	assert.True(t, mock.called)
	assert.Equal(t, "database timeout", err.Error())
}
