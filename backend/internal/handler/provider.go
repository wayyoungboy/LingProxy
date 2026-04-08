package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
	"github.com/xuri/excelize/v2"
)

// LLMResourceHandler LLM资源处理器
type LLMResourceHandler struct {
	storage       *storage.StorageFacade
	openaiService *service.OpenAIService
}

// LLMResourceRequest 资源请求结构（包含类型特定配置）
type LLMResourceRequest struct {
	storage.LLMResource
	EmbeddingConfig *storage.EmbeddingConfig `json:"embedding_config,omitempty"`
	RerankConfig    *storage.RerankConfig    `json:"rerank_config,omitempty"`
	ChatConfig      *storage.ChatConfig       `json:"chat_config,omitempty"`
	ImageConfig     map[string]interface{}   `json:"image_config,omitempty"`
	AudioConfig     map[string]interface{}   `json:"audio_config,omitempty"`
}

// processTypeConfig 处理类型特定配置，将对象转换为JSON字符串存储
func (h *LLMResourceHandler) processTypeConfig(resource *storage.LLMResource, req *LLMResourceRequest) error {
	switch resource.Type {
	case "embedding":
		if req.EmbeddingConfig != nil {
			configJSON, err := json.Marshal(req.EmbeddingConfig)
			if err != nil {
				return fmt.Errorf("failed to marshal embedding config: %w", err)
			}
			resource.EmbeddingConfig = string(configJSON)
		}
	case "rerank":
		if req.RerankConfig != nil {
			configJSON, err := json.Marshal(req.RerankConfig)
			if err != nil {
				return fmt.Errorf("failed to marshal rerank config: %w", err)
			}
			resource.RerankConfig = string(configJSON)
		}
	case "chat":
		if req.ChatConfig != nil {
			configJSON, err := json.Marshal(req.ChatConfig)
			if err != nil {
				return fmt.Errorf("failed to marshal chat config: %w", err)
			}
			resource.ChatConfig = string(configJSON)
		}
	case "image":
		if req.ImageConfig != nil {
			configJSON, err := json.Marshal(req.ImageConfig)
			if err != nil {
				return fmt.Errorf("failed to marshal image config: %w", err)
			}
			resource.ImageConfig = string(configJSON)
		}
	case "audio":
		if req.AudioConfig != nil {
			configJSON, err := json.Marshal(req.AudioConfig)
			if err != nil {
				return fmt.Errorf("failed to marshal audio config: %w", err)
			}
			resource.AudioConfig = string(configJSON)
		}
	}
	return nil
}

// enrichResourceResponse 丰富资源响应，将JSON字符串解析为对象
func (h *LLMResourceHandler) enrichResourceResponse(resource *storage.LLMResource) map[string]interface{} {
	result := map[string]interface{}{
		"id":          resource.ID,
		"name":        resource.Name,
		"type":        resource.Type,
		"driver":      resource.Driver,
		"model":       resource.Model,
		"base_url":    resource.BaseURL,
		"api_key":     resource.APIKey,
		"status":      resource.Status,
		"test_status": resource.TestStatus,
		"created_at":  resource.CreatedAt,
		"updated_at":  resource.UpdatedAt,
	}

	// 根据类型解析对应的配置
	switch resource.Type {
	case "embedding":
		if resource.EmbeddingConfig != "" {
			var config storage.EmbeddingConfig
			if err := json.Unmarshal([]byte(resource.EmbeddingConfig), &config); err == nil {
				result["embedding_config"] = config
			}
		}
	case "rerank":
		if resource.RerankConfig != "" {
			var config storage.RerankConfig
			if err := json.Unmarshal([]byte(resource.RerankConfig), &config); err == nil {
				result["rerank_config"] = config
			}
		}
	case "chat":
		if resource.ChatConfig != "" {
			var config storage.ChatConfig
			if err := json.Unmarshal([]byte(resource.ChatConfig), &config); err == nil {
				result["chat_config"] = config
			}
		}
	case "image":
		if resource.ImageConfig != "" {
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(resource.ImageConfig), &config); err == nil {
				result["image_config"] = config
			}
		}
	case "audio":
		if resource.AudioConfig != "" {
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(resource.AudioConfig), &config); err == nil {
				result["audio_config"] = config
			}
		}
	}

	return result
}

// NewLLMResourceHandler 创建新的LLM资源处理器
func NewLLMResourceHandler(storage *storage.StorageFacade) *LLMResourceHandler {
	return &LLMResourceHandler{
		storage:       storage,
		openaiService: service.NewOpenAIService(),
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
		logger.Error("获取LLM资源列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Info("获取LLM资源列表成功", logger.F("count", len(resources)))
	
	// 丰富响应数据，包含类型特定配置
	enrichedResources := make([]map[string]interface{}, len(resources))
	for i, resource := range resources {
		enrichedResources[i] = h.enrichResourceResponse(resource)
	}
	
	c.JSON(http.StatusOK, gin.H{"data": enrichedResources})
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
		logger.Warn("获取LLM资源失败：资源不存在", logger.F("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
		return
	}
	logger.Info("获取LLM资源成功", logger.F("id", id), logger.F("name", resource.Name))
	
	// 丰富响应数据，包含类型特定配置
	responseData := h.enrichResourceResponse(resource)
	c.JSON(http.StatusOK, gin.H{"data": responseData})
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
	var req LLMResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("创建LLM资源失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resource := req.LLMResource

	// 验证必填字段
	if resource.Name == "" {
		logger.Warn("创建LLM资源失败：名称为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "名称是必填项"})
		return
	}
	if resource.Type == "" {
		logger.Warn("创建LLM资源失败：模型类别为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型类别是必填项"})
		return
	}
	// 验证资源类型
	typeValLower := strings.ToLower(resource.Type)
	validTypes := map[string]bool{
		"chat":       true,
		"completion": true,
		"embedding":  true,
		"image":      true,
		"audio":      true,
		"moderation": true,
		"rerank":     true,
	}
	if !validTypes[typeValLower] {
		if typeValLower == "rank" {
			logger.Warn("创建LLM资源失败：不支持的资源类型", logger.F("type", resource.Type))
			c.JSON(http.StatusBadRequest, gin.H{"error": "资源类型 'rank' 不支持，请使用 'rerank'"})
		} else {
			logger.Warn("创建LLM资源失败：不支持的资源类型", logger.F("type", resource.Type))
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("不支持的资源类型 '%s'，支持的类型: chat, completion, embedding, image, audio, moderation, rerank", resource.Type)})
		}
		return
	}
	resource.Type = typeValLower // 统一转换为小写
	// 驱动固定为openai
	if resource.Driver == "" {
		resource.Driver = "openai"
	} else if resource.Driver != "openai" {
		logger.Warn("创建LLM资源失败：不支持的驱动", logger.F("driver", resource.Driver))
		c.JSON(http.StatusBadRequest, gin.H{"error": "目前仅支持openai驱动"})
		return
	}
	if resource.Model == "" {
		logger.Warn("创建LLM资源失败：模型标识为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型标识是必填项"})
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
	// 默认测试状态为未测试
	if resource.TestStatus == "" {
		resource.TestStatus = "untested"
	}

	// 处理类型特定配置
	if err := h.processTypeConfig(&resource, &req); err != nil {
		logger.Error("处理类型特定配置失败", logger.F("error", err.Error()), logger.F("name", resource.Name), logger.F("type", resource.Type))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.CreateLLMResource(&resource); err != nil {
		logger.Error("创建LLM资源失败", logger.F("error", err.Error()), logger.F("name", resource.Name))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("创建LLM资源成功", logger.F("id", resource.ID), logger.F("name", resource.Name), logger.F("model", resource.Model))
	
	// 返回时丰富响应数据
	responseData := h.enrichResourceResponse(&resource)
	c.JSON(http.StatusCreated, gin.H{"data": responseData})
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
	var req LLMResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("更新LLM资源失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resource := req.LLMResource
	resource.ID = id

	// 验证必填字段
	if resource.Name == "" {
		logger.Warn("更新LLM资源失败：名称为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "名称是必填项"})
		return
	}
	if resource.Type == "" {
		logger.Warn("更新LLM资源失败：模型类别为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型类别是必填项"})
		return
	}
	// 验证资源类型
	typeValLower := strings.ToLower(resource.Type)
	validTypes := map[string]bool{
		"chat":       true,
		"completion": true,
		"embedding":  true,
		"image":      true,
		"audio":      true,
		"moderation": true,
		"rerank":     true,
	}
	if !validTypes[typeValLower] {
		if typeValLower == "rank" {
			logger.Warn("更新LLM资源失败：不支持的资源类型", logger.F("id", id), logger.F("type", resource.Type))
			c.JSON(http.StatusBadRequest, gin.H{"error": "资源类型 'rank' 不支持，请使用 'rerank'"})
		} else {
			logger.Warn("更新LLM资源失败：不支持的资源类型", logger.F("id", id), logger.F("type", resource.Type))
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("不支持的资源类型 '%s'，支持的类型: chat, completion, embedding, image, audio, moderation, rerank", resource.Type)})
		}
		return
	}
	resource.Type = typeValLower // 统一转换为小写
	// 驱动固定为openai
	if resource.Driver == "" {
		resource.Driver = "openai"
	} else if resource.Driver != "openai" {
		logger.Warn("更新LLM资源失败：不支持的驱动", logger.F("id", id), logger.F("driver", resource.Driver))
		c.JSON(http.StatusBadRequest, gin.H{"error": "目前仅支持openai驱动"})
		return
	}
	if resource.Model == "" {
		logger.Warn("更新LLM资源失败：模型标识为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型标识是必填项"})
		return
	}
	if resource.BaseURL == "" {
		logger.Warn("更新LLM资源失败：Base URL为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Base URL是必填项"})
		return
	}
	if resource.APIKey == "" {
		logger.Warn("更新LLM资源失败：API Key为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "API Key是必填项"})
		return
	}

	// 确保资源ID与路径参数一致
	resource.ID = id

	if err := h.storage.UpdateLLMResource(&resource); err != nil {
		if err == storage.ErrNotFound {
			logger.Warn("更新LLM资源失败：资源不存在", logger.F("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
			return
		}
		logger.Error("更新LLM资源失败", logger.F("error", err.Error()), logger.F("id", id), logger.F("name", resource.Name))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("更新LLM资源成功", logger.F("id", id), logger.F("name", resource.Name), logger.F("model", resource.Model))
	
	// 返回时丰富响应数据
	responseData := h.enrichResourceResponse(&resource)
	c.JSON(http.StatusOK, gin.H{"data": responseData})
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
			logger.Warn("删除LLM资源失败：资源不存在", logger.F("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
			return
		}
		logger.Error("删除LLM资源失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("删除LLM资源成功", logger.F("id", id))
	c.Status(http.StatusNoContent)
}

// TestLLMResource 测试LLM资源是否可以正常调用
// @Summary Test LLM resource
// @Description Test if an LLM resource can be called successfully
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param id path string true "LLM Resource ID"
// @Success 200 {object} map[string]interface{} "Test result"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "LLM resource not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources/{id}/test [post]
func (h *LLMResourceHandler) TestLLMResource(c *gin.Context) {
	id := c.Param("id")
	resource, err := h.storage.GetLLMResource(id)
	if err != nil {
		logger.Warn("测试LLM资源失败：资源不存在", logger.F("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
		return
	}

	// 检查资源状态
	if resource.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Resource is not active",
			"message": "资源状态为禁用，无法测试",
		})
		return
	}

	// 根据资源类型进行测试
	var testResult map[string]interface{}
	switch resource.Type {
	case "chat":
		testResult = h.testChatResource(resource)
	case "embedding":
		testResult = h.testEmbeddingResource(resource)
	case "rerank":
		testResult = h.testRerankResource(resource)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Unsupported resource type",
			"message": fmt.Sprintf("资源类型 %s 暂不支持测试", resource.Type),
		})
		return
	}

	// 更新测试状态
	if testResult["success"].(bool) {
		resource.TestStatus = "passed"
		logger.Info("测试LLM资源成功", logger.F("id", id), logger.F("name", resource.Name), logger.F("type", resource.Type))
	} else {
		resource.TestStatus = "failed"
		logger.Warn("测试LLM资源失败", logger.F("id", id), logger.F("name", resource.Name), logger.F("type", resource.Type), logger.F("error", testResult["error"]))
	}

	// 保存更新后的测试状态
	if err := h.storage.UpdateLLMResource(resource); err != nil {
		logger.Error("更新LLM资源测试状态失败", logger.F("id", id), logger.F("error", err.Error()))
		// 即使更新失败也继续返回测试结果
	}

	c.JSON(http.StatusOK, testResult)
}

// testChatResource 测试chat类型资源
func (h *LLMResourceHandler) testChatResource(resource *storage.LLMResource) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Debug("Test chat resource", logger.F("component", "handler"), logger.F("base_url", resource.BaseURL), logger.F("model", resource.Model))

	// 使用统一的服务层进行测试
	return h.openaiService.TestChatResource(ctx, resource)
}

// testEmbeddingResource 测试embedding类型资源
func (h *LLMResourceHandler) testEmbeddingResource(resource *storage.LLMResource) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Debug("Test embedding resource", logger.F("component", "handler"), logger.F("base_url", resource.BaseURL), logger.F("model", resource.Model))

	// 使用统一的服务层进行测试
	return h.openaiService.TestEmbeddingResource(ctx, resource)
}

// testRerankResource 测试rerank类型资源
func (h *LLMResourceHandler) testRerankResource(resource *storage.LLMResource) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Debug("Test rerank resource", logger.F("component", "handler"), logger.F("base_url", resource.BaseURL), logger.F("model", resource.Model))

	// 使用统一的服务层进行测试
	return h.openaiService.TestRerankResource(ctx, resource)
}

// ImportLLMResources 批量导入LLM资源
// @Summary Import LLM resources
// @Description Import multiple LLM resources from Excel file or JSON body
// @Tags llm-resources
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Param file formData file false "Excel file"
// @Param resources body []storage.LLMResource false "LLM resources JSON array"
// @Success 200 {object} map[string]interface{} "Import result"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources/import [post]
func (h *LLMResourceHandler) ImportLLMResources(c *gin.Context) {
	// 优先根据 Content-Type 判断是否为 JSON 导入
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		h.importLLMResourcesFromJSON(c)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.Warn("批量导入LLM资源失败：文件获取失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件获取失败: " + err.Error()})
		return
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		logger.Error("批量导入LLM资源失败：文件打开失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件打开失败: " + err.Error()})
		return
	}
	defer src.Close()

	// 读取Excel文件
	f, err := excelize.OpenReader(src)
	if err != nil {
		logger.Error("批量导入LLM资源失败：Excel文件解析失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Excel文件解析失败: " + err.Error()})
		return
	}
	defer f.Close()

	// 读取第一个工作表
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		logger.Error("批量导入LLM资源失败：读取工作表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取工作表失败: " + err.Error()})
		return
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Excel文件至少需要包含表头和一行数据"})
		return
	}

	// 解析表头
	header := rows[0]
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// 验证必需字段（每个字段至少需要一个别名存在）
	requiredFieldGroups := map[string][]string{
		"资源名称":  {"资源名称", "name"},
		"模型类别":  {"模型类别", "type"},
		"驱动":    {"驱动", "driver"},
		"模型标识":  {"模型标识", "model"},
		"基础URL": {"基础url", "base_url", "baseurl"},
		"API密钥": {"api密钥", "api_key", "apikey"},
	}
	missingFields := []string{}
	for fieldName, aliases := range requiredFieldGroups {
		found := false
		for _, alias := range aliases {
			if _, exists := headerMap[strings.ToLower(alias)]; exists {
				found = true
				break
			}
		}
		if !found {
			missingFields = append(missingFields, fieldName)
		}
	}
	if len(missingFields) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("缺少必需字段: %v", missingFields)})
		return
	}

	// 获取列索引
	nameCol := getColumnIndex(headerMap, []string{"资源名称", "name"})
	typeCol := getColumnIndex(headerMap, []string{"模型类别", "type"})
	driverCol := getColumnIndex(headerMap, []string{"驱动", "driver"})
	modelCol := getColumnIndex(headerMap, []string{"模型标识", "model"})
	baseURLCol := getColumnIndex(headerMap, []string{"基础url", "base_url", "baseurl"})
	apiKeyCol := getColumnIndex(headerMap, []string{"api密钥", "api_key", "apikey"})
	statusCol := getColumnIndex(headerMap, []string{"状态", "status"})

	// 获取所有现有资源，用于重复检查
	existingResources, err := h.storage.ListLLMResources()
	if err != nil {
		logger.Error("批量导入LLM资源失败：获取现有资源列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取现有资源列表失败: " + err.Error()})
		return
	}

	// 解析数据行
	successCount := 0
	failCount := 0
	duplicateCount := 0
	errors := []string{}
	duplicates := []map[string]interface{}{}

	for i, row := range rows[1:] {
		rowNum := i + 2 // Excel行号（从2开始，因为第1行是表头）

		// 跳过空行
		if len(row) == 0 {
			continue
		}

		// 获取字段值并去掉前后空格
		name := strings.TrimSpace(getCellValue(row, nameCol))
		typeVal := strings.TrimSpace(getCellValue(row, typeCol))
		driver := strings.TrimSpace(getCellValue(row, driverCol))
		model := strings.TrimSpace(getCellValue(row, modelCol))
		baseURL := strings.TrimSpace(getCellValue(row, baseURLCol))
		apiKey := strings.TrimSpace(getCellValue(row, apiKeyCol))
		status := strings.TrimSpace(getCellValue(row, statusCol))
		if status == "" {
			status = "active"
		}
		// 驱动默认为openai，如果为空或不是openai则设置为openai
		if driver == "" || strings.ToLower(driver) != "openai" {
			driver = "openai"
		}

		// 验证必填字段
		if name == "" || typeVal == "" || model == "" || baseURL == "" || apiKey == "" {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d行: 必填字段不能为空", rowNum))
			continue
		}

		// 验证资源类型
		typeValLower := strings.ToLower(typeVal)
		validTypes := map[string]bool{
			"chat":       true,
			"completion": true,
			"embedding":  true,
			"image":      true,
			"audio":      true,
			"moderation": true,
			"rerank":     true,
		}
		if !validTypes[typeValLower] {
			failCount++
			if typeValLower == "rank" {
				errors = append(errors, fmt.Sprintf("第%d行: 资源类型 'rank' 不支持，请使用 'rerank'", rowNum))
			} else {
				errors = append(errors, fmt.Sprintf("第%d行: 不支持的资源类型 '%s'，支持的类型: chat, completion, embedding, image, audio, moderation, rerank", rowNum, typeVal))
			}
			continue
		}

		// 创建资源对象（用于重复检查）
		resourceToCheck := &storage.LLMResource{
			Type:    strings.ToLower(typeVal),
			Model:   model,
			BaseURL: baseURL,
			APIKey:  apiKey,
		}

		// 检查是否重复
		if h.isDuplicateResource(resourceToCheck, existingResources) {
			duplicateCount++
			duplicates = append(duplicates, map[string]interface{}{
				"row":      rowNum,
				"name":     name,
				"type":     strings.ToLower(typeVal),
				"model":    model,
				"base_url": baseURL,
			})
			continue
		}

		// 创建资源
		resource := &storage.LLMResource{
			Name:       name,
			Type:       strings.ToLower(typeVal),
			Driver:     strings.ToLower(driver),
			Model:      model,
			BaseURL:    baseURL,
			APIKey:     apiKey,
			Status:     strings.ToLower(status),
			TestStatus: "untested", // 导入的资源默认为未测试状态
		}

		if err := h.storage.CreateLLMResource(resource); err != nil {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d行: %s", rowNum, err.Error()))
			continue
		}

		// 将新创建的资源添加到现有资源列表，避免后续重复检查时重复添加
		existingResources = append(existingResources, resource)
		successCount++
	}

	logger.Info("批量导入LLM资源完成", logger.F("success", successCount), logger.F("fail", failCount), logger.F("duplicate", duplicateCount))
	c.JSON(http.StatusOK, gin.H{
		"message":    "导入完成",
		"success":    successCount,
		"fail":       failCount,
		"duplicate":  duplicateCount,
		"errors":     errors,
		"duplicates": duplicates,
		"total":      successCount + failCount + duplicateCount,
	})
}

// importLLMResourcesFromJSON 通过 JSON 批量导入LLM资源
func (h *LLMResourceHandler) importLLMResourcesFromJSON(c *gin.Context) {
	var resources []storage.LLMResource
	if err := c.ShouldBindJSON(&resources); err != nil {
		logger.Warn("批量导入LLM资源失败：JSON解析失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON解析失败: " + err.Error()})
		return
	}

	if len(resources) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据为空"})
		return
	}

	// 获取所有现有资源，用于重复检查
	existingResources, err := h.storage.ListLLMResources()
	if err != nil {
		logger.Error("批量导入LLM资源失败：获取现有资源列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取现有资源列表失败: " + err.Error()})
		return
	}

	successCount := 0
	failCount := 0
	duplicateCount := 0
	errors := []string{}
	duplicates := []map[string]interface{}{}

	for i, item := range resources {
		rowNum := i + 1 // JSON数组下标，从1开始更易读

		// 去掉所有字段的前后空格
		name := strings.TrimSpace(item.Name)
		typeVal := strings.TrimSpace(item.Type)
		driver := strings.TrimSpace(item.Driver)
		model := strings.TrimSpace(item.Model)
		baseURL := strings.TrimSpace(item.BaseURL)
		apiKey := strings.TrimSpace(item.APIKey)
		status := strings.TrimSpace(item.Status)

		if status == "" {
			status = "active"
		}
		// 驱动默认为openai，如果为空或不是openai则设置为openai
		if driver == "" || strings.ToLower(driver) != "openai" {
			driver = "openai"
		}

		// 验证必填字段
		if name == "" || typeVal == "" || model == "" || baseURL == "" || apiKey == "" {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d条: 必填字段不能为空", rowNum))
			continue
		}

		// 验证资源类型
		typeValLower := strings.ToLower(typeVal)
		validTypes := map[string]bool{
			"chat":       true,
			"completion": true,
			"embedding":  true,
			"image":      true,
			"audio":      true,
			"moderation": true,
			"rerank":     true,
		}
		if !validTypes[typeValLower] {
			failCount++
			if typeValLower == "rank" {
				errors = append(errors, fmt.Sprintf("第%d条: 资源类型 'rank' 不支持，请使用 'rerank'", rowNum))
			} else {
				errors = append(errors, fmt.Sprintf("第%d条: 不支持的资源类型 '%s'，支持的类型: chat, completion, embedding, image, audio, moderation, rerank", rowNum, typeVal))
			}
			continue
		}

		// 创建资源对象（用于重复检查）
		resourceToCheck := &storage.LLMResource{
			Type:       strings.ToLower(typeVal),
			Model:      model,
			BaseURL:    baseURL,
			APIKey:     apiKey,
			TestStatus: "untested", // 导入的资源默认为未测试状态
		}

		// 检查是否重复
		if h.isDuplicateResource(resourceToCheck, existingResources) {
			duplicateCount++
			duplicates = append(duplicates, map[string]interface{}{
				"row":      rowNum,
				"name":     name,
				"type":     strings.ToLower(typeVal),
				"model":    model,
				"base_url": baseURL,
			})
			continue
		}

		resource := &storage.LLMResource{
			Name:       name,
			Type:       strings.ToLower(typeVal),
			Driver:     strings.ToLower(driver),
			Model:      model,
			BaseURL:    baseURL,
			APIKey:     apiKey,
			Status:     strings.ToLower(status),
			TestStatus: "untested", // 导入的资源默认为未测试状态
		}

		if err := h.storage.CreateLLMResource(resource); err != nil {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d条: %s", rowNum, err.Error()))
			continue
		}

		// 将新创建的资源添加到现有资源列表，避免后续重复检查时重复添加
		existingResources = append(existingResources, resource)
		successCount++
	}

	logger.Info("批量导入LLM资源(JSON)完成", logger.F("success", successCount), logger.F("fail", failCount), logger.F("duplicate", duplicateCount))
	c.JSON(http.StatusOK, gin.H{
		"message":    "导入完成",
		"success":    successCount,
		"fail":       failCount,
		"duplicate":  duplicateCount,
		"errors":     errors,
		"duplicates": duplicates,
		"total":      successCount + failCount + duplicateCount,
	})
}

// DownloadImportTemplate 下载导入模板
// @Summary Download import template
// @Description Download Excel template for importing LLM resources
// @Tags llm-resources
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} file "Excel template file"
// @Router /api/v1/llm-resources/import/template [get]
func (h *LLMResourceHandler) DownloadImportTemplate(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "LLM资源导入模板"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		logger.Error("创建Excel模板失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建Excel模板失败"})
		return
	}

	// 先设置新sheet为活动sheet，再删除默认Sheet1
	f.SetActiveSheet(index)

	// 删除默认Sheet1
	if err := f.DeleteSheet("Sheet1"); err != nil {
		logger.Warn("删除默认Sheet失败", logger.F("error", err.Error()))
		// 删除失败不影响文件生成，继续执行
	}

	// 设置表头
	headers := []string{"资源名称", "模型类别", "驱动", "模型标识", "基础URL", "API密钥", "状态"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			logger.Error("设置表头失败", logger.F("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "设置表头失败"})
			return
		}
	}

	// 设置示例数据（只包含核心字段，不包含模型元数据）
	examples := [][]interface{}{
		{"OpenAI GPT-4", "chat", "openai", "gpt-4", "https://api.openai.com/v1", "sk-xxxxxxxxxxxxx", "active"},
		{"OpenAI GPT-3.5", "chat", "openai", "gpt-3.5-turbo", "https://api.openai.com/v1", "sk-yyyyyyyyyyyyy", "active"},
		{"OpenAI GPT-4o", "chat", "openai", "gpt-4o", "https://api.openai.com/v1", "sk-zzzzzzzzzzzzz", "active"},
	}

	for rowIdx, example := range examples {
		for colIdx, value := range example {
			cell := fmt.Sprintf("%c%d", 'A'+colIdx, rowIdx+2)
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				logger.Error("设置示例数据失败", logger.F("error", err.Error()))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "设置示例数据失败"})
				return
			}
		}
	}

	// 设置列宽（根据字段长度优化）
	columnWidths := map[string]float64{
		"A": 20, // 资源名称
		"B": 12, // 模型类别
		"C": 10, // 驱动
		"D": 20, // 模型标识
		"E": 30, // 基础URL
		"F": 25, // API密钥
		"G": 10, // 状态
	}
	for col, width := range columnWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			logger.Warn("设置列宽失败", logger.F("column", col), logger.F("error", err.Error()))
		}
	}

	// 设置表头样式
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", "G1", headerStyle)
	}

	// 设置响应头（必须在写入数据之前设置）
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename="llm_resources_import_template.xlsx"`)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Status(http.StatusOK)

	// 写入响应
	if err := f.Write(c.Writer); err != nil {
		logger.Error("写入Excel模板失败", logger.F("error", err.Error()))
		// 注意：此时响应头已经设置，不能再用c.JSON，需要直接写入错误信息
		if !c.Writer.Written() {
			c.String(http.StatusInternalServerError, "写入Excel模板失败: %v", err)
		}
		return
	}

	logger.Info("Excel模板下载成功")
}

// BailianImportRequest 百炼平台导入请求格式
// 一个 provider (base_url + api_key) + 多个模型
type BailianImportRequest struct {
	BaseURL string                `json:"base_url"`
	APIKey  string                `json:"api_key"`
	Models  []BailianModelConfig  `json:"models"`
}

// BailianModelConfig 百炼平台模型配置
type BailianModelConfig struct {
	Name        string `json:"name"`        // 模型标识，如 qwen-turbo
	Type        string `json:"type"`        // 模型类型: chat, embedding, rerank 等
	DisplayName string `json:"display_name"` // 显示名称（可选）
}

// ImportLLMResourcesFromBailian 从百炼平台格式导入LLM资源
// @Summary Import LLM resources from Bailian format
// @Description Import LLM resources using Bailian platform format (one provider with multiple models)
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param request body BailianImportRequest true "Bailian import request"
// @Success 200 {object} map[string]interface{} "Import result"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources/import/bailian [post]
func (h *LLMResourceHandler) ImportLLMResourcesFromBailian(c *gin.Context) {
	var req BailianImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("百炼格式导入失败：JSON解析失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON解析失败: " + err.Error()})
		return
	}

	// 验证必填字段
	if req.BaseURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "base_url 是必填项"})
		return
	}
	if req.APIKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key 是必填项"})
		return
	}
	if len(req.Models) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "models 列表不能为空"})
		return
	}

	// 获取所有现有资源，用于重复检查
	existingResources, err := h.storage.ListLLMResources()
	if err != nil {
		logger.Error("百炼格式导入失败：获取现有资源列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取现有资源列表失败: " + err.Error()})
		return
	}

	successCount := 0
	failCount := 0
	duplicateCount := 0
	errors := []string{}
	duplicates := []map[string]interface{}{}

	validTypes := map[string]bool{
		"chat":       true,
		"completion": true,
		"embedding":  true,
		"image":      true,
		"audio":      true,
		"moderation": true,
		"rerank":     true,
	}

	for i, modelConfig := range req.Models {
		itemNum := i + 1

		// 验证模型配置
		modelName := strings.TrimSpace(modelConfig.Name)
		modelType := strings.ToLower(strings.TrimSpace(modelConfig.Type))
		displayName := strings.TrimSpace(modelConfig.DisplayName)

		if modelName == "" {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d个模型: 模型名称不能为空", itemNum))
			continue
		}

		if modelType == "" {
			modelType = "chat" // 默认为 chat 类型
		}

		if !validTypes[modelType] {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d个模型: 不支持的模型类型 '%s'", itemNum, modelConfig.Type))
			continue
		}

		// 生成资源名称
		resourceName := displayName
		if resourceName == "" {
			resourceName = modelName
		}

		// 创建资源对象（用于重复检查）
		resourceToCheck := &storage.LLMResource{
			Type:       modelType,
			Model:      modelName,
			BaseURL:    req.BaseURL,
			APIKey:     req.APIKey,
			TestStatus: "untested",
		}

		// 检查是否重复
		if h.isDuplicateResource(resourceToCheck, existingResources) {
			duplicateCount++
			duplicates = append(duplicates, map[string]interface{}{
				"item":     itemNum,
				"name":     resourceName,
				"type":     modelType,
				"model":    modelName,
				"base_url": req.BaseURL,
			})
			continue
		}

		// 创建资源
		resource := &storage.LLMResource{
			Name:       resourceName,
			Type:       modelType,
			Driver:     "openai",
			Model:      modelName,
			BaseURL:    req.BaseURL,
			APIKey:     req.APIKey,
			Status:     "active",
			TestStatus: "untested",
		}

		if err := h.storage.CreateLLMResource(resource); err != nil {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d个模型: %s", itemNum, err.Error()))
			continue
		}

		// 将新创建的资源添加到现有资源列表
		existingResources = append(existingResources, resource)
		successCount++
	}

	logger.Info("百炼格式导入LLM资源完成", logger.F("success", successCount), logger.F("fail", failCount), logger.F("duplicate", duplicateCount))
	c.JSON(http.StatusOK, gin.H{
		"message":    "导入完成",
		"success":    successCount,
		"fail":       failCount,
		"duplicate":  duplicateCount,
		"errors":     errors,
		"duplicates": duplicates,
		"total":      successCount + failCount + duplicateCount,
	})
}

// isDuplicateResource 检查资源是否重复
// 重复的定义：type、model、base_url、api_key 都一样
func (h *LLMResourceHandler) isDuplicateResource(resource *storage.LLMResource, existingResources []*storage.LLMResource) bool {
	for _, existing := range existingResources {
		if existing.Type == resource.Type &&
			existing.Model == resource.Model &&
			existing.BaseURL == resource.BaseURL &&
			existing.APIKey == resource.APIKey {
			return true
		}
	}
	return false
}

// getColumnIndex 获取列索引（支持多个可能的列名）
func getColumnIndex(headerMap map[string]int, possibleNames []string) int {
	for _, name := range possibleNames {
		if idx, exists := headerMap[strings.ToLower(name)]; exists {
			return idx
		}
	}
	return -1
}

// getCellValue 安全获取单元格值
func getCellValue(row []string, colIndex int) string {
	if colIndex < 0 || colIndex >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[colIndex])
}
