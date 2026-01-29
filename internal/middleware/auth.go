package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	storage *storage.StorageFacade
}

// NewAuthMiddleware 创建新的认证中间件
func NewAuthMiddleware(storage *storage.StorageFacade) *AuthMiddleware {
	return &AuthMiddleware{
		storage: storage,
	}
}

// RequireAuth 需要认证
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		apiKey := parts[1]
		user, err := m.storage.GetUserByAPIKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}