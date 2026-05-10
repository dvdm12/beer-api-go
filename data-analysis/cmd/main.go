package main

import (
	"dataanalysis/internal/controllers"
	"dataanalysis/internal/db"
	"dataanalysis/internal/repository"
	"dataanalysis/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	collection := db.Connect()
	repo := repository.NewAnalysisRepository(collection, nil)
	service := services.NewAnalysisService(repo)
	controller := controllers.NewAnalysisController(service)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Rankings
	r.GET("/beers/rankings/strongest", controller.GetStrongest)
	r.GET("/beers/rankings/weakest", controller.GetWeakest)
	r.GET("/beers/rankings/oldest", controller.GetOldest)
	r.GET("/beers/rankings/newest", controller.GetNewest)

	// Statistics
	r.GET("/beers/stats", controller.GetStats)
	r.GET("/beers/stats/brand", controller.GetStatsByBrand)

	r.Run(":8084")
}
