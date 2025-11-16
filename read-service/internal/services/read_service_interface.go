package services

import "readservice/internal/models"

type ReadServiceInterface interface {
	GetBeerByID(id string) (*models.Beer, error)
	GetAllBeers() ([]models.Beer, error)
}
