// Package services provides business logic for the create service.
package services

import (
	"createservice/internal/errors"
	"createservice/internal/models"
	"createservice/internal/repository"
	"fmt"
	"time"
)

// CreateService handles business operations for beer creation.
type CreateService struct {
	repo repository.CreateRepositoryInterface
}

// NewCreateService creates a new service instance.
func NewCreateService(repo repository.CreateRepositoryInterface) *CreateService {
	return &CreateService{repo: repo}
}

// CreateBeer validates and creates a new beer.
func (s *CreateService) CreateBeer(beer models.Beer) error {
	// Validate input data.
	if err := validateBeer(beer); err != nil {
		return err
	}

	// Check for duplicate beer.
	exists, err := s.repo.ExistsBeer(beer.Name)
	if err != nil {
		return errors.Internal(err)
	}
	if exists {
		return errors.NewDuplicateBeerError(
			fmt.Sprintf("beer '%s' already exists", beer.Name),
		)
	}

	// Persist beer.
	_, err = s.repo.CreateBeer(beer)
	if err != nil {
		return errors.Internal(err)
	}

	return nil
}

// Business validation constraints.
const (
	MinAlcohol = 0.0
	MaxAlcohol = 70.0
)

// validateBeer validates beer fields.
func validateBeer(beer models.Beer) error {
	if beer.Name == "" {
		return errors.NewValidationError("name is required")
	}
	if beer.Brand == "" {
		return errors.NewValidationError("brand is required")
	}

	if beer.Alcohol < MinAlcohol || beer.Alcohol > MaxAlcohol {
		return errors.NewValidationError(
			fmt.Sprintf("alcohol percentage must be between %.1f and %.1f", MinAlcohol, MaxAlcohol),
		)
	}

	currentYear := time.Now().Year()

	if beer.Year < 1800 || beer.Year > currentYear {
		return errors.NewValidationError(
			fmt.Sprintf("year must be between 1800 and %d", currentYear),
		)
	}

	return nil
}