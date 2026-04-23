package services

import (
	"context"
	"time"

	"updateservice/internal/errors"
	"updateservice/internal/models"
	"updateservice/internal/repository"
)

// UpdateService orchestrates business logic for beer updates.
type UpdateService struct {
	repo repository.UpdateRepositoryInterface
}

// NewUpdateService creates a new UpdateService instance.
func NewUpdateService(repo repository.UpdateRepositoryInterface) *UpdateService {
	return &UpdateService{repo: repo}
}

// UpdateBeer validates domain rules and invokes the repository layer.
func (s *UpdateService) UpdateBeer(ctx context.Context, id string, beer models.Beer) error {
	// Validate required fields.
	if beer.Name == "" || beer.Brand == "" {
		return errors.ValidationError("Name and brand are required", models.ErrCodeMissingData)
	}

	// Validate alcohol content.
	if beer.Alcohol < 0 {
		return errors.ValidationError("Alcohol content cannot be negative", models.ErrCodeInvalidAlcohol)
	}

	// Validate production year.
	currentYear := time.Now().Year()
	if beer.Year < 1800 || beer.Year > currentYear {
		return errors.ValidationError("Invalid brew year provided", models.ErrCodeInvalidYear)
	}

	// Delegate persistence to the repository.
	// Note: The repository should be updated to accept the context and return only an error.
	return s.repo.UpdateBeer(ctx, id, beer)
}
