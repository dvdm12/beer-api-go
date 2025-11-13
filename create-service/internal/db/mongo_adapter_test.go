package db

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMongoCollection struct {
	called bool
	err    error
}

func (m *mockMongoCollection) InsertOne(ctx context.Context, doc interface{}) (*struct{}, error) {
	m.called = true
	if m.err != nil {
		return nil, m.err
	}
	return &struct{}{}, nil
}

type testAdapter struct {
	mock *mockMongoCollection
}

func (t *testAdapter) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	return t.mock.InsertOne(ctx, document)
}

func TestMongoCollectionAdapter_InsertOne_Success(t *testing.T) {
	mock := &mockMongoCollection{}
	adapter := &testAdapter{mock: mock}

	_, err := adapter.InsertOne(context.Background(), map[string]string{"test": "demo"})

	assert.Nil(t, err)
	assert.True(t, mock.called)
}

func TestMongoCollectionAdapter_InsertOne_Error(t *testing.T) {
	mock := &mockMongoCollection{
		err: errors.New("insert error"),
	}
	adapter := &testAdapter{mock: mock}

	_, err := adapter.InsertOne(context.Background(), map[string]string{"test": "demo"})

	assert.NotNil(t, err)
	assert.Equal(t, "insert error", err.Error())
	assert.True(t, mock.called)
}
