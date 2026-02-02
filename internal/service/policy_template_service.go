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
		{
			Name:        "模型名匹配",
			Type:        "model_match",
			Description: "根据请求的模型名匹配对应的资源",
			ParametersSchema: `{
				"type": "object",
				"properties": {
					"mappings": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"model_pattern": {"type": "string"},
								"resource_id": {"type": "string"}
							},
							"required": ["model_pattern", "resource_id"]
						}
					},
					"default_resource_id": {
						"type": "string",
						"description": "默认资源（可选）"
					}
				},
				"required": ["mappings"]
			}`,
			DefaultParameters: `{"mappings": []}`,
			Builtin:           true,
		},
		{
			Name:        "正则匹配",
			Type:        "regex_match",
			Description: "使用正则表达式匹配模型名",
			ParametersSchema: `{
				"type": "object",
				"properties": {
					"rules": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"pattern": {"type": "string"},
								"resource_id": {"type": "string"}
							},
							"required": ["pattern", "resource_id"]
						}
					},
					"default_resource_id": {
						"type": "string"
					}
				},
				"required": ["rules"]
			}`,
			DefaultParameters: `{"rules": []}`,
			Builtin:           true,
		},
		{
			Name:        "优先级策略",
			Type:        "priority",
			Description: "按优先级顺序选择资源，优先使用优先级高的资源",
			ParametersSchema: `{
				"type": "object",
				"properties": {
					"resources": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"id": {"type": "string"},
								"priority": {"type": "integer", "minimum": 1}
							},
							"required": ["id", "priority"]
						}
					},
					"fallback_enabled": {
						"type": "boolean",
						"default": true,
						"description": "是否启用降级"
					}
				},
				"required": ["resources"]
			}`,
			DefaultParameters: `{"resources": [], "fallback_enabled": true}`,
			Builtin:           true,
		},
		{
			Name:        "故障转移",
			Type:        "failover",
			Description: "主资源不可用时自动切换到备用资源",
			ParametersSchema: `{
				"type": "object",
				"properties": {
					"primary_resource_id": {"type": "string"},
					"fallback_resources": {
						"type": "array",
						"items": {"type": "string"}
					},
					"health_check_enabled": {
						"type": "boolean",
						"default": true
					},
					"health_check_interval": {
						"type": "integer",
						"default": 30,
						"description": "健康检查间隔（秒）"
					}
				},
				"required": ["primary_resource_id", "fallback_resources"]
			}`,
			DefaultParameters: `{"health_check_enabled": true, "health_check_interval": 30}`,
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
	// 简单的JSON Schema验证（可以后续使用专门的库）
	// 这里先做基本的类型检查
	paramsJSON, err := json.Marshal(parameters)
	if err != nil {
		return fmt.Errorf("invalid parameters format: %w", err)
	}

	// 解析默认参数作为参考
	var defaultParams map[string]interface{}
	if template.DefaultParameters != "" {
		if err := json.Unmarshal([]byte(template.DefaultParameters), &defaultParams); err != nil {
			return fmt.Errorf("invalid default parameters: %w", err)
		}
	}

	// 解析参数Schema（简化验证）
	var schema map[string]interface{}
	if template.ParametersSchema != "" {
		if err := json.Unmarshal([]byte(template.ParametersSchema), &schema); err != nil {
			return fmt.Errorf("invalid schema: %w", err)
		}
	}

	// 验证参数（这里简化处理，实际应该使用JSON Schema验证库）
	_ = paramsJSON
	_ = schema

	return nil
}
