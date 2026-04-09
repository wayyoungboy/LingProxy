package middleware

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket 令牌桶
type TokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

// NewTokenBucket 创建新的令牌桶
func NewTokenBucket(maxTokens float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow 检查是否允许请求
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens = math.Min(b.maxTokens, b.tokens+elapsed*b.refillRate)
	b.lastRefill = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// Remaining 返回剩余令牌数
func (b *TokenBucket) Remaining() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	tokens := math.Min(b.maxTokens, b.tokens+elapsed*b.refillRate)
	return int(tokens)
}

// Reset 重置令牌桶
func (b *TokenBucket) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.tokens = b.maxTokens
	b.lastRefill = time.Now()
}

// RateLimiter 全局限流器
type RateLimiter struct {
	mu          sync.RWMutex
	buckets     map[string]*TokenBucket
	maxTokens   float64
	refillRate  float64
	enabled     bool
	concurrency int
	semaphore   chan struct{}
	cleanupTick time.Duration
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(maxTokens float64, refillRate float64, enabled bool, concurrency int) *RateLimiter {
	rl := &RateLimiter{
		buckets:     make(map[string]*TokenBucket),
		maxTokens:   maxTokens,
		refillRate:  refillRate,
		enabled:     enabled,
		concurrency: concurrency,
		semaphore:   make(chan struct{}, concurrency),
		cleanupTick: 10 * time.Minute,
	}
	go rl.cleanup()
	return rl
}

// GetBucket 获取或创建客户端的令牌桶
func (rl *RateLimiter) GetBucket(clientID string) *TokenBucket {
	rl.mu.RLock()
	bucket, exists := rl.buckets[clientID]
	rl.mu.RUnlock()

	if exists {
		return bucket
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 双重检查
	if bucket, exists = rl.buckets[clientID]; exists {
		return bucket
	}

	bucket = NewTokenBucket(rl.maxTokens, rl.refillRate)
	rl.buckets[clientID] = bucket
	return bucket
}

// UpdateConfig 动态更新配置
func (rl *RateLimiter) UpdateConfig(maxTokens float64, refillRate float64, enabled bool) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.maxTokens = maxTokens
	rl.refillRate = refillRate
	rl.enabled = enabled

	// 更新所有现有桶的配置
	for _, bucket := range rl.buckets {
		bucket.mu.Lock()
		bucket.maxTokens = maxTokens
		bucket.refillRate = refillRate
		if bucket.tokens > maxTokens {
			bucket.tokens = maxTokens
		}
		bucket.mu.Unlock()
	}
}

// UpdateConcurrency 动态更新并发限制
func (rl *RateLimiter) UpdateConcurrency(concurrency int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.concurrency = concurrency
	// 替换信号量通道
	rl.semaphore = make(chan struct{}, concurrency)
}

// GetStats 获取限流器当前状态
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"enabled":        rl.enabled,
		"rpm_limit":      int(rl.maxTokens),
		"refill_rate":    rl.refillRate,
		"concurrency":    rl.concurrency,
		"active_clients": len(rl.buckets),
	}
}

// IsEnabled 检查限流是否启用
func (rl *RateLimiter) IsEnabled() bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.enabled
}

// cleanup 定期清理空闲的令牌桶
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupTick)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for id, bucket := range rl.buckets {
			bucket.mu.Lock()
			if time.Since(bucket.lastRefill) > 30*time.Minute {
				delete(rl.buckets, id)
			}
			bucket.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// GlobalRateLimiter 全局限流器实例
var GlobalRateLimiter *RateLimiter

// InitGlobalRateLimiter 初始化全局限流器
func InitGlobalRateLimiter(maxTokens float64, refillRate float64, enabled bool, concurrency int) {
	GlobalRateLimiter = NewRateLimiter(maxTokens, refillRate, enabled, concurrency)
}

// extractClientID 提取客户端标识
func extractClientID(c *gin.Context) string {
	// 优先使用 API Key（如果有的话）
	if apiKey := c.GetHeader("Authorization"); apiKey != "" {
		return "apikey:" + apiKey
	}

	// 使用客户端 IP
	if ip := c.ClientIP(); ip != "" {
		return "ip:" + ip
	}

	// 使用 X-Forwarded-For
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return "xff:" + xff
	}

	return "ip:unknown"
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GlobalRateLimiter == nil || !GlobalRateLimiter.IsEnabled() {
			c.Next()
			return
		}

		clientID := extractClientID(c)
		bucket := GlobalRateLimiter.GetBucket(clientID)

		if !bucket.Allow() {
			remaining := bucket.Remaining()
			c.Header("X-RateLimit-Limit", strconv.Itoa(int(GlobalRateLimiter.maxTokens)))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"message": "Rate limit exceeded. Please try again later.",
					"type":    "rate_limit_error",
					"code":    429,
				},
			})
			return
		}

		remaining := bucket.Remaining()
		c.Header("X-RateLimit-Limit", strconv.Itoa(int(GlobalRateLimiter.maxTokens)))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))

		c.Next()
	}
}

// ConcurrencyLimit 并发限制中间件
func ConcurrencyLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GlobalRateLimiter == nil || !GlobalRateLimiter.IsEnabled() {
			c.Next()
			return
		}

		rl := GlobalRateLimiter
		rl.mu.RLock()
		semaphore := rl.semaphore
		rl.mu.RUnlock()

		if semaphore == nil {
			c.Next()
			return
		}

		select {
		case semaphore <- struct{}{}:
			defer func() { <-semaphore }()
			c.Next()
		default:
			c.Header("X-Concurrency-Limit", strconv.Itoa(rl.concurrency))
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"message": "Server is at capacity. Please try again later.",
					"type":    "rate_limit_error",
					"code":    429,
				},
			})
		}
	}
}

// GetClientIP 获取客户端 IP 地址（辅助函数）
func GetClientIP(c *gin.Context) string {
	// 优先使用 X-Real-IP
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}

	// 使用 gin 的 ClientIP
	return c.ClientIP()
}
