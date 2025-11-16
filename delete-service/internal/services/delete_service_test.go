package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	deleteErr error
}

func (m *mockRepository) DeleteBeer(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	return nil
}

func TestDeleteService_DeleteBeer_Success(t *testing.T) {
	mock := &mockRepository{}
	service := NewDeleteService(mock)

	err := service.DeleteBeer("507f1f77bcf86cd799439011")

	assert.Nil(t, err)
}

func TestDeleteService_DeleteBeer_Error(t *testing.T) {
	mock := &mockRepository{
		deleteErr: errors.New("not found"),
	}
	service := NewDeleteService(mock)

	err := service.DeleteBeer("507f1f77bcf86cd799439011")

	assert.NotNil(t, err)
	assert.Equal(t, "not found", err.Error())
}
