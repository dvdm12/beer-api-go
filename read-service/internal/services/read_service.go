package services

import (
	"readservice/internal/models"
	"readservice/internal/repository"
)

type ReadService struct {
	repo repository.ReadRepositoryInterface
}

func NewReadService(repo repository.ReadRepositoryInterface) *ReadService {
	return &ReadService{repo: repo}
}

func (s *ReadService) GetBeerByID(id string) (*models.Beer, error) {
	return s.repo.GetBeerByID(id)
}

func (s *ReadService) GetAllBeers() ([]models.Beer, error) {
	return s.repo.GetAllBeers()
}
