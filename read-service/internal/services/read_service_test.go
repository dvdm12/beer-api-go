package services

import (
	"errors"
	"testing"

	"readservice/internal/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepository struct {
	getBeerErr     error
	getAllBeersErr error
	beer           *models.Beer
	beers          []models.Beer
}

func (m *mockRepository) GetBeerByID(id string) (*models.Beer, error) {
	if m.getBeerErr != nil {
		return nil, m.getBeerErr
	}
	return m.beer, nil
}

func (m *mockRepository) GetAllBeers() ([]models.Beer, error) {
	if m.getAllBeersErr != nil {
		return nil, m.getAllBeersErr
	}
	return m.beers, nil
}

func TestReadService_GetBeerByID_Success(t *testing.T) {
	beer := &models.Beer{
		ID:      primitive.NewObjectID(),
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 5.0,
		Year:    2023,
	}

	mock := &mockRepository{beer: beer}
	service := NewReadService(mock)

	result, err := service.GetBeerByID("123")

	assert.Nil(t, err)
	assert.Equal(t, beer, result)
}

func TestReadService_GetBeerByID_Error(t *testing.T) {
	mock := &mockRepository{
		getBeerErr: errors.New("not found"),
	}
	service := NewReadService(mock)

	_, err := service.GetBeerByID("123")

	assert.NotNil(t, err)
	assert.Equal(t, "not found", err.Error())
}

func TestReadService_GetAllBeers_Success(t *testing.T) {
	beers := []models.Beer{
		{
			ID:      primitive.NewObjectID(),
			Name:    "Beer1",
			Brand:   "Brand1",
			Alcohol: 5.0,
			Year:    2023,
		},
		{
			ID:      primitive.NewObjectID(),
			Name:    "Beer2",
			Brand:   "Brand2",
			Alcohol: 6.0,
			Year:    2024,
		},
	}

	mock := &mockRepository{beers: beers}
	service := NewReadService(mock)

	result, err := service.GetAllBeers()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, beers, result)
}

func TestReadService_GetAllBeers_Error(t *testing.T) {
	mock := &mockRepository{
		getAllBeersErr: errors.New("database error"),
	}
	service := NewReadService(mock)

	_, err := service.GetAllBeers()

	assert.NotNil(t, err)
	assert.Equal(t, "database error", err.Error())
}
