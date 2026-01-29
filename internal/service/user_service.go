package service

import (
	"github.com/lingproxy/lingproxy/internal/storage"
)

// UserService 用户服务
type UserService struct {
	storage *storage.StorageFacade
}

// NewUserService 创建新的用户服务
func NewUserService(storage *storage.StorageFacade) *UserService {
	return &UserService{
		storage: storage,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *storage.User) error {
	return s.storage.CreateUser(user)
}

// GetUser 获取用户
func (s *UserService) GetUser(id string) (*storage.User, error) {
	return s.storage.GetUser(id)
}

// GetUserByAPIKey 通过API Key获取用户
func (s *UserService) GetUserByAPIKey(apiKey string) (*storage.User, error) {
	return s.storage.GetUserByAPIKey(apiKey)
}

// ListUsers 获取所有用户
func (s *UserService) ListUsers() ([]*storage.User, error) {
	return s.storage.ListUsers()
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *storage.User) error {
	return s.storage.UpdateUser(user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id string) error {
	return s.storage.DeleteUser(id)
}