package services

import "updateservice/internal/models"

type UpdateServiceInterface interface {
    UpdateBeer(id string, beer models.Beer) error
}
