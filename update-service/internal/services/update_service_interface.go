package services

import (
	"context"
	"updateservice/internal/models"
)

type UpdateServiceInterface interface {
	UpdateBeer(ctx context.Context, id string, beer models.Beer) error
}
