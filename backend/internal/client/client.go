package client

// Client 定义了向资源方发送请求的通用接口
type Client interface {
	// Close 关闭客户端连接
	Close() error
}

// OpenAIClient 定义了 OpenAI 客户端接口
type OpenAIClient interface {
	Client

	// GetClient 获取底层客户端实例
	GetClient() interface{}

	// Chat 返回聊天相关的 API
	Chat() interface{}

	// Completions 返回文本补全相关的 API
	Completions() interface{}

	// Embeddings 返回嵌入相关的 API
	Embeddings() interface{}
}

// ClientOptions 定义了创建客户端的选项
type ClientOptions struct {
	APIKey      string
	BaseURL     string
	Model       string
	ExtraConfig map[string]interface{}
}
