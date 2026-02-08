package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lingproxy/lingproxy/internal/client/embedding"
	"github.com/lingproxy/lingproxy/internal/client/openai"
	"github.com/lingproxy/lingproxy/internal/config"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/storage"
	openaiSDK "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

// OpenAIService OpenAI服务，统一管理OpenAI客户端调用
type OpenAIService struct {
	// 可以添加缓存、连接池等
}

// NewOpenAIService 创建新的OpenAI服务
func NewOpenAIService() *OpenAIService {
	return &OpenAIService{}
}

// CreateClient 根据LLM资源创建OpenAI客户端
func (s *OpenAIService) CreateClient(resource *storage.LLMResource) *openai.Client {
	return openai.NewClient(resource.APIKey, resource.BaseURL)
}

// CreateEmbeddingClient 根据LLM资源创建Embedding客户端
func (s *OpenAIService) CreateEmbeddingClient(resource *storage.LLMResource) *embedding.Client {
	model := resource.Model
	if model == "" {
		model = "text-embedding-3-small" // 默认模型
	}
	return embedding.NewClient(resource.APIKey, resource.BaseURL, model)
}

// ChatCompletionRequest 聊天补全请求参数
type ChatCompletionRequest struct {
	Model            string
	Messages         []openaiSDK.ChatCompletionMessageParamUnion
	MaxTokens        int64
	Temperature      float64
	TopP             float64
	Stop             []string
	PresencePenalty  float64
	FrequencyPenalty float64
	User             string
	Stream           bool
}

// ChatCompletionResponse 聊天补全响应
type ChatCompletionResponse struct {
	Response *openaiSDK.ChatCompletion
	Duration time.Duration
	Error    error
}

// CreateChatCompletion 创建聊天补全请求（带自动重试）
func (s *OpenAIService) CreateChatCompletion(ctx context.Context, resource *storage.LLMResource, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	return s.createChatCompletionWithRetry(ctx, resource, req)
}

// createChatCompletionWithRetry 带重试的聊天补全请求
func (s *OpenAIService) createChatCompletionWithRetry(ctx context.Context, resource *storage.LLMResource, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	client := s.CreateClient(resource)
	defer client.Close()

	// 确定使用的模型
	modelToUse := req.Model
	if resource.Model != "" {
		modelToUse = resource.Model
	}

	// 构建请求参数
	params := s.buildChatCompletionParams(modelToUse, req)

	// 获取重试配置
	maxRetries := config.C.Provider.MaxRetries
	retryDelay := config.C.Provider.RetryDelay
	if maxRetries <= 0 {
		maxRetries = 0 // 禁用重试
	}

	// 记录开始时间
	startTime := time.Now()
	var lastErr error
	var lastResponse *openaiSDK.ChatCompletion

	// 重试逻辑
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// 重试前等待
			logger.Debug("Retrying chat completion request", logger.F("component", "service"), logger.F("attempt", attempt), logger.F("max_retries", maxRetries), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
			select {
			case <-ctx.Done():
				return &ChatCompletionResponse{
					Response: nil,
					Duration: time.Since(startTime),
					Error:    ctx.Err(),
				}, ctx.Err()
			case <-time.After(retryDelay * time.Duration(attempt)): // 指数退避
			}
		}

		// 调用API
		response, err := client.Chat().Completions().New(ctx, params)
		duration := time.Since(startTime)

		if err == nil {
			// 成功，返回结果
			return &ChatCompletionResponse{
				Response: response,
				Duration: duration,
				Error:    nil,
			}, nil
		}

		lastErr = err
		lastResponse = response

		// 检查是否可重试
		if !s.isRetryableError(err) {
			logger.Debug("Error is not retryable", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("resource_id", resource.ID))
			break
		}

		// 如果还有重试机会，继续
		if attempt < maxRetries {
			logger.Warn("Chat completion request failed, will retry", logger.F("component", "service"), logger.F("attempt", attempt+1), logger.F("max_retries", maxRetries), logger.F("error", err.Error()), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
		}
	}

	// 所有重试都失败
	duration := time.Since(startTime)
	logger.Error("Chat completion request failed after retries", logger.F("component", "service"), logger.F("attempts", maxRetries+1), logger.F("error", lastErr.Error()), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
	return &ChatCompletionResponse{
		Response: lastResponse,
		Duration: duration,
		Error:    lastErr,
	}, lastErr
}

// buildChatCompletionParams 构建聊天补全请求参数
func (s *OpenAIService) buildChatCompletionParams(modelToUse string, req ChatCompletionRequest) openaiSDK.ChatCompletionNewParams {
	params := openaiSDK.ChatCompletionNewParams{
		Model:    modelToUse,
		Messages: req.Messages,
	}

	// 设置可选参数
	if req.MaxTokens > 0 {
		params.MaxTokens = param.NewOpt(req.MaxTokens)
	}
	if req.Temperature > 0 {
		params.Temperature = param.NewOpt(req.Temperature)
	}
	if req.TopP > 0 {
		params.TopP = param.NewOpt(req.TopP)
	}
	if len(req.Stop) > 0 {
		params.Stop = openaiSDK.ChatCompletionNewParamsStopUnion{
			OfStringArray: req.Stop,
		}
	}
	if req.PresencePenalty != 0 {
		params.PresencePenalty = param.NewOpt(req.PresencePenalty)
	}
	if req.FrequencyPenalty != 0 {
		params.FrequencyPenalty = param.NewOpt(req.FrequencyPenalty)
	}
	if req.User != "" {
		params.User = param.NewOpt(req.User)
	}

	return params
}

// CreateChatCompletionStream 创建流式聊天补全请求（带自动重试）
func (s *OpenAIService) CreateChatCompletionStream(ctx context.Context, resource *storage.LLMResource, req ChatCompletionRequest) (*ssestream.Stream[openaiSDK.ChatCompletionChunk], error) {
	return s.createChatCompletionStreamWithRetry(ctx, resource, req)
}

// createChatCompletionStreamWithRetry 带重试的流式聊天补全请求
func (s *OpenAIService) createChatCompletionStreamWithRetry(ctx context.Context, resource *storage.LLMResource, req ChatCompletionRequest) (*ssestream.Stream[openaiSDK.ChatCompletionChunk], error) {
	// 确定使用的模型
	modelToUse := req.Model
	if resource.Model != "" {
		modelToUse = resource.Model
	}

	// 构建请求参数
	params := s.buildChatCompletionParams(modelToUse, req)

	// 获取重试配置
	maxRetries := config.C.Provider.MaxRetries
	retryDelay := config.C.Provider.RetryDelay
	if maxRetries <= 0 {
		maxRetries = 0 // 禁用重试
	}

	var lastErr error
	var lastStream *ssestream.Stream[openaiSDK.ChatCompletionChunk]

	// 重试逻辑（仅在创建流之前失败时重试）
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// 重试前等待
			logger.Debug("Retrying streaming chat completion request", logger.F("component", "service"), logger.F("attempt", attempt), logger.F("max_retries", maxRetries), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryDelay * time.Duration(attempt)): // 指数退避
			}
		}

		client := s.CreateClient(resource)
		// 注意：流式请求不能立即关闭client，需要在流结束后关闭
		// 这里先不defer，让调用者管理client生命周期

		// 调用流式API
		stream := client.Chat().Completions().NewStreaming(ctx, params)
		if stream == nil {
			lastErr = fmt.Errorf("failed to create stream")
			client.Close()
			// 检查是否可重试
			if !s.isRetryableError(lastErr) || attempt >= maxRetries {
				break
			}
			continue
		}

		// 成功创建流，返回
		return stream, nil
	}

	// 所有重试都失败
	if lastErr == nil {
		lastErr = fmt.Errorf("failed to create stream after %d attempts", maxRetries+1)
	}
	logger.Error("Streaming chat completion request failed after retries", logger.F("component", "service"), logger.F("attempts", maxRetries+1), logger.F("error", lastErr.Error()), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
	return lastStream, lastErr
}

// CompletionRequest 文本补全请求参数
type CompletionRequest struct {
	Model            string
	Prompt           string
	MaxTokens        int64
	Temperature      float64
	TopP             float64
	Stop             []string
	PresencePenalty  float64
	FrequencyPenalty float64
	User             string
}

// CompletionResponse 文本补全响应
type CompletionResponse struct {
	Response *openaiSDK.Completion
	Duration time.Duration
	Error    error
}

// CreateCompletion 创建文本补全请求（带自动重试）
func (s *OpenAIService) CreateCompletion(ctx context.Context, resource *storage.LLMResource, req CompletionRequest) (*CompletionResponse, error) {
	return s.createCompletionWithRetry(ctx, resource, req)
}

// createCompletionWithRetry 带重试的文本补全请求
func (s *OpenAIService) createCompletionWithRetry(ctx context.Context, resource *storage.LLMResource, req CompletionRequest) (*CompletionResponse, error) {
	client := s.CreateClient(resource)
	defer client.Close()

	// 确定使用的模型
	modelToUse := req.Model
	if resource.Model != "" {
		modelToUse = resource.Model
	}

	// 构建请求参数
	params := openaiSDK.CompletionNewParams{
		Model: openaiSDK.CompletionNewParamsModel(modelToUse),
		Prompt: openaiSDK.CompletionNewParamsPromptUnion{
			OfString: param.NewOpt(req.Prompt),
		},
	}

	// 设置可选参数
	if req.MaxTokens > 0 {
		params.MaxTokens = param.NewOpt(req.MaxTokens)
	}
	if req.Temperature > 0 {
		params.Temperature = param.NewOpt(req.Temperature)
	}
	if req.TopP > 0 {
		params.TopP = param.NewOpt(req.TopP)
	}
	if len(req.Stop) > 0 {
		params.Stop = openaiSDK.CompletionNewParamsStopUnion{
			OfStringArray: req.Stop,
		}
	}
	if req.PresencePenalty != 0 {
		params.PresencePenalty = param.NewOpt(req.PresencePenalty)
	}
	if req.FrequencyPenalty != 0 {
		params.FrequencyPenalty = param.NewOpt(req.FrequencyPenalty)
	}
	if req.User != "" {
		params.User = param.NewOpt(req.User)
	}

	// 获取重试配置
	maxRetries := config.C.Provider.MaxRetries
	retryDelay := config.C.Provider.RetryDelay
	if maxRetries <= 0 {
		maxRetries = 0 // 禁用重试
	}

	// 记录开始时间
	startTime := time.Now()
	var lastErr error
	var lastResponse *openaiSDK.Completion

	// 重试逻辑
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// 重试前等待
			logger.Debug("Retrying completion request", logger.F("component", "service"), logger.F("attempt", attempt), logger.F("max_retries", maxRetries), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
			select {
			case <-ctx.Done():
				return &CompletionResponse{
					Response: nil,
					Duration: time.Since(startTime),
					Error:    ctx.Err(),
				}, ctx.Err()
			case <-time.After(retryDelay * time.Duration(attempt)): // 指数退避
			}
		}

		// 调用API
		response, err := client.Completions().New(ctx, params)
		duration := time.Since(startTime)

		if err == nil {
			// 成功，返回结果
			return &CompletionResponse{
				Response: response,
				Duration: duration,
				Error:    nil,
			}, nil
		}

		lastErr = err
		lastResponse = response

		// 检查是否可重试
		if !s.isRetryableError(err) {
			logger.Debug("Error is not retryable", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("resource_id", resource.ID))
			break
		}

		// 如果还有重试机会，继续
		if attempt < maxRetries {
			logger.Warn("Completion request failed, will retry", logger.F("component", "service"), logger.F("attempt", attempt+1), logger.F("max_retries", maxRetries), logger.F("error", err.Error()), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
		}
	}

	// 所有重试都失败
	duration := time.Since(startTime)
	logger.Error("Completion request failed after retries", logger.F("component", "service"), logger.F("attempts", maxRetries+1), logger.F("error", lastErr.Error()), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
	return &CompletionResponse{
		Response: lastResponse,
		Duration: duration,
		Error:    lastErr,
	}, lastErr
}

// EmbeddingRequest 嵌入请求参数
type EmbeddingRequest struct {
	Model string
	Input string
}

// EmbeddingResponse 嵌入响应
type EmbeddingResponse struct {
	Response *openaiSDK.CreateEmbeddingResponse
	Duration time.Duration
	Error    error
}

// CreateEmbedding 创建嵌入请求（带自动重试）
func (s *OpenAIService) CreateEmbedding(ctx context.Context, resource *storage.LLMResource, req EmbeddingRequest) (*EmbeddingResponse, error) {
	return s.createEmbeddingWithRetry(ctx, resource, req)
}

// createEmbeddingWithRetry 带重试的嵌入请求
func (s *OpenAIService) createEmbeddingWithRetry(ctx context.Context, resource *storage.LLMResource, req EmbeddingRequest) (*EmbeddingResponse, error) {
	client := s.CreateEmbeddingClient(resource)
	defer client.Close()

	// 确定使用的模型
	modelToUse := req.Model
	if modelToUse == "" {
		modelToUse = resource.Model
	}
	if modelToUse == "" {
		modelToUse = "text-embedding-3-small" // 默认模型
	}

	// 构建请求参数
	params := openaiSDK.EmbeddingNewParams{
		Model: modelToUse,
		Input: openaiSDK.EmbeddingNewParamsInputUnion{
			OfString: openaiSDK.String(req.Input),
		},
	}

	// 获取重试配置
	maxRetries := config.C.Provider.MaxRetries
	retryDelay := config.C.Provider.RetryDelay
	if maxRetries <= 0 {
		maxRetries = 0 // 禁用重试
	}

	// 记录开始时间
	startTime := time.Now()
	var lastErr error
	var lastResponse *openaiSDK.CreateEmbeddingResponse

	// 重试逻辑
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// 重试前等待
			logger.Debug("Retrying embedding request", logger.F("component", "service"), logger.F("attempt", attempt), logger.F("max_retries", maxRetries), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
			select {
			case <-ctx.Done():
				return &EmbeddingResponse{
					Response: nil,
					Duration: time.Since(startTime),
					Error:    ctx.Err(),
				}, ctx.Err()
			case <-time.After(retryDelay * time.Duration(attempt)): // 指数退避
			}
		}

		// 调用API
		response, err := client.New(ctx, params)
		duration := time.Since(startTime)

		if err == nil {
			// 成功，返回结果
			return &EmbeddingResponse{
				Response: response,
				Duration: duration,
				Error:    nil,
			}, nil
		}

		lastErr = err
		lastResponse = response

		// 检查是否可重试
		if !s.isRetryableError(err) {
			logger.Debug("Error is not retryable", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("resource_id", resource.ID))
			break
		}

		// 如果还有重试机会，继续
		if attempt < maxRetries {
			logger.Warn("Embedding request failed, will retry", logger.F("component", "service"), logger.F("attempt", attempt+1), logger.F("max_retries", maxRetries), logger.F("error", err.Error()), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
		}
	}

	// 所有重试都失败
	duration := time.Since(startTime)
	logger.Error("Embedding request failed after retries", logger.F("component", "service"), logger.F("attempts", maxRetries+1), logger.F("error", lastErr.Error()), logger.F("resource_id", resource.ID), logger.F("model", modelToUse))
	return &EmbeddingResponse{
		Response: lastResponse,
		Duration: duration,
		Error:    lastErr,
	}, lastErr
}

// isRetryableError 判断错误是否可重试
func (s *OpenAIService) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// 检查是否是网络错误
	if strings.Contains(errStr, "network") || strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "timeout") || strings.Contains(errStr, "EOF") ||
		strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset") {
		return true
	}

	// 检查是否是HTTP错误
	// OpenAI SDK 可能会返回包含状态码的错误信息
	if strings.Contains(errStr, "429") || // Too Many Requests
		strings.Contains(errStr, "500") || // Internal Server Error
		strings.Contains(errStr, "502") || // Bad Gateway
		strings.Contains(errStr, "503") || // Service Unavailable
		strings.Contains(errStr, "504") { // Gateway Timeout
		return true
	}

	// 检查是否是临时性错误
	if strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "overloaded") ||
		strings.Contains(errStr, "temporarily unavailable") || strings.Contains(errStr, "try again") {
		return true
	}

	// 检查是否是上下文取消错误（不可重试）
	if strings.Contains(errStr, "context canceled") || strings.Contains(errStr, "context deadline exceeded") {
		return false
	}

	// 检查是否是认证错误（不可重试）
	if strings.Contains(errStr, "401") || strings.Contains(errStr, "unauthorized") ||
		strings.Contains(errStr, "403") || strings.Contains(errStr, "forbidden") ||
		strings.Contains(errStr, "invalid api key") || strings.Contains(errStr, "authentication") {
		return false
	}

	// 检查是否是客户端错误（4xx，除了429，不可重试）
	if strings.Contains(errStr, "400") || strings.Contains(errStr, "404") ||
		strings.Contains(errStr, "bad request") || strings.Contains(errStr, "not found") {
		return false
	}

	// 默认情况下，对于未知错误，不重试（避免无限重试）
	return false
}

// TestChatResource 测试chat类型资源
func (s *OpenAIService) TestChatResource(ctx context.Context, resource *storage.LLMResource) map[string]interface{} {
	modelToUse := resource.Model
	if modelToUse == "" {
		modelToUse = "gpt-3.5-turbo" // 默认模型
	}

	req := ChatCompletionRequest{
		Model: modelToUse,
		Messages: []openaiSDK.ChatCompletionMessageParamUnion{
			openaiSDK.UserMessage("Hello"),
		},
		MaxTokens: 10, // 限制token数量，快速测试
	}

	resp, err := s.CreateChatCompletion(ctx, resource, req)
	if err != nil {
		return map[string]interface{}{
			"success":     false,
			"error":       err.Error(),
			"message":     fmt.Sprintf("测试失败: %s", err.Error()),
			"duration_ms": resp.Duration.Milliseconds(),
		}
	}

	// 检查响应
	if resp.Response == nil || len(resp.Response.Choices) == 0 {
		return map[string]interface{}{
			"success":     false,
			"error":       "Empty response",
			"message":     "API返回了空响应",
			"duration_ms": resp.Duration.Milliseconds(),
		}
	}

	return map[string]interface{}{
		"success":  true,
		"message":  "测试成功",
		"model":    string(resp.Response.Model),
		"response": resp.Response.Choices[0].Message.Content,
		"usage": map[string]interface{}{
			"prompt_tokens":     resp.Response.Usage.PromptTokens,
			"completion_tokens": resp.Response.Usage.CompletionTokens,
			"total_tokens":      resp.Response.Usage.TotalTokens,
		},
		"duration_ms": resp.Duration.Milliseconds(),
	}
}

// TestEmbeddingResource 测试embedding类型资源
func (s *OpenAIService) TestEmbeddingResource(ctx context.Context, resource *storage.LLMResource) map[string]interface{} {
	modelToUse := resource.Model
	if modelToUse == "" {
		modelToUse = "text-embedding-3-small" // 默认模型
	}

	req := EmbeddingRequest{
		Model: modelToUse,
		Input: "test",
	}

	resp, err := s.CreateEmbedding(ctx, resource, req)
	if err != nil {
		return map[string]interface{}{
			"success":     false,
			"error":       err.Error(),
			"message":     fmt.Sprintf("测试失败: %s", err.Error()),
			"duration_ms": resp.Duration.Milliseconds(),
		}
	}

	// 检查响应
	if resp.Response == nil || len(resp.Response.Data) == 0 {
		return map[string]interface{}{
			"success":     false,
			"error":       "Empty response",
			"message":     "API返回了空响应",
			"duration_ms": resp.Duration.Milliseconds(),
		}
	}

	return map[string]interface{}{
		"success":             true,
		"message":             "测试成功",
		"model":               resp.Response.Model,
		"embedding_dimension": len(resp.Response.Data[0].Embedding),
		"usage": map[string]interface{}{
			"prompt_tokens": resp.Response.Usage.PromptTokens,
			"total_tokens":  resp.Response.Usage.TotalTokens,
		},
		"duration_ms": resp.Duration.Milliseconds(),
	}
}
