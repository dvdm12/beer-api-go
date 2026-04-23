package controllers

import (
	"net/http"

	"updateservice/internal/errors"
	"updateservice/internal/models"
	"updateservice/internal/services"

	"github.com/gin-gonic/gin"
)

// UpdateController handles HTTP requests for beer updates.
type UpdateController struct {
	service services.UpdateServiceInterface
}

// NewUpdateController creates a new UpdateController instance.
func NewUpdateController(service services.UpdateServiceInterface) *UpdateController {
	return &UpdateController{service: service}
}

// UpdateBeer processes the PUT request to update a beer's information.
func (c *UpdateController) UpdateBeer(ctx *gin.Context) {
	// Extract the beer ID from the URL path.
	id := ctx.Param("id")
	if id == "" {
		appErr := errors.BadRequest("Beer ID parameter is required")
		ctx.JSON(appErr.StatusCode(), gin.H{"error": appErr.Error(), "code": appErr.Code()})
		return
	}

	// Bind the incoming JSON payload to the Beer model.
	var beer models.Beer
	if err := ctx.ShouldBindJSON(&beer); err != nil {
		appErr := errors.BadRequest("Invalid JSON format or structure")
		ctx.JSON(appErr.StatusCode(), gin.H{"error": appErr.Error(), "code": appErr.Code()})
		return
	}

	// Forward the request to the service layer using the request context.
	// This allows cancellation signals to propagate through the architecture.
	if err := c.service.UpdateBeer(ctx.Request.Context(), id, beer); err != nil {
		// Map any error to the custom AppError contract.
		appErr, _ := errors.FromError(err)

		ctx.JSON(appErr.StatusCode(), gin.H{
			"error": appErr.Error(),
			"code":  appErr.Code(),
		})
		return
	}

	// Return a successful response.
	ctx.JSON(http.StatusOK, gin.H{"message": "Beer updated successfully"})
}
