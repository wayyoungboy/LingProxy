package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
)

const (
	// RequestIDKey 请求ID在gin.Context中的key
	RequestIDKey = "request_id"
	// RequestIDHeader 请求ID的HTTP头名称
	RequestIDHeader = "X-Request-ID"
)

// RequestID 请求ID中间件，为每个请求生成唯一ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取Request ID
		requestID := c.GetHeader(RequestIDHeader)

		// 如果没有，生成一个新的
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 设置到上下文和响应头
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// generateRequestID 生成请求ID（16字节随机数，32位十六进制字符串）
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GetRequestID 从gin.Context获取Request ID
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(RequestIDKey); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}

// RequestLogger HTTP请求日志中间件（优化版）
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		requestID := GetRequestID(c)

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 构建基础字段
		baseFields := []logger.Field{
			logger.F("component", "middleware"),
			logger.F("request_id", requestID),
			logger.F("method", method),
			logger.F("path", path),
			logger.F("status", statusCode),
			logger.F("latency_ms", latency.Milliseconds()),
			logger.F("ip", clientIP),
		}

		// 根据状态码记录不同级别的日志
		if statusCode >= 500 {
			if errorMessage != "" {
				baseFields = append(baseFields, logger.F("error", errorMessage))
			}
			logger.Error("HTTP request failed", baseFields...)
		} else if statusCode >= 400 {
			if errorMessage != "" {
				baseFields = append(baseFields, logger.F("error", errorMessage))
			}
			logger.Warn("HTTP request client error", baseFields...)
		} else {
			// 成功请求使用DEBUG级别，减少日志量（可通过配置调整）
			logger.Debug("HTTP request", baseFields...)
		}

		if raw != "" {
			path = path + "?" + raw
		}
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
