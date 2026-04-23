// Package services contains business logic.
package services

import (
	stderrors "errors"
	"readservice/internal/errors"
	"readservice/internal/models"
	"readservice/internal/repository"
)

// ReadService handles read operations.
type ReadService struct {
	repo repository.ReadRepositoryInterface
}

// NewReadService creates a service instance.
func NewReadService(repo repository.ReadRepositoryInterface) *ReadService {
	return &ReadService{repo: repo}
}

// GetBeerByID retrieves a beer by ID.
func (s *ReadService) GetBeerByID(id string) (*models.Beer, error) {
	// Validate ID
	if id == "" {
		return nil, errors.NewInvalidIDError(id)
	}

	// Call repository
	beer, err := s.repo.GetBeerByID(id)
	if err == nil {
		return beer, nil
	}

	// Map errors
	switch {
	case stderrors.Is(err, repository.ErrInvalidID):
		return nil, errors.NewInvalidIDError(id)
	case stderrors.Is(err, repository.ErrBeerNotFound):
		return nil, errors.NewBeerNotFoundError(id)
	default:
		return nil, errors.Internal(err)
	}
}

// GetAllBeers retrieves all beers.
func (s *ReadService) GetAllBeers() ([]models.Beer, error) {
	// Call repository
	beers, err := s.repo.GetAllBeers()
	if err != nil {
		return nil, errors.Internal(err)
	}

	// Ensure non-nil slice
	if beers == nil {
		return []models.Beer{}, nil
	}

	return beers, nil
}
