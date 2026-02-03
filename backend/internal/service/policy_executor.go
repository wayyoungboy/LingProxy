package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/storage"
)

var (
	ErrNoResourcesAvailable = errors.New("no resources available")
	ErrPolicyNotFound       = errors.New("policy not found")
	ErrInvalidPolicyParams  = errors.New("invalid policy parameters")
)

// PolicyExecutor 策略执行器接口
type PolicyExecutor interface {
	Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error)
}

// RandomPolicyExecutor 随机选择策略执行器
type RandomPolicyExecutor struct{}

func (e *RandomPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Info("执行随机策略", logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model_name", modelName))

	var params struct {
		Resources      []string `json:"resources"`
		FilterByStatus bool     `json:"filter_by_status"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("解析随机策略参数失败", logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 过滤资源
	filtered := filterResources(resources, params.Resources, params.FilterByStatus)
	logger.Info("随机策略资源过滤结果", logger.F("policy_id", policy.ID), logger.F("original_count", len(resources)), logger.F("filtered_count", len(filtered)))

	if len(filtered) == 0 {
		logger.Warn("随机策略无可用资源", logger.F("policy_id", policy.ID), logger.F("model_name", modelName))
		return nil, ErrNoResourcesAvailable
	}

	// 随机选择（Go 1.20+ 全局随机数生成器已自动初始化）
	selected := filtered[rand.Intn(len(filtered))]
	logger.Info("随机策略选择结果", logger.F("policy_id", policy.ID), logger.F("resource_id", selected.ID), logger.F("resource_name", selected.Name), logger.F("resource_name_base_url", selected.BaseURL))

	return selected, nil
}

// RoundRobinPolicyExecutor 轮询负载均衡策略执行器
type RoundRobinPolicyExecutor struct {
	mu           sync.Mutex
	currentIndex map[string]int // 每个策略的当前索引
}

func NewRoundRobinPolicyExecutor() *RoundRobinPolicyExecutor {
	return &RoundRobinPolicyExecutor{
		currentIndex: make(map[string]int),
	}
}

func (e *RoundRobinPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Info("执行轮询策略", logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model_name", modelName))

	var params struct {
		Resources      []string `json:"resources"`
		FilterByStatus bool     `json:"filter_by_status"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("解析轮询策略参数失败", logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 过滤资源
	filtered := filterResources(resources, params.Resources, params.FilterByStatus)
	logger.Info("轮询策略资源过滤结果", logger.F("policy_id", policy.ID), logger.F("original_count", len(resources)), logger.F("filtered_count", len(filtered)))

	if len(filtered) == 0 {
		logger.Warn("轮询策略无可用资源", logger.F("policy_id", policy.ID), logger.F("model_name", modelName))
		return nil, ErrNoResourcesAvailable
	}

	// 轮询选择
	e.mu.Lock()
	index := e.currentIndex[policy.ID]
	nextIndex := (index + 1) % len(filtered)
	e.currentIndex[policy.ID] = nextIndex
	e.mu.Unlock()

	selected := filtered[index]
	logger.Info("轮询策略选择结果", logger.F("policy_id", policy.ID), logger.F("current_index", index), logger.F("next_index", nextIndex), logger.F("resource_id", selected.ID), logger.F("resource_name", selected.Name))

	return selected, nil
}

// WeightedPolicyExecutor 加权负载均衡策略执行器
type WeightedPolicyExecutor struct{}

func (e *WeightedPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Info("执行加权策略", logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model_name", modelName))

	var params struct {
		Resources []struct {
			ID     string `json:"id"`
			Weight int    `json:"weight"`
		} `json:"resources"`
		FilterByStatus bool `json:"filter_by_status"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("解析加权策略参数失败", logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建资源映射
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if params.FilterByStatus && r.Status != "active" {
			continue
		}
		resourceMap[r.ID] = r
	}

	// 计算总权重
	totalWeight := 0
	validResources := make([]struct {
		Resource *storage.LLMResource
		Weight   int
	}, 0)

	for _, wr := range params.Resources {
		if r, exists := resourceMap[wr.ID]; exists {
			totalWeight += wr.Weight
			validResources = append(validResources, struct {
				Resource *storage.LLMResource
				Weight   int
			}{Resource: r, Weight: wr.Weight})
		}
	}

	logger.Info("加权策略资源处理结果", logger.F("policy_id", policy.ID), logger.F("original_count", len(resources)), logger.F("valid_count", len(validResources)), logger.F("total_weight", totalWeight))

	if len(validResources) == 0 {
		logger.Warn("加权策略无可用资源", logger.F("policy_id", policy.ID), logger.F("model_name", modelName))
		return nil, ErrNoResourcesAvailable
	}

	// 根据权重随机选择（Go 1.20+ 全局随机数生成器已自动初始化）
	random := rand.Intn(totalWeight)
	currentWeight := 0

	for _, wr := range validResources {
		currentWeight += wr.Weight
		if random < currentWeight {
			logger.Info("加权策略选择结果", logger.F("policy_id", policy.ID), logger.F("resource_id", wr.Resource.ID), logger.F("resource_name", wr.Resource.Name), logger.F("weight", wr.Weight))
			return wr.Resource, nil
		}
	}

	// 兜底返回第一个资源
	selected := validResources[0].Resource
	logger.Info("加权策略兜底选择", logger.F("policy_id", policy.ID), logger.F("resource_id", selected.ID), logger.F("resource_name", selected.Name))
	return selected, nil
}

// ModelMatchPolicyExecutor 模型名匹配策略执行器
type ModelMatchPolicyExecutor struct{}

func (e *ModelMatchPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Info("执行模型匹配策略", logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model_name", modelName))

	var params struct {
		Mappings []struct {
			ModelPattern string `json:"model_pattern"`
			ResourceID   string `json:"resource_id"`
		} `json:"mappings"`
		DefaultResourceID string `json:"default_resource_id"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("解析模型匹配策略参数失败", logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建资源映射
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if r.Status == "active" {
			resourceMap[r.ID] = r
		}
	}

	// 匹配模型名（支持通配符）
	for _, mapping := range params.Mappings {
		if matchPattern(modelName, mapping.ModelPattern) {
			logger.Info("模型匹配策略找到匹配", logger.F("policy_id", policy.ID), logger.F("model_name", modelName), logger.F("pattern", mapping.ModelPattern), logger.F("resource_id", mapping.ResourceID))
			if resource, exists := resourceMap[mapping.ResourceID]; exists {
				logger.Info("模型匹配策略选择结果", logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name))
				return resource, nil
			}
			logger.Warn("模型匹配策略资源不存在", logger.F("policy_id", policy.ID), logger.F("resource_id", mapping.ResourceID))
		}
	}

	// 使用默认资源
	if params.DefaultResourceID != "" {
		logger.Info("模型匹配策略使用默认资源", logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
		if resource, exists := resourceMap[params.DefaultResourceID]; exists {
			logger.Info("模型匹配策略默认资源选择结果", logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name))
			return resource, nil
		}
		logger.Warn("模型匹配策略默认资源不存在", logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
	}

	logger.Warn("模型匹配策略无可用资源", logger.F("policy_id", policy.ID), logger.F("model_name", modelName))
	return nil, ErrNoResourcesAvailable
}

// RegexMatchPolicyExecutor 正则匹配策略执行器
type RegexMatchPolicyExecutor struct{}

func (e *RegexMatchPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Info("执行正则匹配策略", logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model_name", modelName))

	var params struct {
		Rules []struct {
			Pattern    string `json:"pattern"`
			ResourceID string `json:"resource_id"`
		} `json:"rules"`
		DefaultResourceID string `json:"default_resource_id"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("解析正则匹配策略参数失败", logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建资源映射
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if r.Status == "active" {
			resourceMap[r.ID] = r
		}
	}

	// 正则匹配
	for _, rule := range params.Rules {
		matched, err := regexp.MatchString(rule.Pattern, modelName)
		if err != nil {
			logger.Warn("正则匹配策略规则无效", logger.F("policy_id", policy.ID), logger.F("pattern", rule.Pattern), logger.F("error", err.Error()))
			continue // 跳过无效的正则表达式
		}
		if matched {
			logger.Info("正则匹配策略找到匹配", logger.F("policy_id", policy.ID), logger.F("model_name", modelName), logger.F("pattern", rule.Pattern), logger.F("resource_id", rule.ResourceID))
			if resource, exists := resourceMap[rule.ResourceID]; exists {
				logger.Info("正则匹配策略选择结果", logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name))
				return resource, nil
			}
			logger.Warn("正则匹配策略资源不存在", logger.F("policy_id", policy.ID), logger.F("resource_id", rule.ResourceID))
		}
	}

	// 使用默认资源
	if params.DefaultResourceID != "" {
		logger.Info("正则匹配策略使用默认资源", logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
		if resource, exists := resourceMap[params.DefaultResourceID]; exists {
			logger.Info("正则匹配策略默认资源选择结果", logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name))
			return resource, nil
		}
		logger.Warn("正则匹配策略默认资源不存在", logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
	}

	logger.Warn("正则匹配策略无可用资源", logger.F("policy_id", policy.ID), logger.F("model_name", modelName))
	return nil, ErrNoResourcesAvailable
}

// PriorityPolicyExecutor 优先级策略执行器
type PriorityPolicyExecutor struct{}

func (e *PriorityPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Info("执行优先级策略", logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model_name", modelName))

	var params struct {
		Resources []struct {
			ID       string `json:"id"`
			Priority int    `json:"priority"`
		} `json:"resources"`
		FallbackEnabled bool `json:"fallback_enabled"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("解析优先级策略参数失败", logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建资源映射
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if r.Status == "active" {
			resourceMap[r.ID] = r
		}
	}

	// 按优先级排序
	sort.Slice(params.Resources, func(i, j int) bool {
		return params.Resources[i].Priority < params.Resources[j].Priority
	})

	// 选择第一个可用的资源
	for _, pr := range params.Resources {
		if resource, exists := resourceMap[pr.ID]; exists {
			logger.Info("优先级策略选择结果", logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name), logger.F("priority", pr.Priority))
			return resource, nil
		}
		logger.Warn("优先级策略资源不可用", logger.F("policy_id", policy.ID), logger.F("resource_id", pr.ID), logger.F("priority", pr.Priority))
	}

	if params.FallbackEnabled {
		// 降级：使用任何可用的资源
		logger.Info("优先级策略启用降级", logger.F("policy_id", policy.ID))
		for _, r := range resources {
			if r.Status == "active" {
				logger.Info("优先级策略降级选择结果", logger.F("policy_id", policy.ID), logger.F("resource_id", r.ID), logger.F("resource_name", r.Name))
				return r, nil
			}
		}
	}

	logger.Warn("优先级策略无可用资源", logger.F("policy_id", policy.ID), logger.F("model_name", modelName))
	return nil, ErrNoResourcesAvailable
}

// FailoverPolicyExecutor 故障转移策略执行器
type FailoverPolicyExecutor struct {
	mu             sync.RWMutex
	resourceStatus map[string]bool // 资源健康状态
	lastCheckTime  map[string]time.Time
}

func NewFailoverPolicyExecutor() *FailoverPolicyExecutor {
	return &FailoverPolicyExecutor{
		resourceStatus: make(map[string]bool),
		lastCheckTime:  make(map[string]time.Time),
	}
}

func (e *FailoverPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Info("执行故障转移策略", logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model_name", modelName))

	var params struct {
		PrimaryResourceID   string   `json:"primary_resource_id"`
		FallbackResources   []string `json:"fallback_resources"`
		HealthCheckEnabled  bool     `json:"health_check_enabled"`
		HealthCheckInterval int      `json:"health_check_interval"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("解析故障转移策略参数失败", logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建资源映射
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if r.Status == "active" {
			resourceMap[r.ID] = r
		}
	}

	// 检查主资源
	if primary, exists := resourceMap[params.PrimaryResourceID]; exists {
		isHealthy := !params.HealthCheckEnabled || e.isResourceHealthy(params.PrimaryResourceID, params.HealthCheckInterval)
		logger.Info("故障转移策略检查主资源", logger.F("policy_id", policy.ID), logger.F("primary_resource_id", params.PrimaryResourceID), logger.F("exists", exists), logger.F("is_healthy", isHealthy))

		if isHealthy {
			logger.Info("故障转移策略选择主资源", logger.F("policy_id", policy.ID), logger.F("resource_id", primary.ID), logger.F("resource_name", primary.Name))
			return primary, nil
		}
		logger.Warn("故障转移策略主资源不健康", logger.F("policy_id", policy.ID), logger.F("primary_resource_id", params.PrimaryResourceID))
	} else {
		logger.Warn("故障转移策略主资源不存在", logger.F("policy_id", policy.ID), logger.F("primary_resource_id", params.PrimaryResourceID))
	}

	// 使用备用资源
	logger.Info("故障转移策略尝试备用资源", logger.F("policy_id", policy.ID), logger.F("fallback_count", len(params.FallbackResources)))
	for _, fallbackID := range params.FallbackResources {
		if fallback, exists := resourceMap[fallbackID]; exists {
			isHealthy := !params.HealthCheckEnabled || e.isResourceHealthy(fallbackID, params.HealthCheckInterval)
			logger.Info("故障转移策略检查备用资源", logger.F("policy_id", policy.ID), logger.F("fallback_resource_id", fallbackID), logger.F("exists", exists), logger.F("is_healthy", isHealthy))

			if isHealthy {
				logger.Info("故障转移策略选择备用资源", logger.F("policy_id", policy.ID), logger.F("resource_id", fallback.ID), logger.F("resource_name", fallback.Name))
				return fallback, nil
			}
			logger.Warn("故障转移策略备用资源不健康", logger.F("policy_id", policy.ID), logger.F("fallback_resource_id", fallbackID))
		} else {
			logger.Warn("故障转移策略备用资源不存在", logger.F("policy_id", policy.ID), logger.F("fallback_resource_id", fallbackID))
		}
	}

	logger.Warn("故障转移策略无可用资源", logger.F("policy_id", policy.ID), logger.F("model_name", modelName))
	return nil, ErrNoResourcesAvailable
}

func (e *FailoverPolicyExecutor) isResourceHealthy(resourceID string, checkInterval int) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	lastCheck, exists := e.lastCheckTime[resourceID]
	if !exists || time.Since(lastCheck) > time.Duration(checkInterval)*time.Second {
		// 需要检查，这里简化处理，假设资源是健康的
		// 实际应该进行健康检查
		logger.Debug("资源健康检查（需要检查）", logger.F("resource_id", resourceID), logger.F("last_check", lastCheck), logger.F("check_interval", checkInterval))
		return true
	}

	isHealthy := e.resourceStatus[resourceID]
	logger.Debug("资源健康检查结果", logger.F("resource_id", resourceID), logger.F("is_healthy", isHealthy), logger.F("last_check", lastCheck))
	return isHealthy
}

// PolicyExecutorFactory 策略执行器工厂
type PolicyExecutorFactory struct {
	roundRobinExecutor *RoundRobinPolicyExecutor
	failoverExecutor   *FailoverPolicyExecutor
}

func NewPolicyExecutorFactory() *PolicyExecutorFactory {
	return &PolicyExecutorFactory{
		roundRobinExecutor: NewRoundRobinPolicyExecutor(),
		failoverExecutor:   NewFailoverPolicyExecutor(),
	}
}

func (f *PolicyExecutorFactory) GetExecutor(policyType string) PolicyExecutor {
	logger.Info("获取策略执行器", logger.F("policy_type", policyType))

	switch policyType {
	case "random":
		logger.Info("使用随机策略执行器", logger.F("policy_type", policyType))
		return &RandomPolicyExecutor{}
	case "round_robin":
		logger.Info("使用轮询策略执行器", logger.F("policy_type", policyType))
		return f.roundRobinExecutor
	case "weighted":
		logger.Info("使用加权策略执行器", logger.F("policy_type", policyType))
		return &WeightedPolicyExecutor{}
	case "model_match":
		logger.Info("使用模型匹配策略执行器", logger.F("policy_type", policyType))
		return &ModelMatchPolicyExecutor{}
	case "regex_match":
		logger.Info("使用正则匹配策略执行器", logger.F("policy_type", policyType))
		return &RegexMatchPolicyExecutor{}
	case "priority":
		logger.Info("使用优先级策略执行器", logger.F("policy_type", policyType))
		return &PriorityPolicyExecutor{}
	case "failover":
		logger.Info("使用故障转移策略执行器", logger.F("policy_type", policyType))
		return f.failoverExecutor
	default:
		logger.Warn("未知策略类型，使用默认随机策略", logger.F("policy_type", policyType))
		return &RandomPolicyExecutor{} // 默认使用随机策略
	}
}

// filterResources 过滤资源
func filterResources(resources []*storage.LLMResource, resourceIDs []string, filterByStatus bool) []*storage.LLMResource {
	filtered := make([]*storage.LLMResource, 0)

	logger.Debug("开始过滤资源", logger.F("original_count", len(resources)), logger.F("filter_by_status", filterByStatus), logger.F("resource_ids_count", len(resourceIDs)))

	for _, r := range resources {
		// 状态过滤
		if filterByStatus && r.Status != "active" {
			logger.Debug("资源过滤：状态不活跃", logger.F("resource_id", r.ID), logger.F("resource_name", r.Name), logger.F("status", r.Status))
			continue
		}

		// ID过滤
		if len(resourceIDs) > 0 {
			found := false
			for _, id := range resourceIDs {
				if r.ID == id {
					found = true
					break
				}
			}
			if !found {
				logger.Debug("资源过滤：ID不在列表中", logger.F("resource_id", r.ID), logger.F("resource_name", r.Name))
				continue
			}
		}

		filtered = append(filtered, r)
		logger.Debug("资源过滤：通过", logger.F("resource_id", r.ID), logger.F("resource_name", r.Name))
	}

	logger.Debug("资源过滤完成", logger.F("filtered_count", len(filtered)))
	return filtered
}

// matchPattern 匹配模式（支持通配符）
func matchPattern(text, pattern string) bool {
	// 简单的通配符匹配：* 匹配任意字符
	if strings.Contains(pattern, "*") {
		regexPattern := "^" + strings.ReplaceAll(strings.ReplaceAll(pattern, ".", "\\."), "*", ".*") + "$"
		matched, _ := regexp.MatchString(regexPattern, text)
		logger.Debug("模式匹配结果（通配符）", logger.F("text", text), logger.F("pattern", pattern), logger.F("regex_pattern", regexPattern), logger.F("matched", matched))
		return matched
	}
	matched := text == pattern
	logger.Debug("模式匹配结果（精确）", logger.F("text", text), logger.F("pattern", pattern), logger.F("matched", matched))
	return matched
}
