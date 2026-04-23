package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"readservice/internal/errors"
	"readservice/internal/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockService is a test double for the service layer.
// It records method calls and returns predefined results.
type mockService struct {
	getByIDCalled bool
	getAllCalled  bool

	getByIDResult *models.Beer
	getByIDErr    error

	getAllResult []models.Beer
	getAllErr    error
}

// GetBeerByID mocks retrieving a beer by ID.
func (m *mockService) GetBeerByID(id string) (*models.Beer, error) {
	m.getByIDCalled = true
	return m.getByIDResult, m.getByIDErr
}

// GetAllBeers mocks retrieving all beers.
func (m *mockService) GetAllBeers() ([]models.Beer, error) {
	m.getAllCalled = true
	return m.getAllResult, m.getAllErr
}

// setupRouter initializes a Gin engine with test routes.
func setupRouter(mock *mockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	controller := NewReadController(mock)
	r := gin.Default()

	r.GET("/beers/:id", controller.GetBeerByID)
	r.GET("/beers", controller.GetAllBeers)

	return r
}

// performGetByID executes a GET /beers/:id request.
func performGetByID(r *gin.Engine, id string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodGet, "/beers/"+id, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// performGetAll executes a GET /beers request.
func performGetAll(r *gin.Engine) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodGet, "/beers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// parseResponse unmarshals a JSON object response.
func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return resp
}

// validID represents a valid MongoDB ObjectID string.
const validID = "507f1f77bcf86cd799439011"

// validBeer returns a sample Beer entity.
func validBeer() *models.Beer {
	return &models.Beer{
		ID:      primitive.NewObjectID(),
		Name:    "Corona",
		Brand:   "AB InBev",
		Alcohol: 4.5,
		Year:    2021,
	}
}

// TestReadController_GetBeerByID_Success verifies successful retrieval.
func TestReadController_GetBeerByID_Success(t *testing.T) {
	mock := &mockService{getByIDResult: validBeer()}
	r := setupRouter(mock)

	w := performGetByID(r, validID)
	resp := parseResponse(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Corona", resp["name"])
	assert.True(t, mock.getByIDCalled)
}

// TestReadController_GetBeerByID_InvalidID verifies invalid ID handling.
func TestReadController_GetBeerByID_InvalidID(t *testing.T) {
	mock := &mockService{getByIDErr: errors.NewInvalidIDError("bad-id")}
	r := setupRouter(mock)

	w := performGetByID(r, "bad-id")
	resp := parseResponse(w)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errors.CodeInvalidID, resp["code"])
	assert.True(t, mock.getByIDCalled)
}

// TestReadController_GetBeerByID_NotFound verifies not found behavior.
func TestReadController_GetBeerByID_NotFound(t *testing.T) {
	mock := &mockService{getByIDErr: errors.NewBeerNotFoundError(validID)}
	r := setupRouter(mock)

	w := performGetByID(r, validID)
	resp := parseResponse(w)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, errors.CodeBeerNotFound, resp["code"])
	assert.Contains(t, resp["message"], validID)
	assert.True(t, mock.getByIDCalled)
}

// TestReadController_GetBeerByID_InternalError verifies internal error handling.
func TestReadController_GetBeerByID_InternalError(t *testing.T) {
	mock := &mockService{getByIDErr: errors.Internal(assert.AnError)}
	r := setupRouter(mock)

	w := performGetByID(r, validID)
	resp := parseResponse(w)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, errors.ErrInternal, resp["code"])
	assert.True(t, mock.getByIDCalled)
}

// TestReadController_GetAllBeers_Success verifies successful retrieval of all beers.
func TestReadController_GetAllBeers_Success(t *testing.T) {
	mock := &mockService{getAllResult: []models.Beer{*validBeer()}}
	r := setupRouter(mock)

	w := performGetAll(r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, mock.getAllCalled)

	var result []interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &result)
	assert.Len(t, result, 1)
}

// TestReadController_GetAllBeers_Empty verifies empty result handling.
func TestReadController_GetAllBeers_Empty(t *testing.T) {
	mock := &mockService{getAllResult: []models.Beer{}}
	r := setupRouter(mock)

	w := performGetAll(r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, mock.getAllCalled)

	var result []interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &result)
	assert.Empty(t, result)
}

// TestReadController_GetAllBeers_InternalError verifies error handling.
func TestReadController_GetAllBeers_InternalError(t *testing.T) {
	mock := &mockService{getAllErr: errors.Internal(assert.AnError)}
	r := setupRouter(mock)

	w := performGetAll(r)
	resp := parseResponse(w)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, errors.ErrInternal, resp["code"])
	assert.True(t, mock.getAllCalled)
}
