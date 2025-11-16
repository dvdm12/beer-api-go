package controllers

import "github.com/gin-gonic/gin"

type DeleteControllerInterface interface {
	DeleteBeer(ctx *gin.Context)
}
