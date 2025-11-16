package repository

import "readservice/internal/models"

type ReadRepositoryInterface interface {
	GetBeerByID(id string) (*models.Beer, error)
	GetAllBeers() ([]models.Beer, error)
}
