package repository

import (
	"context"
	"createservice/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMongoCollection struct {
	called bool
	err    error
}

func (m *mockMongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	m.called = true
	return nil, m.err
}


func TestCreateRepository_CreateBeer_Success(t *testing.T) {
	mockCollection := &mockMongoCollection{}
	repo := NewCreateRepository(mockCollection)

	beer := models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 5.5,
		Year:    2023,
	}

	_, err := repo.CreateBeer(beer)

	assert.Nil(t, err)
	assert.True(t, mockCollection.called)
}

func TestCreateRepository_CreateBeer_Error(t *testing.T) {
	mockCollection := &mockMongoCollection{
		err: errors.New("insert failed"),
	}

	repo := NewCreateRepository(mockCollection)

	beer := models.Beer{
		Name:    "ErrorBeer",
		Brand:   "BrandX",
		Alcohol: 3.0,
		Year:    2020,
	}

	_, err := repo.CreateBeer(beer)

	assert.NotNil(t, err)
	assert.Equal(t, "insert failed", err.Error())
	assert.True(t, mockCollection.called)
}
