package router

import (
	"github.com/gin-gonic/gin"
)

// NewRouter initializes a new Gin engine
func NewRouter() *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
		})
	})

	return r
}
