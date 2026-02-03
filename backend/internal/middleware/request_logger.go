package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
)

// RequestLogger HTTP请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 记录日志
		if statusCode >= 500 {
			logger.Error("HTTP Request",
				logger.F("method", method),
				logger.F("path", path),
				logger.F("status", statusCode),
				logger.F("latency", latency),
				logger.F("ip", clientIP),
				logger.F("error", errorMessage),
			)
		} else if statusCode >= 400 {
			logger.Warn("HTTP Request",
				logger.F("method", method),
				logger.F("path", path),
				logger.F("status", statusCode),
				logger.F("latency", latency),
				logger.F("ip", clientIP),
			)
		} else {
			logger.Info("HTTP Request",
				logger.F("method", method),
				logger.F("path", path),
				logger.F("status", statusCode),
				logger.F("latency", latency),
				logger.F("ip", clientIP),
			)
		}

		if raw != "" {
			path = path + "?" + raw
		}
	}
}

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
