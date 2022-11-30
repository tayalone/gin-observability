package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// // ---------- router ------------------------
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "get User data",
		})
	})

	r.PATCH("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "edit user data",
		})
	})

	r.GET("/:id/todo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "get all user's todo",
		})
	})

	r.POST("/:id/todo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "create new  todo user data",
		})
	})

	r.PATCH("/:id/todo/:todoId", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "edit  todo user data",
		})
	})

	r.DELETE("/:id/todo/:todoId", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "delete  todo user data",
		})
	})
	// // -------------------------------------------------
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
