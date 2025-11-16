package controllers

import "github.com/gin-gonic/gin"

type ReadControllerInterface interface {
	GetBeerByID(ctx *gin.Context)
	GetAllBeers(ctx *gin.Context)
}
