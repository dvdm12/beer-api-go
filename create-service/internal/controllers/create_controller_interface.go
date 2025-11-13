package controllers

import "github.com/gin-gonic/gin"

type CreateControllerInterface interface {
    CreateBeer(ctx *gin.Context)
}
