package db

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockMongoCollection struct {
	called bool
	err    error
}

func (m *mockMongoCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	m.called = true
	if m.err != nil {
		return nil, m.err
	}
	return &mongo.DeleteResult{DeletedCount: 1}, nil
}

type testAdapter struct {
	mock *mockMongoCollection
}

func (t *testAdapter) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return t.mock.DeleteOne(ctx, filter)
}

func TestMongoCollectionAdapter_DeleteOne_Success(t *testing.T) {
	mock := &mockMongoCollection{}
	adapter := &testAdapter{mock: mock}

	result, err := adapter.DeleteOne(context.Background(), map[string]string{"test": "demo"})

	assert.Nil(t, err)
	assert.True(t, mock.called)
	assert.Equal(t, int64(1), result.DeletedCount)
}

func TestMongoCollectionAdapter_DeleteOne_Error(t *testing.T) {
	mock := &mockMongoCollection{
		err: errors.New("delete error"),
	}
	adapter := &testAdapter{mock: mock}

	_, err := adapter.DeleteOne(context.Background(), map[string]string{"test": "demo"})

	assert.NotNil(t, err)
	assert.Equal(t, "delete error", err.Error())
	assert.True(t, mock.called)
}
