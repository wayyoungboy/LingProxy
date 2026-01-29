package storage

import (
	"time"
)

// Model 模型配置 - 深度原子化模型管理
type Model struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	LLMResourceID  string            `json:"llm_resource_id"`
	ModelID        string            `json:"model_id"`     // 提供商内部的模型标识
	Type        string            `json:"type"`         // 模型类型: chat, completion, embedding, image
	Category    string            `json:"category"`     // 模型分类: gpt, claude, gemini, llama, etc.
	Version     string            `json:"version"`      // 模型版本
	Description string            `json:"description"`
	Capabilities []string         `json:"capabilities"` // 模型能力: text-generation, code-completion, image-analysis, etc.
	Pricing     ModelPricing      `json:"pricing"`      // 定价信息
	Limits      ModelLimits       `json:"limits"`       // 使用限制
	Parameters  ModelParameters   `json:"parameters"`   // 默认参数
	Features    ModelFeatures     `json:"features"`     // 功能特性
	Status      string            `json:"status"`       // active, inactive, deprecated
	Metadata    map[string]interface{} `json:"metadata"` // 扩展元数据
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ModelPricing 模型定价信息
type ModelPricing struct {
	InputTokenPrice  float64 `json:"input_token_price"`  // 每1000个输入token的价格（美元）
	OutputTokenPrice float64 `json:"output_token_price"` // 每1000个输出token的价格（美元）
	ImagePrice       float64 `json:"image_price,omitempty"` // 每张图片的价格（美元）
	AudioPrice       float64 `json:"audio_price,omitempty"` // 每分钟音频的价格（美元）
	Currency         string  `json:"currency"`           // 货币单位
}

// ModelLimits 模型使用限制
type ModelLimits struct {
	MaxTokens        int `json:"max_tokens"`         // 最大token数
	MaxInputTokens   int `json:"max_input_tokens"`   // 最大输入token数
	MaxOutputTokens  int `json:"max_output_tokens"`  // 最大输出token数
	ContextWindow    int `json:"context_window"`     // 上下文窗口大小
	MaxTemperature   float64 `json:"max_temperature"` // 最大温度值
	MinTemperature   float64 `json:"min_temperature"` // 最小温度值
	MaxTopP          float64 `json:"max_top_p"`        // 最大top-p值
	MaxPresencePenalty float64 `json:"max_presence_penalty"`
	MaxFrequencyPenalty float64 `json:"max_frequency_penalty"`
}

// ModelParameters 模型默认参数
type ModelParameters struct {
	Temperature      float64  `json:"temperature"`
	TopP             float64  `json:"top_p"`
	PresencePenalty  float64  `json:"presence_penalty"`
	FrequencyPenalty float64  `json:"frequency_penalty"`
	MaxTokens        int      `json:"max_tokens"`
	StopSequences    []string `json:"stop_sequences,omitempty"`
	SystemPrompt     string   `json:"system_prompt,omitempty"`
}

// ModelFeatures 模型功能特性
type ModelFeatures struct {
	Streaming        bool `json:"streaming"`         // 是否支持流式响应
	FunctionCalling  bool `json:"function_calling"`  // 是否支持函数调用
	Vision           bool `json:"vision"`            // 是否支持视觉理解
	CodeCompletion   bool `json:"code_completion"`   // 是否支持代码补全
	Multilingual     bool `json:"multilingual"`      // 是否支持多语言
	FineTuning       bool `json:"fine_tuning"`       // 是否支持微调
	Embeddings       bool `json:"embeddings"`        // 是否支持嵌入
	Moderation       bool `json:"moderation"`        // 是否支持内容审核
}

// ModelEndpoint 模型端点映射
type ModelEndpoint struct {
	ID             string `json:"id"`
	ModelID        string `json:"model_id"`
	LLMResourceID  string `json:"llm_resource_id"`
	Endpoint       string `json:"endpoint"` // 完整的API端点路径
	Method   string `json:"method"`
	Headers  map[string]string `json:"headers"`
	Status   string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ModelVersion 模型版本管理
type ModelVersion struct {
	ID        string    `json:"id"`
	ModelID   string    `json:"model_id"`
	Version   string    `json:"version"`
	ReleaseDate time.Time `json:"release_date"`
	Changelog string    `json:"changelog"`
	Status    string    `json:"status"` // stable, beta, alpha, deprecated
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}