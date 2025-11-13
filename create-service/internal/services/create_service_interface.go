package services

import "createservice/internal/models"

type CreateServiceInterface interface {
    CreateBeer(beer models.Beer) error
}
