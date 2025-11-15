package main

import (
    "updateservice/internal/controllers"
    "updateservice/internal/db"
    "updateservice/internal/repository"
    "updateservice/internal/services"
    "github.com/gin-gonic/gin"
)

func main() {

    collection := db.Connect()

    repo := repository.NewUpdateRepository(collection)
    service := services.NewUpdateService(repo)
    controller := controllers.NewUpdateController(service)

    r := gin.Default()

    r.PUT("/beers/update", controller.UpdateBeer)

    r.Run(":8081")
}
