package service

import (
	"errors"
	"strings"
	"time"

	"github.com/lingproxy/lingproxy/internal/pkg/password"
	"github.com/lingproxy/lingproxy/internal/storage"
)

var (
	ErrAPIKeyNotFound   = errors.New("API key not found")
	ErrAPIKeyExpired    = errors.New("API key has expired")
	ErrAPIKeyInactive   = errors.New("API key is inactive")
	ErrAPIKeyNameExists = errors.New("API key name already exists")
)

// 保持向后兼容的错误变量
var (
	ErrTokenNotFound   = ErrAPIKeyNotFound
	ErrTokenExpired    = ErrAPIKeyExpired
	ErrTokenInactive   = ErrAPIKeyInactive
	ErrTokenNameExists = ErrAPIKeyNameExists
)

// APIKeyService API Key服务
type APIKeyService struct {
	storage *storage.StorageFacade
}

// NewAPIKeyService 创建新的API Key服务
func NewAPIKeyService(storage *storage.StorageFacade) *APIKeyService {
	return &APIKeyService{
		storage: storage,
	}
}

// TokenService 保持向后兼容的类型别名
// Deprecated: 使用 APIKeyService 代替
type TokenService = APIKeyService

// NewTokenService 保持向后兼容的函数别名
// Deprecated: 使用 NewAPIKeyService 代替
func NewTokenService(storage *storage.StorageFacade) *APIKeyService {
	return NewAPIKeyService(storage)
}

// CreateAPIKey 创建API Key
func (s *APIKeyService) CreateAPIKey(name string, expiresAt *time.Time) (*storage.APIKey, error) {
	// 检查名称是否已存在
	apiKeys, err := s.storage.ListAPIKeys()
	if err == nil {
		for _, k := range apiKeys {
			if k.Name == name {
				return nil, ErrAPIKeyNameExists
			}
		}
	}

	// 生成API Key值（默认生成ling-开头的API Key）
	// password.GenerateAPIKey()生成的是"ling_"开头，我们需要改为"ling-"
	baseKey := password.GenerateAPIKey()
	apiKeyValue := strings.Replace(baseKey, "ling_", "ling-", 1)

	// 提取前缀（前12个字符）
	prefix := apiKeyValue
	if len(apiKeyValue) > 12 {
		prefix = apiKeyValue[:12] + "..."
	}

	apiKey := &storage.APIKey{
		Name:      name,
		APIKey:    apiKeyValue,
		Prefix:    prefix,
		Status:    "active",
		ExpiresAt: expiresAt,
	}

	if err := s.storage.CreateAPIKey(apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// CreateToken 保持向后兼容的函数别名
// Deprecated: 使用 CreateAPIKey 代替
func (s *APIKeyService) CreateToken(name string, expiresAt *time.Time) (*storage.APIKey, error) {
	return s.CreateAPIKey(name, expiresAt)
}

// GetAPIKey 获取API Key
func (s *APIKeyService) GetAPIKey(id string) (*storage.APIKey, error) {
	return s.storage.GetAPIKey(id)
}

// GetToken 保持向后兼容的函数别名
// Deprecated: 使用 GetAPIKey 代替
func (s *APIKeyService) GetToken(id string) (*storage.APIKey, error) {
	return s.GetAPIKey(id)
}

// ListAPIKeys 获取API Key列表
func (s *APIKeyService) ListAPIKeys() ([]*storage.APIKey, error) {
	return s.storage.ListAPIKeys()
}

// ListTokens 保持向后兼容的函数别名
// Deprecated: 使用 ListAPIKeys 代替
func (s *APIKeyService) ListTokens() ([]*storage.APIKey, error) {
	return s.ListAPIKeys()
}

// UpdateAPIKey 更新API Key
func (s *APIKeyService) UpdateAPIKey(id string, name *string, status *string) (*storage.APIKey, error) {
	apiKey, err := s.storage.GetAPIKey(id)
	if err != nil {
		return nil, ErrAPIKeyNotFound
	}

	if name != nil {
		// 检查名称是否与其他API Key冲突
		apiKeys, err := s.storage.ListAPIKeys()
		if err == nil {
			for _, k := range apiKeys {
				if k.ID != id && k.Name == *name {
					return nil, ErrAPIKeyNameExists
				}
			}
		}
		apiKey.Name = *name
	}

	if status != nil {
		apiKey.Status = *status
	}

	if err := s.storage.UpdateAPIKey(apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// UpdateToken 保持向后兼容的函数别名
// Deprecated: 使用 UpdateAPIKey 代替
func (s *APIKeyService) UpdateToken(id string, name *string, status *string) (*storage.APIKey, error) {
	return s.UpdateAPIKey(id, name, status)
}

// DeleteAPIKey 删除API Key
func (s *APIKeyService) DeleteAPIKey(id string) error {
	_, err := s.storage.GetAPIKey(id)
	if err != nil {
		return ErrAPIKeyNotFound
	}
	return s.storage.DeleteAPIKey(id)
}

// DeleteToken 保持向后兼容的函数别名
// Deprecated: 使用 DeleteAPIKey 代替
func (s *APIKeyService) DeleteToken(id string) error {
	return s.DeleteAPIKey(id)
}

// ResetAPIKey 重置API Key（生成新的API Key值）
func (s *APIKeyService) ResetAPIKey(id string) (*storage.APIKey, error) {
	apiKey, err := s.storage.GetAPIKey(id)
	if err != nil {
		return nil, ErrAPIKeyNotFound
	}

	// 生成新的API Key值（保持原有前缀）
	oldPrefix := ""
	if strings.HasPrefix(apiKey.APIKey, "ling-") {
		oldPrefix = "ling-"
	} else if strings.HasPrefix(apiKey.APIKey, "ling_") {
		oldPrefix = "ling-"
	}

	baseKey := password.GenerateAPIKey()
	newAPIKeyValue := strings.Replace(baseKey, "ling_", oldPrefix, 1)
	if oldPrefix == "" {
		newAPIKeyValue = baseKey
	}

	// 提取前缀
	prefix := newAPIKeyValue
	if len(newAPIKeyValue) > 12 {
		prefix = newAPIKeyValue[:12] + "..."
	}

	apiKey.APIKey = newAPIKeyValue
	apiKey.Prefix = prefix

	if err := s.storage.UpdateAPIKey(apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// ResetToken 保持向后兼容的函数别名
// Deprecated: 使用 ResetAPIKey 代替
func (s *APIKeyService) ResetToken(id string) (*storage.APIKey, error) {
	return s.ResetAPIKey(id)
}

// ValidateAPIKey 验证API Key（用于认证中间件）
func (s *APIKeyService) ValidateAPIKey(apiKeyValue string) (*storage.APIKey, error) {
	apiKey, err := s.storage.GetAPIKeyByValue(apiKeyValue)
	if err != nil {
		// 如果API Key不存在，尝试通过User的APIKey验证（向后兼容）
		user, err := s.storage.GetUserByAPIKey(apiKeyValue)
		if err != nil {
			return nil, ErrAPIKeyNotFound
		}
		// 将User转换为API Key格式（用于兼容）
		return &storage.APIKey{
			ID:     user.ID,
			Name:   user.Username,
			APIKey: user.APIKey,
			Prefix: getPrefix(user.APIKey),
			Status: user.Status,
		}, nil
	}

	// 检查API Key状态
	if apiKey.Status != "active" {
		return nil, ErrAPIKeyInactive
	}

	// 检查过期时间
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, ErrAPIKeyExpired
	}

	// 更新最后使用时间
	now := time.Now()
	apiKey.LastUsedAt = &now
	s.storage.UpdateAPIKey(apiKey)

	return apiKey, nil
}

// ValidateToken 保持向后兼容的函数别名
// Deprecated: 使用 ValidateAPIKey 代替
func (s *APIKeyService) ValidateToken(tokenValue string) (*storage.APIKey, error) {
	return s.ValidateAPIKey(tokenValue)
}

// UpdateAPIKeyPolicy 更新API Key的策略
func (s *APIKeyService) UpdateAPIKeyPolicy(apiKeyID, policyID string) (*storage.APIKey, error) {
	apiKey, err := s.storage.GetAPIKey(apiKeyID)
	if err != nil {
		return nil, ErrAPIKeyNotFound
	}

	apiKey.PolicyID = policyID
	if err := s.storage.UpdateAPIKey(apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// UpdateTokenPolicy 保持向后兼容的函数别名
// Deprecated: 使用 UpdateAPIKeyPolicy 代替
func (s *APIKeyService) UpdateTokenPolicy(tokenID, policyID string) (*storage.APIKey, error) {
	return s.UpdateAPIKeyPolicy(tokenID, policyID)
}

// RemoveAPIKeyPolicy 移除API Key的策略
func (s *APIKeyService) RemoveAPIKeyPolicy(apiKeyID string) (*storage.APIKey, error) {
	apiKey, err := s.storage.GetAPIKey(apiKeyID)
	if err != nil {
		return nil, ErrAPIKeyNotFound
	}

	apiKey.PolicyID = ""
	if err := s.storage.UpdateAPIKey(apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// RemoveTokenPolicy 保持向后兼容的函数别名
// Deprecated: 使用 RemoveAPIKeyPolicy 代替
func (s *APIKeyService) RemoveTokenPolicy(tokenID string) (*storage.APIKey, error) {
	return s.RemoveAPIKeyPolicy(tokenID)
}

// getPrefix 获取API Key前缀
func getPrefix(tokenValue string) string {
	if len(tokenValue) > 12 {
		return tokenValue[:12] + "..."
	}
	return tokenValue
}
