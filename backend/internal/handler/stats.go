package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/middleware"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
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

	requests, err := h.storage.ListRequests(&storage.RequestQueryParams{Limit: 1000})
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
	requests, err := h.storage.ListRequests(&storage.RequestQueryParams{Limit: 100000}) // 获取足够多的记录
	if err != nil {
		logger.Error("Failed to get requests for usage stats", logger.F("component", "handler"), logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get requests: " + err.Error()})
		return
	}

	logger.Debug("Retrieved requests for usage stats", logger.F("component", "handler"), logger.F("request_count", len(requests)))

	// 获取所有LLM资源
	resources, err := h.storage.ListLLMResources()
	if err != nil {
		logger.Error("Failed to get llm resources for usage stats", logger.F("component", "handler"), logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get llm resources: " + err.Error()})
		return
	}

	logger.Debug("Retrieved llm resources for usage stats", logger.F("component", "handler"), logger.F("resource_count", len(resources)))

	// 创建资源映射
	resourceMap := make(map[string]*storage.LLMResource)
	for _, resource := range resources {
		resourceMap[resource.ID] = resource
	}

	// 按资源分组统计
	usageMap := make(map[string]map[string]interface{})
	skippedRequests := 0
	for _, request := range requests {
		resourceID := request.LLMResourceID
		if resourceID == "" {
			skippedRequests++
			continue // 跳过没有关联资源的请求
		}

		if _, exists := usageMap[resourceID]; !exists {
			resource, ok := resourceMap[resourceID]
			if !ok {
				skippedRequests++
				logger.Debug("Request references non-existent resource", logger.F("component", "handler"), logger.F("resource_id", resourceID), logger.F("request_id", request.ID))
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

	if skippedRequests > 0 {
		logger.Debug("Skipped requests without valid resource", logger.F("component", "handler"), logger.F("skipped_count", skippedRequests))
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

	logger.Debug("Usage stats calculated", logger.F("component", "handler"), logger.F("usage_items_count", len(usageList)))

	// 即使没有数据也返回空数组，而不是错误
	c.JSON(http.StatusOK, gin.H{"data": usageList})
}

// GetMonitorStats 获取监控数据（时间序列，供 ECharts 使用）
func (h *StatsHandler) GetMonitorStats(c *gin.Context) {
	period := c.DefaultQuery("period", "1h")
	limit := 100000 // 获取足够多的请求记录

	requests, err := h.storage.ListRequests(&storage.RequestQueryParams{Limit: limit})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get requests"})
		return
	}

	// 解析时间范围
	var duration time.Duration
	switch period {
	case "1h":
		duration = 1 * time.Hour
	case "6h":
		duration = 6 * time.Hour
	case "24h":
		duration = 24 * time.Hour
	case "7d":
		duration = 7 * 24 * time.Hour
	default:
		duration = 1 * time.Hour
	}

	now := time.Now()
	windowStart := now.Add(-duration)

	// 过滤在时间窗口内的请求
	var filteredRequests []*storage.Request
	for _, req := range requests {
		if req.CreatedAt.After(windowStart) {
			filteredRequests = append(filteredRequests, req)
		}
	}

	// 按时间粒度聚合数据
	interval := h.calculateInterval(duration)
	timeBuckets := make(map[int64]*TimeBucket)

	for _, req := range filteredRequests {
		bucketKey := req.CreatedAt.Truncate(interval).UnixMilli()
		if _, exists := timeBuckets[bucketKey]; !exists {
			timeBuckets[bucketKey] = &TimeBucket{
				Timestamp:    time.UnixMilli(bucketKey).Format(time.RFC3339),
				TotalReqs:    0,
				SuccessReqs:  0,
				FailedReqs:   0,
				TotalTokens:  0,
				TotalDurMs:   0,
				EndpointHits: make(map[string]int),
			}
		}
		bucket := timeBuckets[bucketKey]
		bucket.TotalReqs++
		bucket.TotalTokens += req.Tokens
		bucket.TotalDurMs += req.Duration
		if req.Status == "success" {
			bucket.SuccessReqs++
		} else {
			bucket.FailedReqs++
		}
		if req.Endpoint != "" {
			bucket.EndpointHits[req.Endpoint]++
		}
	}

	// 转换为有序数组
	var timeline []TimeBucket
	for _, bucket := range timeBuckets {
		if bucket.TotalReqs > 0 {
			bucket.AvgDurMs = bucket.TotalDurMs / int64(bucket.TotalReqs)
			bucket.SuccessRate = float64(bucket.SuccessReqs) / float64(bucket.TotalReqs) * 100
		}
		timeline = append(timeline, *bucket)
	}

	// 按时间排序
	for i := 0; i < len(timeline); i++ {
		for j := i + 1; j < len(timeline); j++ {
			if timeline[i].Timestamp > timeline[j].Timestamp {
				timeline[i], timeline[j] = timeline[j], timeline[i]
			}
		}
	}

	// 填充缺失的时间点（让图表更平滑）
	timeline = h.fillMissingBuckets(timeline, windowStart, now, interval)

	// 获取限流器统计
	var rateLimiterStats map[string]interface{}
	if middleware.GlobalRateLimiter != nil {
		rateLimiterStats = middleware.GlobalRateLimiter.GetStats()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"period":        period,
			"interval_ms":   interval.Milliseconds(),
			"timeline":      timeline,
			"rate_limiter":  rateLimiterStats,
			"total_points":  len(timeline),
			"total_requests": len(filteredRequests),
		},
	})
}

// TimeBucket 时间聚合桶
type TimeBucket struct {
	Timestamp    string         `json:"timestamp"`
	TotalReqs    int            `json:"total_requests"`
	SuccessReqs  int            `json:"success_requests"`
	FailedReqs   int            `json:"failed_requests"`
	SuccessRate  float64        `json:"success_rate"`
	TotalTokens  int            `json:"total_tokens"`
	TotalDurMs   int64          `json:"total_duration_ms"`
	AvgDurMs     int64          `json:"avg_duration_ms"`
	EndpointHits map[string]int `json:"endpoint_hits"`
}

// TopEndpoint 高频端点
type TopEndpoint struct {
	Endpoint string `json:"endpoint"`
	Count    int    `json:"count"`
}

func (h *StatsHandler) calculateInterval(duration time.Duration) time.Duration {
	switch {
	case duration <= time.Hour:
		return 1 * time.Minute
	case duration <= 6*time.Hour:
		return 5 * time.Minute
	case duration <= 24*time.Hour:
		return 15 * time.Minute
	default:
		return 1 * time.Hour
	}
}

func (h *StatsHandler) fillMissingBuckets(timeline []TimeBucket, windowStart, now time.Time, interval time.Duration) []TimeBucket {
	if len(timeline) == 0 {
		return timeline
	}

	filled := make([]TimeBucket, 0)
	bucketSet := make(map[int64]bool)
	for _, b := range timeline {
		t, _ := time.Parse(time.RFC3339, b.Timestamp)
		bucketSet[t.UnixMilli()] = true
	}

	start := windowStart.Truncate(interval)
	end := now.Truncate(interval)

	for t := start; !t.After(end); t = t.Add(interval) {
		key := t.UnixMilli()
		if !bucketSet[key] {
			filled = append(filled, TimeBucket{
				Timestamp:   t.Format(time.RFC3339),
				TotalReqs:   0,
				SuccessReqs: 0,
				FailedReqs:  0,
				SuccessRate: 0,
				TotalTokens: 0,
				AvgDurMs:    0,
			})
		}
	}

	filled = append(filled, timeline...)

	// 重新排序
	for i := 0; i < len(filled); i++ {
		for j := i + 1; j < len(filled); j++ {
			if filled[i].Timestamp > filled[j].Timestamp {
				filled[i], filled[j] = filled[j], filled[i]
			}
		}
	}

	return filled
}
