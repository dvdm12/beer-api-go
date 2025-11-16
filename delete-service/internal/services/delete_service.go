package services

import (
	"deleteservice/internal/repository"
)

type DeleteService struct {
	repo repository.DeleteRepositoryInterface
}

func NewDeleteService(repo repository.DeleteRepositoryInterface) *DeleteService {
	return &DeleteService{repo: repo}
}

func (s *DeleteService) DeleteBeer(id string) error {
	return s.repo.DeleteBeer(id)
}
