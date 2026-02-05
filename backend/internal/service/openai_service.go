package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lingproxy/lingproxy/internal/client/embedding"
	"github.com/lingproxy/lingproxy/internal/client/openai"
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

// CreateChatCompletion 创建聊天补全请求
func (s *OpenAIService) CreateChatCompletion(ctx context.Context, resource *storage.LLMResource, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	client := s.CreateClient(resource)
	defer client.Close()

	// 确定使用的模型
	modelToUse := req.Model
	if resource.Model != "" {
		modelToUse = resource.Model
	}

	// 构建请求参数
	params := s.buildChatCompletionParams(modelToUse, req)

	// 记录开始时间
	startTime := time.Now()

	// 调用API
	response, err := client.Chat().Completions().New(ctx, params)
	duration := time.Since(startTime)

	if err != nil {
		return &ChatCompletionResponse{
			Response: nil,
			Duration: duration,
			Error:    err,
		}, err
	}

	return &ChatCompletionResponse{
		Response: response,
		Duration: duration,
		Error:    nil,
	}, nil
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

// CreateChatCompletionStream 创建流式聊天补全请求
func (s *OpenAIService) CreateChatCompletionStream(ctx context.Context, resource *storage.LLMResource, req ChatCompletionRequest) (*ssestream.Stream[openaiSDK.ChatCompletionChunk], error) {
	client := s.CreateClient(resource)
	defer client.Close()

	// 确定使用的模型
	modelToUse := req.Model
	if resource.Model != "" {
		modelToUse = resource.Model
	}

	// 构建请求参数
	params := s.buildChatCompletionParams(modelToUse, req)

	// 调用流式API
	stream := client.Chat().Completions().NewStreaming(ctx, params)
	if stream == nil {
		return nil, fmt.Errorf("failed to create stream")
	}
	return stream, nil
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

// CreateCompletion 创建文本补全请求
func (s *OpenAIService) CreateCompletion(ctx context.Context, resource *storage.LLMResource, req CompletionRequest) (*CompletionResponse, error) {
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

	// 记录开始时间
	startTime := time.Now()

	// 调用API
	response, err := client.Completions().New(ctx, params)
	duration := time.Since(startTime)

	if err != nil {
		return &CompletionResponse{
			Response: nil,
			Duration: duration,
			Error:    err,
		}, err
	}

	return &CompletionResponse{
		Response: response,
		Duration: duration,
		Error:    nil,
	}, nil
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

// CreateEmbedding 创建嵌入请求
func (s *OpenAIService) CreateEmbedding(ctx context.Context, resource *storage.LLMResource, req EmbeddingRequest) (*EmbeddingResponse, error) {
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

	// 记录开始时间
	startTime := time.Now()

	// 调用API
	response, err := client.New(ctx, params)
	duration := time.Since(startTime)

	if err != nil {
		return &EmbeddingResponse{
			Response: nil,
			Duration: duration,
			Error:    err,
		}, err
	}

	return &EmbeddingResponse{
		Response: response,
		Duration: duration,
		Error:    nil,
	}, nil
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
