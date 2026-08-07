package response

import "github.com/gin-gonic/gin"

// Success returns a 200 response with the given data.
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"data": data})
}

// Error returns an error response with the given status code and message.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
