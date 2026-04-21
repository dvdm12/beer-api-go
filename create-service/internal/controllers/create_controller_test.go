package controllers

import (
	"bytes"
	"createservice/internal/errors"
	"createservice/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	called bool
	err    error
}

func (m *mockService) CreateBeer(beer models.Beer) error {
	m.called = true
	return m.err
}

func setupRouter(mock *mockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	controller := NewCreateController(mock)
	r := gin.Default()
	r.POST("/beers", controller.CreateBeer) // ← ruta correcta
	return r
}

func performRequest(r *gin.Engine, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(http.MethodPost, "/beers", &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performRawRequest(r *gin.Engine, rawBody string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, "/beers", bytes.NewBufferString(rawBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return resp
}

func validBeer() models.Beer {
	return models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 4.5,
		Year:    2021,
	}
}

func TestCreateController_CreateBeer_Success(t *testing.T) {
	mock := &mockService{}
	r := setupRouter(mock)

	w := performRequest(r, validBeer())
	resp := parseResponse(w)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.True(t, mock.called)
	assert.Equal(t, "Beer has been created successfully", resp["message"]) // ← mensaje correcto
}

func TestCreateController_CreateBeer_InvalidJSON(t *testing.T) {
	mock := &mockService{}
	r := setupRouter(mock)

	w := performRawRequest(r, `{invalid json}`)
	resp := parseResponse(w)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, errors.CodeInvalidInput, resp["code"])
	assert.False(t, mock.called)
}

func TestCreateController_CreateBeer_EmptyBody(t *testing.T) {
	mock := &mockService{}
	r := setupRouter(mock)

	w := performRawRequest(r, ``)
	resp := parseResponse(w)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, errors.CodeInvalidInput, resp["code"])
	assert.False(t, mock.called)
}

func TestCreateController_CreateBeer_Duplicate(t *testing.T) {
	mock := &mockService{err: errors.NewDuplicateBeerError("beer 'TestBeer' already exists")}
	r := setupRouter(mock)

	w := performRequest(r, validBeer())
	resp := parseResponse(w)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, errors.CodeDuplicateBeer, resp["code"])
	assert.Equal(t, "beer 'TestBeer' already exists", resp["message"])
	assert.True(t, mock.called)
}

func TestCreateController_CreateBeer_ValidationError(t *testing.T) {
	mock := &mockService{err: errors.NewValidationError("alcohol must be between 0 and 100")}
	r := setupRouter(mock)

	w := performRequest(r, validBeer())
	resp := parseResponse(w)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, errors.CodeInvalidInput, resp["code"])
	assert.Equal(t, "alcohol must be between 0 and 100", resp["message"])
	assert.True(t, mock.called)
}

func TestCreateController_CreateBeer_InternalError(t *testing.T) {
	mock := &mockService{err: errors.Internal(assert.AnError)}
	r := setupRouter(mock)

	w := performRequest(r, validBeer())
	resp := parseResponse(w)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, errors.ErrInternal, resp["code"])
	assert.True(t, mock.called)
}

func TestCreateController_CreateBeer_ServiceNotCalledOnBadPayload(t *testing.T) {
	mock := &mockService{}
	r := setupRouter(mock)

	w := performRawRequest(r, `{invalid}`)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.False(t, mock.called, "service should not be called when binding fails")
}