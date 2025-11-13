package controllers

import (
	"bytes"
	"createservice/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock que implementa CreateServiceInterface
type mockService struct {
	called bool
	err    error
}

func (m *mockService) CreateBeer(beer models.Beer) error {
	m.called = true
	return m.err
}

func TestCreateController_CreateBeer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{}
	controller := NewCreateController(mock)

	r := gin.Default()
	r.POST("/beers/create", controller.CreateBeer)

	beer := models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 4.5,
		Year:    2021,
	}

	body, _ := json.Marshal(beer)

	req, _ := http.NewRequest("POST", "/beers/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.True(t, mock.called)
}
