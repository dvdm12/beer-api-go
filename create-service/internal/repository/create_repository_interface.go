package repository

import "createservice/internal/models"

type CreateRepositoryInterface interface {
    CreateBeer(beer models.Beer) (*struct{}, error)
}
