package ctxutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/pkg/response"
)

// GetUserID retrieves the authenticated user's ID from the Gin context.
// Set by the Auth middleware after JWT verification.
func GetUserID(c *gin.Context) string {
	return c.GetString("userID")
}

// GetEmail retrieves the authenticated user's email from the Gin context.
func GetEmail(c *gin.Context) string {
	return c.GetString("email")
}

// BindJSON parses the request body into the given struct.
// Returns false and writes a 400 response if binding fails.
// Usage: if !ctxutil.BindJSON(c, &req) { return }
func BindJSON(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
