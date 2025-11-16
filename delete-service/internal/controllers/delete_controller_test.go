package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	deleteErr error
}

func (m *mockService) DeleteBeer(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	return nil
}

func TestDeleteController_DeleteBeer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{}
	controller := NewDeleteController(mock)

	r := gin.Default()
	r.DELETE("/beers/:id", controller.DeleteBeer)

	req, _ := http.NewRequest("DELETE", "/beers/507f1f77bcf86cd799439011", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestDeleteController_DeleteBeer_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockService{
		deleteErr: errors.New("not found"),
	}
	controller := NewDeleteController(mock)

	r := gin.Default()
	r.DELETE("/beers/:id", controller.DeleteBeer)

	req, _ := http.NewRequest("DELETE", "/beers/507f1f77bcf86cd799439011", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
