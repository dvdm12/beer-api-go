package controllers

import (
    "createservice/internal/models"
    "createservice/internal/services"
    "github.com/gin-gonic/gin"
    "net/http"
)

type CreateController struct {
    service *services.CreateService
}

func NewCreateController(service *services.CreateService) *CreateController {
    return &CreateController{service: service}
}

func (c *CreateController) CreateBeer(ctx *gin.Context) {
    var beer models.Beer

    if err := ctx.ShouldBindJSON(&beer); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := c.service.CreateBeer(beer)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create beer"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Beer created successfully"})
}
