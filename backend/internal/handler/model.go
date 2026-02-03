package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// ModelHandler 模型处理器
type ModelHandler struct {
	storage *storage.StorageFacade
}

// NewModelHandler 创建新的模型处理器
func NewModelHandler(storage *storage.StorageFacade) *ModelHandler {
	return &ModelHandler{storage: storage}
}

// ListModels godoc
// @Summary List all models
// @Description Get a list of all AI models with filtering options
// @Tags models
// @Accept json
// @Produce json
// @Param provider query string false "Filter by provider ID"
// @Param type query string false "Filter by model type (chat, completion, embedding, image)"
// @Param category query string false "Filter by model category"
// @Param status query string false "Filter by status (active, inactive, deprecated)"
// @Success 200 {object} map[string]interface{} "List of models"
// @Router /api/v1/models [get]
func (h *ModelHandler) ListModels(c *gin.Context) {
	llmResource := c.Query("llm_resource")

	var models []*storage.Model
	var err error

	if llmResource != "" {
		models, err = h.storage.ListModelsByLLMResource(llmResource)
	} else {
		models, err = h.storage.ListModels()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": models})
}

// GetModel godoc
// @Summary Get model by ID
// @Description Get detailed information about a specific model
// @Tags models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Success 200 {object} map[string]interface{} "Model details"
// @Failure 404 {object} map[string]string "Model not found"
// @Router /api/v1/models/{id} [get]
func (h *ModelHandler) GetModel(c *gin.Context) {
	id := c.Param("id")
	model, err := h.storage.GetModel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": model})
}

// CreateModel godoc
// @Summary Create a new model
// @Description Create a new AI model configuration with detailed parameters
// @Tags models
// @Accept json
// @Produce json
// @Param model body storage.Model true "Model configuration"
// @Success 201 {object} map[string]interface{} "Created model"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/models [post]
func (h *ModelHandler) CreateModel(c *gin.Context) {
	var model storage.Model
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if model.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "model name is required"})
		return
	}
	if model.LLMResourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "llm_resource_id is required"})
		return
	}
	if model.ModelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "model_id is required"})
		return
	}

	// 验证LLM资源是否存在
	_, err := h.storage.GetLLMResource(model.LLMResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "llm resource not found"})
		return
	}

	// 设置默认值
	if model.Type == "" {
		model.Type = "chat"
	}
	if model.Status == "" {
		model.Status = "active"
	}

	if err := h.storage.CreateModel(&model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": model})
}

// UpdateModel godoc
// @Summary Update model
// @Description Update an existing model configuration
// @Tags models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Param model body storage.Model true "Updated model configuration"
// @Success 200 {object} map[string]interface{} "Updated model"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/models/{id} [put]
func (h *ModelHandler) UpdateModel(c *gin.Context) {
	id := c.Param("id")
	var model storage.Model
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model.ID = id
	if err := h.storage.UpdateModel(&model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": model})
}

// DeleteModel godoc
// @Summary Delete model
// @Description Delete a model by ID
// @Tags models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/models/{id} [delete]
func (h *ModelHandler) DeleteModel(c *gin.Context) {
	id := c.Param("id")
	if err := h.storage.DeleteModel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "model deleted successfully"})
}

// ListModelTypes godoc
// @Summary List model types
// @Description Get available model types and categories
// @Tags models
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Model types and categories"
// @Router /api/v1/models/types [get]
func (h *ModelHandler) ListModelTypes(c *gin.Context) {
	types := map[string][]string{
		"types":      {"chat", "completion", "embedding", "image", "audio", "moderation"},
		"categories": {"gpt", "claude", "gemini", "llama", "mistral", "custom"},
	}
	c.JSON(http.StatusOK, types)
}

// GetModelPricing godoc
// @Summary Get model pricing
// @Description Get pricing information for a specific model
// @Tags models
// @Accept json
// @Produce json
// @Param id path string true "Model ID"
// @Success 200 {object} map[string]interface{} "Model pricing"
// @Router /api/v1/models/{id}/pricing [get]
func (h *ModelHandler) GetModelPricing(c *gin.Context) {
	id := c.Param("id")
	model, err := h.storage.GetModel(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
		return
	}

	// 解析JSON字符串为对象
	var pricing storage.ModelPricing
	if model.Pricing != "" {
		if err := json.Unmarshal([]byte(model.Pricing), &pricing); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse pricing data"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": pricing})
}

// ListModelsByLLMResource godoc
// @Summary List models by LLM resource
// @Description Get all models for a specific LLM resource
// @Tags models
// @Accept json
// @Produce json
// @Param llm_resource_id path string true "LLM Resource ID"
// @Success 200 {object} map[string]interface{} "List of models"
// @Router /api/v1/llm-resources/{llm_resource_id}/models [get]
func (h *ModelHandler) ListModelsByLLMResource(c *gin.Context) {
	llmResourceID := c.Param("id")
	models, err := h.storage.ListModelsByLLMResource(llmResourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": models})
}
