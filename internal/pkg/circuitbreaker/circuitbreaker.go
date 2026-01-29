package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State 熔断器状态
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Config 熔断器配置
type Config struct {
	Timeout               time.Duration
	MaxConcurrentRequests int
	ErrorThreshold        int
	RequestVolume         int
	SleepWindow           time.Duration
}

// CircuitBreaker 熔断器接口
type CircuitBreaker interface {
	Execute(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error)
	GetState() State
	GetName() string
}

// DefaultCircuitBreaker 默认熔断器实现
type DefaultCircuitBreaker struct {
	name               string
	maxFailures        int
	resetTimeout       time.Duration
	state              State
	failureCount       int
	successCount       int
	lastFailureTime    time.Time
	halfOpenMaxRequests int
	halfOpenRequests   int
	mutex              sync.RWMutex
}

// NewCircuitBreaker 创建新的熔断器
func NewCircuitBreaker(config *Config) *DefaultCircuitBreaker {
	return &DefaultCircuitBreaker{
		name:               "default",
		maxFailures:        config.ErrorThreshold,
		resetTimeout:       config.SleepWindow,
		state:              StateClosed,
		halfOpenMaxRequests: 3,
	}
}

// Execute 执行函数
func (cb *DefaultCircuitBreaker) Execute(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	cb.mutex.Lock()
	
	// 检查熔断器状态
	switch cb.state {
	case StateOpen:
		// 检查是否可以进入半开状态
		if time.Since(cb.lastFailureTime) >= cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.halfOpenRequests = 0
		} else {
			cb.mutex.Unlock()
			return nil, errors.New("circuit breaker is open")
		}
	case StateHalfOpen:
		if cb.halfOpenRequests >= cb.halfOpenMaxRequests {
			cb.mutex.Unlock()
			return nil, errors.New("circuit breaker is half-open and max requests reached")
		}
		cb.halfOpenRequests++
	}
	
	cb.mutex.Unlock()

	// 执行函数
	result, err := fn(ctx)
	
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		
		if cb.state == StateHalfOpen {
			cb.state = StateOpen
		} else if cb.failureCount >= cb.maxFailures {
			cb.state = StateOpen
		}
		
		return nil, err
	}
	
	// 成功
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		cb.failureCount = 0
	}
	
	return result, nil
}

// GetState 获取当前状态
func (cb *DefaultCircuitBreaker) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetName 获取熔断器名称
func (cb *DefaultCircuitBreaker) GetName() string {
	return cb.name
}