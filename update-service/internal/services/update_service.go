package services

import (
	"updateservice/internal/models"
	"updateservice/internal/repository"
)

type UpdateService struct {
	repo repository.UpdateRepositoryInterface
}

func NewUpdateService(repo repository.UpdateRepositoryInterface) *UpdateService {
	return &UpdateService{repo: repo}
}

func (s *UpdateService) UpdateBeer(id string, beer models.Beer) error {
	_, err := s.repo.UpdateBeer(id, beer)
	return err
}
