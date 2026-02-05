package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// RequestHandler 请求处理器
type RequestHandler struct {
	storage *storage.StorageFacade
}

// NewRequestHandler 创建新的请求处理器
func NewRequestHandler(storage *storage.StorageFacade) *RequestHandler {
	return &RequestHandler{
		storage: storage,
	}
}

// ListRequests godoc
// @Summary List all requests
// @Description Get a list of API request records with optional filters (time range, endpoint, status) and limit
// @Tags requests
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results (default: 100)"
// @Param start_time query string false "Start time (RFC3339 format, e.g., 2024-01-01T00:00:00Z)"
// @Param end_time query string false "End time (RFC3339 format, e.g., 2024-01-31T23:59:59Z)"
// @Param endpoint query string false "Endpoint path (supports partial match)"
// @Param status query string false "Request status"
// @Success 200 {object} map[string]interface{} "List of requests"
// @Router /api/v1/requests [get]
func (h *RequestHandler) ListRequests(c *gin.Context) {
	params := &storage.RequestQueryParams{}
	
	// 解析 limit 参数
	limitStr := c.Query("limit")
	params.Limit = 100 // 默认值
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			params.Limit = l
		}
	}
	
	// 解析时间范围参数
	startTimeStr := c.Query("start_time")
	if startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			params.StartTime = &t
		}
	}
	
	endTimeStr := c.Query("end_time")
	if endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			params.EndTime = &t
		}
	}
	
	// 解析请求路径参数
	params.Endpoint = c.Query("endpoint")
	
	// 解析状态参数
	params.Status = c.Query("status")

	requests, err := h.storage.ListRequests(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": requests})
}

// GetRequest 获取单个请求
func (h *RequestHandler) GetRequest(c *gin.Context) {
	id := c.Param("id")
	request, err := h.storage.GetRequest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": request})
}

// CreateRequest 创建请求记录
func (h *RequestHandler) CreateRequest(c *gin.Context) {
	var request storage.Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.CreateRequest(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": request})
}