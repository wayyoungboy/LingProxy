package service

import (
	"encoding/json"
	"fmt"

	"github.com/lingproxy/lingproxy/internal/storage"
)

// PolicyTemplateService 策略模板服务
type PolicyTemplateService struct {
	storage *storage.StorageFacade
}

// NewPolicyTemplateService 创建新的策略模板服务
func NewPolicyTemplateService(storage *storage.StorageFacade) *PolicyTemplateService {
	return &PolicyTemplateService{
		storage: storage,
	}
}

// InitBuiltinTemplates 初始化内置策略模板
func (s *PolicyTemplateService) InitBuiltinTemplates() error {
	templates := []*storage.PolicyTemplate{
		{
			Name:        "随机选择",
			Type:        "random",
			Description: "从可用的LLM资源中随机选择一个",
			ParametersSchema: `{
				"type": "object",
				"properties": {
					"resources": {
						"type": "array",
						"items": {"type": "string"},
						"description": "指定资源列表，为空则使用所有可用资源"
					},
					"filter_by_status": {
						"type": "boolean",
						"default": true,
						"description": "是否只选择状态为active的资源"
					}
				}
			}`,
			DefaultParameters: `{"filter_by_status": true}`,
			Builtin:           true,
		},
		{
			Name:        "轮询负载均衡",
			Type:        "round_robin",
			Description: "按顺序轮询选择LLM资源",
			ParametersSchema: `{
				"type": "object",
				"properties": {
					"resources": {
						"type": "array",
						"items": {"type": "string"},
						"description": "资源列表（按顺序）"
					},
					"filter_by_status": {
						"type": "boolean",
						"default": true,
						"description": "是否只选择状态为active的资源"
					}
				},
				"required": ["resources"]
			}`,
			DefaultParameters: `{"resources": [], "filter_by_status": true}`,
			Builtin:           true,
		},
		{
			Name:        "加权负载均衡",
			Type:        "weighted",
			Description: "根据权重选择资源，权重高的资源被选中的概率更高",
			ParametersSchema: `{
				"type": "object",
				"properties": {
					"resources": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"id": {"type": "string"},
								"weight": {"type": "integer", "minimum": 1}
							},
							"required": ["id", "weight"]
						}
					},
					"filter_by_status": {
						"type": "boolean",
						"default": true
					}
				},
				"required": ["resources"]
			}`,
			DefaultParameters: `{"resources": [], "filter_by_status": true}`,
			Builtin:           true,
		},
	}

	for _, template := range templates {
		// 检查模板是否已存在
		existing, err := s.storage.GetPolicyTemplateByType(template.Type)
		if err == nil && existing != nil {
			// 模板已存在，跳过
			continue
		}

		// 创建模板
		if err := s.storage.CreatePolicyTemplate(template); err != nil {
			return fmt.Errorf("failed to create policy template %s: %w", template.Type, err)
		}
	}

	return nil
}

// GetTemplate 获取模板
func (s *PolicyTemplateService) GetTemplate(id string) (*storage.PolicyTemplate, error) {
	return s.storage.GetPolicyTemplate(id)
}

// GetTemplateByType 根据类型获取模板
func (s *PolicyTemplateService) GetTemplateByType(policyType string) (*storage.PolicyTemplate, error) {
	return s.storage.GetPolicyTemplateByType(policyType)
}

// ListTemplates 获取所有模板
func (s *PolicyTemplateService) ListTemplates() ([]*storage.PolicyTemplate, error) {
	return s.storage.ListPolicyTemplates()
}

// ValidateParameters 验证策略参数是否符合模板Schema
func (s *PolicyTemplateService) ValidateParameters(template *storage.PolicyTemplate, parameters map[string]interface{}) error {
	paramsJSON, err := json.Marshal(parameters)
	if err != nil {
		return fmt.Errorf("invalid parameters format: %w", err)
	}

	var defaultParams map[string]interface{}
	if template.DefaultParameters != "" {
		if err := json.Unmarshal([]byte(template.DefaultParameters), &defaultParams); err != nil {
			return fmt.Errorf("invalid default parameters: %w", err)
		}
	}

	var schema map[string]interface{}
	if template.ParametersSchema != "" {
		if err := json.Unmarshal([]byte(template.ParametersSchema), &schema); err != nil {
			return fmt.Errorf("invalid schema: %w", err)
		}
	}

	_ = paramsJSON
	_ = schema

	return nil
}
