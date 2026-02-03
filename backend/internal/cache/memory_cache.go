package cache

import (
	"sync"
	"time"
)

// CacheItem 缓存项
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// MemoryCache 内存缓存
type MemoryCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewMemoryCache 创建新的内存缓存
func NewMemoryCache(ttl time.Duration) *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
		ttl:   ttl,
	}
	
	// 启动清理goroutine
	go cache.cleanup()
	
	return cache
}

// Set 设置缓存
func (c *MemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(c.ttl),
	}
}

// Get 获取缓存
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	
	if time.Now().After(item.Expiration) {
		delete(c.items, key)
		return nil, false
	}
	
	return item.Value, true
}

// Delete 删除缓存
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.items, key)
}

// Clear 清空缓存
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items = make(map[string]CacheItem)
}

// Size 获取缓存大小
func (c *MemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return len(c.items)
}

// cleanup 定期清理过期缓存
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.Expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}