package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/config"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	storage      *storage.StorageFacade
	tokenService *service.TokenService
	config       *config.Config
}

// NewAuthMiddleware 创建新的认证中间件
func NewAuthMiddleware(storage *storage.StorageFacade, tokenService *service.TokenService, cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		storage:      storage,
		tokenService: tokenService,
		config:       cfg,
	}
}

// RequireAuth 需要认证（如果认证开关关闭，则跳过认证检查）
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查认证开关
		if !m.config.Security.Auth.Enabled {
			// 认证已禁用，直接通过
			c.Next()
			return
		}

		// 认证已启用，执行认证检查
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// 解析 Bearer API Key
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenValue := parts[1]

		// 首先尝试通过API Key验证
		token, err := m.tokenService.ValidateToken(tokenValue)
		if err == nil {
			// API Key验证成功
			c.Set("token", token)
			c.Next()
			return
		}

		// 如果API Key验证失败，尝试通过User的APIKey验证（向后兼容）
		user, err := m.storage.GetUserByAPIKey(tokenValue)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key or token"})
			c.Abort()
			return
		}

		// 检查用户状态
		if user.Status != "active" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user account is not active"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
