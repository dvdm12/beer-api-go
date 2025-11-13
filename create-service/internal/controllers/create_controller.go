package controllers

import (
    "createservice/internal/models"
    "createservice/internal/services"
    "github.com/gin-gonic/gin"
    "net/http"
)

type CreateController struct {
    service services.CreateServiceInterface
}

func NewCreateController(service services.CreateServiceInterface) *CreateController {
    return &CreateController{service: service}
}

func (c *CreateController) CreateBeer(ctx *gin.Context) {
    var beer models.Beer

    if err := ctx.ShouldBindJSON(&beer); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.service.CreateBeer(beer); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create beer"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Beer created successfully"})
}
