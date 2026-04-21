// Package controllers provides HTTP handlers for the create service.
package controllers

import (
	"createservice/internal/errors"
	"createservice/internal/models"
	"createservice/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateController handles beer creation requests.
type CreateController struct {
	service services.CreateServiceInterface
}

// NewCreateController creates a new controller instance.
func NewCreateController(service services.CreateServiceInterface) *CreateController {
	return &CreateController{service: service}
}

// CreateBeer handles the HTTP request to create a beer.
func (c *CreateController) CreateBeer(ctx *gin.Context) {
	var beer models.Beer

	// Parse and validate request body.
	if err := ctx.ShouldBindJSON(&beer); err != nil {
		appErr := errors.NewValidationError(err.Error())
		ctx.JSON(appErr.StatusCode(), errorResponse(appErr))
		return
	}

	// Execute business logic.
	if err := c.service.CreateBeer(beer); err != nil {
		appErr, _ := errors.FromError(err)
		ctx.JSON(appErr.StatusCode(), errorResponse(appErr))
		return
	}

	// Return success response.
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Beer has been created successfully",
	})
}

// errorResponse formats an AppError for HTTP responses.
func errorResponse(err errors.AppError) gin.H {
	return gin.H{
		"code":    err.Code(),
		"message": err.Error(),
	}
}