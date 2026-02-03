package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
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
}

// NewOpenAIHandler 创建新的OpenAI处理器
func NewOpenAIHandler(storage *storage.StorageFacade, policyService *service.PolicyService, tokenService *service.TokenService) *OpenAIHandler {
	return &OpenAIHandler{
		storage:       storage,
		policyService: policyService,
		tokenService:  tokenService,
	}
}

// findLLMResourceByModel 根据模型名称和Token查找对应的LLM资源
func (h *OpenAIHandler) findLLMResourceByModel(c *gin.Context, modelName string) (*storage.LLMResource, error) {
	// 获取所有LLM资源
	resources, err := h.storage.ListLLMResources()
	if err != nil {
		return nil, err
	}

	if len(resources) == 0 {
		return nil, fmt.Errorf("no LLM resources available")
	}

	// 从上下文获取Token（由认证中间件设置）
	tokenValue := h.getTokenFromContext(c)
	if tokenValue == "" {
		// 没有Token，使用默认策略
		return h.selectResourceWithDefaultPolicy(modelName, resources)
	}

	// 获取Token信息
	token, err := h.tokenService.ValidateToken(tokenValue)
	if err != nil {
		// Token验证失败，使用默认策略
		return h.selectResourceWithDefaultPolicy(modelName, resources)
	}

	// 检查Token是否有策略
	if token.PolicyID == "" {
		// Token没有配置策略
		// 如果Token是"ling-"开头，建议配置策略，但这里使用默认策略
		// 注意：ling-开头的Token建议配置策略，但这里先使用默认策略以保证向后兼容
		return h.selectResourceWithDefaultPolicy(modelName, resources)
	}

	// 使用Token的策略选择资源
	resource, err := h.policyService.ExecutePolicy(token.PolicyID, modelName, resources)
	if err != nil {
		// 策略执行失败，降级到默认策略
		return h.selectResourceWithDefaultPolicy(modelName, resources)
	}

	return resource, nil
}

// selectResourceWithDefaultPolicy 使用默认策略选择资源
func (h *OpenAIHandler) selectResourceWithDefaultPolicy(modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	// 默认策略：返回第一个可用的LLM资源
	for _, resource := range resources {
		if resource.Status == "active" {
			return resource, nil
		}
	}

	// 如果没有找到活跃的LLM资源，返回第一个可用的
	if len(resources) > 0 {
		return resources[0], nil
	}

	return nil, fmt.Errorf("no LLM resources available")
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

	// 查找对应的LLM资源
	llmResource, err := h.findLLMResourceByModel(c, req.Model)
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

	// 模拟响应，实际应用中需要使用client发送请求
	// 这里暂时使用模拟响应，因为需要修复客户端实现
	// 构建响应
	response := ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []struct {
			Index        int         `json:"index"`
			Message      ChatMessage `json:"message"`
			FinishReason string      `json:"finish_reason"`
		}{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: fmt.Sprintf("This is a response from the AI provider using LLM resource: %s. In production, this would be the actual response from the LLM resource.", llmResource.Name),
				},
				FinishReason: "stop",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     len(req.Messages) * 10,
			CompletionTokens: 20,
			TotalTokens:      len(req.Messages)*10 + 20,
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

	// 查找对应的LLM资源
	llmResource, err := h.findLLMResourceByModel(c, req.Model)
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

	// 模拟响应，实际应用中需要使用client发送请求
	// 这里暂时使用模拟响应，因为需要修复客户端实现
	response := CompletionResponse{
		ID:      fmt.Sprintf("cmpl-%d", time.Now().Unix()),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []struct {
			Text         string `json:"text"`
			Index        int    `json:"index"`
			FinishReason string `json:"finish_reason"`
		}{
			{
				Text:         fmt.Sprintf("This is a completion response from the AI provider using LLM resource: %s. In production, this would be the actual response from the LLM resource.", llmResource.Name),
				Index:        0,
				FinishReason: "stop",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     len(req.Prompt) / 4, // 粗略估计token数
			CompletionTokens: 20,
			TotalTokens:      len(req.Prompt)/4 + 20,
		},
	}

	c.JSON(http.StatusOK, response)
}

// ModelsResponse 模型列表响应
type ModelsResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

// ListModels godoc
// @Summary List models
// @Description List all available AI models
// @Tags openai
// @Accept json
// @Produce json
// @Success 200 {object} ModelsResponse "Models list"
// @Router /llm/v1/models [get]
func (h *OpenAIHandler) ListModels(c *gin.Context) {
	response := ModelsResponse{
		Object: "list",
		Data: []struct {
			ID      string `json:"id"`
			Object  string `json:"object"`
			Created int64  `json:"created"`
			OwnedBy string `json:"owned_by"`
		}{
			{ID: "gpt-4", Object: "model", Created: 1687882411, OwnedBy: "openai"},
			{ID: "gpt-3.5-turbo", Object: "model", Created: 1677610602, OwnedBy: "openai"},
			{ID: "claude-3-opus", Object: "model", Created: 1698959748, OwnedBy: "anthropic"},
			{ID: "gemini-pro", Object: "model", Created: 1701879923, OwnedBy: "google"},
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetModel godoc
// @Summary Get model
// @Description Get information about a specific model
// @Tags openai
// @Accept json
// @Produce json
// @Param model path string true "Model ID"
// @Success 200 {object} map[string]interface{} "Model information"
// @Router /llm/v1/models/{model} [get]
func (h *OpenAIHandler) GetModel(c *gin.Context) {
	model := c.Param("model")

	models := map[string]struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	}{
		"gpt-4":         {ID: "gpt-4", Object: "model", Created: 1687882411, OwnedBy: "openai"},
		"gpt-3.5-turbo": {ID: "gpt-3.5-turbo", Object: "model", Created: 1677610602, OwnedBy: "openai"},
		"claude-3-opus": {ID: "claude-3-opus", Object: "model", Created: 1698959748, OwnedBy: "anthropic"},
		"gemini-pro":    {ID: "gemini-pro", Object: "model", Created: 1701879923, OwnedBy: "google"},
	}

	if modelInfo, exists := models[model]; exists {
		c.JSON(http.StatusOK, modelInfo)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": map[string]interface{}{
			"message": fmt.Sprintf("The model '%s' does not exist", model),
			"type":    "invalid_request_error",
			"param":   "model",
			"code":    "model_not_found",
		},
	})
}
