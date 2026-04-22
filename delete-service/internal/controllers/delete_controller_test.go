// Package controllers contains HTTP handler tests.
package controllers

import (
	"deleteservice/internal/errors"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockService simulates service behavior.
type mockService struct {
	called    bool
	deleteErr error
}

// DeleteBeer mock implementation.
func (m *mockService) DeleteBeer(id string) error {
	m.called = true
	return m.deleteErr
}

// setupRouter configures test routes.
func setupRouter(mock *mockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	controller := NewDeleteController(mock)
	r := gin.Default()
	r.DELETE("/beers/:id", controller.DeleteBeer)
	return r
}

// performRequest executes HTTP request.
func performRequest(r *gin.Engine, id string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodDelete, "/beers/"+id, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// parseResponse decodes JSON response.
func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return resp
}

// Test successful deletion.
func TestDeleteController_DeleteBeer_Success(t *testing.T) {
	mock := &mockService{}
	r := setupRouter(mock)

	w := performRequest(r, "507f1f77bcf86cd799439011")
	resp := parseResponse(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Beer deleted successfully", resp["message"])
	assert.True(t, mock.called)
}

// Test invalid ID error.
func TestDeleteController_DeleteBeer_InvalidID(t *testing.T) {
	mock := &mockService{deleteErr: errors.NewInvalidIDError("bad-id")}
	r := setupRouter(mock)

	w := performRequest(r, "bad-id")
	resp := parseResponse(w)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errors.CodeInvalidID, resp["code"])
	assert.True(t, mock.called)
}

// Test not found error.
func TestDeleteController_DeleteBeer_NotFound(t *testing.T) {
	mock := &mockService{deleteErr: errors.NewBeerNotFoundError("507f1f77bcf86cd799439011")}
	r := setupRouter(mock)

	w := performRequest(r, "507f1f77bcf86cd799439011")
	resp := parseResponse(w)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, errors.CodeBeerNotFound, resp["code"])
	assert.Contains(t, resp["message"], "507f1f77bcf86cd799439011")
	assert.True(t, mock.called)
}

// Test internal error case.
func TestDeleteController_DeleteBeer_InternalError(t *testing.T) {
	mock := &mockService{deleteErr: errors.Internal(assert.AnError)}
	r := setupRouter(mock)

	w := performRequest(r, "507f1f77bcf86cd799439011")
	resp := parseResponse(w)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, errors.ErrInternal, resp["code"])
	assert.True(t, mock.called)
}