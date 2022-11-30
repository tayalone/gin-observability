package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	ess "github.com/tayalone/go-ess-package/otel"
)

const (
	service     = "todo"
	environment = "dev"
)

func main() {
	tp, err := ess.JaegertracerProvider(os.Getenv("JEAGER_ENDPOINT"), service, environment)
	if err != nil {
		log.Fatal(err)
	}

	otelCtx := context.Background()
	defer func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(otelCtx)

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
