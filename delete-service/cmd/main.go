package main

import (
	"deleteservice/internal/controllers"
	"deleteservice/internal/db"
	"deleteservice/internal/repository"
	"deleteservice/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {

	collection := db.Connect()

	repo := repository.NewDeleteRepository(collection)
	service := services.NewDeleteService(repo)
	controller := controllers.NewDeleteController(service)

	r := gin.Default()

	r.DELETE("/beers/:id", controller.DeleteBeer)

	r.Run(":8083")
}
