package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors" // Standard library for mock errors
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "updateservice/internal/errors" // Aliased to avoid conflicts
	"updateservice/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockService implements the UpdateServiceInterface contract.
type mockService struct {
	called bool
	err    error
}

// UpdateBeer simulates the service layer execution.
func (m *mockService) UpdateBeer(ctx context.Context, id string, beer models.Beer) error {
	m.called = true
	return m.err
}

// setupRouter initializes the Gin engine and binds the controller for testing.
func setupRouter(mock *mockService) (*gin.Engine, *UpdateController) {
	gin.SetMode(gin.TestMode)
	controller := NewUpdateController(mock)
	r := gin.Default()

	// Register both routes to test standard requests and the missing ID edge case.
	r.PUT("/beers/:id", controller.UpdateBeer)
	r.PUT("/beers", controller.UpdateBeer)

	return r, controller
}

func TestUpdateController_UpdateBeer_Success(t *testing.T) {
	mock := &mockService{}
	r, _ := setupRouter(mock)

	beer := models.Beer{Name: "TestBeer", Brand: "TestBrand", Alcohol: 4.5, Year: 2021}
	body, _ := json.Marshal(beer)

	req, _ := http.NewRequest(http.MethodPut, "/beers/507f1f77bcf86cd799439011", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, mock.called)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Beer updated successfully", response["message"])
}

func TestUpdateController_UpdateBeer_MissingID(t *testing.T) {
	mock := &mockService{}
	r, _ := setupRouter(mock)

	beer := models.Beer{Name: "TestBeer", Brand: "TestBrand"}
	body, _ := json.Marshal(beer)

	// Sending request to the root path without an ID.
	req, _ := http.NewRequest(http.MethodPut, "/beers", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, mock.called, "Service should not be called if ID is missing")

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Beer ID parameter is required", response["error"])
	assert.Equal(t, apperrors.ErrBadRequest, response["code"])
}

func TestUpdateController_UpdateBeer_InvalidJSON(t *testing.T) {
	mock := &mockService{}
	r, _ := setupRouter(mock)

	// Sending malformed JSON.
	req, _ := http.NewRequest(http.MethodPut, "/beers/123", bytes.NewBuffer([]byte("{invalid-json: true}")))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, mock.called)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, apperrors.ErrBadRequest, response["code"])
	assert.Contains(t, response["error"], "Invalid JSON format")
}

func TestUpdateController_UpdateBeer_ServiceValidationError(t *testing.T) {
	// Simulate a domain validation error (e.g., negative alcohol).
	mockErr := apperrors.ValidationError("Alcohol content cannot be negative", models.ErrCodeInvalidAlcohol)
	mock := &mockService{err: mockErr}
	r, _ := setupRouter(mock)

	beer := models.Beer{Name: "TestBeer", Brand: "Brand", Alcohol: -1.0, Year: 2021}
	body, _ := json.Marshal(beer)

	req, _ := http.NewRequest(http.MethodPut, "/beers/123", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // 422 HTTP Code
	assert.True(t, mock.called)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Alcohol content cannot be negative", response["error"])
	assert.Equal(t, models.ErrCodeInvalidAlcohol, response["code"])
}

func TestUpdateController_UpdateBeer_InternalError(t *testing.T) {
	// Simulate an unhandled database error.
	mock := &mockService{err: errors.New("database connection timeout")}
	r, _ := setupRouter(mock)

	beer := models.Beer{Name: "TestBeer", Brand: "Brand", Alcohol: 5.0, Year: 2021}
	body, _ := json.Marshal(beer)

	req, _ := http.NewRequest(http.MethodPut, "/beers/123", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code) // 500 HTTP Code
	assert.True(t, mock.called)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, apperrors.ErrInternal, response["code"])
	assert.Contains(t, response["error"], "database connection timeout")
}
