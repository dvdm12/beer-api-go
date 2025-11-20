package main

import (
	"createservice/internal/controllers"
	"createservice/internal/db"
	"createservice/internal/repository"
	"createservice/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {

	collection := db.Connect()

	repo := repository.NewCreateRepository(collection)
	service := services.NewCreateService(repo)
	controller := controllers.NewCreateController(service)

	r := gin.Default()

	r.POST("/beers", controller.CreateBeer)

	r.Run(":8080")
}
