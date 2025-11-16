package repository

import (
	"context"
	"updateservice/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockMongoCollection struct {
	called bool
	err    error
}

func (m *mockMongoCollection) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	m.called = true
	return &mongo.UpdateResult{}, m.err
}

func TestUpdateRepository_UpdateBeer_Success(t *testing.T) {
	mockCollection := &mockMongoCollection{}
	repo := NewUpdateRepository(mockCollection)

	beer := models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 5.5,
		Year:    2023,
	}

	validID := "507f1f77bcf86cd799439011"
	_, err := repo.UpdateBeer(validID, beer)

	assert.Nil(t, err)
	assert.True(t, mockCollection.called)
}

func TestUpdateRepository_UpdateBeer_Error(t *testing.T) {
	mockCollection := &mockMongoCollection{
		err: errors.New("update failed"),
	}

	repo := NewUpdateRepository(mockCollection)

	beer := models.Beer{
		Name:    "ErrorBeer",
		Brand:   "BrandX",
		Alcohol: 3.0,
		Year:    2020,
	}

	validID := "507f1f77bcf86cd799439011"
	_, err := repo.UpdateBeer(validID, beer)

	assert.NotNil(t, err)
	assert.Equal(t, "update failed", err.Error())
	assert.True(t, mockCollection.called)
}

func TestUpdateRepository_UpdateBeer_InvalidID(t *testing.T) {
	mockCollection := &mockMongoCollection{}
	repo := NewUpdateRepository(mockCollection)

	beer := models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 5.5,
		Year:    2023,
	}

	invalidID := "invalid-id"
	_, err := repo.UpdateBeer(invalidID, beer)

	assert.NotNil(t, err)
	assert.False(t, mockCollection.called)
}
