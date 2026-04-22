package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"webhook-tower/internal/config"
	"webhook-tower/internal/executor"

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
		// 1. API Key auth
		if route.APIKey != "" {
			key := c.GetHeader("X-API-Key")
			if key == "" {
				key = c.Query("api_key")
			}

			if key != route.APIKey {
				c.JSON(401, gin.H{"error": "unauthorized"})
				return
			}
		}

		// 2. Header Matching
		for _, h := range route.Headers {
			if c.GetHeader(h.Key) != h.Value {
				c.JSON(403, gin.H{"error": "header mismatch"})
				return
			}
		}

		var body []byte
		var err error
		if c.Request.Body != nil {
			body, err = io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		// 3. Payload matching
		if len(route.Rules) > 0 {
			if err != nil {
				c.JSON(400, gin.H{"error": "failed to read body"})
				return
			}

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

		// 4. Execution
		if route.Command.Execute != "" {
			var payload map[string]interface{}
			if body != nil {
				json.Unmarshal(body, &payload)
			}

			cmdStr, err := executor.PrepareCommand(route.Command.Execute, payload)
			if err != nil {
				c.JSON(500, gin.H{"error": "failed to prepare command"})
				return
			}

			if route.Command.Async {
				go executor.Execute(cmdStr)
				c.JSON(202, gin.H{"message": "accepted"})
				return
			} else {
				result := executor.Execute(cmdStr)
				status := 200
				if result.ExitCode != 0 {
					status = 500
				}
				c.JSON(status, gin.H{
					"message":   "executed",
					"stdout":    result.Stdout,
					"stderr":    result.Stderr,
					"exit_code": result.ExitCode,
				})
				return
			}
		}

		c.JSON(200, gin.H{
			"message": "matched",
		})
	}
}
