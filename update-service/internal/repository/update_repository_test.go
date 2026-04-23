package repository

import (
	"context"
	"errors" // Standard library errors
	"testing"

	apperrors "updateservice/internal/errors"
	"updateservice/internal/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mockMongoCollection simulates MongoDB operations.
type mockMongoCollection struct {
	called       bool
	err          error
	matchedCount int64
}

func (m *mockMongoCollection) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	m.called = true
	return &mongo.UpdateResult{MatchedCount: m.matchedCount}, m.err
}

func TestUpdateRepository_UpdateBeer_Success(t *testing.T) {
	// Setup mock to return 1 matched document.
	mockCollection := &mockMongoCollection{matchedCount: 1}
	repo := NewUpdateRepository(mockCollection)

	beer := models.Beer{Name: "TestBeer", Brand: "TestBrand", Alcohol: 5.5, Year: 2023}
	validID := "507f1f77bcf86cd799439011"

	// Pass context.Background() as required by the new signature.
	err := repo.UpdateBeer(context.Background(), validID, beer)

	assert.NoError(t, err)
	assert.True(t, mockCollection.called)
}

func TestUpdateRepository_UpdateBeer_InternalError(t *testing.T) {
	// Setup mock to simulate a database failure.
	mockCollection := &mockMongoCollection{
		err: errors.New("connection lost"),
	}
	repo := NewUpdateRepository(mockCollection)

	beer := models.Beer{Name: "ErrorBeer"}
	validID := "507f1f77bcf86cd799439011"

	err := repo.UpdateBeer(context.Background(), validID, beer)

	assert.Error(t, err)
	assert.True(t, mockCollection.called)

	// Cast and assert custom AppError properties.
	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok, "error should be of type AppError")
	assert.Equal(t, apperrors.ErrInternal, appErr.Code())
}

func TestUpdateRepository_UpdateBeer_NotFound(t *testing.T) {
	// Setup mock to simulate 0 matched documents (no error, just not found).
	mockCollection := &mockMongoCollection{matchedCount: 0}
	repo := NewUpdateRepository(mockCollection)

	beer := models.Beer{Name: "GhostBeer"}
	validID := "507f1f77bcf86cd799439011"

	err := repo.UpdateBeer(context.Background(), validID, beer)

	assert.Error(t, err)
	assert.True(t, mockCollection.called)

	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrNotFound, appErr.Code())
}

func TestUpdateRepository_UpdateBeer_InvalidID(t *testing.T) {
	mockCollection := &mockMongoCollection{}
	repo := NewUpdateRepository(mockCollection)

	beer := models.Beer{Name: "TestBeer"}
	invalidID := "invalid-id-format"

	err := repo.UpdateBeer(context.Background(), invalidID, beer)

	assert.Error(t, err)
	assert.False(t, mockCollection.called, "DB should not be called with an invalid ID")

	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrBadRequest, appErr.Code())
}
