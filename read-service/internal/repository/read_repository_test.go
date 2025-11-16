package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockCollection struct {
	findOneErr error
	findErr    error
}

func (m *mockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return nil
}

func (m *mockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return nil, nil
}

func TestReadRepository_GetBeerByID_InvalidID(t *testing.T) {
	mock := &mockCollection{}
	repo := NewReadRepository(mock)

	_, err := repo.GetBeerByID("invalid-id")

	assert.NotNil(t, err)
}

func TestReadRepository_GetAllBeers_Error(t *testing.T) {
	mock := &mockCollection{
		findErr: errors.New("find error"),
	}
	repo := NewReadRepository(mock)

	_, err := repo.GetAllBeers()

	assert.NotNil(t, err)
	assert.Equal(t, "find error", err.Error())
}

func TestReadRepository_Creation(t *testing.T) {
	mock := &mockCollection{}
	repo := NewReadRepository(mock)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.collection)
}
