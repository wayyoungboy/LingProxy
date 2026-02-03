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
	var params struct {
		Resources      []string `json:"resources"`
		FilterByStatus bool     `json:"filter_by_status"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 过滤资源
	filtered := filterResources(resources, params.Resources, params.FilterByStatus)
	if len(filtered) == 0 {
		return nil, ErrNoResourcesAvailable
	}

	// 随机选择（Go 1.20+ 全局随机数生成器已自动初始化）
	return filtered[rand.Intn(len(filtered))], nil
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
	var params struct {
		Resources      []string `json:"resources"`
		FilterByStatus bool     `json:"filter_by_status"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPolicyParams, err)
	}

	// 过滤资源
	filtered := filterResources(resources, params.Resources, params.FilterByStatus)
	if len(filtered) == 0 {
		return nil, ErrNoResourcesAvailable
	}

	// 轮询选择
	e.mu.Lock()
	index := e.currentIndex[policy.ID]
	e.currentIndex[policy.ID] = (index + 1) % len(filtered)
	e.mu.Unlock()

	return filtered[index], nil
}

// WeightedPolicyExecutor 加权负载均衡策略执行器
type WeightedPolicyExecutor struct{}

func (e *WeightedPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	var params struct {
		Resources []struct {
			ID     string `json:"id"`
			Weight int    `json:"weight"`
		} `json:"resources"`
		FilterByStatus bool `json:"filter_by_status"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
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

	if len(validResources) == 0 {
		return nil, ErrNoResourcesAvailable
	}

	// 根据权重随机选择（Go 1.20+ 全局随机数生成器已自动初始化）
	random := rand.Intn(totalWeight)
	currentWeight := 0

	for _, wr := range validResources {
		currentWeight += wr.Weight
		if random < currentWeight {
			return wr.Resource, nil
		}
	}

	return validResources[0].Resource, nil
}

// ModelMatchPolicyExecutor 模型名匹配策略执行器
type ModelMatchPolicyExecutor struct{}

func (e *ModelMatchPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	var params struct {
		Mappings []struct {
			ModelPattern string `json:"model_pattern"`
			ResourceID   string `json:"resource_id"`
		} `json:"mappings"`
		DefaultResourceID string `json:"default_resource_id"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
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
			if resource, exists := resourceMap[mapping.ResourceID]; exists {
				return resource, nil
			}
		}
	}

	// 使用默认资源
	if params.DefaultResourceID != "" {
		if resource, exists := resourceMap[params.DefaultResourceID]; exists {
			return resource, nil
		}
	}

	return nil, ErrNoResourcesAvailable
}

// RegexMatchPolicyExecutor 正则匹配策略执行器
type RegexMatchPolicyExecutor struct{}

func (e *RegexMatchPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	var params struct {
		Rules []struct {
			Pattern    string `json:"pattern"`
			ResourceID string `json:"resource_id"`
		} `json:"rules"`
		DefaultResourceID string `json:"default_resource_id"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
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
			continue // 跳过无效的正则表达式
		}
		if matched {
			if resource, exists := resourceMap[rule.ResourceID]; exists {
				return resource, nil
			}
		}
	}

	// 使用默认资源
	if params.DefaultResourceID != "" {
		if resource, exists := resourceMap[params.DefaultResourceID]; exists {
			return resource, nil
		}
	}

	return nil, ErrNoResourcesAvailable
}

// PriorityPolicyExecutor 优先级策略执行器
type PriorityPolicyExecutor struct{}

func (e *PriorityPolicyExecutor) Execute(policy *storage.Policy, modelName string, resources []*storage.LLMResource) (*storage.LLMResource, error) {
	var params struct {
		Resources []struct {
			ID       string `json:"id"`
			Priority int    `json:"priority"`
		} `json:"resources"`
		FallbackEnabled bool `json:"fallback_enabled"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
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
			return resource, nil
		}
	}

	if params.FallbackEnabled {
		// 降级：使用任何可用的资源
		for _, r := range resources {
			if r.Status == "active" {
				return r, nil
			}
		}
	}

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
	var params struct {
		PrimaryResourceID   string   `json:"primary_resource_id"`
		FallbackResources   []string `json:"fallback_resources"`
		HealthCheckEnabled  bool     `json:"health_check_enabled"`
		HealthCheckInterval int      `json:"health_check_interval"`
	}

	if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
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
		if !params.HealthCheckEnabled || e.isResourceHealthy(params.PrimaryResourceID, params.HealthCheckInterval) {
			return primary, nil
		}
	}

	// 使用备用资源
	for _, fallbackID := range params.FallbackResources {
		if fallback, exists := resourceMap[fallbackID]; exists {
			if !params.HealthCheckEnabled || e.isResourceHealthy(fallbackID, params.HealthCheckInterval) {
				return fallback, nil
			}
		}
	}

	return nil, ErrNoResourcesAvailable
}

func (e *FailoverPolicyExecutor) isResourceHealthy(resourceID string, checkInterval int) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	lastCheck, exists := e.lastCheckTime[resourceID]
	if !exists || time.Since(lastCheck) > time.Duration(checkInterval)*time.Second {
		// 需要检查，这里简化处理，假设资源是健康的
		// 实际应该进行健康检查
		return true
	}

	return e.resourceStatus[resourceID]
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
		return &RandomPolicyExecutor{} // 默认使用随机策略
	}
}

// filterResources 过滤资源
func filterResources(resources []*storage.LLMResource, resourceIDs []string, filterByStatus bool) []*storage.LLMResource {
	filtered := make([]*storage.LLMResource, 0)

	for _, r := range resources {
		// 状态过滤
		if filterByStatus && r.Status != "active" {
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
		return matched
	}
	return text == pattern
}
