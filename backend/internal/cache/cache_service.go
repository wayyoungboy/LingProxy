package cache

import (
	"time"

	"github.com/lingproxy/lingproxy/internal/storage"
)

// CacheService 缓存服务
type CacheService struct {
	userCache        *MemoryCache
	llmResourceCache *MemoryCache
	endpointCache    *MemoryCache
	statsCache       *MemoryCache
}

// NewCacheService 创建新的缓存服务
func NewCacheService() *CacheService {
	return &CacheService{
		userCache:        NewMemoryCache(5 * time.Minute),
		llmResourceCache: NewMemoryCache(10 * time.Minute),
		endpointCache:    NewMemoryCache(5 * time.Minute),
		statsCache:       NewMemoryCache(1 * time.Minute),
	}
}

// CacheUser 缓存用户信息
func (s *CacheService) CacheUser(user *storage.User) {
	s.userCache.Set(user.ID, user)
	s.userCache.Set("api_key_"+user.APIKey, user)
}

// GetUser 获取缓存的用户
func (s *CacheService) GetUser(id string) (*storage.User, bool) {
	if val, found := s.userCache.Get(id); found {
		if user, ok := val.(*storage.User); ok {
			return user, true
		}
	}
	return nil, false
}

// GetUserByAPIKey 通过API Key获取缓存的用户
func (s *CacheService) GetUserByAPIKey(apiKey string) (*storage.User, bool) {
	if val, found := s.userCache.Get("api_key_" + apiKey); found {
		if user, ok := val.(*storage.User); ok {
			return user, true
		}
	}
	return nil, false
}

// CacheLLMResource 缓存LLM资源信息
func (s *CacheService) CacheLLMResource(resource *storage.LLMResource) {
	s.llmResourceCache.Set(resource.ID, resource)
}

// GetLLMResource 获取缓存的LLM资源
func (s *CacheService) GetLLMResource(id string) (*storage.LLMResource, bool) {
	if val, found := s.llmResourceCache.Get(id); found {
		if resource, ok := val.(*storage.LLMResource); ok {
			return resource, true
		}
	}
	return nil, false
}

// CacheEndpoint 缓存端点信息
func (s *CacheService) CacheEndpoint(endpoint *storage.Endpoint) {
	s.endpointCache.Set(endpoint.ID, endpoint)
}

// GetEndpoint 获取缓存的端点
func (s *CacheService) GetEndpoint(id string) (*storage.Endpoint, bool) {
	if val, found := s.endpointCache.Get(id); found {
		if endpoint, ok := val.(*storage.Endpoint); ok {
			return endpoint, true
		}
	}
	return nil, false
}

// CacheStats 缓存统计信息
func (s *CacheService) CacheStats(key string, stats interface{}) {
	s.statsCache.Set(key, stats)
}

// GetStats 获取缓存的统计信息
func (s *CacheService) GetStats(key string) (interface{}, bool) {
	return s.statsCache.Get(key)
}

// InvalidateUserCache 使用户缓存失效
func (s *CacheService) InvalidateUserCache(userID string) {
	s.userCache.Delete(userID)
}

// InvalidateLLMResourceCache 使LLM资源缓存失效
func (s *CacheService) InvalidateLLMResourceCache(llmResourceID string) {
	s.llmResourceCache.Delete(llmResourceID)
}

// InvalidateEndpointCache 使端点缓存失效
func (s *CacheService) InvalidateEndpointCache(endpointID string) {
	s.endpointCache.Delete(endpointID)
}

// GetCacheStats 获取缓存统计
func (s *CacheService) GetCacheStats() map[string]int {
	return map[string]int{
		"users":         s.userCache.Size(),
		"llm_resources": s.llmResourceCache.Size(),
		"endpoints":     s.endpointCache.Size(),
		"stats":         s.statsCache.Size(),
	}
}
