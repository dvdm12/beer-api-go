package repository

import (
	"context"
	"errors"
	"testing"

	"readservice/internal/db"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mockSingleResult simulates a MongoDB single result.
// It allows controlling the Decode behavior in tests.
type mockSingleResult struct {
	err error
}

// Decode simulates decoding a MongoDB document.
func (m *mockSingleResult) Decode(v interface{}) error {
	return m.err
}

// mockCollection simulates a MongoDB collection.
// It returns predefined errors for testing scenarios.
type mockCollection struct {
	findOneErr error
	findErr    error
}

// FindOne mocks retrieving a single document.
// It returns a SingleResultInterface for testability.
func (m *mockCollection) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOneOptions,
) db.SingleResultInterface {
	return &mockSingleResult{err: m.findOneErr}
}

// Find mocks retrieving multiple documents.
func (m *mockCollection) Find(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (*mongo.Cursor, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return nil, nil
}

// TestReadRepository_Creation verifies repository initialization.
func TestReadRepository_Creation(t *testing.T) {
	mock := &mockCollection{}
	repo := NewReadRepository(mock)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.collection)
}

// TestReadRepository_GetBeerByID_InvalidID verifies invalid ID handling.
func TestReadRepository_GetBeerByID_InvalidID(t *testing.T) {
	mock := &mockCollection{}
	repo := NewReadRepository(mock)

	_, err := repo.GetBeerByID("invalid-id")

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidID, err)
}

// TestReadRepository_GetBeerByID_NotFound verifies mapping of not found errors.
func TestReadRepository_GetBeerByID_NotFound(t *testing.T) {
	mock := &mockCollection{
		findOneErr: mongo.ErrNoDocuments,
	}
	repo := NewReadRepository(mock)

	validID := primitive.NewObjectID().Hex()
	_, err := repo.GetBeerByID(validID)

	assert.NotNil(t, err)
	assert.Equal(t, ErrBeerNotFound, err)
}

// TestReadRepository_GetAllBeers_Error verifies Find error handling.
func TestReadRepository_GetAllBeers_Error(t *testing.T) {
	mock := &mockCollection{
		findErr: errors.New("find error"),
	}
	repo := NewReadRepository(mock)

	_, err := repo.GetAllBeers()

	assert.NotNil(t, err)
	assert.Equal(t, "find error", err.Error())
}

func TestReadRepository_GetBeerByID_Success(t *testing.T) {
	mock := &mockCollection{
		findOneErr: nil,
	}
	repo := NewReadRepository(mock)

	validID := primitive.NewObjectID().Hex()

	beer, err := repo.GetBeerByID(validID)

	assert.Nil(t, err)
	assert.NotNil(t, beer)
}
