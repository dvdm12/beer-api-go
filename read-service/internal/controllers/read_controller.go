package controllers

import (
	"net/http"
	"readservice/internal/services"

	"github.com/gin-gonic/gin"
)

type ReadController struct {
	service services.ReadServiceInterface
}

func NewReadController(service services.ReadServiceInterface) *ReadController {
	return &ReadController{service: service}
}

func (c *ReadController) GetBeerByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Beer ID is required"})
		return
	}

	beer, err := c.service.GetBeerByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Beer not found"})
		return
	}

	ctx.JSON(http.StatusOK, beer)
}

func (c *ReadController) GetAllBeers(ctx *gin.Context) {
	beers, err := c.service.GetAllBeers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch beers"})
		return
	}

	ctx.JSON(http.StatusOK, beers)
}
