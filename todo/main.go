package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Get Todo by Id",
		})
	})

	r.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.PATCH("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "edit to do by id",
		})
	})

	r.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "delete to do by id",
		})
	})

	r.DELETE("user/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "delete to do by id",
		})
	})

	r.Run(":3001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
