package services

import (
	"deleteservice/internal/errors"
	"deleteservice/internal/repository"
	stderrors "errors"
)

// DeleteService handles delete operations.
type DeleteService struct {
	repo repository.DeleteRepositoryInterface
}

// NewDeleteService creates a DeleteService instance.
func NewDeleteService(repo repository.DeleteRepositoryInterface) *DeleteService {
	return &DeleteService{repo: repo}
}

// DeleteBeer deletes a beer by ID.
func (s *DeleteService) DeleteBeer(id string) error {
	// Validate ID
	if id == "" {
		return errors.NewInvalidIDError(id)
	}

	// Call repository
	err := s.repo.DeleteBeer(id)
	if err == nil {
		return nil
	}

	// Map errors
	switch {
	case stderrors.Is(err, repository.ErrInvalidID):
		return errors.NewInvalidIDError(id)
	case stderrors.Is(err, repository.ErrBeerNotFound):
		return errors.NewBeerNotFoundError(id)
	default:
		return errors.Internal(err)
	}
}