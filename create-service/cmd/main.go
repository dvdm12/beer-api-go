package main

import (
    "createservice/internal/controllers"
    "createservice/internal/db"
    "createservice/internal/repository"
    "createservice/internal/services"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    client := db.Connect()
    repo := repository.NewCreateRepository(client)
    service := services.NewCreateService(repo)
    controller := controllers.NewCreateController(service)

    router.POST("/beers/create", controller.CreateBeer)

    router.Run(":8080")
}
