package storage

import (
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
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`     // 模型类别: chat(对话), image(生图), embedding(嵌入), rerank(重排序), audio(语音), video(视频)
	Provider  string    `json:"provider"` // 服务提供商: openai, zai, anthropic, google, azure, custom等
	Model     string    `json:"model"`    // 模型标识：此资源对应的模型标识（如gpt-4, gpt-3.5-turbo），与base_url和api_key共同构成可调用资源
	BaseURL   string    `json:"base_url"` // API基础URL：此资源的API端点地址
	APIKey    string    `json:"api_key"`  // API密钥：此资源的认证密钥
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Endpoint  string    `json:"endpoint"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	Duration  int64     `json:"duration"`
	Tokens    int       `json:"tokens"`
	CreatedAt time.Time `json:"created_at"`
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

// Token Token模型（用于API请求认证）
type Token struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`                // Token名称/描述
	Token      string     `json:"token"`               // Token值（API Key）
	Prefix     string     `json:"prefix"`              // Token前缀（用于显示）
	Status     string     `json:"status"`              // active/inactive
	PolicyID   string     `json:"policy_id,omitempty"` // 关联的策略ID
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

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
