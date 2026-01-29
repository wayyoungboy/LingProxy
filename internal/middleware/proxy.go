package middleware

import (
	"github.com/gin-gonic/gin"
)

// ProxyMiddleware 代理中间件
type ProxyMiddleware struct{}

// NewProxyMiddleware 创建新的代理中间件
func NewProxyMiddleware() *ProxyMiddleware {
	return &ProxyMiddleware{}
}

// Proxy 代理处理
func (m *ProxyMiddleware) Proxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Proxy endpoint",
			"path":    c.Param("proxyPath"),
			"method":  c.Request.Method,
		})
	}
}