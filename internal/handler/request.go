package handler

import (
	"net/http"
	"strconv"

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
// @Description Get a list of API request records with optional limit
// @Tags requests
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results (default: 100)"
// @Success 200 {object} map[string]interface{} "List of requests"
// @Router /api/v1/requests [get]
func (h *RequestHandler) ListRequests(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	requests, err := h.storage.ListRequests(limit)
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