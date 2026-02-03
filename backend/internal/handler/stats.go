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
	// 从数据库获取真实数据
	users, err := h.storage.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	resources, err := h.storage.ListLLMResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get llm resources"})
		return
	}

	requests, err := h.storage.ListRequests(1000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get requests"})
		return
	}

	// 计算统计数据
	totalUsers := len(users)
	totalLLMResources := len(resources)
	totalRequests := len(requests)

	// 计算成功率
	successCount := 0
	totalDuration := int64(0)
	for _, req := range requests {
		if req.Status == "success" {
			successCount++
		}
		totalDuration += req.Duration
	}

	var successRate float64
	if totalRequests > 0 {
		successRate = float64(successCount) / float64(totalRequests) * 100
	}

	var avgResponseTime float64
	if totalRequests > 0 {
		avgResponseTime = float64(totalDuration) / float64(totalRequests)
	}

	stats := gin.H{
		"total_requests":      totalRequests,
		"total_users":         totalUsers,
		"total_llm_resources": totalLLMResources,
		"success_rate":        successRate,
		"avg_response_time":   avgResponseTime,
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetLLMResourceStats 获取LLM资源统计
func (h *StatsHandler) GetLLMResourceStats(c *gin.Context) {
	llmResourceID := c.Param("id")
	stats := gin.H{
		"llm_resource_id":   llmResourceID,
		"total_requests":    500,
		"success_rate":      99.2,
		"avg_response_time": 110,
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetUserStats 获取用户统计
func (h *StatsHandler) GetUserStats(c *gin.Context) {
	userID := c.Param("id")
	stats := gin.H{
		"user_id":           userID,
		"total_requests":    100,
		"total_tokens":      50000,
		"success_rate":      97.8,
		"avg_response_time": 130,
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
