package services

import (
    "createservice/internal/models"
    "createservice/internal/repository"
)

type CreateService struct {
    repo *repository.CreateRepository
}

func NewCreateService(repo *repository.CreateRepository) *CreateService {
    return &CreateService{repo: repo}
}

func (s *CreateService) CreateBeer(beer models.Beer) error {
    _, err := s.repo.CreateBeer(beer)
    return err
}
