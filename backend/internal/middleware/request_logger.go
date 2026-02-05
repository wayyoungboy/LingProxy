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

// RequestLogger HTTP请求日志中间件（记录所有请求）
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		requestID := GetRequestID(c)
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()
		contentType := c.Request.Header.Get("Content-Type")
		contentLength := c.Request.ContentLength

		// 构建完整路径（包含查询参数）
		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		// 记录请求开始日志
		requestFields := []logger.Field{
			logger.F("component", "middleware"),
			logger.F("request_id", requestID),
			logger.F("method", method),
			logger.F("path", fullPath),
			logger.F("ip", clientIP),
		}

		// 添加可选字段
		if userAgent != "" {
			requestFields = append(requestFields, logger.F("user_agent", userAgent))
		}
		if contentType != "" {
			requestFields = append(requestFields, logger.F("content_type", contentType))
		}
		if contentLength > 0 {
			requestFields = append(requestFields, logger.F("content_length", contentLength))
		}

		logger.Info("Incoming request", requestFields...)

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		responseSize := c.Writer.Size()

		// 构建响应日志字段
		responseFields := []logger.Field{
			logger.F("component", "middleware"),
			logger.F("request_id", requestID),
			logger.F("method", method),
			logger.F("path", fullPath),
			logger.F("status", statusCode),
			logger.F("latency_ms", latency.Milliseconds()),
			logger.F("ip", clientIP),
		}

		if responseSize > 0 {
			responseFields = append(responseFields, logger.F("response_size", responseSize))
		}

		// 根据状态码记录不同级别的日志，但都记录
		if statusCode >= 500 {
			if errorMessage != "" {
				responseFields = append(responseFields, logger.F("error", errorMessage))
			}
			logger.Error("Request completed with server error", responseFields...)
		} else if statusCode >= 400 {
			if errorMessage != "" {
				responseFields = append(responseFields, logger.F("error", errorMessage))
			}
			logger.Warn("Request completed with client error", responseFields...)
		} else {
			// 成功请求使用INFO级别，确保所有请求都被记录
			logger.Info("Request completed successfully", responseFields...)
		}
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
