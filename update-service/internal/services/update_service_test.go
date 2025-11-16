package services

import (
	"updateservice/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	called bool
	err    error
}

func (m *mockRepo) UpdateBeer(id string, beer models.Beer) (*struct{}, error) {
	m.called = true
	return nil, m.err
}

func TestUpdateService_UpdateBeer_Success(t *testing.T) {
	mock := &mockRepo{}

	service := NewUpdateService(mock)

	beer := models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 4.8,
		Year:    2022,
	}

	err := service.UpdateBeer("507f1f77bcf86cd799439011", beer)

	assert.Nil(t, err)
	assert.True(t, mock.called)
}
