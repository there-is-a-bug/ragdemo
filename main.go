package main

import (
	"MyMemory/logs"
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logs.InitLogger(ctx)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	api := r.Group("/api")
	{
		api.POST("/retrieve", Retrieve)
		api.POST("/upload", Insert)
	}
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("gin start failed: %v", err)
	}
}
