package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pkgjwt "github.com/wingc34/mini-commerce-api/pkg/jwt"
)

// Auth validates the JWT token in the Authorization header.
// On success, sets "userID" and "email" in the gin context for downstream handlers.
// On failure, aborts with 401 Unauthorized.
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 第一步：拿到 Authorization header
		// 格式是 "Bearer <token>"
		authHeader := c.GetHeader("Authorization")

		// 第二步：檢查 header 是否存在，格式是否正確
		// Authorization: Bearer <token>
		// strings.HasPrefix 檢查開頭
		// strings.TrimPrefix 去掉 "Bearer " 拿到 token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 第三步：驗證 JWT token
		claims, err := pkgjwt.Verify(token, jwtSecret)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 第四步：把 userID 存進 context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		// 第五步：繼續往下
		c.Next()
	}
}
