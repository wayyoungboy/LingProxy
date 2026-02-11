package handler

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
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

// findLLMResourceByModel 根据模型名称和API Key查找对应的LLM资源
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

	// 从上下文获取API Key（由认证中间件设置）
	tokenValue := h.getTokenFromContext(c)
	if tokenValue == "" {
		// 没有API Key，使用默认策略
		logger.Debug("No API key found, using default policy", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("model", modelName), logger.F("resource_type", resourceType))
		return h.selectResourceWithDefaultPolicy(c, modelName, typeFiltered)
	}

	// 获取API Key信息
	token, err := h.tokenService.ValidateToken(tokenValue)
	if err != nil {
		// API Key验证失败，使用默认策略
		logger.Warn("API key validation failed, using default policy", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("error", err.Error()), logger.F("model", modelName), logger.F("resource_type", resourceType))
		return h.selectResourceWithDefaultPolicy(c, modelName, typeFiltered)
	}

	// 检查模型许可
	if !token.IsModelAllowed(modelName) {
		logger.Warn("Model not allowed for API key", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("token_id", token.ID), logger.F("model", modelName))
		return nil, fmt.Errorf("model '%s' is not allowed for this API key", modelName)
	}

	// 获取按类型配置的策略ID
	policyID := token.GetPolicyIDByType(resourceType)
	if policyID == "" {
		// API Key没有配置策略，使用默认策略
		logger.Debug("API key has no policy for type, using default policy", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("token_id", token.ID), logger.F("model", modelName), logger.F("resource_type", resourceType))
		return h.selectResourceWithDefaultPolicy(c, modelName, typeFiltered)
	}

	// 使用API Key的策略选择资源
	logger.Debug("Executing policy for API key", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("token_id", token.ID), logger.F("policy_id", policyID), logger.F("model", modelName), logger.F("resource_type", resourceType))
	
	// 对于 embedding 类型，需要检查是否有 dimensions 参数
	var dimensions *int
	if resourceType == "embedding" {
		// 尝试从请求中获取 dimensions 参数
		if dims, exists := c.Get("embedding_dimensions"); exists {
			if d, ok := dims.(*int); ok {
				dimensions = d
			}
		}
	}
	
	resource, err := h.policyService.ExecutePolicyWithDimensions(policyID, modelName, typeFiltered, dimensions)
	if err != nil {
		// 策略执行失败，降级到默认策略
		logger.Warn("Policy execution failed, falling back to default", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("error", err.Error()), logger.F("token_id", token.ID), logger.F("policy_id", policyID))
		return h.selectResourceWithDefaultPolicy(c, modelName, typeFiltered)
	}
	logger.Debug("Resource selected successfully", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("token_id", token.ID), logger.F("policy_id", policyID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name), logger.F("resource_type", resource.Type))

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
		logger.Debug("Default policy (random) selected resource", logger.F("component", "handler"), logger.F("request_id", middleware.GetRequestID(c)), logger.F("model", modelName), logger.F("resource_id", selected.ID), logger.F("resource_name", selected.Name), logger.F("total_active", len(activeResources)))
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

// getTokenFromContext 从上下文获取API Key
func (h *OpenAIHandler) getTokenFromContext(c *gin.Context) string {
	// 尝试从上下文获取API Key（由认证中间件设置）
	if token, exists := c.Get("token"); exists {
		if t, ok := token.(*storage.APIKey); ok {
			return t.APIKey
		}
		if t, ok := token.(*storage.Token); ok {
			// 向后兼容：如果存储的是Token类型别名，转换为APIKey
			return t.APIKey
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

// ContentPart 内容部分（用于多模态消息）
type ContentPart struct {
	Type     string `json:"type"` // "text" 或 "image_url"
	Text     string `json:"text,omitempty"`
	ImageURL struct {
		URL string `json:"url"`
	} `json:"image_url,omitempty"`
}

// ChatMessage 聊天消息
// @Description 单条聊天消息，支持多模态（文本和图片）
// @Param role body string true "角色: system, user, assistant"
// @Param content body string|array true "消息内容，可以是字符串（纯文本）或数组（多模态：文本+图片）"
// @Param name body string false "消息发送者名称"
// @Example content as string: "Hello, how are you?"
// @Example content as array: [{"type": "text", "text": "What is in this image?"}, {"type": "image_url", "image_url": {"url": "https://example.com/image.png"}}]
type ChatMessage struct {
	Role    string      `json:"role" binding:"required"`
	Content interface{} `json:"content" binding:"required"` // 可以是 string 或 []ContentPart
	Name    string      `json:"name,omitempty"`
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

	// 获取 requestID（用于日志）
	requestID := middleware.GetRequestID(c)

	// 记录请求的 messages（用于调试）
	logger.Debug("Chat completion request received", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("model", req.Model), logger.F("message_count", len(req.Messages)), logger.F("messages", req.Messages))

	// 转换消息格式（支持多模态）
	messages := make([]openaiSDK.ChatCompletionMessageParamUnion, 0, len(req.Messages))
	for _, msg := range req.Messages {
		// 处理 content 字段（支持 string 或 array）
		var messageParam openaiSDK.ChatCompletionMessageParamUnion

		// 判断 content 类型
		switch content := msg.Content.(type) {
		case string:
			// 纯文本消息
			switch strings.ToLower(msg.Role) {
			case "system":
				messageParam = openaiSDK.SystemMessage(content)
			case "user":
				messageParam = openaiSDK.UserMessage(content)
			case "assistant":
				messageParam = openaiSDK.AssistantMessage(content)
			default:
				messageParam = openaiSDK.UserMessage(content)
			}
		case []interface{}:
			// 多模态消息（数组格式）- 仅 user 角色支持多模态
			if strings.ToLower(msg.Role) != "user" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": map[string]interface{}{
						"message": "Multimodal content (array) is only supported for 'user' role",
						"type":    "invalid_request_error",
						"param":   nil,
						"code":    "invalid_request",
					},
				})
				return
			}

			contentParts := make([]openaiSDK.ChatCompletionContentPartUnionParam, 0, len(content))
			for _, part := range content {
				partMap, ok := part.(map[string]interface{})
				if !ok {
					logger.Warn("Invalid content part format", logger.F("component", "handler"), logger.F("request_id", requestID))
					continue
				}

				partType, ok := partMap["type"].(string)
				if !ok {
					logger.Warn("Content part missing type", logger.F("component", "handler"), logger.F("request_id", requestID))
					continue
				}

				switch partType {
				case "text":
					if text, ok := partMap["text"].(string); ok {
						contentParts = append(contentParts, openaiSDK.TextContentPart(text))
					}
				case "image_url":
					if imageURLMap, ok := partMap["image_url"].(map[string]interface{}); ok {
						if url, ok := imageURLMap["url"].(string); ok {
							imageURLParam := openaiSDK.ChatCompletionContentPartImageImageURLParam{
								URL: url,
							}
							// 可选：支持 detail 参数
							if detail, ok := imageURLMap["detail"].(string); ok {
								imageURLParam.Detail = detail
							}
							contentParts = append(contentParts, openaiSDK.ImageContentPart(imageURLParam))
						}
					}
				default:
					logger.Warn("Unsupported content part type", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("type", partType))
				}
			}

			if len(contentParts) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": map[string]interface{}{
						"message": "Invalid content array: no valid content parts",
						"type":    "invalid_request_error",
						"param":   nil,
						"code":    "invalid_request",
					},
				})
				return
			}

			// User 消息支持多模态
			messageParam = openaiSDK.UserMessage(contentParts)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": map[string]interface{}{
					"message": "Invalid content type: must be string or array",
					"type":    "invalid_request_error",
					"param":   nil,
					"code":    "invalid_request",
				},
			})
			return
		}

		messages = append(messages, messageParam)
	}

	// 使用资源中配置的model覆写请求中的model（如果资源配置了model）
	modelToUse := req.Model
	if llmResource.Model != "" {
		modelToUse = llmResource.Model
		logger.Debug("Model overridden by resource", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("request_model", req.Model), logger.F("resource_model", llmResource.Model), logger.F("final_model", modelToUse))
	}
	logger.Debug("Request parameters prepared", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("model", modelToUse), logger.F("base_url", llmResource.BaseURL), logger.F("message_count", len(messages)))

	// 获取用户ID（从API Key或user中）
	userID := ""
	if token, exists := c.Get("token"); exists {
		if t, ok := token.(*storage.APIKey); ok {
			userID = t.ID // 使用API Key ID作为用户标识
		} else if t, ok := token.(*storage.Token); ok {
			// 向后兼容：如果存储的是Token类型别名
			userID = t.ID
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
		Stream:           req.Stream,
	}

	// 调用统一的API服务
	ctx := c.Request.Context()

	// 如果请求流式响应，使用流式处理
	if req.Stream {
		h.handleStreamingChatCompletion(c, ctx, llmResource, serviceReq, requestID, userID, modelToUse)
		return
	}

	// 非流式响应处理
	serviceResp, err := h.openaiService.CreateChatCompletion(ctx, llmResource, serviceReq)
	var duration int64
	var openaiResponse *openaiSDK.ChatCompletion
	
	if serviceResp != nil {
		duration = serviceResp.Duration.Milliseconds()
		openaiResponse = serviceResp.Response
	}

	// 记录API调用结果
	if err != nil {
		logger.Error("OpenAI API call failed", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", err.Error()))
	} else if openaiResponse == nil {
		logger.Error("OpenAI API returned nil response", logger.F("component", "handler"), logger.F("request_id", requestID))
	} else {
		logger.Debug("OpenAI API call succeeded", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("response_id", openaiResponse.ID), logger.F("model", openaiResponse.Model), logger.F("choices_count", len(openaiResponse.Choices)))
		if len(openaiResponse.Choices) > 0 {
			logger.Debug("Chat completion response", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("response_content", openaiResponse.Choices[0].Message.Content), logger.F("finish_reason", openaiResponse.Choices[0].FinishReason))
		}
	}

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
		} else {
			logger.Debug("Failed request record saved", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("resource_id", llmResource.ID), logger.F("user_id", userID))
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
	if openaiResponse != nil && openaiResponse.Usage.TotalTokens > 0 {
		requestRecord.Tokens = int(openaiResponse.Usage.TotalTokens)
		logger.Debug("Token usage recorded", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("tokens", requestRecord.Tokens))
	}

	// 保存请求记录
	if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
		logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
	} else {
		logger.Debug("Request record saved successfully", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("resource_id", llmResource.ID), logger.F("user_id", userID), logger.F("tokens", requestRecord.Tokens), logger.F("status", requestRecord.Status))
	}

	// 检查响应是否有效
	if openaiResponse == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Received nil response from OpenAI API",
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	if len(openaiResponse.Choices) == 0 {
		logger.Warn("OpenAI API returned empty choices", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("response_id", openaiResponse.ID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Received empty choices from OpenAI API",
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 直接返回从资源侧收到的响应，不做任何转换
	// 记录最终响应（用于调试）
	logger.Debug("Sending chat completion response", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("response_id", openaiResponse.ID), logger.F("choices_count", len(openaiResponse.Choices)))
	if len(openaiResponse.Choices) > 0 {
		logger.Debug("Response content", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("content", openaiResponse.Choices[0].Message.Content), logger.F("content_type", fmt.Sprintf("%T", openaiResponse.Choices[0].Message.Content)))
	}

	// 尝试序列化响应以检查是否有错误
	responseJSON, err := json.Marshal(openaiResponse)
	if err != nil {
		logger.Error("Failed to marshal response to JSON", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Failed to serialize response: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	previewLen := len(responseJSON)
	if previewLen > 500 {
		previewLen = 500
	}
	logger.Debug("Response JSON serialized successfully", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("json_size", len(responseJSON)), logger.F("json_preview", string(responseJSON[:previewLen])))

	// 直接返回从资源侧收到的原始响应
	c.JSON(http.StatusOK, openaiResponse)

	// 记录响应发送完成
	logger.Debug("Response sent to client", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("status_code", c.Writer.Status()), logger.F("response_written", c.Writer.Written()))
}

// handleStreamingChatCompletion 处理流式聊天补全请求
func (h *OpenAIHandler) handleStreamingChatCompletion(c *gin.Context, ctx context.Context, llmResource *storage.LLMResource, req service.ChatCompletionRequest, requestID string, userID string, modelToUse string) {
	logger.Debug("Starting streaming chat completion", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("model", modelToUse))

	// 设置响应头为 Server-Sent Events (SSE)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // 禁用 nginx 缓冲

	// 记录开始时间
	startTime := time.Now()

	// 获取流式响应
	stream, err := h.openaiService.CreateChatCompletionStream(ctx, llmResource, req)
	if err != nil {
		logger.Error("Failed to create streaming chat completion", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", err.Error()))
		
		// 记录失败的请求
		duration := time.Since(startTime)
		requestRecord := &storage.Request{
			UserID:        userID,
			LLMResourceID: llmResource.ID,
			Endpoint:      "/llm/v1/chat/completions",
			Method:        "POST",
			Duration:      duration.Milliseconds(),
			Status:        "error",
			Tokens:        0,
			CreatedAt:     time.Now(),
		}
		if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
			logger.Error("Failed to save failed streaming request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
		}
		
		c.SSEvent("error", gin.H{
			"error": map[string]interface{}{
				"message": "Failed to create streaming chat completion: " + err.Error(),
				"type":    "internal_server_error",
				"code":    "internal_server_error",
			},
		})
		c.Writer.Flush()
		return
	}
	var totalTokens int
	var promptTokens int
	var completionTokens int
	var responseID string
	var responseModel string

	// 处理流式响应
	defer stream.Close()

	// 迭代流式响应
	for stream.Next() {
		chunk := stream.Current()
		if chunk.ID == "" {
			continue
		}

		// 提取 response ID 和 model（通常在第一个 chunk）
		if responseID == "" {
			responseID = chunk.ID
			responseModel = string(chunk.Model)
		}

		// 提取 usage（通常在最后一个 chunk）
		if chunk.Usage.PromptTokens > 0 || chunk.Usage.CompletionTokens > 0 {
			promptTokens = int(chunk.Usage.PromptTokens)
			completionTokens = int(chunk.Usage.CompletionTokens)
			totalTokens = int(chunk.Usage.TotalTokens)
		}

		// 清理 chunk 数据以符合 OpenAI API 规范
		cleanedChunk := h.cleanChatCompletionChunk(chunk)

		// 将清理后的 chunk 转换为 JSON
		chunkJSON, err := json.Marshal(cleanedChunk)
		if err != nil {
			logger.Warn("Failed to marshal chunk", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", err.Error()))
			continue
		}

		// 发送 SSE 事件
		c.Writer.WriteString("data: ")
		c.Writer.Write(chunkJSON)
		c.Writer.WriteString("\n\n")
		c.Writer.Flush()
	}

	// 检查流式响应错误
	if err := stream.Err(); err != nil {
		logger.Error("Stream iteration error", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", err.Error()))
		
		// 记录失败的请求
		duration := time.Since(startTime)
		requestRecord := &storage.Request{
			UserID:        userID,
			LLMResourceID: llmResource.ID,
			Endpoint:      "/llm/v1/chat/completions",
			Method:        "POST",
			Duration:      duration.Milliseconds(),
			Status:        "error",
			Tokens:        totalTokens,
			CreatedAt:     time.Now(),
		}
		if totalTokens == 0 {
			requestRecord.Tokens = promptTokens + completionTokens
		}
		if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
			logger.Error("Failed to save failed streaming request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
		}
		
		errorJSON, _ := json.Marshal(gin.H{
			"error": map[string]interface{}{
				"message": "Stream error: " + err.Error(),
				"type":    "internal_server_error",
				"code":    "internal_server_error",
			},
		})
		c.Writer.WriteString("data: ")
		c.Writer.Write(errorJSON)
		c.Writer.WriteString("\n\n")
		c.Writer.Flush()
		return
	}

	// 发送结束标记
	c.Writer.WriteString("data: [DONE]\n\n")
	c.Writer.Flush()

	duration := time.Since(startTime)

	// 记录请求到数据库
	requestRecord := &storage.Request{
		UserID:        userID,
		LLMResourceID: llmResource.ID,
		Endpoint:      "/llm/v1/chat/completions",
		Method:        "POST",
		Duration:      duration.Milliseconds(),
		Status:        "success",
		Tokens:        totalTokens,
		CreatedAt:     time.Now(),
	}

	if totalTokens == 0 {
		requestRecord.Tokens = promptTokens + completionTokens
	}

	// 保存请求记录
	if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
		logger.Error("Failed to save streaming request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
	} else {
		logger.Debug("Streaming request record saved successfully", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("resource_id", llmResource.ID), logger.F("user_id", userID), logger.F("tokens", totalTokens), logger.F("status", requestRecord.Status))
	}

	logger.Debug("Streaming chat completion completed", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("response_id", responseID), logger.F("model", responseModel), logger.F("duration_ms", duration.Milliseconds()), logger.F("tokens", totalTokens))
}

// cleanChatCompletionChunk 清理 ChatCompletionChunk 以符合 OpenAI API 规范
// 移除空字符串字段，确保符合 OpenAI SDK 的验证要求
func (h *OpenAIHandler) cleanChatCompletionChunk(chunk openaiSDK.ChatCompletionChunk) map[string]interface{} {
	chunkJSON, err := json.Marshal(chunk)
	if err != nil {
		logger.Warn("Failed to marshal chunk for cleaning", logger.F("error", err.Error()))
		return map[string]interface{}{
			"id":      chunk.ID,
			"object":  chunk.Object,
			"created": chunk.Created,
			"model":   chunk.Model,
			"choices": chunk.Choices,
		}
	}

	var chunkMap map[string]interface{}
	if err := json.Unmarshal(chunkJSON, &chunkMap); err != nil {
		logger.Warn("Failed to unmarshal chunk for cleaning", logger.F("error", err.Error()))
		return chunkMap
	}

	// 清理 choices 中的 delta 字段
	if choices, ok := chunkMap["choices"].([]interface{}); ok {
		for _, choice := range choices {
			if choiceMap, ok := choice.(map[string]interface{}); ok {
				if delta, ok := choiceMap["delta"].(map[string]interface{}); ok {
					// 移除 role 字段（OpenAI API 规范中 delta 不应该包含 role）
					delete(delta, "role")
					h.removeEmptyStringField(delta, "refusal")
					// 清理 function_call
					if functionCall, ok := delta["function_call"].(map[string]interface{}); ok {
						h.removeEmptyStringField(functionCall, "name")
						h.removeEmptyStringField(functionCall, "arguments")
						if len(functionCall) == 0 {
							delete(delta, "function_call")
						}
					}
					// 移除 null 的 tool_calls
					if toolCalls, ok := delta["tool_calls"]; ok && toolCalls == nil {
						delete(delta, "tool_calls")
					}
				}
				// 移除空的 finish_reason
				h.removeEmptyStringField(choiceMap, "finish_reason")
			}
		}
	}

	// 移除其他空字段
	h.removeEmptyStringField(chunkMap, "system_fingerprint")
	h.removeEmptyStringField(chunkMap, "service_tier")

	return chunkMap
}

// removeEmptyStringField 辅助函数：如果字段是空字符串则移除
func (h *OpenAIHandler) removeEmptyStringField(m map[string]interface{}, key string) {
	if val, ok := m[key].(string); ok && val == "" {
		delete(m, key)
	}
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
		logger.Debug("Model overridden by resource", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("request_model", req.Model), logger.F("resource_model", llmResource.Model), logger.F("final_model", modelToUse))
	}
	logger.Debug("Request parameters prepared", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("model", modelToUse), logger.F("base_url", llmResource.BaseURL))

	// 获取用户ID（从API Key或user中）
	userID := ""
	if token, exists := c.Get("token"); exists {
		if t, ok := token.(*storage.APIKey); ok {
			userID = t.ID // 使用API Key ID作为用户标识
		} else if t, ok := token.(*storage.Token); ok {
			// 向后兼容：如果存储的是Token类型别名
			userID = t.ID
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
	var duration int64
	var openaiResponse *openaiSDK.Completion
	
	if serviceResp != nil {
		duration = serviceResp.Duration.Milliseconds()
		openaiResponse = serviceResp.Response
	}

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
		} else {
			logger.Debug("Failed completion request record saved", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("resource_id", llmResource.ID), logger.F("user_id", userID))
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
	if openaiResponse != nil && openaiResponse.Usage.TotalTokens > 0 {
		requestRecord.Tokens = int(openaiResponse.Usage.TotalTokens)
		logger.Debug("Completion token usage recorded", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("tokens", requestRecord.Tokens))
	}

	// 保存请求记录
	if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
		logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
	} else {
		logger.Debug("Completion request record saved successfully", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("resource_id", llmResource.ID), logger.F("user_id", userID), logger.F("tokens", requestRecord.Tokens), logger.F("status", requestRecord.Status))
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

// RerankRequest 重排序请求
type RerankRequest struct {
	Model     string   `json:"model" binding:"required"`
	Query     string   `json:"query" binding:"required"`
	Documents []string `json:"documents" binding:"required"`
	TopN      int      `json:"top_n,omitempty"` // 返回前N个结果，可选
}

// CreateRerank godoc
// @Summary Create rerank
// @Description Rerank documents based on a query
// @Tags openai
// @Accept json
// @Produce json
// @Param request body RerankRequest true "Rerank request"
// @Success 200 {object} map[string]interface{} "Rerank response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /llm/v1/reranks [post]
func (h *OpenAIHandler) CreateRerank(c *gin.Context) {
	var req RerankRequest
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

	// 验证必填字段
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"message": "query is required",
				"type":    "invalid_request_error",
				"param":   "query",
				"code":    "invalid_request",
			},
		})
		return
	}
	if len(req.Documents) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"message": "documents cannot be empty",
				"type":    "invalid_request_error",
				"param":   "documents",
				"code":    "invalid_request",
			},
		})
		return
	}

	// 查找对应的LLM资源（rerank类型）
	llmResource, err := h.findLLMResourceByModel(c, req.Model, "rerank")
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

	requestID := middleware.GetRequestID(c)

	// 使用资源中配置的model覆写请求中的model（如果资源配置了model）
	modelToUse := req.Model
	if llmResource.Model != "" {
		modelToUse = llmResource.Model
		logger.Debug("Model overridden by resource", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("request_model", req.Model), logger.F("resource_model", llmResource.Model), logger.F("final_model", modelToUse))
	}

	// 获取用户ID
	userID := ""
	if token, exists := c.Get("token"); exists {
		if t, ok := token.(*storage.APIKey); ok {
			userID = t.ID
		} else if t, ok := token.(*storage.Token); ok {
			// 向后兼容：如果存储的是Token类型别名
			userID = t.ID
		}
	} else if user, exists := c.Get("user"); exists {
		if u, ok := user.(*storage.User); ok {
			userID = u.ID
		}
	}

	// 构建服务层请求
	serviceReq := service.RerankRequest{
		Model:     modelToUse,
		Query:     req.Query,
		Documents: req.Documents,
		TopN:      req.TopN,
	}

	// 调用服务层
	ctx := c.Request.Context()
	serviceResp, err := h.openaiService.CreateRerank(ctx, llmResource, serviceReq)
	var duration int64
	var rerankResponse map[string]interface{}

	if serviceResp != nil {
		duration = serviceResp.Duration.Milliseconds()
		rerankResponse = serviceResp.Response
	}

	// 记录请求到数据库
	requestRecord := &storage.Request{
		UserID:        userID,
		LLMResourceID: llmResource.ID,
		Endpoint:      "/llm/v1/reranks",
		Method:        "POST",
		Duration:      duration,
		Status:        "success",
		Tokens:        0,
		CreatedAt:     time.Now(),
	}

	if err != nil {
		requestRecord.Status = "error"
		logger.Error("Rerank failed", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", err.Error()), logger.F("request_model", req.Model), logger.F("actual_model", modelToUse), logger.F("resource_id", llmResource.ID), logger.F("resource_name", llmResource.Name))
		if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
			logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Failed to create rerank: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 保存请求记录
	if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
		logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
	} else {
		logger.Debug("Rerank request record saved successfully", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("resource_id", llmResource.ID), logger.F("user_id", userID), logger.F("status", requestRecord.Status))
	}

	// 直接返回从资源侧收到的响应
	if rerankResponse == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Received nil response from rerank API",
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, rerankResponse)
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

// EmbeddingRequest 嵌入请求
type EmbeddingRequest struct {
	Model      string      `json:"model" binding:"required"`
	Input      interface{} `json:"input" binding:"required"` // 支持 string 或 []string
	Dimensions *int        `json:"dimensions,omitempty"`     // 向量维度（可选）
	User       string      `json:"user,omitempty"`
}

// EmbeddingResponse 嵌入响应
type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// CreateEmbedding godoc
// @Summary Create embeddings
// @Description Create embeddings for text input
// @Tags openai
// @Accept json
// @Produce json
// @Param request body EmbeddingRequest true "Embedding request"
// @Success 200 {object} EmbeddingResponse "Embedding response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /llm/v1/embeddings [post]
func (h *OpenAIHandler) CreateEmbedding(c *gin.Context) {
	var req EmbeddingRequest
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

	// 如果有 dimensions 参数，设置到 context 中供策略执行器使用
	if req.Dimensions != nil {
		c.Set("embedding_dimensions", req.Dimensions)
	}

	// 查找对应的LLM资源（embedding类型）
	llmResource, err := h.findLLMResourceByModel(c, req.Model, "embedding")
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

	requestID := middleware.GetRequestID(c)

	// 处理输入：支持 string 或 []string
	var inputStr string
	switch v := req.Input.(type) {
	case string:
		inputStr = v
	case []interface{}:
		// 如果是数组，只取第一个元素（简化处理）
		if len(v) > 0 {
			if str, ok := v[0].(string); ok {
				inputStr = str
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": map[string]interface{}{
						"message": "Input array must contain strings",
						"type":    "invalid_request_error",
						"param":   "input",
						"code":    "invalid_request",
					},
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": map[string]interface{}{
					"message": "Input array cannot be empty",
					"type":    "invalid_request_error",
					"param":   "input",
					"code":    "invalid_request",
				},
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"message": "Input must be a string or array of strings",
				"type":    "invalid_request_error",
				"param":   "input",
				"code":    "invalid_request",
			},
		})
		return
	}

	if inputStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"message": "Input cannot be empty",
				"type":    "invalid_request_error",
				"param":   "input",
				"code":    "invalid_request",
			},
		})
		return
	}

	// 使用资源中配置的model覆写请求中的model（如果资源配置了model）
	modelToUse := req.Model
	if llmResource.Model != "" {
		modelToUse = llmResource.Model
		logger.Debug("Model overridden by resource", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("request_model", req.Model), logger.F("resource_model", llmResource.Model), logger.F("final_model", modelToUse))
	}

	// 获取用户ID
	userID := ""
	if token, exists := c.Get("token"); exists {
		if t, ok := token.(*storage.APIKey); ok {
			userID = t.ID
		} else if t, ok := token.(*storage.Token); ok {
			// 向后兼容：如果存储的是Token类型别名
			userID = t.ID
		}
	} else if user, exists := c.Get("user"); exists {
		if u, ok := user.(*storage.User); ok {
			userID = u.ID
		}
	}

	// 构建服务层请求
	serviceReq := service.EmbeddingRequest{
		Model: modelToUse,
		Input: inputStr,
	}

	// 调用服务层
	ctx := c.Request.Context()
	serviceResp, err := h.openaiService.CreateEmbedding(ctx, llmResource, serviceReq)
	var duration int64
	var embeddingResponse *openaiSDK.CreateEmbeddingResponse

	if serviceResp != nil {
		duration = serviceResp.Duration.Milliseconds()
		embeddingResponse = serviceResp.Response
	}

	// 记录请求到数据库
	requestRecord := &storage.Request{
		UserID:        userID,
		LLMResourceID: llmResource.ID,
		Endpoint:      "/llm/v1/embeddings",
		Method:        "POST",
		Duration:      duration,
		Status:        "success",
		Tokens:        0,
		CreatedAt:     time.Now(),
	}

	if err != nil {
		requestRecord.Status = "error"
		logger.Error("Embedding failed", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", err.Error()), logger.F("request_model", req.Model), logger.F("actual_model", modelToUse), logger.F("resource_id", llmResource.ID), logger.F("resource_name", llmResource.Name))
		if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
			logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Failed to create embedding: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 记录token使用量
	if embeddingResponse != nil && embeddingResponse.Usage.TotalTokens > 0 {
		requestRecord.Tokens = int(embeddingResponse.Usage.TotalTokens)
		logger.Debug("Embedding token usage recorded", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("tokens", requestRecord.Tokens))
	}

	// 保存请求记录
	if saveErr := h.storage.CreateRequest(requestRecord); saveErr != nil {
		logger.Error("Failed to save request record", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("error", saveErr.Error()))
	} else {
		logger.Debug("Embedding request record saved successfully", logger.F("component", "handler"), logger.F("request_id", requestID), logger.F("resource_id", llmResource.ID), logger.F("user_id", userID), logger.F("tokens", requestRecord.Tokens), logger.F("status", requestRecord.Status))
	}

	if embeddingResponse == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"message": "Received nil response from embedding API",
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "internal_server_error",
			},
		})
		return
	}

	// 转换响应格式
	response := EmbeddingResponse{
		Object: string(embeddingResponse.Object),
		Model:  string(embeddingResponse.Model),
		Data:   make([]struct {
			Object    string    `json:"object"`
			Embedding []float64 `json:"embedding"`
			Index     int       `json:"index"`
		}, len(embeddingResponse.Data)),
		Usage: struct {
			PromptTokens int `json:"prompt_tokens"`
			TotalTokens  int `json:"total_tokens"`
		}{
			PromptTokens: int(embeddingResponse.Usage.PromptTokens),
			TotalTokens:  int(embeddingResponse.Usage.TotalTokens),
		},
	}

	for i, data := range embeddingResponse.Data {
		response.Data[i] = struct {
			Object    string    `json:"object"`
			Embedding []float64 `json:"embedding"`
			Index     int       `json:"index"`
		}{
			Object:    string(data.Object),
			Embedding: data.Embedding,
			Index:     i,
		}
	}

	c.JSON(http.StatusOK, response)
}
