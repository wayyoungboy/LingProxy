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
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	APIKey    string    `json:"api_key"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LLMResource LLM资源模型
type LLMResource struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Model     string    `json:"model"`
	BaseURL   string    `json:"base_url"`
	APIKey    string    `json:"api_key"`
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
