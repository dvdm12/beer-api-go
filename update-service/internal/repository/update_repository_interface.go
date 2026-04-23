package repository

import (
	"context"
	"updateservice/internal/models"
)

type UpdateRepositoryInterface interface {
	UpdateBeer(ctx context.Context, id string, beer models.Beer) error
}
