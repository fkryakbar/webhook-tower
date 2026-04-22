package router

import (
	"bytes"
	"fmt"
	"io"
	"webhook-tower/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
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

		// 2. Payload matching
		if len(route.Rules) > 0 {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(400, gin.H{"error": "failed to read body"})
				return
			}
			// Restore body for further use if needed
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			for _, rule := range route.Rules {
				result := gjson.GetBytes(body, rule.Field)
				if !result.Exists() {
					c.JSON(403, gin.H{"error": fmt.Sprintf("field %s not found in payload", rule.Field)})
					return
				}

				// Simple operator matching
				if rule.Operator == "==" {
					if result.String() != fmt.Sprintf("%v", rule.Value) {
						c.JSON(403, gin.H{"error": "payload rule mismatch"})
						return
					}
				}
				// Support more operators if needed
			}
		}

		// 3. TODO: API Key auth

		c.JSON(200, gin.H{
			"message": "matched",
		})
	}
}
