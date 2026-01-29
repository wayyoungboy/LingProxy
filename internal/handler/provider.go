package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/storage"
	"github.com/lingproxy/lingproxy/pkg/logger"
)

// LLMResourceHandler LLM资源处理器
type LLMResourceHandler struct {
	storage *storage.StorageFacade
}

// NewLLMResourceHandler 创建新的LLM资源处理器
func NewLLMResourceHandler(storage *storage.StorageFacade) *LLMResourceHandler {
	return &LLMResourceHandler{
		storage: storage,
	}
}

// ListLLMResources godoc
// @Summary List all LLM resources
// @Description Get a list of all AI service LLM resources
// @Tags llm-resources
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of LLM resources"
// @Router /api/v1/llm-resources [get]
func (h *LLMResourceHandler) ListLLMResources(c *gin.Context) {
	resources, err := h.storage.ListLLMResources()
	if err != nil {
		logger.Error("获取LLM资源列表失败", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Info("获取LLM资源列表成功", "count", len(resources))
	c.JSON(http.StatusOK, gin.H{"data": resources})
}

// GetLLMResource godoc
// @Summary Get LLM resource by ID
// @Description Get a specific AI service LLM resource by ID
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param id path string true "LLM Resource ID"
// @Success 200 {object} map[string]interface{} "LLM resource details"
// @Failure 404 {object} map[string]string "LLM resource not found"
// @Router /api/v1/llm-resources/{id} [get]
func (h *LLMResourceHandler) GetLLMResource(c *gin.Context) {
	id := c.Param("id")
	resource, err := h.storage.GetLLMResource(id)
	if err != nil {
		logger.Warn("获取LLM资源失败：资源不存在", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
		return
	}
	logger.Info("获取LLM资源成功", "id", id, "name", resource.Name)
	c.JSON(http.StatusOK, gin.H{"data": resource})
}

// CreateLLMResource godoc
// @Summary Create a new LLM resource
// @Description Create a new AI service LLM resource configuration
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param resource body storage.LLMResource true "LLM resource configuration"
// @Success 201 {object} map[string]interface{} "Created LLM resource"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources [post]
func (h *LLMResourceHandler) CreateLLMResource(c *gin.Context) {
	var resource storage.LLMResource
	if err := c.ShouldBindJSON(&resource); err != nil {
		logger.Error("创建LLM资源失败", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if resource.Name == "" {
		logger.Warn("创建LLM资源失败：名称为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "名称是必填项"})
		return
	}
	if resource.Type == "" {
		logger.Warn("创建LLM资源失败：类型为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "类型是必填项"})
		return
	}
	if resource.Model == "" {
		logger.Warn("创建LLM资源失败：模型为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型是必填项"})
		return
	}
	if resource.BaseURL == "" {
		logger.Warn("创建LLM资源失败：Base URL为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Base URL是必填项"})
		return
	}
	if resource.APIKey == "" {
		logger.Warn("创建LLM资源失败：API Key为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "API Key是必填项"})
		return
	}

	if err := h.storage.CreateLLMResource(&resource); err != nil {
		logger.Error("创建LLM资源失败", "error", err.Error(), "name", resource.Name)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("创建LLM资源成功", "id", resource.ID, "name", resource.Name, "model", resource.Model)
	c.JSON(http.StatusCreated, gin.H{"data": resource})
}

// UpdateLLMResource godoc
// @Summary Update an existing LLM resource
// @Description Update an existing AI service LLM resource configuration
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param id path string true "LLM Resource ID"
// @Param resource body storage.LLMResource true "LLM resource configuration"
// @Success 200 {object} map[string]interface{} "Updated LLM resource"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "LLM resource not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources/{id} [put]
func (h *LLMResourceHandler) UpdateLLMResource(c *gin.Context) {
	id := c.Param("id")
	var resource storage.LLMResource
	if err := c.ShouldBindJSON(&resource); err != nil {
		logger.Error("更新LLM资源失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if resource.Name == "" {
		logger.Warn("更新LLM资源失败：名称为空", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "名称是必填项"})
		return
	}
	if resource.Type == "" {
		logger.Warn("更新LLM资源失败：类型为空", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "类型是必填项"})
		return
	}
	if resource.Model == "" {
		logger.Warn("更新LLM资源失败：模型为空", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型是必填项"})
		return
	}
	if resource.BaseURL == "" {
		logger.Warn("更新LLM资源失败：Base URL为空", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Base URL是必填项"})
		return
	}
	if resource.APIKey == "" {
		logger.Warn("更新LLM资源失败：API Key为空", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "API Key是必填项"})
		return
	}

	// 确保资源ID与路径参数一致
	resource.ID = id

	if err := h.storage.UpdateLLMResource(&resource); err != nil {
		if err == storage.ErrNotFound {
			logger.Warn("更新LLM资源失败：资源不存在", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
			return
		}
		logger.Error("更新LLM资源失败", "error", err.Error(), "id", id, "name", resource.Name)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("更新LLM资源成功", "id", id, "name", resource.Name, "model", resource.Model)
	c.JSON(http.StatusOK, gin.H{"data": resource})
}

// DeleteLLMResource godoc
// @Summary Delete an LLM resource
// @Description Delete an AI service LLM resource configuration
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param id path string true "LLM Resource ID"
// @Success 204 {object} nil "LLM resource deleted"
// @Failure 404 {object} map[string]string "LLM resource not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources/{id} [delete]
func (h *LLMResourceHandler) DeleteLLMResource(c *gin.Context) {
	id := c.Param("id")

	if err := h.storage.DeleteLLMResource(id); err != nil {
		if err == storage.ErrNotFound {
			logger.Warn("删除LLM资源失败：资源不存在", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
			return
		}
		logger.Error("删除LLM资源失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("删除LLM资源成功", "id", id)
	c.Status(http.StatusNoContent)
}
