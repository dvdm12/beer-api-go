// Package repository contains persistence tests.
package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mockCollection simulates MongoDB behavior.
type mockCollection struct {
	deleteErr    error
	deletedCount int64
}

// DeleteOne mock implementation.
func (m *mockCollection) DeleteOne(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
	if m.deleteErr != nil {
		return nil, m.deleteErr
	}
	return &mongo.DeleteResult{DeletedCount: m.deletedCount}, nil
}

// Test constants.
const (
	validObjectID   = "507f1f77bcf86cd799439011"
	invalidObjectID = "invalid-id"
)

// Test repository creation.
func TestDeleteRepository_Creation(t *testing.T) {
	mock := &mockCollection{}
	repo := NewDeleteRepository(mock)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.collection)
}

// Test successful deletion.
func TestDeleteRepository_DeleteBeer_Success(t *testing.T) {
	mock := &mockCollection{deletedCount: 1}
	repo := NewDeleteRepository(mock)

	err := repo.DeleteBeer(validObjectID)

	assert.Nil(t, err)
}

// Test invalid ID format.
func TestDeleteRepository_DeleteBeer_InvalidID(t *testing.T) {
	mock := &mockCollection{}
	repo := NewDeleteRepository(mock)

	err := repo.DeleteBeer(invalidObjectID)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidID, err)
}

// Test beer not found.
func TestDeleteRepository_DeleteBeer_NotFound(t *testing.T) {
	mock := &mockCollection{deletedCount: 0}
	repo := NewDeleteRepository(mock)

	err := repo.DeleteBeer(validObjectID)

	assert.NotNil(t, err)
	assert.Equal(t, ErrBeerNotFound, err)
}

// Test database error.
func TestDeleteRepository_DeleteBeer_DBError(t *testing.T) {
	mock := &mockCollection{deleteErr: errors.New("database error")}
	repo := NewDeleteRepository(mock)

	err := repo.DeleteBeer(validObjectID)

	assert.NotNil(t, err)
	assert.Equal(t, "database error", err.Error())
}