package handler

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/middleware"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
	openaiSDK "github.com/openai/openai-go/v3"
)

// OpenAIHandler 处理OpenAI兼容的API请求
// @title OpenAI Compatible API
// @version 1.0.0
// @description OpenAI标准API接口，支持聊天补全、文本补全、嵌入等功能
// @host localhost:8080
// @BasePath /llm/v1
// @schemes http https

type OpenAIHandler struct {
	storage       *storage.StorageFacade
	policyService *service.PolicyService
	tokenService  *service.TokenService
	openaiService *service.OpenAIService
}

// NewOpenAIHandler 创建新的OpenAI处理器
func NewOpenAIHandler(storage *storage.StorageFacade, policyService *service.PolicyService, tokenService *service.TokenService) *OpenAIHandler {
	return &OpenAIHandler{
		storage:       storage,
		policyService: policyService,
		tokenService:  tokenService,
		openaiService: service.NewOpenAIService(),
	}
}

// findLLMResourceByModel 根据模型名称和Token查找对应的LLM资源
func (h *OpenAIHandler) findLLMResourceByModel(c *gin.Context, modelName string, resourceType string) (*storage.LLMResource, error) {
	// 获取所有LLM资源
	resources, err := h.storage.ListLLMResources()
	if err != nil {
		return nil, err
	}

	if len(resources) == 0 {
		return nil, fmt.Errorf("no LLM resources available")
	}

	// 根据资源类型过滤资源（只选择同类型的资源）
	typeFiltered := make([]*storage.LLMResource, 0)
	for _, r := range resources {
		if r.Type == resourceType {
			typeFiltered = append(typeFiltered, r)
		}
	}

	if len(typeFiltered) == 0 {
		logger.Warn("No resources found for type", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("resource_type", resourceType), logger.F("model", modelName))
		return nil, fmt.Errorf("no LLM resources available for type: %s", resourceType)
	}

	logger.Debug("Filtered resources by type", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("resource_type", resourceType), logger.F("original_count", len(resources)), logger.F("filtered_count", len(typeFiltered)))

	// 从上下文获取Token（由认证中间件设置）
	tokenValue := h.getTokenFromContext(c)
	if tokenValue == "" {
		// 没有Token，使用默认策略
		logger.Debug("No token found, using default policy", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("model", modelName), logger.F("resource_type", resourceType))
		return h.selectResourceWithDefaultPolicy(c, modelName, typeFiltered)
	}

	// 获取Token信息
	token, err := h.tokenService.ValidateToken(tokenValue)
	if err != nil {
		// Token验证失败，使用默认策略
		logger.Warn("Token validation failed, using default policy", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("error", err.Error()), logger.F("model", modelName), logger.F("resource_type", resourceType))
		return h.selectResourceWithDefaultPolicy(c, modelName, typeFiltered)
	}

	// 检查Token是否有策略
	if token.PolicyID == "" {
		// Token没有配置策略，使用默认策略
		logger.Debug("Token has no policy, using default policy", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("token_id", token.ID), logger.F("model", modelName), logger.F("resource_type", resourceType))
		return h.selectResourceWithDefaultPolicy(c, modelName, typeFiltered)
	}

	// 使用Token的策略选择资源
	logger.Debug("Executing policy for token", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("token_id", token.ID), logger.F("policy_id", token.PolicyID), logger.F("model", modelName), logger.F("resource_type", resourceType))
	resource, err := h.policyService.ExecutePolicy(token.PolicyID, modelName, typeFiltered)
	if err != nil {
		// 策略执行失败，降级到默认策略
		logger.Warn("Policy execution failed, falling back to default", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("error", err.Error()), logger.F("token_id", token.ID), logger.F("policy_id", token.PolicyID))
		return h.selectResourceWithDefaultPolicy(c, modelName, typeFiltered)
	}
	logger.Info("Resource selected successfully", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("token_id", token.ID), logger.F("policy_id", token.PolicyID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name), logger.F("resource_type", resource.Type))

	return resource, nil
}

// selectResourceWithDefaultPolicy 使用默认策略选择资源（随机策略）
func (h *OpenAIHandler) selectResourceWithDefaultPolicy(c *gin.Context, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	// 默认策略：从所有可用资源中随机选择
	activeResources := make([]*storage.LLMResource, 0)
	for _, resource := range resources {
		if resource.Status == "active" {
			activeResources = append(activeResources, resource)
		}
	}

	// 如果有活跃资源，随机选择一个
	if len(activeResources) > 0 {
		// 使用crypto/rand确保真正的随机性
		randomIndex, err := h.getRandomInt(len(activeResources))
		if err != nil {
			// 如果crypto/rand失败，使用math/rand作为后备
			randomIndex = int(time.Now().UnixNano()) % len(activeResources)
		}
		selected := activeResources[randomIndex]
		logger.Info("Default policy (random) selected resource", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("model", modelName), logger.F("resource_id", selected.ID), logger.F("resource_name", selected.Name), logger.F("total_active", len(activeResources)))
		return selected, nil
	}

	// 如果没有活跃资源，返回第一个可用的（降级处理）
	if len(resources) > 0 {
		logger.Warn("Default policy selected inactive resource (no active resources available)", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("model", modelName), logger.F("resource_id", resources[0].ID), logger.F("resource_name", resources[0].Name), logger.F("resource_status", resources[0].Status))
		return resources[0], nil
	}

	return nil, fmt.Errorf("no LLM resources available")
}

// getRandomInt 使用crypto/rand生成随机整数
func (h *OpenAIHandler) getRandomInt(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("max must be positive")
	}
	var n uint32
	if err := binary.Read(rand.Reader, binary.BigEndian, &n); err != nil {
		return 0, err
	}
	return int(n) % max, nil
}

// getTokenFromContext 从上下文获取Token
func (h *OpenAIHandler) getTokenFromContext(c *gin.Context) string {
	// 尝试从上下文获取Token（由认证中间件设置）
	if token, exists := c.Get("token"); exists {
		if t, ok := token.(*storage.Token); ok {
			return t.Token
		}
	}

	// 尝试从User获取APIKey（向后兼容）
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*storage.User); ok {
			return u.APIKey
		}
	}

	// 从请求头获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	return ""
}

// ChatCompletionRequest 聊天补全请求
// @Description OpenAI聊天补全请求结构
// @Param messages body []ChatMessage true "对话消息列表"
// @Param model body string true "模型名称"
// @Param max_tokens body int false "最大token数"
// @Param temperature body float64 false "温度参数(0-2)"
// @Param top_p body float64 false "top-p参数"
// @Param stream body bool false "是否流式响应"
// @Param stop body []string false "停止序列"
// @Param presence_penalty body float64 false "存在惩罚"
// @Param frequency_penalty body float64 false "频率惩罚"
type ChatCompletionRequest struct {
	Model            string        `json:"model" binding:"required"`
	Messages         []ChatMessage `json:"messages" binding:"required"`
	MaxTokens        int           `json:"max_tokens,omitempty"`
	Temperature      float64       `json:"temperature,omitempty"`
	TopP             float64       `json:"top_p,omitempty"`
	Stream           bool          `json:"stream,omitempty"`
	Stop             []string      `json:"stop,omitempty"`
	PresencePenalty  float64       `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64       `json:"frequency_penalty,omitempty"`
	User             string        `json:"user,omitempty"`
}

// ChatMessage 聊天消息
// @Description 单条聊天消息
// @Param role body string true "角色: system, user, assistant"
// @Param content body string true "消息内容"
// @Param name body string false "消息发送者名称"
type ChatMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
	Name    string `json:"name,omitempty"`
}

// ChatCompletionResponse 聊天补全响应
// @Description OpenAI聊天补全响应结构
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// CreateChatCompletion godoc
// @Summary Create chat completion
// @Description Create a chat completion using OpenAI-compatible API
// @Tags openai
// @Accept json
// @Produce json
// @Param request body ChatCompletionRequest true "Chat completion request"
// @Success 200 {object} ChatCompletionResponse "Chat completion response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /llm/v1/chat/completions [post]
func (h *OpenAIHandler) CreateChatCompletion(c *gin.Context) {
	var req ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"message": err.Error(),
				"type":    "invalid_request_error",
				"param":   nil,
				"code":    "invalid_request",
			},
		})
		return
	}

	// 查找对应的LLM资源（chat类型）
	llmResource, err := h.findLLMResourceByModel(c, req.Model, "chat")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Failed to find LLM resource: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 转换消息格式
	messages := make([]openaiSDK.ChatCompletionMessageParamUnion, 0, len(req.Messages))
	for _, msg := range req.Messages {
		switch strings.ToLower(msg.Role) {
		case "system":
			messages = append(messages, openaiSDK.SystemMessage(msg.Content))
		case "user":
			messages = append(messages, openaiSDK.UserMessage(msg.Content))
		case "assistant":
			messages = append(messages, openaiSDK.AssistantMessage(msg.Content))
		default:
			messages = append(messages, openaiSDK.UserMessage(msg.Content))
		}
	}

	// 使用资源中配置的model覆写请求中的model（如果资源配置了model）
	requestID := middleware.GetRequestID(c)
	modelToUse := req.Model
	if llmResource.Model != "" {
		modelToUse = llmResource.Model
		logger.Info("Model overridden by resource", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("request_model", req.Model), logger.F("resource_model", llmResource.Model), logger.F("final_model", modelToUse))
	}
	logger.Debug("Request parameters prepared", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("model", modelToUse), logger.F("base_url", llmResource.BaseURL), logger.F("message_count", len(messages)))

	// 获取用户ID（从token或user中）
	userID := ""
	if token, exists := c.Get("token"); exists {
		if t, ok := token.(*storage.Token); ok {
			userID = t.ID // 使用Token ID作为用户标识
		}
	} else if user, exists := c.Get("user"); exists {
		if u, ok := user.(*storage.User); ok {
			userID = u.ID
		}
	}

	// 构建服务层请求
	serviceReq := service.ChatCompletionRequest{
		Model:            modelToUse,
		Messages:         messages,
		MaxTokens:        int64(req.MaxTokens),
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		Stop:             req.Stop,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		User:             req.User,
	}

	// 调用统一的API服务
	ctx := c.Request.Context()
	serviceResp, err := h.openaiService.CreateChatCompletion(ctx, llmResource, serviceReq)
	duration := serviceResp.Duration.Milliseconds()
	openaiResponse := serviceResp.Response

	// 记录请求到数据库
	requestRecord := &storage.Request{
		UserID:        userID,
		LLMResourceID: llmResource.ID,
		Endpoint:      "/llm/v1/chat/completions",
		Method:        "POST",
		Duration:      duration,
		Status:        "success",
		Tokens:        0,
		CreatedAt:     time.Now(),
	}

	if err != nil {
		requestRecord.Status = "error"
		logger.Error("Chat completion failed", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", err.Error()), logger.F("request_model", req.Model), logger.F("actual_model", modelToUse), logger.F("resource_id", llmResource.ID), logger.F("resource_name", llmResource.Name))
		// 保存失败的请求记录
		if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
			logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Failed to create chat completion: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 记录token使用量
	if openaiResponse != nil {
		requestRecord.Tokens = int(openaiResponse.Usage.TotalTokens)
	}

	// 保存请求记录
	if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
		logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
	}

	// 转换响应格式
	choices := make([]struct {
		Index        int         `json:"index"`
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	}, len(openaiResponse.Choices))
	for i, choice := range openaiResponse.Choices {
		choices[i] = struct {
			Index        int         `json:"index"`
			Message      ChatMessage `json:"message"`
			FinishReason string      `json:"finish_reason"`
		}{
			Index: i,
			Message: ChatMessage{
				Role:    string(choice.Message.Role),
				Content: choice.Message.Content,
			},
			FinishReason: string(choice.FinishReason),
		}
	}

	response := ChatCompletionResponse{
		ID:      openaiResponse.ID,
		Object:  string(openaiResponse.Object),
		Created: int64(openaiResponse.Created),
		Model:   openaiResponse.Model,
		Choices: choices,
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     int(openaiResponse.Usage.PromptTokens),
			CompletionTokens: int(openaiResponse.Usage.CompletionTokens),
			TotalTokens:      int(openaiResponse.Usage.TotalTokens),
		},
	}

	c.JSON(http.StatusOK, response)
}

// CompletionRequest 文本补全请求
// @Description OpenAI文本补全请求结构
type CompletionRequest struct {
	Model            string   `json:"model" binding:"required"`
	Prompt           string   `json:"prompt" binding:"required"`
	MaxTokens        int      `json:"max_tokens,omitempty"`
	Temperature      float64  `json:"temperature,omitempty"`
	TopP             float64  `json:"top_p,omitempty"`
	Stream           bool     `json:"stream,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	PresencePenalty  float64  `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64  `json:"frequency_penalty,omitempty"`
	User             string   `json:"user,omitempty"`
}

// CompletionResponse 文本补全响应
type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// CreateCompletion godoc
// @Summary Create completion
// @Description Create a text completion using OpenAI-compatible API
// @Tags openai
// @Accept json
// @Produce json
// @Param request body CompletionRequest true "Completion request"
// @Success 200 {object} CompletionResponse "Completion response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /llm/v1/completions [post]
func (h *OpenAIHandler) CreateCompletion(c *gin.Context) {
	var req CompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"message": err.Error(),
				"type":    "invalid_request_error",
				"param":   nil,
				"code":    "invalid_request",
			},
		})
		return
	}

	// 查找对应的LLM资源（chat类型，completion也使用chat类型）
	llmResource, err := h.findLLMResourceByModel(c, req.Model, "chat")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Failed to find LLM resource: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 使用资源中配置的model覆写请求中的model（如果资源配置了model）
	requestID := middleware.GetRequestID(c)
	modelToUse := req.Model
	if llmResource.Model != "" {
		modelToUse = llmResource.Model
		logger.Info("Model overridden by resource", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("request_model", req.Model), logger.F("resource_model", llmResource.Model), logger.F("final_model", modelToUse))
	}
	logger.Debug("Request parameters prepared", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("model", modelToUse), logger.F("base_url", llmResource.BaseURL))

	// 获取用户ID（从token或user中）
	userID := ""
	if token, exists := c.Get("token"); exists {
		if t, ok := token.(*storage.Token); ok {
			userID = t.ID // 使用Token ID作为用户标识
		}
	} else if user, exists := c.Get("user"); exists {
		if u, ok := user.(*storage.User); ok {
			userID = u.ID
		}
	}

	// 构建服务层请求
	serviceReq := service.CompletionRequest{
		Model:            modelToUse,
		Prompt:           req.Prompt,
		MaxTokens:        int64(req.MaxTokens),
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		Stop:             req.Stop,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		User:             req.User,
	}

	// 调用统一的API服务
	ctx := c.Request.Context()
	serviceResp, err := h.openaiService.CreateCompletion(ctx, llmResource, serviceReq)
	duration := serviceResp.Duration.Milliseconds()
	openaiResponse := serviceResp.Response

	// 记录请求到数据库
	requestRecord := &storage.Request{
		UserID:        userID,
		LLMResourceID: llmResource.ID,
		Endpoint:      "/llm/v1/completions",
		Method:        "POST",
		Duration:      duration,
		Status:        "success",
		Tokens:        0,
		CreatedAt:     time.Now(),
	}

	if err != nil {
		requestRecord.Status = "error"
		logger.Error("Completion failed", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", err.Error()), logger.F("request_model", req.Model), logger.F("actual_model", modelToUse), logger.F("resource_id", llmResource.ID), logger.F("resource_name", llmResource.Name))
		// 保存失败的请求记录
		if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
			logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Failed to create completion: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 记录token使用量
	if openaiResponse != nil {
		requestRecord.Tokens = int(openaiResponse.Usage.TotalTokens)
	}

	// 保存请求记录
	if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
		logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
	}

	// 转换响应格式
	choices := make([]struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	}, len(openaiResponse.Choices))
	for i, choice := range openaiResponse.Choices {
		choices[i] = struct {
			Text         string `json:"text"`
			Index        int    `json:"index"`
			FinishReason string `json:"finish_reason"`
		}{
			Text:         choice.Text,
			Index:        i,
			FinishReason: string(choice.FinishReason),
		}
	}

	response := CompletionResponse{
		ID:      openaiResponse.ID,
		Object:  string(openaiResponse.Object),
		Created: int64(openaiResponse.Created),
		Model:   string(openaiResponse.Model),
		Choices: choices,
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     int(openaiResponse.Usage.PromptTokens),
			CompletionTokens: int(openaiResponse.Usage.CompletionTokens),
			TotalTokens:      int(openaiResponse.Usage.TotalTokens),
		},
	}

	c.JSON(http.StatusOK, response)
}

// ModelInfo 模型信息
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ModelsResponse 模型列表响应
type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

// ListModels godoc
// @Summary List models
// @Description List all available AI models from database
// @Tags openai
// @Accept json
// @Produce json
// @Success 200 {object} ModelsResponse "Models list"
// @Router /llm/v1/models [get]
func (h *OpenAIHandler) ListModels(c *gin.Context) {
	// 从数据库获取所有模型
	models, err := h.storage.ListModels()
	if err != nil {
		logger.Error("Failed to list models", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Failed to list models: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 如果没有模型，尝试从LLM资源中获取
	if len(models) == 0 {
		resources, err := h.storage.ListLLMResources()
		if err == nil && len(resources) > 0 {
			// 从LLM资源中提取模型信息
			data := make([]ModelInfo, 0, len(resources))
			for _, resource := range resources {
				if resource.Status == "active" && resource.Model != "" {
					data = append(data, ModelInfo{
						ID:      resource.Model,
						Object:  "model",
						Created: resource.CreatedAt.Unix(),
						OwnedBy: resource.Driver,
					})
				}
			}
			c.JSON(http.StatusOK, ModelsResponse{
				Object: "list",
				Data:   data,
			})
			return
		}
	}

	// 转换模型数据为OpenAI格式
	data := make([]ModelInfo, 0, len(models))
	for _, model := range models {
		if model.Status == "active" {
			modelID := model.ModelID
			if modelID == "" {
				// 如果没有ModelID，使用Name
				modelID = model.Name
			}
			ownedBy := model.Category
			if ownedBy == "" {
				// 如果没有Category，尝试从LLMResource获取
				if resource, err := h.storage.GetLLMResource(model.LLMResourceID); err == nil {
					ownedBy = resource.Driver
				} else {
					ownedBy = "unknown"
				}
			}
			data = append(data, ModelInfo{
				ID:      modelID,
				Object:  "model",
				Created: model.CreatedAt.Unix(),
				OwnedBy: ownedBy,
			})
		}
	}

	c.JSON(http.StatusOK, ModelsResponse{
		Object: "list",
		Data:   data,
	})
}

// GetModel godoc
// @Summary Get model
// @Description Get information about a specific model from database
// @Tags openai
// @Accept json
// @Produce json
// @Param model path string true "Model ID"
// @Success 200 {object} map[string]interface{} "Model information"
// @Router /llm/v1/models/{model} [get]
func (h *OpenAIHandler) GetModel(c *gin.Context) {
	modelID := c.Param("model")

	// 从数据库获取所有模型
	models, err := h.storage.ListModels()
	if err != nil {
		logger.Error("Failed to list models", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Failed to get model: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 查找匹配的模型（通过ModelID或Name）
	var foundModel *storage.Model
	for _, model := range models {
		if model.Status == "active" && (model.ModelID == modelID || model.Name == modelID) {
			foundModel = model
			break
		}
	}

	// 如果没找到，尝试从LLM资源中查找
	if foundModel == nil {
		resources, err := h.storage.ListLLMResources()
		if err == nil {
			for _, resource := range resources {
				if resource.Status == "active" && resource.Model == modelID {
					// 返回LLM资源的模型信息
					c.JSON(http.StatusOK, map[string]interface{}{
						"id":       resource.Model,
						"object":   "model",
						"created":  resource.CreatedAt.Unix(),
						"owned_by": resource.Driver,
					})
					return
				}
			}
		}
	}

	// 如果找到了模型
	if foundModel != nil {
		modelIDValue := foundModel.ModelID
		if modelIDValue == "" {
			modelIDValue = foundModel.Name
		}
		ownedBy := foundModel.Category
		if ownedBy == "" {
			if resource, err := h.storage.GetLLMResource(foundModel.LLMResourceID); err == nil {
				ownedBy = resource.Driver
			} else {
				ownedBy = "unknown"
			}
		}
		c.JSON(http.StatusOK, map[string]interface{}{
			"id":       modelIDValue,
			"object":   "model",
			"created":  foundModel.CreatedAt.Unix(),
			"owned_by": ownedBy,
		})
		return
	}

	// 模型不存在
	c.JSON(http.StatusNotFound, gin.H{
		"error": map[string]interface{}{
			"message": fmt.Sprintf("The model '%s' does not exist", modelID),
			"type":    "invalid_request_error",
			"param":   "model",
			"code":    "model_not_found",
		},
	})
}
