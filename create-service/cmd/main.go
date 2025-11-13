package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/greeting", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "hi, how are you"})
    })

    r.Run(":8080")
}
