package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	Enabled   bool      `json:"enabled" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LLMResource LLM资源模型
type LLMResource struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // openai, anthropic, google, etc.
	BaseURL     string    `json:"base_url"`
	APIKey      string    `json:"api_key"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled" gorm:"default:true"`
	Weight      int       `json:"weight" gorm:"default:1"`
	Timeout     int       `json:"timeout" gorm:"default:30"`
	MaxRetries  int       `json:"max_retries" gorm:"default:3"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Endpoint 端点模型
type Endpoint struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	LLMResourceID uint      `json:"llm_resource_id" gorm:"not null"`
	Path          string    `json:"path" gorm:"not null"`
	Method        string    `json:"method" gorm:"default:POST"`
	Enabled       bool      `json:"enabled" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	LLMResource LLMResource `json:"llm_resource" gorm:"foreignKey:LLMResourceID"`
}

// Request 请求模型
type Request struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null"`
	LLMResourceID uint      `json:"llm_resource_id" gorm:"not null"`
	EndpointID    uint      `json:"endpoint_id" gorm:"not null"`
	Path          string    `json:"path"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`   // pending, success, failed, timeout
	Duration      int64     `json:"duration"` // milliseconds
	RequestAt     time.Time `json:"request_at"`
	CreatedAt     time.Time `json:"created_at"`

	User        User        `json:"user" gorm:"foreignKey:UserID"`
	LLMResource LLMResource `json:"llm_resource" gorm:"foreignKey:LLMResourceID"`
	Endpoint    Endpoint    `json:"endpoint" gorm:"foreignKey:EndpointID"`
}

// Response 响应模型
type Response struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RequestID uint      `json:"request_id" gorm:"not null"`
	Status    int       `json:"status"`
	Body      string    `json:"body" gorm:"type:text"`
	Headers   string    `json:"headers" gorm:"type:text"`
	Error     string    `json:"error" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`

	Request Request `json:"request" gorm:"foreignKey:RequestID"`
}

// Log 日志模型
type Log struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Level     string    `json:"level" gorm:"not null"`
	Message   string    `json:"message" gorm:"type:text"`
	Context   string    `json:"context" gorm:"type:text"`
	UserID    *uint     `json:"user_id"`
	RequestID *uint     `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`

	User    *User    `json:"user" gorm:"foreignKey:UserID"`
	Request *Request `json:"request" gorm:"foreignKey:RequestID"`
}
