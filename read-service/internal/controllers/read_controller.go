// Package controllers handles HTTP requests.
package controllers

import (
	"net/http"
	"readservice/internal/errors"
	"readservice/internal/services"

	"github.com/gin-gonic/gin"
)

// ReadController handles read endpoints.
type ReadController struct {
	service services.ReadServiceInterface
}

// NewReadController creates a controller instance.
func NewReadController(service services.ReadServiceInterface) *ReadController {
	return &ReadController{service: service}
}

// GetBeerByID handles GET /beers/:id.
func (c *ReadController) GetBeerByID(ctx *gin.Context) {
	// Get ID from path
	id := ctx.Param("id")

	// Call service
	beer, err := c.service.GetBeerByID(id)
	if err != nil {
		appErr, _ := errors.FromError(err)
		ctx.JSON(appErr.StatusCode(), errorResponse(appErr))
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, beer)
}

// GetAllBeers handles GET /beers.
func (c *ReadController) GetAllBeers(ctx *gin.Context) {
	// Call service
	beers, err := c.service.GetAllBeers()
	if err != nil {
		appErr, _ := errors.FromError(err)
		ctx.JSON(appErr.StatusCode(), errorResponse(appErr))
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, beers)
}

// errorResponse formats error output.
func errorResponse(err errors.AppError) gin.H {
	return gin.H{
		"code":    err.Code(),
		"message": err.Error(),
	}
}
