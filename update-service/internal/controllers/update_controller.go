package controllers

import (
    "updateservice/internal/models"
    "updateservice/internal/services"
    "github.com/gin-gonic/gin"
    "net/http"
)

type UpdateController struct {
    service services.UpdateServiceInterface
}

func NewUpdateController(service services.UpdateServiceInterface) *UpdateController {
    return &UpdateController{service: service}
}

func (c *UpdateController) UpdateBeer(ctx *gin.Context) {
    id := ctx.Param("id")

    if id == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Beer ID is required"})
        return
    }

    var beer models.Beer
    if err := ctx.ShouldBindJSON(&beer); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.service.UpdateBeer(id, beer); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update beer"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Beer updated successfully"})
}
