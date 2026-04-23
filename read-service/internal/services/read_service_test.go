package services

import (
	"errors"
	"readservice/internal/models"
	"readservice/internal/repository"
	"testing"

	apperrors "readservice/internal/errors"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockRepository is a test double that implements the repository interface.
// It records calls and returns predefined results.
type mockRepository struct {
	getByIDCalled bool
	getAllCalled  bool

	getByIDResult *models.Beer
	getByIDErr    error

	getAllResult []models.Beer
	getAllErr    error
}

// GetBeerByID mocks fetching a beer by ID.
func (m *mockRepository) GetBeerByID(id string) (*models.Beer, error) {
	m.getByIDCalled = true
	return m.getByIDResult, m.getByIDErr
}

// GetAllBeers mocks fetching all beers.
func (m *mockRepository) GetAllBeers() ([]models.Beer, error) {
	m.getAllCalled = true
	return m.getAllResult, m.getAllErr
}

// validID represents a valid MongoDB ObjectID string used in tests.
const validID = "507f1f77bcf86cd799439011"

// validBeer returns a sample valid Beer entity.
func validBeer() *models.Beer {
	return &models.Beer{
		ID:      primitive.NewObjectID(),
		Name:    "Corona",
		Brand:   "AB InBev",
		Alcohol: 4.5,
		Year:    2021,
	}
}

// TestReadService_GetBeerByID_Success verifies successful retrieval by ID.
func TestReadService_GetBeerByID_Success(t *testing.T) {
	mock := &mockRepository{getByIDResult: validBeer()}
	service := NewReadService(mock)

	beer, err := service.GetBeerByID(validID)

	assert.Nil(t, err)
	assert.NotNil(t, beer)
	assert.Equal(t, "Corona", beer.Name)
	assert.True(t, mock.getByIDCalled)
}

// TestReadService_GetBeerByID_EmptyID verifies validation for empty ID.
func TestReadService_GetBeerByID_EmptyID(t *testing.T) {
	mock := &mockRepository{}
	service := NewReadService(mock)

	beer, err := service.GetBeerByID("")

	assert.Nil(t, beer)
	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.CodeInvalidID, appErr.Code())
	assert.Equal(t, 400, appErr.StatusCode())
	assert.False(t, mock.getByIDCalled)
}

// TestReadService_GetBeerByID_InvalidIDFormat verifies handling of invalid ID format.
func TestReadService_GetBeerByID_InvalidIDFormat(t *testing.T) {
	mock := &mockRepository{getByIDErr: repository.ErrInvalidID}
	service := NewReadService(mock)

	beer, err := service.GetBeerByID("bad-id")

	assert.Nil(t, beer)
	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.CodeInvalidID, appErr.Code())
	assert.Equal(t, 400, appErr.StatusCode())
	assert.True(t, mock.getByIDCalled)
}

// TestReadService_GetBeerByID_NotFound verifies behavior when beer is not found.
func TestReadService_GetBeerByID_NotFound(t *testing.T) {
	mock := &mockRepository{getByIDErr: repository.ErrBeerNotFound}
	service := NewReadService(mock)

	beer, err := service.GetBeerByID(validID)

	assert.Nil(t, beer)
	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.CodeBeerNotFound, appErr.Code())
	assert.Equal(t, 404, appErr.StatusCode())
	assert.Contains(t, appErr.Error(), validID)
	assert.True(t, mock.getByIDCalled)
}

// TestReadService_GetBeerByID_InternalError verifies handling of unexpected errors.
func TestReadService_GetBeerByID_InternalError(t *testing.T) {
	mock := &mockRepository{getByIDErr: errors.New("mongo connection lost")}
	service := NewReadService(mock)

	beer, err := service.GetBeerByID(validID)

	assert.Nil(t, beer)
	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrInternal, appErr.Code())
	assert.Equal(t, 500, appErr.StatusCode())
	assert.True(t, mock.getByIDCalled)
}

// TestReadService_GetAllBeers_Success verifies successful retrieval of all beers.
func TestReadService_GetAllBeers_Success(t *testing.T) {
	beers := []models.Beer{
		*validBeer(),
		{Name: "Heineken", Brand: "Heineken", Alcohol: 5, Year: 1873},
	}
	mock := &mockRepository{getAllResult: beers}
	service := NewReadService(mock)

	result, err := service.GetAllBeers()

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.True(t, mock.getAllCalled)
}

// TestReadService_GetAllBeers_Empty verifies behavior when no beers are found.
func TestReadService_GetAllBeers_Empty(t *testing.T) {
	mock := &mockRepository{getAllResult: []models.Beer{}}
	service := NewReadService(mock)

	result, err := service.GetAllBeers()

	assert.Nil(t, err)
	assert.Empty(t, result)
	assert.True(t, mock.getAllCalled)
}

// TestReadService_GetAllBeers_NilResult ensures nil results are normalized to empty slices.
func TestReadService_GetAllBeers_NilResult(t *testing.T) {
	mock := &mockRepository{getAllResult: nil}
	service := NewReadService(mock)

	result, err := service.GetAllBeers()

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// TestReadService_GetAllBeers_InternalError verifies handling of repository errors.
func TestReadService_GetAllBeers_InternalError(t *testing.T) {
	mock := &mockRepository{getAllErr: errors.New("mongo find failed")}
	service := NewReadService(mock)

	result, err := service.GetAllBeers()

	assert.Nil(t, result)
	appErr, ok := err.(apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrInternal, appErr.Code())
	assert.Equal(t, 500, appErr.StatusCode())
}
