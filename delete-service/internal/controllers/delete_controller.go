package controllers

import (
	"deleteservice/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteController struct {
	service services.DeleteServiceInterface
}

func NewDeleteController(service services.DeleteServiceInterface) *DeleteController {
	return &DeleteController{service: service}
}

func (c *DeleteController) DeleteBeer(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Beer ID is required"})
		return
	}

	if err := c.service.DeleteBeer(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Beer not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Beer deleted successfully"})
}
