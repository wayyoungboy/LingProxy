package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// StatsHandler 统计处理器
type StatsHandler struct {
	storage *storage.StorageFacade
}

// NewStatsHandler 创建新的统计处理器
func NewStatsHandler(storage *storage.StorageFacade) *StatsHandler {
	return &StatsHandler{
		storage: storage,
	}
}

// GetSystemStats 获取系统统计
func (h *StatsHandler) GetSystemStats(c *gin.Context) {
	// 返回模拟数据
	stats := gin.H{
		"total_requests":   1000,
		"total_users":      50,
		"total_llm_resources": 3,
		"success_rate":     98.5,
		"avg_response_time": 120,
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetLLMResourceStats 获取LLM资源统计
func (h *StatsHandler) GetLLMResourceStats(c *gin.Context) {
	llmResourceID := c.Param("id")
	stats := gin.H{
		"llm_resource_id": llmResourceID,
		"total_requests":  500,
		"success_rate":    99.2,
		"avg_response_time": 110,
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetUserStats 获取用户统计
func (h *StatsHandler) GetUserStats(c *gin.Context) {
	userID := c.Param("id")
	stats := gin.H{
		"user_id":        userID,
		"total_requests": 100,
		"total_tokens":   50000,
		"success_rate":   97.8,
		"avg_response_time": 130,
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}