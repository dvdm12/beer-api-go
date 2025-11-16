package services

type DeleteServiceInterface interface {
	DeleteBeer(id string) error
}
