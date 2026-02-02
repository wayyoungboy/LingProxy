package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lingproxy/lingproxy/internal/storage"
)

var (
	ErrPolicyDisabled   = errors.New("policy is disabled")
	ErrTemplateNotFound = errors.New("policy template not found")
)

// PolicyService 策略服务
type PolicyService struct {
	storage         *storage.StorageFacade
	templateService *PolicyTemplateService
	executorFactory *PolicyExecutorFactory
}

// NewPolicyService 创建新的策略服务
func NewPolicyService(storage *storage.StorageFacade) *PolicyService {
	return &PolicyService{
		storage:         storage,
		templateService: NewPolicyTemplateService(storage),
		executorFactory: NewPolicyExecutorFactory(),
	}
}

// CreatePolicy 创建策略
func (s *PolicyService) CreatePolicy(name, templateID string, parameters map[string]interface{}) (*storage.Policy, error) {
	// 获取模板
	template, err := s.storage.GetPolicyTemplate(templateID)
	if err != nil {
		return nil, ErrTemplateNotFound
	}

	// 验证参数
	if err := s.templateService.ValidateParameters(template, parameters); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// 序列化参数
	paramsJSON, err := json.Marshal(parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	policy := &storage.Policy{
		Name:       name,
		TemplateID: templateID,
		Type:       template.Type,
		Parameters: string(paramsJSON),
		Enabled:    true,
	}

	if err := s.storage.CreatePolicy(policy); err != nil {
		return nil, err
	}

	return policy, nil
}

// GetPolicy 获取策略
func (s *PolicyService) GetPolicy(id string) (*storage.Policy, error) {
	return s.storage.GetPolicy(id)
}

// ListPolicies 获取策略列表
func (s *PolicyService) ListPolicies() ([]*storage.Policy, error) {
	return s.storage.ListPolicies()
}

// UpdatePolicy 更新策略
func (s *PolicyService) UpdatePolicy(id string, name *string, parameters map[string]interface{}, enabled *bool) (*storage.Policy, error) {
	policy, err := s.storage.GetPolicy(id)
	if err != nil {
		return nil, ErrPolicyNotFound
	}

	if name != nil {
		policy.Name = *name
	}

	if parameters != nil {
		// 获取模板验证参数
		template, err := s.storage.GetPolicyTemplate(policy.TemplateID)
		if err != nil {
			return nil, ErrTemplateNotFound
		}

		if err := s.templateService.ValidateParameters(template, parameters); err != nil {
			return nil, fmt.Errorf("invalid parameters: %w", err)
		}

		paramsJSON, err := json.Marshal(parameters)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parameters: %w", err)
		}
		policy.Parameters = string(paramsJSON)
	}

	if enabled != nil {
		policy.Enabled = *enabled
	}

	if err := s.storage.UpdatePolicy(policy); err != nil {
		return nil, err
	}

	return policy, nil
}

// DeletePolicy 删除策略
func (s *PolicyService) DeletePolicy(id string) error {
	_, err := s.storage.GetPolicy(id)
	if err != nil {
		return ErrPolicyNotFound
	}
	return s.storage.DeletePolicy(id)
}

// ExecutePolicy 执行策略，选择资源
func (s *PolicyService) ExecutePolicy(policyID, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	// 获取策略
	policy, err := s.storage.GetPolicy(policyID)
	if err != nil {
		return nil, ErrPolicyNotFound
	}

	// 检查策略是否启用
	if !policy.Enabled {
		return nil, ErrPolicyDisabled
	}

	// 获取执行器
	executor := s.executorFactory.GetExecutor(policy.Type)

	// 执行策略
	return executor.Execute(policy, modelName, resources)
}

// GetDefaultPolicyExecutor 获取默认策略执行器（用于没有配置策略的情况）
func (s *PolicyService) GetDefaultPolicyExecutor() PolicyExecutor {
	return NewRoundRobinPolicyExecutor()
}
