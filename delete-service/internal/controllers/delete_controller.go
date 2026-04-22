// Package controllers handles HTTP requests.
package controllers

import (
	"deleteservice/internal/errors"
	"deleteservice/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteController handles delete endpoints.
type DeleteController struct {
	service services.DeleteServiceInterface
}

// NewDeleteController creates a controller instance.
func NewDeleteController(service services.DeleteServiceInterface) *DeleteController {
	return &DeleteController{service: service}
}

// DeleteBeer handles DELETE /beer/:id.
func (c *DeleteController) DeleteBeer(ctx *gin.Context) {
	// Get ID from path
	id := ctx.Param("id")

	// Call service
	if err := c.service.DeleteBeer(id); err != nil {
		appErr, _ := errors.FromError(err)
		ctx.JSON(appErr.StatusCode(), errorResponse(appErr))
		return
	}

	// Success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Beer deleted successfully",
	})
}

// errorResponse formats error output.
func errorResponse(err errors.AppError) gin.H {
	return gin.H{
		"code":    err.Code(),
		"message": err.Error(),
	}
}