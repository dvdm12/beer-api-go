package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"readservice/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockService struct {
	getBeerErr     error
	getAllBeersErr error
	beer           *models.Beer
	beers          []models.Beer
}

func (m *mockService) GetBeerByID(id string) (*models.Beer, error) {
	if m.getBeerErr != nil {
		return nil, m.getBeerErr
	}
	return m.beer, nil
}

func (m *mockService) GetAllBeers() ([]models.Beer, error) {
	if m.getAllBeersErr != nil {
		return nil, m.getAllBeersErr
	}
	return m.beers, nil
}

func TestReadController_GetBeerByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	beer := &models.Beer{
		ID:      primitive.NewObjectID(),
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 5.0,
		Year:    2023,
	}

	mock := &mockService{beer: beer}
	controller := NewReadController(mock)

	r := gin.Default()
	r.GET("/beers/:id", controller.GetBeerByID)

	req, _ := http.NewRequest("GET", "/beers/"+beer.ID.Hex(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestReadController_GetBeerByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{
		getBeerErr: errors.New("not found"),
	}
	controller := NewReadController(mock)

	r := gin.Default()
	r.GET("/beers/:id", controller.GetBeerByID)

	req, _ := http.NewRequest("GET", "/beers/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestReadController_GetAllBeers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	beers := []models.Beer{
		{
			ID:      primitive.NewObjectID(),
			Name:    "Beer1",
			Brand:   "Brand1",
			Alcohol: 5.0,
			Year:    2023,
		},
	}

	mock := &mockService{beers: beers}
	controller := NewReadController(mock)

	r := gin.Default()
	r.GET("/beers", controller.GetAllBeers)

	req, _ := http.NewRequest("GET", "/beers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestReadController_GetAllBeers_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{
		getAllBeersErr: errors.New("database error"),
	}
	controller := NewReadController(mock)

	r := gin.Default()
	r.GET("/beers", controller.GetAllBeers)

	req, _ := http.NewRequest("GET", "/beers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
}
