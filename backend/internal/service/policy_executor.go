package service

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	mathrand "math/rand"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// init 初始化随机数生成器，使用crypto/rand作为种子源以确保真正的随机性
func init() {
	// 使用crypto/rand生成种子，确保每次程序启动时都有不同的随机序列
	var seed int64
	if err := binary.Read(rand.Reader, binary.BigEndian, &seed); err != nil {
		// 如果crypto/rand失败，使用时间戳作为后备
		seed = time.Now().UnixNano()
	}
	mathrand.Seed(seed)
}

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
	logger.Debug("Executing random policy", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model", modelName))

	var params struct {
		Resources      []string `json:"resources"`
		FilterByStatus bool     `json:"filter_by_status"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("Failed to parse random policy parameters", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 过滤资源
	filtered := filterResources(resources, params.Resources, params.FilterByStatus)
	logger.Debug("Random policy filtered resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("original_count", len(resources)), logger.F("filtered_count", len(filtered)))

	if len(filtered) == 0 {
		logger.Warn("Random policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
		return nil, ErrNoResourcesAvailable
	}

	// 随机选择（使用已初始化的随机数生成器）
	selected := filtered[mathrand.Intn(len(filtered))]
		logger.Debug("Random policy selected resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", selected.ID), logger.F("resource_name", selected.Name), logger.F("base_url", selected.BaseURL))

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
	logger.Debug("Executing round-robin policy", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model", modelName))

	var params struct {
		Resources      []string `json:"resources"`
		FilterByStatus bool     `json:"filter_by_status"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("Failed to parse round-robin policy parameters", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 过滤资源
	filtered := filterResources(resources, params.Resources, params.FilterByStatus)
	logger.Debug("Round-robin policy filtered resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("original_count", len(resources)), logger.F("filtered_count", len(filtered)))

	if len(filtered) == 0 {
		logger.Warn("Round-robin policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
		return nil, ErrNoResourcesAvailable
	}

	// 轮询选择
	e.mu.Lock()
	index := e.currentIndex[policy.ID]
	nextIndex := (index + 1) % len(filtered)
	e.currentIndex[policy.ID] = nextIndex
	e.mu.Unlock()

	selected := filtered[index]
		logger.Debug("Round-robin policy selected resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("current_index", index), logger.F("next_index", nextIndex), logger.F("resource_id", selected.ID), logger.F("resource_name", selected.Name))

	return selected, nil
}

// WeightedPolicyExecutor 加权负载均衡策略执行器
type WeightedPolicyExecutor struct{}

func (e *WeightedPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Debug("Executing weighted policy", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model", modelName))

	var params struct {
		Resources []struct {
			ID     string `json:"id"`
			Weight int    `json:"weight"`
		} `json:"resources"`
		FilterByStatus bool `json:"filter_by_status"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("Failed to parse weighted policy parameters", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建资源映射（包含所有资源，用于存在性验证）
	allResourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		allResourceMap[r.ID] = r
	}

	// 构建可用资源映射（根据FilterByStatus过滤）
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if params.FilterByStatus && r.Status != "active" {
			continue
		}
		resourceMap[r.ID] = r
	}

	// 验证配置的资源是否存在，并计算总权重
	totalWeight := 0
	validResources := make([]struct {
		Resource *storage.LLMResource
		Weight   int
	}, 0)
	missingResources := make([]string, 0)

	for _, wr := range params.Resources {
		// 首先验证资源是否存在（在所有资源中）
		if _, exists := allResourceMap[wr.ID]; !exists {
			missingResources = append(missingResources, wr.ID)
			continue
		}
		// 然后检查资源是否在可用资源映射中
		if r, exists := resourceMap[wr.ID]; exists {
			totalWeight += wr.Weight
			validResources = append(validResources, struct {
				Resource *storage.LLMResource
				Weight   int
			}{Resource: r, Weight: wr.Weight})
		}
	}

	// 记录缺失资源的警告
	if len(missingResources) > 0 {
		logger.Warn("Weighted policy: some configured resources are not found", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("missing_resource_ids", missingResources), logger.F("total_configured", len(params.Resources)))
	}

	logger.Debug("Weighted policy processed resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("original_count", len(resources)), logger.F("valid_count", len(validResources)), logger.F("total_weight", totalWeight))

	if len(validResources) == 0 {
		logger.Warn("Weighted policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
		return nil, ErrNoResourcesAvailable
	}

	// 根据权重随机选择（使用已初始化的随机数生成器）
	random := mathrand.Intn(totalWeight)
	currentWeight := 0

	for _, wr := range validResources {
		currentWeight += wr.Weight
		if random < currentWeight {
			logger.Debug("Weighted policy selected resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", wr.Resource.ID), logger.F("resource_name", wr.Resource.Name), logger.F("weight", wr.Weight))
			return wr.Resource, nil
		}
	}

	// 兜底返回第一个资源
	selected := validResources[0].Resource
	logger.Debug("Weighted policy fallback selection", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", selected.ID), logger.F("resource_name", selected.Name))
	return selected, nil
}

// ModelMatchPolicyExecutor 模型名匹配策略执行器
type ModelMatchPolicyExecutor struct{}

func (e *ModelMatchPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Debug("Executing model-match policy", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model", modelName))

	var params struct {
		Mappings []struct {
			ModelPattern string `json:"model_pattern"`
			ResourceID   string `json:"resource_id"`
		} `json:"mappings"`
		DefaultResourceID string `json:"default_resource_id"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("Failed to parse model-match policy parameters", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建所有资源映射（用于存在性验证）
	allResourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		allResourceMap[r.ID] = r
	}

	// 构建可用资源映射（只包含active状态）
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if r.Status == "active" {
			resourceMap[r.ID] = r
		}
	}

	// 匹配模型名（支持通配符）
	for _, mapping := range params.Mappings {
		if matchPattern(modelName, mapping.ModelPattern) {
			logger.Debug("Model-match policy found match", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName), logger.F("pattern", mapping.ModelPattern), logger.F("resource_id", mapping.ResourceID))
			// 首先验证资源是否存在
			if _, exists := allResourceMap[mapping.ResourceID]; !exists {
				logger.Warn("Model-match policy: configured resource not found", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", mapping.ResourceID), logger.F("pattern", mapping.ModelPattern))
				continue
			}
			// 然后检查资源是否可用
			if resource, exists := resourceMap[mapping.ResourceID]; exists {
				logger.Debug("Model-match policy selected resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name))
				return resource, nil
			}
			logger.Warn("Model-match policy: resource exists but not active", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", mapping.ResourceID))
		}
	}

	// 使用默认资源
	if params.DefaultResourceID != "" {
		logger.Debug("Model-match policy using default resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
		// 首先验证默认资源是否存在
		if _, exists := allResourceMap[params.DefaultResourceID]; !exists {
			logger.Warn("Model-match policy: configured default resource not found", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
		} else if resource, exists := resourceMap[params.DefaultResourceID]; exists {
			logger.Debug("Model-match policy selected default resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name))
			return resource, nil
		} else {
			logger.Warn("Model-match policy: default resource exists but not active", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
		}
	}

	logger.Warn("Model-match policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
	return nil, ErrNoResourcesAvailable
}

// RegexMatchPolicyExecutor 正则匹配策略执行器
type RegexMatchPolicyExecutor struct{}

func (e *RegexMatchPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Debug("Executing regex-match policy", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model", modelName))

	var params struct {
		Rules []struct {
			Pattern    string `json:"pattern"`
			ResourceID string `json:"resource_id"`
		} `json:"rules"`
		DefaultResourceID string `json:"default_resource_id"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("Failed to parse regex-match policy parameters", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建所有资源映射（用于存在性验证）
	allResourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		allResourceMap[r.ID] = r
	}

	// 构建可用资源映射（只包含active状态）
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
			logger.Warn("Regex-match policy: invalid pattern", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("pattern", rule.Pattern), logger.F("error", err.Error()))
			continue // 跳过无效的正则表达式
		}
		if matched {
			logger.Debug("Regex-match policy found match", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName), logger.F("pattern", rule.Pattern), logger.F("resource_id", rule.ResourceID))
			// 首先验证资源是否存在
			if _, exists := allResourceMap[rule.ResourceID]; !exists {
				logger.Warn("Regex-match policy: configured resource not found", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", rule.ResourceID), logger.F("pattern", rule.Pattern))
				continue
			}
			// 然后检查资源是否可用
			if resource, exists := resourceMap[rule.ResourceID]; exists {
				logger.Debug("Regex-match policy selected resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name))
				return resource, nil
			}
			logger.Warn("Regex-match policy: resource exists but not active", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", rule.ResourceID))
		}
	}

	// 使用默认资源
	if params.DefaultResourceID != "" {
		logger.Debug("Regex-match policy using default resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
		// 首先验证默认资源是否存在
		if _, exists := allResourceMap[params.DefaultResourceID]; !exists {
			logger.Warn("Regex-match policy: configured default resource not found", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
		} else if resource, exists := resourceMap[params.DefaultResourceID]; exists {
			logger.Debug("Regex-match policy selected default resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name))
			return resource, nil
		} else {
			logger.Warn("Regex-match policy: default resource exists but not active", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("default_resource_id", params.DefaultResourceID))
		}
	}

	logger.Warn("Regex-match policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
	return nil, ErrNoResourcesAvailable
}

// PriorityPolicyExecutor 优先级策略执行器
type PriorityPolicyExecutor struct{}

func (e *PriorityPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	logger.Debug("Executing priority policy", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model", modelName))

	var params struct {
		Resources []struct {
			ID       string `json:"id"`
			Priority int    `json:"priority"`
		} `json:"resources"`
		FallbackEnabled bool `json:"fallback_enabled"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("Failed to parse priority policy parameters", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建所有资源映射（用于存在性验证）
	allResourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		allResourceMap[r.ID] = r
	}

	// 构建可用资源映射（只包含active状态）
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if r.Status == "active" {
			resourceMap[r.ID] = r
		}
	}

	// 验证配置的资源是否存在
	missingResources := make([]string, 0)
	for _, pr := range params.Resources {
		if _, exists := allResourceMap[pr.ID]; !exists {
			missingResources = append(missingResources, pr.ID)
		}
	}
	if len(missingResources) > 0 {
		logger.Warn("Priority policy: some configured resources are not found", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("missing_resource_ids", missingResources), logger.F("total_configured", len(params.Resources)))
	}

	// 按优先级排序
	sort.Slice(params.Resources, func(i, j int) bool {
		return params.Resources[i].Priority < params.Resources[j].Priority
	})

	// 选择第一个可用的资源
	for _, pr := range params.Resources {
		// 跳过不存在的资源
		if _, exists := allResourceMap[pr.ID]; !exists {
			continue
		}
		if resource, exists := resourceMap[pr.ID]; exists {
			logger.Debug("Priority policy selected resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", resource.ID), logger.F("resource_name", resource.Name), logger.F("priority", pr.Priority))
			return resource, nil
		}
		logger.Debug("Priority policy: resource exists but not active", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", pr.ID), logger.F("priority", pr.Priority))
	}

	if params.FallbackEnabled {
		// 降级：使用任何可用的资源
		logger.Debug("Priority policy: fallback enabled", logger.F("component", "service"), logger.F("policy_id", policy.ID))
		for _, r := range resources {
			if r.Status == "active" {
				logger.Debug("Priority policy selected fallback resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", r.ID), logger.F("resource_name", r.Name))
				return r, nil
			}
		}
	}

	logger.Warn("Priority policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
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
	logger.Debug("Executing failover policy", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("policy_name", policy.Name), logger.F("model", modelName))

	var params struct {
		PrimaryResourceID   string   `json:"primary_resource_id"`
		FallbackResources   []string `json:"fallback_resources"`
		HealthCheckEnabled  bool     `json:"health_check_enabled"`
		HealthCheckInterval int      `json:"health_check_interval"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		logger.Error("Failed to parse failover policy parameters", logger.F("component", "service"), logger.F("error", err.Error()), logger.F("policy_id", policy.ID))
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 构建所有资源映射（用于存在性验证）
	allResourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		allResourceMap[r.ID] = r
	}

	// 构建可用资源映射（只包含active状态）
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if r.Status == "active" {
			resourceMap[r.ID] = r
		}
	}

	// 验证主资源是否存在
	if params.PrimaryResourceID != "" {
		if _, exists := allResourceMap[params.PrimaryResourceID]; !exists {
			logger.Warn("Failover policy: configured primary resource not found", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("primary_resource_id", params.PrimaryResourceID))
		}
	}

	// 验证备用资源是否存在
	missingFallbackResources := make([]string, 0)
	for _, fallbackID := range params.FallbackResources {
		if _, exists := allResourceMap[fallbackID]; !exists {
			missingFallbackResources = append(missingFallbackResources, fallbackID)
		}
	}
	if len(missingFallbackResources) > 0 {
		logger.Warn("Failover policy: some configured fallback resources are not found", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("missing_fallback_resource_ids", missingFallbackResources), logger.F("total_configured", len(params.FallbackResources)))
	}

	// 检查主资源
	if params.PrimaryResourceID != "" {
		if _, exists := allResourceMap[params.PrimaryResourceID]; exists {
			if primary, exists := resourceMap[params.PrimaryResourceID]; exists {
				isHealthy := !params.HealthCheckEnabled || e.isResourceHealthy(params.PrimaryResourceID, params.HealthCheckInterval)
				logger.Debug("Failover policy checking primary resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("primary_resource_id", params.PrimaryResourceID), logger.F("is_healthy", isHealthy))

				if isHealthy {
					logger.Debug("Failover policy selected primary resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", primary.ID), logger.F("resource_name", primary.Name))
					return primary, nil
				}
				logger.Warn("Failover policy: primary resource unhealthy", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("primary_resource_id", params.PrimaryResourceID))
			} else {
				logger.Warn("Failover policy: primary resource exists but not active", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("primary_resource_id", params.PrimaryResourceID))
			}
		}
	}

	// 使用备用资源
	logger.Debug("Failover policy trying fallback resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("fallback_count", len(params.FallbackResources)))
	for _, fallbackID := range params.FallbackResources {
		// 跳过不存在的资源
		if _, exists := allResourceMap[fallbackID]; !exists {
			continue
		}
		if fallback, exists := resourceMap[fallbackID]; exists {
			isHealthy := !params.HealthCheckEnabled || e.isResourceHealthy(fallbackID, params.HealthCheckInterval)
			logger.Debug("Failover policy checking fallback resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("fallback_resource_id", fallbackID), logger.F("is_healthy", isHealthy))

			if isHealthy {
				logger.Debug("Failover policy selected fallback resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", fallback.ID), logger.F("resource_name", fallback.Name))
				return fallback, nil
			}
			logger.Debug("Failover policy: fallback resource unhealthy", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("fallback_resource_id", fallbackID))
		} else {
			logger.Debug("Failover policy: fallback resource exists but not active", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("fallback_resource_id", fallbackID))
		}
	}

	logger.Warn("Failover policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
	return nil, ErrNoResourcesAvailable
}

func (e *FailoverPolicyExecutor) isResourceHealthy(resourceID string, checkInterval int) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	lastCheck, exists := e.lastCheckTime[resourceID]
	if !exists || time.Since(lastCheck) > time.Duration(checkInterval)*time.Second {
		// 需要检查，这里简化处理，假设资源是健康的
		// 实际应该进行健康检查
		logger.Debug("Resource health check (needs check)", logger.F("component", "service"), logger.F("resource_id", resourceID), logger.F("last_check", lastCheck), logger.F("check_interval", checkInterval))
		return true
	}

	isHealthy := e.resourceStatus[resourceID]
	logger.Debug("Resource health check result", logger.F("component", "service"), logger.F("resource_id", resourceID), logger.F("is_healthy", isHealthy), logger.F("last_check", lastCheck))
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
	logger.Debug("Getting policy executor", logger.F("component", "service"), logger.F("policy_type", policyType))

	switch policyType {
	case "random":
		return &RandomPolicyExecutor{}
	case "round_robin":
		return f.roundRobinExecutor
	case "weighted":
		return &WeightedPolicyExecutor{}
	case "model_match":
		return &ModelMatchPolicyExecutor{}
	case "regex_match":
		return &RegexMatchPolicyExecutor{}
	case "priority":
		return &PriorityPolicyExecutor{}
	case "failover":
		return f.failoverExecutor
	default:
		logger.Warn("Unknown policy type, using random executor", logger.F("component", "service"), logger.F("policy_type", policyType))
		return &RandomPolicyExecutor{} // 默认使用随机策略
	}
}

// filterResources 过滤资源，并验证资源是否存在
// 如果resourceIDs不为空，只返回在resourceIDs中且存在于resources列表中的资源
// 如果策略参数中配置的资源ID不存在，会记录警告日志
func filterResources(resources []*storage.LLMResource, resourceIDs []string, filterByStatus bool) []*storage.LLMResource {
	filtered := make([]*storage.LLMResource, 0)

	// 构建资源映射，用于快速查找
	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		resourceMap[r.ID] = r
	}

	// 如果指定了资源ID列表，验证每个ID是否存在
	if len(resourceIDs) > 0 {
		missingResources := make([]string, 0)
		for _, id := range resourceIDs {
			if _, exists := resourceMap[id]; !exists {
				missingResources = append(missingResources, id)
			}
		}
		if len(missingResources) > 0 {
			logger.Warn("Some resources configured in policy are not found", logger.F("component", "service"), logger.F("missing_resource_ids", missingResources), logger.F("total_configured", len(resourceIDs)), logger.F("available_resources", len(resources)))
		}
	}

	// 过滤资源
	for _, r := range resources {
		// 状态过滤
		if filterByStatus && r.Status != "active" {
			continue
		}

		// ID过滤：如果指定了resourceIDs，只返回匹配的资源
		if len(resourceIDs) > 0 {
			found := false
			for _, id := range resourceIDs {
				if r.ID == id {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, r)
	}

	return filtered
}

// matchPattern 匹配模式（支持通配符）
func matchPattern(text, pattern string) bool {
	// 简单的通配符匹配：* 匹配任意字符
	if strings.Contains(pattern, "*") {
		regexPattern := "^" + strings.ReplaceAll(strings.ReplaceAll(pattern, ".", "\\."), "*", ".*") + "$"
		matched, _ := regexp.MatchString(regexPattern, text)
		logger.Debug("Pattern match result (wildcard)", logger.F("component", "service"), logger.F("text", text), logger.F("pattern", pattern), logger.F("regex_pattern", regexPattern), logger.F("matched", matched))
		return matched
	}
	matched := text == pattern
	logger.Debug("Pattern match result (exact)", logger.F("component", "service"), logger.F("text", text), logger.F("pattern", pattern), logger.F("matched", matched))
	return matched
}
