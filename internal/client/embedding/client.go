package embedding

import (

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Client 封装了向量化模型客户端
type Client struct {
	client openai.Client
	model  string
}

// NewClient 创建新的向量化模型客户端
func NewClient(apiKey string, baseURL string, model string) *Client {
	options := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	if baseURL != "" {
		options = append(options, option.WithBaseURL(baseURL))
	}

	if model == "" {
		model = "text-embedding-3-small"
	}

	return &Client{
		client: openai.NewClient(options...),
		model:  model,
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

// New 创建嵌入请求

func (c *Client) New(ctx context.Context, params openai.EmbeddingNewParams) (*openai.CreateEmbeddingResponse, error) {
	if params.Model == "" {
		params.Model = c.model
	}
	return c.client.Embeddings.New(ctx, params)
}

// CreateEmbedding 创建嵌入请求（兼容旧接口）
func (c *Client) CreateEmbedding(ctx context.Context, input interface{}) (*openai.CreateEmbeddingResponse, error) {
	// 只支持字符串输入
	inputStr, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("unsupported input type: %T, only string is supported", input)
	}
	
	// 创建嵌入请求
	response, err := c.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: c.model,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(inputStr),
		},
	})
	
	return response, err
}
