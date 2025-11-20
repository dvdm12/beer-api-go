package main

import (
	"readservice/internal/controllers"
	"readservice/internal/db"
	"readservice/internal/repository"
	"readservice/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {

	collection := db.Connect()

	repo := repository.NewReadRepository(collection)
	service := services.NewReadService(repo)
	controller := controllers.NewReadController(service)

	r := gin.Default()

	r.GET("/beers/:id", controller.GetBeerByID)
	r.GET("/beers", controller.GetAllBeers)

	r.GET("/health", func(c *gin.Context) {
    	c.JSON(200, gin.H{"status": "ok"})
	})

	r.Run(":8081")
}
