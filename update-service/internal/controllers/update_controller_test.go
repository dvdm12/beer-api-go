package controllers

import (
	"bytes"
	"updateservice/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock que implementa UpdateServiceInterface
type mockService struct {
	called bool
	err    error
}

func (m *mockService) UpdateBeer(id string, beer models.Beer) error {
	m.called = true
	return m.err
}

func TestUpdateController_UpdateBeer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{}
	controller := NewUpdateController(mock)

	r := gin.Default()
	r.PUT("/beers/:id", controller.UpdateBeer)

	beer := models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 4.5,
		Year:    2021,
	}

	body, _ := json.Marshal(beer)

	req, _ := http.NewRequest("PUT", "/beers/507f1f77bcf86cd799439011", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.True(t, mock.called)
}

func TestUpdateController_UpdateBeer_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{}
	controller := NewUpdateController(mock)

	r := gin.Default()
	r.PUT("/beers/:id", controller.UpdateBeer)

	beer := models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 4.5,
		Year:    2021,
	}

	body, _ := json.Marshal(beer)

	req, _ := http.NewRequest("PUT", "/beers/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.False(t, mock.called)
}
