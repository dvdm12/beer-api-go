package db

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockMongoCollection struct {
	called bool
	err    error
}

func (m *mockMongoCollection) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	m.called = true
	if m.err != nil {
		return nil, m.err
	}
	return &mongo.UpdateResult{}, nil
}

type testAdapter struct {
	mock *mockMongoCollection
}

func (t *testAdapter) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	return t.mock.UpdateOne(ctx, filter, update, opts...)
}

func TestMongoCollectionAdapter_UpdateOne_Success(t *testing.T) {
	mock := &mockMongoCollection{}
	adapter := &testAdapter{mock: mock}

	_, err := adapter.UpdateOne(
		context.Background(),
		map[string]string{"_id": "123"},
		map[string]interface{}{"$set": map[string]string{"name": "Updated Beer"}},
	)

	assert.Nil(t, err)
	assert.True(t, mock.called)
}

func TestMongoCollectionAdapter_UpdateOne_Error(t *testing.T) {
	mock := &mockMongoCollection{
		err: errors.New("update error"),
	}
	adapter := &testAdapter{mock: mock}

	_, err := adapter.UpdateOne(
		context.Background(),
		map[string]string{"_id": "123"},
		map[string]interface{}{"$set": map[string]string{"name": "Updated Beer"}},
	)

	assert.NotNil(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.True(t, mock.called)
}
