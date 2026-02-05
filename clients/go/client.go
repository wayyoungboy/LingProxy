package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// LingProxyClient is a standard Go client for LingProxy AI API Gateway
type LingProxyClient struct {
	client *openai.Client
}

// ClientOptions contains options for initializing the client
type ClientOptions struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// NewClient creates a new LingProxy client
// If APIKey is empty, it will try to get from LINGPROXY_API_KEY environment variable
// If BaseURL is empty, it defaults to "http://localhost:8080/llm/v1"
func NewClient(options *ClientOptions) (*LingProxyClient, error) {
	apiKey := options.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("LINGPROXY_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("API key is required. Either pass it as parameter or set LINGPROXY_API_KEY environment variable")
		}
	}

	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:8080/llm/v1"
	}

	timeout := options.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
		option.WithTimeout(timeout),
	)

	return &LingProxyClient{
		client: client,
	}, nil
}

// GetClient returns the underlying OpenAI client for direct access
func (c *LingProxyClient) GetClient() *openai.Client {
	return c.client
}

// ListModels lists all available models
func (c *LingProxyClient) ListModels(ctx context.Context) ([]openai.Model, error) {
	models, err := c.client.Models.List(ctx)
	if err != nil {
		return nil, err
	}
	return models.Data, nil
}

// CreateChatCompletion creates a chat completion
func (c *LingProxyClient) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(req.Messages))
	for i, msg := range req.Messages {
		switch msg.Role {
		case "system":
			messages[i] = openai.SystemMessage(msg.Content)
		case "user":
			messages[i] = openai.UserMessage(msg.Content)
		case "assistant":
			messages[i] = openai.AssistantMessage(msg.Content)
		default:
			messages[i] = openai.UserMessage(msg.Content)
		}
	}

	params := openai.ChatCompletionNewParams{
		Model:    req.Model,
		Messages: messages,
	}

	if req.Temperature != nil {
		params.Temperature = req.Temperature
	}
	if req.MaxTokens != nil {
		params.MaxTokens = req.MaxTokens
	}
	if req.TopP != nil {
		params.TopP = req.TopP
	}
	if req.Stream != nil && *req.Stream {
		params.Stream = req.Stream
	}

	response, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}

	return &ChatCompletionResponse{
		ID:      response.ID,
		Object:  response.Object,
		Created: response.Created,
		Model:   response.Model,
		Choices: response.Choices,
		Usage:   response.Usage,
	}, nil
}

// CreateCompletion creates a text completion
func (c *LingProxyClient) CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	params := openai.CompletionNewParams{
		Model:  openai.CompletionNewParamsModel(req.Model),
		Prompt: openai.CompletionNewParamsPromptUnion{OfString: &req.Prompt},
	}

	if req.Temperature != nil {
		params.Temperature = req.Temperature
	}
	if req.MaxTokens != nil {
		params.MaxTokens = req.MaxTokens
	}
	if req.TopP != nil {
		params.TopP = req.TopP
	}

	response, err := c.client.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}

	return &CompletionResponse{
		ID:      response.ID,
		Object:  response.Object,
		Created: response.Created,
		Model:   response.Model,
		Choices: response.Choices,
		Usage:   response.Usage,
	}, nil
}

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Model       string
	Messages    []ChatMessage
	Temperature *float64
	MaxTokens   *int64
	TopP        *float64
	Stream      *bool
}

// ChatMessage represents a single chat message
type ChatMessage struct {
	Role    string
	Content string
}

// ChatCompletionResponse represents a chat completion response
type ChatCompletionResponse struct {
	ID      string
	Object  string
	Created int64
	Model   string
	Choices []openai.ChatCompletionChoice
	Usage   openai.CompletionUsage
}

// CompletionRequest represents a text completion request
type CompletionRequest struct {
	Model       string
	Prompt      string
	Temperature *float64
	MaxTokens   *int64
	TopP        *float64
}

// CompletionResponse represents a text completion response
type CompletionResponse struct {
	ID      string
	Object  string
	Created int64
	Model   string
	Choices []openai.CompletionChoice
	Usage   openai.CompletionUsage
}
