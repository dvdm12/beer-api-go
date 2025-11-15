package controllers

import "github.com/gin-gonic/gin"

type UpdateControllerInterface interface {
    UpdateBeer(ctx *gin.Context)
}
