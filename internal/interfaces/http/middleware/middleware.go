package middleware

import (
	"anyzzapp/pkg/domain"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Phone-Number-ID"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	return cors.New(config)
}

// ErrorHandler middleware for handling panics and errors
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(500, domain.ErrorResponse{
				Error:   "internal_server_error",
				Message: err,
				Code:    500,
			})
		} else {
			c.JSON(500, domain.ErrorResponse{
				Error:   "internal_server_error",
				Message: "An unexpected error occurred",
				Code:    500,
			})
		}
	})
}
