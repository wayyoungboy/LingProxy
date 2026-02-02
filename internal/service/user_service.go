package service

import (
	"errors"
	"time"

	"github.com/lingproxy/lingproxy/internal/pkg/password"
	"github.com/lingproxy/lingproxy/internal/storage"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserExists        = errors.New("user already exists")
	ErrInvalidAPIKey     = errors.New("invalid API key")
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
	// 检查用户名是否已存在
	existingUsers, err := s.storage.ListUsers()
	if err == nil {
		for _, u := range existingUsers {
			if u.Username == user.Username {
				return ErrUserExists
			}
		}
	}

	// 如果提供了密码，加密密码
	if user.Password != "" {
		hash, err := password.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.PasswordHash = hash
		user.Password = "" // 清除明文密码
	}

	// 如果没有API Key，生成一个
	if user.APIKey == "" {
		user.APIKey = password.GenerateAPIKey()
	}

	// 设置默认值
	if user.Role == "" {
		user.Role = "user"
	}
	if user.Status == "" {
		user.Status = "active"
	}

	return s.storage.CreateUser(user)
}

// GetUser 获取用户
func (s *UserService) GetUser(id string) (*storage.User, error) {
	return s.storage.GetUser(id)
}

// GetUserByUsername 通过用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*storage.User, error) {
	users, err := s.storage.ListUsers()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

// GetUserByAPIKey 通过API Key获取用户
func (s *UserService) GetUserByAPIKey(apiKey string) (*storage.User, error) {
	user, err := s.storage.GetUserByAPIKey(apiKey)
	if err != nil {
		return nil, ErrInvalidAPIKey
	}
	return user, nil
}

// Authenticate 验证用户名和密码
func (s *UserService) Authenticate(username, pwd string) (*storage.User, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, ErrInvalidPassword // 不暴露用户是否存在
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("user account is not active")
	}

	// 验证密码
	if user.PasswordHash == "" {
		return nil, ErrInvalidPassword
	}

	valid, err := password.VerifyPassword(pwd, user.PasswordHash)
	if err != nil || !valid {
		return nil, ErrInvalidPassword
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	s.storage.UpdateUser(user)

	return user, nil
}

// ResetAPIKey 重置用户的API Key
func (s *UserService) ResetAPIKey(userID string) (string, error) {
	user, err := s.storage.GetUser(userID)
	if err != nil {
		return "", ErrUserNotFound
	}

	newAPIKey := password.GenerateAPIKey()
	user.APIKey = newAPIKey

	if err := s.storage.UpdateUser(user); err != nil {
		return "", err
	}

	return newAPIKey, nil
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(userID, oldPassword, newPassword string) error {
	user, err := s.storage.GetUser(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// 验证旧密码
	if user.PasswordHash != "" {
		valid, err := password.VerifyPassword(oldPassword, user.PasswordHash)
		if err != nil || !valid {
			return ErrInvalidPassword
		}
	}

	// 加密新密码
	hash, err := password.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hash
	return s.storage.UpdateUser(user)
}

// ListUsers 获取所有用户
func (s *UserService) ListUsers() ([]*storage.User, error) {
	return s.storage.ListUsers()
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *storage.User) error {
	// 如果提供了新密码，加密它
	if user.Password != "" {
		hash, err := password.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.PasswordHash = hash
		user.Password = ""
	}
	return s.storage.UpdateUser(user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id string) error {
	return s.storage.DeleteUser(id)
}