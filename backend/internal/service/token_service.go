package service

import (
	"errors"
	"strings"
	"time"

	"github.com/lingproxy/lingproxy/internal/pkg/password"
	"github.com/lingproxy/lingproxy/internal/storage"
)

var (
	ErrTokenNotFound   = errors.New("token not found")
	ErrTokenExpired    = errors.New("token has expired")
	ErrTokenInactive   = errors.New("token is inactive")
	ErrTokenNameExists = errors.New("token name already exists")
)

// TokenService Token服务
type TokenService struct {
	storage *storage.StorageFacade
}

// NewTokenService 创建新的Token服务
func NewTokenService(storage *storage.StorageFacade) *TokenService {
	return &TokenService{
		storage: storage,
	}
}

// CreateToken 创建Token
func (s *TokenService) CreateToken(name string, expiresAt *time.Time) (*storage.Token, error) {
	// 检查名称是否已存在
	tokens, err := s.storage.ListTokens()
	if err == nil {
		for _, t := range tokens {
			if t.Name == name {
				return nil, ErrTokenNameExists
			}
		}
	}

	// 生成Token值（默认生成ling-开头的Token）
	// password.GenerateAPIKey()生成的是"ling_"开头，我们需要改为"ling-"
	baseKey := password.GenerateAPIKey()
	tokenValue := strings.Replace(baseKey, "ling_", "ling-", 1)

	// 提取前缀（前12个字符）
	prefix := tokenValue
	if len(tokenValue) > 12 {
		prefix = tokenValue[:12] + "..."
	}

	token := &storage.Token{
		Name:      name,
		Token:     tokenValue,
		Prefix:    prefix,
		Status:    "active",
		ExpiresAt: expiresAt,
	}

	if err := s.storage.CreateToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

// GetToken 获取Token
func (s *TokenService) GetToken(id string) (*storage.Token, error) {
	return s.storage.GetToken(id)
}

// ListTokens 获取Token列表
func (s *TokenService) ListTokens() ([]*storage.Token, error) {
	return s.storage.ListTokens()
}

// UpdateToken 更新Token
func (s *TokenService) UpdateToken(id string, name *string, status *string) (*storage.Token, error) {
	token, err := s.storage.GetToken(id)
	if err != nil {
		return nil, ErrTokenNotFound
	}

	if name != nil {
		// 检查名称是否与其他Token冲突
		tokens, err := s.storage.ListTokens()
		if err == nil {
			for _, t := range tokens {
				if t.ID != id && t.Name == *name {
					return nil, ErrTokenNameExists
				}
			}
		}
		token.Name = *name
	}

	if status != nil {
		token.Status = *status
	}

	if err := s.storage.UpdateToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

// DeleteToken 删除Token
func (s *TokenService) DeleteToken(id string) error {
	_, err := s.storage.GetToken(id)
	if err != nil {
		return ErrTokenNotFound
	}
	return s.storage.DeleteToken(id)
}

// ResetToken 重置Token（生成新的Token值）
func (s *TokenService) ResetToken(id string) (*storage.Token, error) {
	token, err := s.storage.GetToken(id)
	if err != nil {
		return nil, ErrTokenNotFound
	}

	// 生成新的Token值（保持原有前缀）
	oldPrefix := ""
	if strings.HasPrefix(token.Token, "ling-") {
		oldPrefix = "ling-"
	} else if strings.HasPrefix(token.Token, "ling_") {
		oldPrefix = "ling-"
	}

	baseKey := password.GenerateAPIKey()
	newTokenValue := strings.Replace(baseKey, "ling_", oldPrefix, 1)
	if oldPrefix == "" {
		newTokenValue = baseKey
	}

	// 提取前缀
	prefix := newTokenValue
	if len(newTokenValue) > 12 {
		prefix = newTokenValue[:12] + "..."
	}

	token.Token = newTokenValue
	token.Prefix = prefix

	if err := s.storage.UpdateToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

// ValidateToken 验证Token（用于认证中间件）
func (s *TokenService) ValidateToken(tokenValue string) (*storage.Token, error) {
	token, err := s.storage.GetTokenByValue(tokenValue)
	if err != nil {
		// 如果Token不存在，尝试通过User的APIKey验证（向后兼容）
		user, err := s.storage.GetUserByAPIKey(tokenValue)
		if err != nil {
			return nil, ErrTokenNotFound
		}
		// 将User转换为Token格式（用于兼容）
		return &storage.Token{
			ID:     user.ID,
			Name:   user.Username,
			Token:  user.APIKey,
			Prefix: getPrefix(user.APIKey),
			Status: user.Status,
		}, nil
	}

	// 检查Token状态
	if token.Status != "active" {
		return nil, ErrTokenInactive
	}

	// 检查过期时间
	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// 更新最后使用时间
	now := time.Now()
	token.LastUsedAt = &now
	s.storage.UpdateToken(token)

	return token, nil
}

// UpdateTokenPolicy 更新Token的策略
func (s *TokenService) UpdateTokenPolicy(tokenID, policyID string) (*storage.Token, error) {
	token, err := s.storage.GetToken(tokenID)
	if err != nil {
		return nil, ErrTokenNotFound
	}

	token.PolicyID = policyID
	if err := s.storage.UpdateToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

// RemoveTokenPolicy 移除Token的策略
func (s *TokenService) RemoveTokenPolicy(tokenID string) (*storage.Token, error) {
	token, err := s.storage.GetToken(tokenID)
	if err != nil {
		return nil, ErrTokenNotFound
	}

	token.PolicyID = ""
	if err := s.storage.UpdateToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

// getPrefix 获取Token前缀
func getPrefix(tokenValue string) string {
	if len(tokenValue) > 12 {
		return tokenValue[:12] + "..."
	}
	return tokenValue
}
