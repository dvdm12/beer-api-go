package db

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockMongoCollection struct {
	called     bool
	err        error
	findCalled bool
}

func (m *mockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*struct{}) *mongo.SingleResult {
	m.called = true
	return nil
}

func (m *mockMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*struct{}) (*mongo.Cursor, error) {
	m.findCalled = true
	if m.err != nil {
		return nil, m.err
	}
	return nil, nil
}

type testAdapter struct {
	mock *mockMongoCollection
}

func (t *testAdapter) FindOne(ctx context.Context, filter interface{}, opts ...*struct{}) *mongo.SingleResult {
	return t.mock.FindOne(ctx, filter, opts...)
}

func (t *testAdapter) Find(ctx context.Context, filter interface{}, opts ...*struct{}) (*mongo.Cursor, error) {
	return t.mock.Find(ctx, filter, opts...)
}

func TestMongoCollectionAdapter_FindOne_Success(t *testing.T) {
	mock := &mockMongoCollection{}
	adapter := &testAdapter{mock: mock}

	_ = adapter.FindOne(context.Background(), map[string]string{"test": "demo"})

	assert.True(t, mock.called)
}

func TestMongoCollectionAdapter_Find_Success(t *testing.T) {
	mock := &mockMongoCollection{}
	adapter := &testAdapter{mock: mock}

	_, err := adapter.Find(context.Background(), map[string]string{"test": "demo"})

	assert.Nil(t, err)
	assert.True(t, mock.findCalled)
}

func TestMongoCollectionAdapter_Find_Error(t *testing.T) {
	mock := &mockMongoCollection{
		err: errors.New("find error"),
	}
	adapter := &testAdapter{mock: mock}

	_, err := adapter.Find(context.Background(), map[string]string{"test": "demo"})

	assert.NotNil(t, err)
	assert.Equal(t, "find error", err.Error())
	assert.True(t, mock.findCalled)
}
