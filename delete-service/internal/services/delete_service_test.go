// Package services contains business logic tests.
package services

import (
	"deleteservice/internal/errors"
	"deleteservice/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRepository simulates repository behavior.
type mockRepository struct {
	called    bool
	deleteErr error
}

// DeleteBeer mock implementation.
func (m *mockRepository) DeleteBeer(id string) error {
	m.called = true
	return m.deleteErr
}

// Test successful deletion.
func TestDeleteService_DeleteBeer_Success(t *testing.T) {
	mock := &mockRepository{}
	service := NewDeleteService(mock)

	err := service.DeleteBeer("507f1f77bcf86cd799439011")

	assert.Nil(t, err)
	assert.True(t, mock.called)
}

// Test empty ID validation.
func TestDeleteService_DeleteBeer_EmptyID(t *testing.T) {
	mock := &mockRepository{}
	service := NewDeleteService(mock)

	err := service.DeleteBeer("")

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeInvalidID, appErr.Code())
	assert.Equal(t, 400, appErr.StatusCode())
	assert.False(t, mock.called)
}

// Test invalid ID format from repository.
func TestDeleteService_DeleteBeer_InvalidIDFormat(t *testing.T) {
	mock := &mockRepository{deleteErr: repository.ErrInvalidID}
	service := NewDeleteService(mock)

	err := service.DeleteBeer("invalid-id")

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeInvalidID, appErr.Code())
	assert.Equal(t, 400, appErr.StatusCode())
	assert.True(t, mock.called)
}

// Test beer not found case.
func TestDeleteService_DeleteBeer_NotFound(t *testing.T) {
	mock := &mockRepository{deleteErr: repository.ErrBeerNotFound}
	service := NewDeleteService(mock)

	err := service.DeleteBeer("507f1f77bcf86cd799439011")

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.CodeBeerNotFound, appErr.Code())
	assert.Equal(t, 404, appErr.StatusCode())
	assert.Contains(t, appErr.Error(), "507f1f77bcf86cd799439011")
	assert.True(t, mock.called)
}

// Test internal error fallback.
func TestDeleteService_DeleteBeer_InternalError(t *testing.T) {
	mock := &mockRepository{deleteErr: assert.AnError}
	service := NewDeleteService(mock)

	err := service.DeleteBeer("507f1f77bcf86cd799439011")

	assert.NotNil(t, err)
	appErr, ok := err.(errors.AppError)
	assert.True(t, ok)
	assert.Equal(t, errors.ErrInternal, appErr.Code())
	assert.Equal(t, 500, appErr.StatusCode())
	assert.True(t, mock.called)
}