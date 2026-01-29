package service

import (
	"github.com/lingproxy/lingproxy/internal/storage"
)

// PolicyService 策略服务
type PolicyService struct {
	storage *storage.StorageFacade
}

// NewPolicyService 创建新的策略服务
func NewPolicyService(storage *storage.StorageFacade) *PolicyService {
	return &PolicyService{
		storage: storage,
	}
}

// CreateLLMResource 创建LLM资源
func (s *PolicyService) CreateLLMResource(resource *storage.LLMResource) error {
	return s.storage.CreateLLMResource(resource)
}

// GetLLMResource 获取LLM资源
func (s *PolicyService) GetLLMResource(id string) (*storage.LLMResource, error) {
	return s.storage.GetLLMResource(id)
}

// ListLLMResources 获取所有LLM资源
func (s *PolicyService) ListLLMResources() ([]*storage.LLMResource, error) {
	return s.storage.ListLLMResources()
}

// CreateEndpoint 创建端点
func (s *PolicyService) CreateEndpoint(endpoint *storage.Endpoint) error {
	return s.storage.CreateEndpoint(endpoint)
}

// GetEndpoint 获取端点
func (s *PolicyService) GetEndpoint(id string) (*storage.Endpoint, error) {
	return s.storage.GetEndpoint(id)
}

// ListEndpoints 获取所有端点
func (s *PolicyService) ListEndpoints() ([]*storage.Endpoint, error) {
	return s.storage.ListEndpoints()
}