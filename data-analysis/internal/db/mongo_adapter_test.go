package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mockCollection is a test double for MongoCollectionInterface.
type mockCollection struct {
	findOneCalled bool
	findCalled    bool
	findErr       error
}

// FindOne marks the method as called and returns an empty result.
func (m *mockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	m.findOneCalled = true
	return &mongo.SingleResult{}
}

// Find marks the method as called and returns a preset error.
func (m *mockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	m.findCalled = true
	return nil, m.findErr
}

// TestMongoCollectionAdapter_FindOne_Called verifies FindOne is invoked.
func TestMongoCollectionAdapter_FindOne_Called(t *testing.T) {
	mock := &mockCollection{}

	mock.FindOne(context.Background(), nil)

	assert.True(t, mock.findOneCalled)
}

// TestMongoCollectionAdapter_Find_Success verifies Find without error.
func TestMongoCollectionAdapter_Find_Success(t *testing.T) {
	mock := &mockCollection{}

	cursor, err := mock.Find(context.Background(), nil)

	assert.Nil(t, err)
	assert.Nil(t, cursor)
	assert.True(t, mock.findCalled)
}

// TestMongoCollectionAdapter_Find_Error verifies Find returns an error.
func TestMongoCollectionAdapter_Find_Error(t *testing.T) {
	mock := &mockCollection{findErr: assert.AnError}

	cursor, err := mock.Find(context.Background(), nil)

	assert.NotNil(t, err)
	assert.Nil(t, cursor)
	assert.True(t, mock.findCalled)
}
