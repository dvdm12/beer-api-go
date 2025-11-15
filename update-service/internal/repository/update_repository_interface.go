package repository

import "updateservice/internal/models"

type UpdateRepositoryInterface interface {
    UpdateBeer(id string, beer models.Beer) (*struct{}, error)
}
