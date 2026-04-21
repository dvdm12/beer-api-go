package repository

import (
	"context"
	"createservice/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)


type mockMongoCollection struct {
	insertCalled bool
	countCalled  bool
	insertErr    error
	countErr     error
	countResult  int64
}

func (m *mockMongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	m.insertCalled = true
	return nil, m.insertErr
}

func (m *mockMongoCollection) CountDocuments(ctx context.Context, filter bson.M) (int64, error) {
	m.countCalled = true
	return m.countResult, m.countErr
}

func validBeer() models.Beer {
	return models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 5.5,
		Year:    2023,
	}
}

func TestCreateRepository_CreateBeer_Success(t *testing.T) {
	mock := &mockMongoCollection{}
	repo := NewCreateRepository(mock)

	_, err := repo.CreateBeer(validBeer())

	assert.Nil(t, err)
	assert.True(t, mock.insertCalled)
}

func TestCreateRepository_CreateBeer_Error(t *testing.T) {
	mock := &mockMongoCollection{insertErr: errors.New("insert failed")}
	repo := NewCreateRepository(mock)

	_, err := repo.CreateBeer(validBeer())

	assert.NotNil(t, err)
	assert.Equal(t, "insert failed", err.Error())
	assert.True(t, mock.insertCalled)
}


func TestCreateRepository_ExistsBeer_Exists(t *testing.T) {
	mock := &mockMongoCollection{countResult: 1}
	repo := NewCreateRepository(mock)

	exists, err := repo.ExistsBeer("TestBeer")

	assert.Nil(t, err)
	assert.True(t, exists)
	assert.True(t, mock.countCalled)
}

func TestCreateRepository_ExistsBeer_NotExists(t *testing.T) {
	mock := &mockMongoCollection{countResult: 0}
	repo := NewCreateRepository(mock)

	exists, err := repo.ExistsBeer("GhostBeer")

	assert.Nil(t, err)
	assert.False(t, exists)
	assert.True(t, mock.countCalled)
}

func TestCreateRepository_ExistsBeer_DBError(t *testing.T) {
	mock := &mockMongoCollection{countErr: errors.New("mongo connection lost")}
	repo := NewCreateRepository(mock)

	exists, err := repo.ExistsBeer("TestBeer")

	assert.NotNil(t, err)
	assert.False(t, exists)
	assert.Equal(t, "mongo connection lost", err.Error())
	assert.True(t, mock.countCalled)
}