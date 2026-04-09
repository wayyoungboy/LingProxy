package service

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	mathrand "math/rand"
	"sync"
	"time"

	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// init 初始化随机数生成器，使用crypto/rand作为种子源以确保真正的随机性
func init() {
	var seed int64
	if err := binary.Read(rand.Reader, binary.BigEndian, &seed); err != nil {
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

	filtered := filterResources(resources, params.Resources, params.FilterByStatus)
	logger.Debug("Random policy filtered resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("original_count", len(resources)), logger.F("filtered_count", len(filtered)))

	if len(filtered) == 0 {
		logger.Warn("Random policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
		return nil, ErrNoResourcesAvailable
	}

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

	filtered := filterResources(resources, params.Resources, params.FilterByStatus)
	logger.Debug("Round-robin policy filtered resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("original_count", len(resources)), logger.F("filtered_count", len(filtered)))

	if len(filtered) == 0 {
		logger.Warn("Round-robin policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
		return nil, ErrNoResourcesAvailable
	}

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

	allResourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		allResourceMap[r.ID] = r
	}

	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		if params.FilterByStatus && r.Status != "active" {
			continue
		}
		resourceMap[r.ID] = r
	}

	totalWeight := 0
	validResources := make([]struct {
		Resource *storage.LLMResource
		Weight   int
	}, 0)
	missingResources := make([]string, 0)

	for _, wr := range params.Resources {
		if _, exists := allResourceMap[wr.ID]; !exists {
			missingResources = append(missingResources, wr.ID)
			continue
		}
		if r, exists := resourceMap[wr.ID]; exists {
			totalWeight += wr.Weight
			validResources = append(validResources, struct {
				Resource *storage.LLMResource
				Weight   int
			}{Resource: r, Weight: wr.Weight})
		}
	}

	if len(missingResources) > 0 {
		logger.Warn("Weighted policy: some configured resources are not found", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("missing_resource_ids", missingResources), logger.F("total_configured", len(params.Resources)))
	}

	logger.Debug("Weighted policy processed resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("original_count", len(resources)), logger.F("valid_count", len(validResources)), logger.F("total_weight", totalWeight))

	if len(validResources) == 0 {
		logger.Warn("Weighted policy: no available resources", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("model", modelName))
		return nil, ErrNoResourcesAvailable
	}

	random := mathrand.Intn(totalWeight)
	currentWeight := 0

	for _, wr := range validResources {
		currentWeight += wr.Weight
		if random < currentWeight {
			logger.Debug("Weighted policy selected resource", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", wr.Resource.ID), logger.F("resource_name", wr.Resource.Name), logger.F("weight", wr.Weight))
			return wr.Resource, nil
		}
	}

	selected := validResources[0].Resource
	logger.Debug("Weighted policy fallback selection", logger.F("component", "service"), logger.F("policy_id", policy.ID), logger.F("resource_id", selected.ID), logger.F("resource_name", selected.Name))
	return selected, nil
}

// PolicyExecutorFactory 策略执行器工厂
type PolicyExecutorFactory struct {
	roundRobinExecutor *RoundRobinPolicyExecutor
}

func NewPolicyExecutorFactory() *PolicyExecutorFactory {
	return &PolicyExecutorFactory{
		roundRobinExecutor: NewRoundRobinPolicyExecutor(),
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
	default:
		logger.Warn("Unknown policy type, using random executor", logger.F("component", "service"), logger.F("policy_type", policyType))
		return &RandomPolicyExecutor{}
	}
}

// filterResources 过滤资源
// 如果resourceIDs不为空，只返回在resourceIDs中且存在于resources列表中的资源
func filterResources(resources []*storage.LLMResource, resourceIDs []string, filterByStatus bool) []*storage.LLMResource {
	filtered := make([]*storage.LLMResource, 0)

	resourceMap := make(map[string]*storage.LLMResource)
	for _, r := range resources {
		resourceMap[r.ID] = r
	}

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

	for _, r := range resources {
		if filterByStatus && r.Status != "active" {
			continue
		}

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
