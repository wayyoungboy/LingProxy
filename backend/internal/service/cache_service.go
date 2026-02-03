package service

import (
	"github.com/lingproxy/lingproxy/internal/cache"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// CachedUserService 带缓存的用户服务
type CachedUserService struct {
	storage *storage.StorageFacade
	cache   *cache.CacheService
}

// NewCachedUserService 创建新的带缓存的用户服务
func NewCachedUserService(storage *storage.StorageFacade, cache *cache.CacheService) *CachedUserService {
	return &CachedUserService{
		storage: storage,
		cache:   cache,
	}
}

// CreateUser 创建用户
func (s *CachedUserService) CreateUser(user *storage.User) error {
	err := s.storage.CreateUser(user)
	if err == nil {
		s.cache.CacheUser(user)
	}
	return err
}

// GetUser 获取用户
func (s *CachedUserService) GetUser(id string) (*storage.User, error) {
	if user, found := s.cache.GetUser(id); found {
		return user, nil
	}
	
	user, err := s.storage.GetUser(id)
	if err == nil {
		s.cache.CacheUser(user)
	}
	return user, err
}

// GetUserByAPIKey 通过API Key获取用户
func (s *CachedUserService) GetUserByAPIKey(apiKey string) (*storage.User, error) {
	if user, found := s.cache.GetUserByAPIKey(apiKey); found {
		return user, nil
	}
	
	user, err := s.storage.GetUserByAPIKey(apiKey)
	if err == nil {
		s.cache.CacheUser(user)
	}
	return user, err
}

// ListUsers 获取所有用户
func (s *CachedUserService) ListUsers() ([]*storage.User, error) {
	return s.storage.ListUsers()
}

// UpdateUser 更新用户
func (s *CachedUserService) UpdateUser(user *storage.User) error {
	err := s.storage.UpdateUser(user)
	if err == nil {
		s.cache.CacheUser(user)
		s.cache.InvalidateUserCache(user.ID)
	}
	return err
}

// DeleteUser 删除用户
func (s *CachedUserService) DeleteUser(id string) error {
	err := s.storage.DeleteUser(id)
	if err == nil {
		s.cache.InvalidateUserCache(id)
	}
	return err
}