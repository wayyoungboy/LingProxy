package handler

import (
	"net/http"
	"time"

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

// GetLLMResourceUsageStats 获取LLM资源使用统计（按资源分组）
func (h *StatsHandler) GetLLMResourceUsageStats(c *gin.Context) {
	// 获取所有请求记录
	requests, err := h.storage.ListRequests(100000) // 获取足够多的记录
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get requests: " + err.Error()})
		return
	}

	// 获取所有LLM资源
	resources, err := h.storage.ListLLMResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get llm resources: " + err.Error()})
		return
	}

	// 创建资源映射
	resourceMap := make(map[string]*storage.LLMResource)
	for _, resource := range resources {
		resourceMap[resource.ID] = resource
	}

	// 按资源分组统计
	usageMap := make(map[string]map[string]interface{})
	for _, request := range requests {
		resourceID := request.LLMResourceID
		if resourceID == "" {
			continue // 跳过没有关联资源的请求
		}

		if _, exists := usageMap[resourceID]; !exists {
			resource, ok := resourceMap[resourceID]
			if !ok {
				continue // 资源不存在，跳过
			}
			usageMap[resourceID] = map[string]interface{}{
				"resource_id":      resourceID,
				"resource_name":    resource.Name,
				"resource_type":    resource.Type,
				"model":            resource.Model,
				"total_tokens":     0,
				"total_requests":   0,
				"success_requests": 0,
				"failed_requests": 0,
				"last_request_time": nil,
			}
		}

		usage := usageMap[resourceID]
		usage["total_tokens"] = usage["total_tokens"].(int) + request.Tokens
		usage["total_requests"] = usage["total_requests"].(int) + 1

		if request.Status == "success" {
			usage["success_requests"] = usage["success_requests"].(int) + 1
		} else {
			usage["failed_requests"] = usage["failed_requests"].(int) + 1
		}

		// 更新最后请求时间
		lastTime, ok := usage["last_request_time"].(time.Time)
		if !ok || request.CreatedAt.After(lastTime) {
			usage["last_request_time"] = request.CreatedAt
		}
	}

	// 转换为数组并计算成功率
	var usageList []map[string]interface{}
	for _, usage := range usageMap {
		totalRequests := usage["total_requests"].(int)
		successRequests := usage["success_requests"].(int)
		totalTokens := usage["total_tokens"].(int)

		var successRate float64
		if totalRequests > 0 {
			successRate = float64(successRequests) / float64(totalRequests) * 100
		}

		var avgTokensPerRequest float64
		if totalRequests > 0 {
			avgTokensPerRequest = float64(totalTokens) / float64(totalRequests)
		}

		usage["success_rate"] = successRate
		usage["avg_tokens_per_request"] = avgTokensPerRequest
		usageList = append(usageList, usage)
	}

	c.JSON(http.StatusOK, gin.H{"data": usageList})
}
