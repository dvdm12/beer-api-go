package services

import (
	"createservice/internal/errors"
	"createservice/internal/models"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	createCalled bool
	existsCalled bool
	createErr    error
	existsResult bool
	existsErr    error
}

func (m *mockRepo) CreateBeer(beer models.Beer) (*struct{}, error) {
	m.createCalled = true
	return nil, m.createErr
}

func (m *mockRepo) ExistsBeer(name string) (bool, error) {
	m.existsCalled = true
	return m.existsResult, m.existsErr
}

func validBeer() models.Beer {
	return models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 4.8,
		Year:    2022,
	}
}

func TestCreateService_CreateBeer_Success(t *testing.T) {
	mock := &mockRepo{}
	service := NewCreateService(mock)

	err := service.CreateBeer(validBeer())

	assert.Nil(t, err)
	assert.True(t, mock.existsCalled)
	assert.True(t, mock.createCalled)
}


func TestCreateService_CreateBeer_EmptyName(t *testing.T) {
	mock := &mockRepo{}
	service := NewCreateService(mock)

	beer := validBeer()
	beer.Name = ""

	err := service.CreateBeer(beer)

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeInvalidInput, appErr.Code())
	assert.Equal(t, 422, appErr.StatusCode())
	assert.False(t, mock.existsCalled)
	assert.False(t, mock.createCalled)
}

func TestCreateService_CreateBeer_EmptyBrand(t *testing.T) {
	mock := &mockRepo{}
	service := NewCreateService(mock)

	beer := validBeer()
	beer.Brand = ""

	err := service.CreateBeer(beer)

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeInvalidInput, appErr.Code())
	assert.False(t, mock.createCalled)
}

func TestCreateService_CreateBeer_AlcoholOutOfRange(t *testing.T) {
	mock := &mockRepo{}
	service := NewCreateService(mock)

	beer := validBeer()
	beer.Alcohol = 999

	err := service.CreateBeer(beer)

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeInvalidInput, appErr.Code())
	assert.Contains(t, appErr.Error(), "alcohol must be between 0 and 100")
	assert.False(t, mock.createCalled)
}

func TestCreateService_CreateBeer_NegativeAlcohol(t *testing.T) {
	mock := &mockRepo{}
	service := NewCreateService(mock)

	beer := validBeer()
	beer.Alcohol = -1

	err := service.CreateBeer(beer)

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeInvalidInput, appErr.Code())
	assert.False(t, mock.createCalled)
}

func TestCreateService_CreateBeer_YearTooOld(t *testing.T) {
	mock := &mockRepo{}
	service := NewCreateService(mock)

	beer := validBeer()
	beer.Year = 1799

	err := service.CreateBeer(beer)

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeInvalidInput, appErr.Code())
	assert.Contains(t, appErr.Error(), "year must be between 1800")
	assert.False(t, mock.createCalled)
}

func TestCreateService_CreateBeer_YearInFuture(t *testing.T) {
	mock := &mockRepo{}
	service := NewCreateService(mock)

	beer := validBeer()
	beer.Year = time.Now().Year() + 1

	err := service.CreateBeer(beer)

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeInvalidInput, appErr.Code())
	assert.Contains(t, appErr.Error(), fmt.Sprintf("year must be between 1800 and %d", time.Now().Year()))
	assert.False(t, mock.createCalled)
}

func TestCreateService_CreateBeer_Duplicate(t *testing.T) {
	mock := &mockRepo{existsResult: true}
	service := NewCreateService(mock)

	err := service.CreateBeer(validBeer())

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeDuplicateBeer, appErr.Code())
	assert.Equal(t, 409, appErr.StatusCode())
	assert.Contains(t, appErr.Error(), "TestBeer")
	assert.True(t, mock.existsCalled)
	assert.False(t, mock.createCalled)
}

func TestCreateService_CreateBeer_ExistsDBError(t *testing.T) {
	mock := &mockRepo{existsErr: fmt.Errorf("mongo connection lost")}
	service := NewCreateService(mock)

	err := service.CreateBeer(validBeer())

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.ErrInternal, appErr.Code())
	assert.Equal(t, 500, appErr.StatusCode())
	assert.False(t, mock.createCalled)
}

func TestCreateService_CreateBeer_InsertDBError(t *testing.T) {
	mock := &mockRepo{createErr: fmt.Errorf("mongo write failed")}
	service := NewCreateService(mock)

	err := service.CreateBeer(validBeer())

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.ErrInternal, appErr.Code())
	assert.Equal(t, 500, appErr.StatusCode())
	assert.True(t, mock.createCalled)
}