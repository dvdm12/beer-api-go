package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockCollection struct {
	deleteErr    error
	deletedCount int64
}

func (m *mockCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	if m.deleteErr != nil {
		return nil, m.deleteErr
	}
	return &mongo.DeleteResult{DeletedCount: m.deletedCount}, nil
}

func TestDeleteRepository_DeleteBeer_Success(t *testing.T) {
	mock := &mockCollection{deletedCount: 1}
	repo := NewDeleteRepository(mock)

	err := repo.DeleteBeer("507f1f77bcf86cd799439011")

	assert.Nil(t, err)
}

func TestDeleteRepository_DeleteBeer_InvalidID(t *testing.T) {
	mock := &mockCollection{deletedCount: 1}
	repo := NewDeleteRepository(mock)

	err := repo.DeleteBeer("invalid-id")

	assert.NotNil(t, err)
}

func TestDeleteRepository_DeleteBeer_NotFound(t *testing.T) {
	mock := &mockCollection{deletedCount: 0}
	repo := NewDeleteRepository(mock)

	err := repo.DeleteBeer("507f1f77bcf86cd799439011")

	assert.NotNil(t, err)
}

func TestDeleteRepository_DeleteBeer_Error(t *testing.T) {
	mock := &mockCollection{
		deleteErr: errors.New("database error"),
	}
	repo := NewDeleteRepository(mock)

	err := repo.DeleteBeer("507f1f77bcf86cd799439011")

	assert.NotNil(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestDeleteRepository_Creation(t *testing.T) {
	mock := &mockCollection{}
	repo := NewDeleteRepository(mock)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.collection)
}
