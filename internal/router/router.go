package router

import (
	"webhook-tower/internal/config"

	"github.com/gin-gonic/gin"
)

// NewRouter initializes a new Gin engine with routes from the config
func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	// Register routes from config
	for _, route := range cfg.Routes {
		r.Handle(route.Method, route.Path, createHandler(route))
	}

	return r
}

func createHandler(route config.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Header Matching
		for _, h := range route.Headers {
			if c.GetHeader(h.Key) != h.Value {
				c.JSON(403, gin.H{"error": "header mismatch"})
				return
			}
		}

		// 2. TODO: Payload matching

		// 3. TODO: API Key auth

		c.JSON(200, gin.H{
			"message": "matched",
		})
	}
}
