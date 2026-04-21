package main

import (
	"log"

	"ez-admin-gin/server/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    cfg.App.Env,
		})
	})

	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
