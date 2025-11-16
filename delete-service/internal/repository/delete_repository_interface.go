package repository

type DeleteRepositoryInterface interface {
	DeleteBeer(id string) error
}
