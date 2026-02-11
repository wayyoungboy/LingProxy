package storage

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

// User 用户模型
type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Password     string     `json:"-" gorm:"-"` // 不序列化到JSON，不存储到数据库（存储哈希值）
	PasswordHash string     `json:"-"`          // 密码哈希值
	APIKey       string     `json:"api_key"`
	Role         string     `json:"role"`   // admin, user
	Status       string     `json:"status"` // active, inactive, suspended
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// LLMResource LLM资源模型
// 一个LLM资源代表一个可调用的AI服务配置，包含base_url、api_key和model三个核心要素
// 这三个要素共同构成一个完整的可调用资源
type LLMResource struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`        // 模型类别: chat(对话), image(生图), embedding(嵌入), rerank(重排序), audio(语音), video(视频)
	Driver     string    `json:"driver"`      // 驱动: openai（目前仅支持openai）
	Model      string    `json:"model"`       // 模型标识：此资源对应的模型标识（如gpt-4, gpt-3.5-turbo），与base_url和api_key共同构成可调用资源
	BaseURL    string    `json:"base_url"`    // API基础URL：此资源的API端点地址
	APIKey     string    `json:"api_key"`     // API密钥：此资源的认证密钥
	Status     string    `json:"status"`      // 资源状态: active(活跃), inactive(禁用)
	TestStatus string    `json:"test_status"` // 测试状态: untested(未测试), passed(测试通过), failed(测试失败)
	
	// 类型特定配置（临时使用 JSON 字段存储，重构后将迁移到扩展表）
	EmbeddingConfig string `json:"embedding_config,omitempty" gorm:"type:text"` // JSON格式的embedding配置
	RerankConfig    string `json:"rerank_config,omitempty" gorm:"type:text"`    // JSON格式的rerank配置
	ChatConfig      string `json:"chat_config,omitempty" gorm:"type:text"`       // JSON格式的chat配置
	ImageConfig     string `json:"image_config,omitempty" gorm:"type:text"`      // JSON格式的image配置
	AudioConfig     string `json:"audio_config,omitempty" gorm:"type:text"`      // JSON格式的audio配置
	
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// EmbeddingConfig Embedding类型配置结构
type EmbeddingConfig struct {
	SupportedDimensions []int `json:"supported_dimensions"` // 支持的维度列表
	DefaultDimension    *int  `json:"default_dimension,omitempty"`
	Normalize           bool  `json:"normalize"`
	MaxInputTokens      *int  `json:"max_input_tokens,omitempty"`
}

// RerankConfig Rerank类型配置结构
type RerankConfig struct {
	DefaultTopN       *int `json:"default_top_n,omitempty"`
	MaxDocuments      *int `json:"max_documents,omitempty"`
	MaxQueryLength    *int `json:"max_query_length,omitempty"`
	MaxDocumentLength *int `json:"max_document_length,omitempty"`
}

// ChatConfig Chat类型配置结构
type ChatConfig struct {
	MaxTokens            *int   `json:"max_tokens,omitempty"`
	ContextWindow        *int   `json:"context_window,omitempty"`
	SupportsStreaming    bool   `json:"supports_streaming"`
	SupportsFunctionCalling bool `json:"supports_function_calling"`
	SupportsVision       bool   `json:"supports_vision"`
}

// Endpoint 端点模型
type Endpoint struct {
	ID            string    `json:"id"`
	LLMResourceID string    `json:"llm_resource_id"`
	Path          string    `json:"path"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Request 请求模型
type Request struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	LLMResourceID string    `json:"llm_resource_id"` // LLM资源ID，用于关联资源
	Endpoint      string    `json:"endpoint"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	Duration      int64     `json:"duration"`
	Tokens        int       `json:"tokens"`
	CreatedAt     time.Time `json:"created_at"`
}

// RequestQueryParams 请求查询参数
type RequestQueryParams struct {
	Limit     int       // 限制返回数量
	StartTime *time.Time // 开始时间（可选）
	EndTime   *time.Time // 结束时间（可选）
	Endpoint  string    // 请求路径（可选，支持模糊匹配）
	Status    string    // 状态（可选）
}

// Response 响应模型
type Response struct {
	ID        string    `json:"id"`
	RequestID string    `json:"request_id"`
	Status    int       `json:"status"`
	Body      string    `json:"body"`
	Headers   string    `json:"headers"`
	CreatedAt time.Time `json:"created_at"`
}

// Quota 配额模型
type Quota struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"user_id"`
	RPM                  int       `json:"rpm"`
	TPM                  int       `json:"tpm"`
	DailyTokens          int       `json:"daily_tokens"`
	DailyRequests        int       `json:"daily_requests"`
	CurrentRPM           int       `json:"current_rpm"`
	CurrentTPM           int       `json:"current_tpm"`
	CurrentDailyTokens   int       `json:"current_daily_tokens"`
	CurrentDailyRequests int       `json:"current_daily_requests"`
	ResetAt              time.Time `json:"reset_at"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Statistics 统计模型
type Statistics struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	TotalRequests int64     `json:"total_requests"`
	TotalTokens   int64     `json:"total_tokens"`
	SuccessCount  int64     `json:"success_count"`
	ErrorCount    int64     `json:"error_count"`
	LastRequestAt time.Time `json:"last_request_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// APIKey API Key模型（用于API请求认证）
type APIKey struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`                // API Key名称/描述
	APIKey     string     `json:"api_key"`             // API Key值
	Prefix     string     `json:"prefix"`             // API Key前缀（用于显示）
	Status     string     `json:"status"`             // active/inactive
	PolicyID   string     `json:"policy_id,omitempty"` // 关联的策略ID（向后兼容，已废弃，使用按类型策略）
	
	// 模型许可：允许使用的模型ID列表（空列表表示允许所有模型）
	AllowedModels string `json:"allowed_models,omitempty" gorm:"type:text"` // JSON数组格式: ["gpt-4", "gpt-3.5-turbo"]
	
	// 按类型配置的策略（优先级高于 PolicyID）
	ChatPolicyID     string `json:"chat_policy_id,omitempty"`     // Chat类型模型的策略ID
	EmbeddingPolicyID string `json:"embedding_policy_id,omitempty"` // Embedding类型模型的策略ID
	RerankPolicyID    string `json:"rerank_policy_id,omitempty"`   // Rerank类型模型的策略ID
	ImagePolicyID     string `json:"image_policy_id,omitempty"`     // Image类型模型的策略ID
	AudioPolicyID     string `json:"audio_policy_id,omitempty"`     // Audio类型模型的策略ID
	VideoPolicyID     string `json:"video_policy_id,omitempty"`    // Video类型模型的策略ID
	
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TableName 指定数据库表名
func (APIKey) TableName() string {
	return "api_keys"
}

// GetAllowedModels 获取允许的模型列表
func (k *APIKey) GetAllowedModels() []string {
	if k.AllowedModels == "" {
		return nil // nil 表示允许所有模型
	}
	var models []string
	if err := json.Unmarshal([]byte(k.AllowedModels), &models); err != nil {
		return nil
	}
	return models
}

// SetAllowedModels 设置允许的模型列表
func (k *APIKey) SetAllowedModels(models []string) error {
	if len(models) == 0 {
		k.AllowedModels = ""
		return nil
	}
	data, err := json.Marshal(models)
	if err != nil {
		return err
	}
	k.AllowedModels = string(data)
	return nil
}

// GetPolicyIDByType 根据模型类型获取策略ID
// 优先返回按类型配置的策略，如果没有则返回通用策略ID（向后兼容）
func (k *APIKey) GetPolicyIDByType(modelType string) string {
	switch modelType {
	case "chat", "completion":
		if k.ChatPolicyID != "" {
			return k.ChatPolicyID
		}
	case "embedding":
		if k.EmbeddingPolicyID != "" {
			return k.EmbeddingPolicyID
		}
	case "rerank":
		if k.RerankPolicyID != "" {
			return k.RerankPolicyID
		}
	case "image":
		if k.ImagePolicyID != "" {
			return k.ImagePolicyID
		}
	case "audio":
		if k.AudioPolicyID != "" {
			return k.AudioPolicyID
		}
	case "video":
		if k.VideoPolicyID != "" {
			return k.VideoPolicyID
		}
	}
	// 向后兼容：如果没有按类型配置的策略，使用通用策略ID
	return k.PolicyID
}

// IsModelAllowed 检查模型是否被允许使用
// 如果 AllowedModels 为空，表示允许所有模型
func (k *APIKey) IsModelAllowed(modelID string) bool {
	allowedModels := k.GetAllowedModels()
	if len(allowedModels) == 0 {
		return true // 空列表表示允许所有模型
	}
	for _, allowed := range allowedModels {
		if allowed == modelID {
			return true
		}
	}
	return false
}

// Token 保持向后兼容的类型别名
// Deprecated: 使用 APIKey 代替
type Token = APIKey

// PolicyTemplate 策略模板模型
type PolicyTemplate struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Type              string    `json:"type"` // random, round_robin, weighted, model_match, regex_match, priority, failover
	Description       string    `json:"description"`
	ParametersSchema  string    `json:"parameters_schema"`  // JSON Schema
	DefaultParameters string    `json:"default_parameters"` // JSON
	Builtin           bool      `json:"builtin"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Policy 策略实例模型
type Policy struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	TemplateID string    `json:"template_id"`
	Type       string    `json:"type"`
	Parameters string    `json:"parameters"` // JSON
	Enabled    bool      `json:"enabled"`
	Builtin    bool      `json:"builtin"` // 是否为内置策略
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Model 模型配置 - 深度原子化模型管理
type Model struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	LLMResourceID string    `json:"llm_resource_id"`
	ModelID       string    `json:"model_id"` // 提供商内部的模型标识
	Type          string    `json:"type"`     // 模型类型: chat, completion, embedding, image
	Category      string    `json:"category"` // 模型分类: gpt, claude, gemini, llama, etc.
	Version       string    `json:"version"`  // 模型版本
	Description   string    `json:"description"`
	Capabilities  string    `json:"capabilities"` // 模型能力: JSON字符串格式的能力数组
	Pricing       string    `json:"pricing"`      // 定价信息: JSON字符串格式
	Limits        string    `json:"limits"`       // 使用限制: JSON字符串格式
	Parameters    string    `json:"parameters"`   // 默认参数: JSON字符串格式
	Features      string    `json:"features"`     // 功能特性: JSON字符串格式
	Status        string    `json:"status"`       // active, inactive, deprecated
	Metadata      string    `json:"metadata"`     // 扩展元数据: JSON字符串格式
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ModelPricing 模型定价信息
type ModelPricing struct {
	InputTokenPrice  float64 `json:"input_token_price"`     // 每1000个输入token的价格（美元）
	OutputTokenPrice float64 `json:"output_token_price"`    // 每1000个输出token的价格（美元）
	ImagePrice       float64 `json:"image_price,omitempty"` // 每张图片的价格（美元）
	AudioPrice       float64 `json:"audio_price,omitempty"` // 每分钟音频的价格（美元）
	Currency         string  `json:"currency"`              // 货币单位
}

// ModelLimits 模型使用限制
type ModelLimits struct {
	MaxTokens           int     `json:"max_tokens"`        // 最大token数
	MaxInputTokens      int     `json:"max_input_tokens"`  // 最大输入token数
	MaxOutputTokens     int     `json:"max_output_tokens"` // 最大输出token数
	ContextWindow       int     `json:"context_window"`    // 上下文窗口大小
	MaxTemperature      float64 `json:"max_temperature"`   // 最大温度值
	MinTemperature      float64 `json:"min_temperature"`   // 最小温度值
	MaxTopP             float64 `json:"max_top_p"`         // 最大top-p值
	MaxPresencePenalty  float64 `json:"max_presence_penalty"`
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
	Streaming       bool `json:"streaming"`        // 是否支持流式响应
	FunctionCalling bool `json:"function_calling"` // 是否支持函数调用
	Vision          bool `json:"vision"`           // 是否支持视觉理解
	CodeCompletion  bool `json:"code_completion"`  // 是否支持代码补全
	Multilingual    bool `json:"multilingual"`     // 是否支持多语言
}
