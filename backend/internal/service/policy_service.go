package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lingproxy/lingproxy/internal/pkg/logger"
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
		return nil, fmt.Errorf("policy not found: %w", err)
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
	policy, err := s.storage.GetPolicy(id)
	if err != nil {
		return fmt.Errorf("policy not found: %w", err)
	}
	// 不允许删除内置策略
	if policy.Builtin {
		return fmt.Errorf("cannot delete builtin policy")
	}
	return s.storage.DeletePolicy(id)
}

// ExecutePolicy 执行策略，选择资源
func (s *PolicyService) ExecutePolicy(policyID, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	return s.ExecutePolicyWithDimensions(policyID, modelName, resources, nil)
}

// ExecutePolicyWithDimensions 执行策略，选择资源（支持维度过滤）
func (s *PolicyService) ExecutePolicyWithDimensions(policyID, modelName string, resources []*storage.LLMResource, dimensions *int) (*storage.LLMResource, error) {
	// 获取策略
	policy, err := s.storage.GetPolicy(policyID)
	if err != nil {
		return nil, fmt.Errorf("policy not found: %w", err)
	}

	// 检查策略是否启用
	if !policy.Enabled {
		return nil, ErrPolicyDisabled
	}

	// 对于 embedding 类型，如果指定了 dimensions，先过滤出支持该维度的资源
	filteredResources := resources
	if dimensions != nil {
		filteredResources = filterResourcesByDimensions(resources, *dimensions)
		if len(filteredResources) == 0 {
			return nil, fmt.Errorf("no resources available that support dimension %d", *dimensions)
		}
		logger.Debug("Filtered resources by dimensions", logger.F("component", "service"), logger.F("policy_id", policyID), logger.F("dimension", *dimensions), logger.F("original_count", len(resources)), logger.F("filtered_count", len(filteredResources)))
	}

	// 获取执行器
	executor := s.executorFactory.GetExecutor(policy.Type)

	// 执行策略
	return executor.Execute(policy, modelName, filteredResources)
}

// filterResourcesByDimensions 根据维度过滤 embedding 资源
func filterResourcesByDimensions(resources []*storage.LLMResource, dimension int) []*storage.LLMResource {
	filtered := make([]*storage.LLMResource, 0)
	for _, resource := range resources {
		if resource.Type != "embedding" {
			continue
		}
		// 解析 embedding_config
		if resource.EmbeddingConfig == "" {
			// 如果没有配置，允许通过（向后兼容）
			filtered = append(filtered, resource)
			continue
		}
		var config storage.EmbeddingConfig
		if err := json.Unmarshal([]byte(resource.EmbeddingConfig), &config); err != nil {
			// 解析失败，允许通过（向后兼容）
			filtered = append(filtered, resource)
			continue
		}
		// 检查是否支持该维度
		for _, supportedDim := range config.SupportedDimensions {
			if supportedDim == dimension {
				filtered = append(filtered, resource)
				break
			}
		}
	}
	return filtered
}

// GetDefaultPolicyExecutor 获取默认策略执行器（用于没有配置策略的情况）
func (s *PolicyService) GetDefaultPolicyExecutor() PolicyExecutor {
	return NewRoundRobinPolicyExecutor()
}

// InitBuiltinPolicies 初始化内置策略
func (s *PolicyService) InitBuiltinPolicies() error {
	// 获取所有策略模板
	templates, err := s.storage.ListPolicyTemplates()
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	// 定义内置策略配置
	builtinPolicies := []struct {
		name         string
		templateType string
		parameters   map[string]interface{}
	}{
		{
			name:         "默认随机策略",
			templateType: "random",
			parameters: map[string]interface{}{
				"filter_by_status": true,
			},
		},
		{
			name:         "默认轮询策略",
			templateType: "round_robin",
			parameters: map[string]interface{}{
				"resources":        []interface{}{},
				"filter_by_status": true,
			},
		},
		{
			name:         "默认加权策略",
			templateType: "weighted",
			parameters: map[string]interface{}{
				"resources":        []interface{}{},
				"filter_by_status": true,
			},
		},
	}

	for _, bp := range builtinPolicies {
		// 查找对应的模板
		var template *storage.PolicyTemplate
		for _, t := range templates {
			if t.Type == bp.templateType {
				template = t
				break
			}
		}
		if template == nil {
			logger.Warn("Template not found for builtin policy", logger.F("type", bp.templateType))
			continue
		}

		// 检查策略是否已存在（通过名称和内置标记）
		policies, err := s.storage.ListPolicies()
		if err == nil {
			exists := false
			for _, p := range policies {
				if p.Name == bp.name && p.Builtin {
					exists = true
					break
				}
			}
			if exists {
				continue
			}
		}

		// 序列化参数
		paramsJSON, err := json.Marshal(bp.parameters)
		if err != nil {
			logger.Warn("Failed to marshal parameters", logger.F("error", err.Error()), logger.F("policy", bp.name))
			continue
		}

		// 创建内置策略
		policy := &storage.Policy{
			Name:       bp.name,
			TemplateID: template.ID,
			Type:       template.Type,
			Parameters: string(paramsJSON),
			Enabled:    true,
			Builtin:    true,
		}

		if err := s.storage.CreatePolicy(policy); err != nil {
			logger.Warn("Failed to create builtin policy", logger.F("error", err.Error()), logger.F("policy", bp.name))
			continue
		}
	}

	return nil
}
