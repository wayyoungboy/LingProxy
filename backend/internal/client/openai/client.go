package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

// Client 封装了 OpenAI 官方客户端
type Client struct {
	client openai.Client
}

// NewClient 创建新的 OpenAI 客户端
func NewClient(apiKey string, baseURL string) *Client {
	options := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	if baseURL != "" {
		options = append(options, option.WithBaseURL(baseURL))
	}

	return &Client{
		client: openai.NewClient(options...),
	}
}

// NewClientWithOptions 使用自定义选项创建新的 OpenAI 客户端
func NewClientWithOptions(options ...option.RequestOption) *Client {
	return &Client{
		client: openai.NewClient(options...),
	}
}

// GetClient 获取底层官方客户端
func (c *Client) GetClient() interface{} {
	return c.client
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	// OpenAI 官方客户端不需要特殊关闭操作
	return nil
}

// Chat 返回聊天相关的 API
func (c *Client) Chat() *Chat {
	return &Chat{client: c.client}
}

// Completions 返回文本补全相关的 API
func (c *Client) Completions() *Completions {
	return &Completions{client: c.client}
}

// Chat 提供聊天相关的 API
type Chat struct {
	client openai.Client
}

// Completions 提供聊天补全相关的 API
func (c *Chat) Completions() *ChatCompletions {
	return &ChatCompletions{client: c.client}
}

// ChatCompletions 提供聊天补全相关的 API
type ChatCompletions struct {
	client openai.Client
}

// New 创建聊天补全请求
func (c *ChatCompletions) New(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	return c.client.Chat.Completions.New(ctx, params)
}

// NewStreaming 创建流式聊天补全请求
func (c *ChatCompletions) NewStreaming(ctx context.Context, params openai.ChatCompletionNewParams) *ssestream.Stream[openai.ChatCompletionChunk] {
	return c.client.Chat.Completions.NewStreaming(ctx, params)
}

// GetAccumulator 创建聊天补全累加器
func (c *ChatCompletions) GetAccumulator() interface{} {
	return nil
}

// Completions 提供文本补全相关的 API
type Completions struct {
	client openai.Client
}

// New 创建文本补全请求
func (c *Completions) New(ctx context.Context, params openai.CompletionNewParams) (*openai.Completion, error) {
	return c.client.Completions.New(ctx, params)
}

// NewStreaming 创建流式文本补全请求
func (c *Completions) NewStreaming(ctx context.Context, params openai.CompletionNewParams) interface{} {
	return c.client.Completions.NewStreaming(ctx, params)
}

// Helper functions for creating message params

// SystemMessage 创建系统消息
func SystemMessage(content string) openai.ChatCompletionMessageParamUnion {
	return openai.SystemMessage(content)
}

// UserMessage 创建用户消息
func UserMessage(content string) openai.ChatCompletionMessageParamUnion {
	return openai.UserMessage(content)
}

// AssistantMessage 创建助手消息
func AssistantMessage(content string) openai.ChatCompletionMessageParamUnion {
	return openai.AssistantMessage(content)
}

// Helper functions for creating tool params

// CreateFunctionTool 创建函数工具
func CreateFunctionTool(name string, description string, parameters map[string]interface{}) openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        name,
		Description: openai.String(description),
		Parameters:  openai.FunctionParameters(parameters),
	})
}

// Error 定义客户端错误
type Error struct {
	Message string
	Err     error
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return fmt.Sprintf("openai client error: %s: %v", e.Message, e.Err)
}

// Unwrap 返回底层错误
func (e *Error) Unwrap() error {
	return e.Err
}
