package services

import (
	"createservice/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	called bool
	err    error
}

func (m *mockRepo) CreateBeer(beer models.Beer) (*struct{}, error) {
	m.called = true
	return nil, m.err
}

func TestCreateService_CreateBeer_Success(t *testing.T) {
	mock := &mockRepo{}

	service := NewCreateService(mock)

	beer := models.Beer{
		Name:    "TestBeer",
		Brand:   "TestBrand",
		Alcohol: 4.8,
		Year:    2022,
	}

	err := service.CreateBeer(beer)

	assert.Nil(t, err)
	assert.True(t, mock.called)
}
