package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS configures Cross-Origin Resource Sharing for the API.
// Allows requests from the Next.js frontend during local development.
func CORS(frontendURL string) gin.HandlerFunc {
	allowOrigins := []string{"http://localhost:3000"}
	if frontendURL != "" {
		allowOrigins = append(allowOrigins, frontendURL)
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
