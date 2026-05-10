package controllers

import (
	stderrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "dataanalysis/internal/errors"
	"dataanalysis/internal/services"
)

// AnalysisController handles beer analysis requests.
type AnalysisController struct {
	service services.AnalysisServiceInterface
}

// NewAnalysisController returns a new AnalysisController.
func NewAnalysisController(service services.AnalysisServiceInterface) *AnalysisController {
	return &AnalysisController{service: service}
}

// GetStrongest handles GET /beers/rankings/strongest.
func (c *AnalysisController) GetStrongest(ctx *gin.Context) {
	beer, err := c.service.GetStrongest(ctx.Request.Context())
	if err != nil {
		c.writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, beer)
}

// GetWeakest handles GET /beers/rankings/weakest.
func (c *AnalysisController) GetWeakest(ctx *gin.Context) {
	beer, err := c.service.GetWeakest(ctx.Request.Context())
	if err != nil {
		c.writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, beer)
}

// GetOldest handles GET /beers/rankings/oldest.
func (c *AnalysisController) GetOldest(ctx *gin.Context) {
	beer, err := c.service.GetOldest(ctx.Request.Context())
	if err != nil {
		c.writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, beer)
}

// GetNewest handles GET /beers/rankings/newest.
func (c *AnalysisController) GetNewest(ctx *gin.Context) {
	beer, err := c.service.GetNewest(ctx.Request.Context())
	if err != nil {
		c.writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, beer)
}

// GetStats handles GET /beers/stats.
func (c *AnalysisController) GetStats(ctx *gin.Context) {
	stats, err := c.service.GetStats(ctx.Request.Context())
	if err != nil {
		c.writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, stats)
}

// GetStatsByBrand handles GET /beers/stats/brand.
func (c *AnalysisController) GetStatsByBrand(ctx *gin.Context) {
	stats, err := c.service.GetStatsByBrand(ctx.Request.Context())
	if err != nil {
		c.writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, stats)
}

// errorResponse defines the error payload.
type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeError writes an HTTP error response.
func (c *AnalysisController) writeError(ctx *gin.Context, err error) {
	var appErr apperrors.AppError
	if stderrors.As(err, &appErr) {
		ctx.JSON(appErr.StatusCode(), errorResponse{
			Code:    appErr.Code(),
			Message: appErr.Error(),
		})
		return
	}

	if stderrors.Is(err, services.ErrNoResults) {
		ctx.JSON(http.StatusNotFound, errorResponse{
			Code:    "EMPTY_COLLECTION",
			Message: "no beers found in database",
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, errorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
	})
}
