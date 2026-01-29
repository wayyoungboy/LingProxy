package storage

import (
	"fmt"
	"sync"
	"time"
)

// MemoryStorage 内存存储实现
type MemoryStorage struct {
	mu sync.RWMutex

	// 数据存储
	users        map[string]*User
	llmResources map[string]*LLMResource
	models       map[string]*Model
	endpoints    map[string]*Endpoint
	requests     map[string]*Request
	responses    map[string]*Response
	quotas       map[string]*Quota
	statistics   map[string]*Statistics
}

// NewMemoryStorage 创建新的内存存储实例
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users:        make(map[string]*User),
		llmResources: make(map[string]*LLMResource),
		models:       make(map[string]*Model),
		endpoints:    make(map[string]*Endpoint),
		requests:     make(map[string]*Request),
		responses:    make(map[string]*Response),
		quotas:       make(map[string]*Quota),
		statistics:   make(map[string]*Statistics),
	}
}

// User methods
func (m *MemoryStorage) CreateUser(user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.ID = generateID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MemoryStorage) GetUser(id string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) GetUserByAPIKey(apiKey string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.APIKey == apiKey {
			return user, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) UpdateUser(user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.ID]; !exists {
		return ErrNotFound
	}

	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MemoryStorage) DeleteUser(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[id]; !exists {
		return ErrNotFound
	}

	delete(m.users, id)
	return nil
}

func (m *MemoryStorage) ListUsers() ([]*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]*User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

// LLMResource methods
func (m *MemoryStorage) CreateLLMResource(resource *LLMResource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	resource.ID = generateID()
	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()
	m.llmResources[resource.ID] = resource
	return nil
}

func (m *MemoryStorage) GetLLMResource(id string) (*LLMResource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if resource, exists := m.llmResources[id]; exists {
		return resource, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) UpdateLLMResource(resource *LLMResource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.llmResources[resource.ID]; !exists {
		return ErrNotFound
	}

	resource.UpdatedAt = time.Now()
	m.llmResources[resource.ID] = resource
	return nil
}

func (m *MemoryStorage) DeleteLLMResource(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.llmResources[id]; !exists {
		return ErrNotFound
	}

	delete(m.llmResources, id)
	return nil
}

func (m *MemoryStorage) ListLLMResources() ([]*LLMResource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	resources := make([]*LLMResource, 0, len(m.llmResources))
	for _, resource := range m.llmResources {
		resources = append(resources, resource)
	}
	return resources, nil
}

// Model methods
func (m *MemoryStorage) CreateModel(model *Model) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	model.ID = generateID()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	m.models[model.ID] = model
	return nil
}

func (m *MemoryStorage) GetModel(id string) (*Model, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if model, exists := m.models[id]; exists {
		return model, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) UpdateModel(model *Model) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.models[model.ID]; !exists {
		return ErrNotFound
	}

	model.UpdatedAt = time.Now()
	m.models[model.ID] = model
	return nil
}

func (m *MemoryStorage) DeleteModel(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.models[id]; !exists {
		return ErrNotFound
	}

	delete(m.models, id)
	return nil
}

func (m *MemoryStorage) ListModels() ([]*Model, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	models := make([]*Model, 0, len(m.models))
	for _, model := range m.models {
		models = append(models, model)
	}
	return models, nil
}

func (m *MemoryStorage) ListModelsByLLMResource(llmResourceID string) ([]*Model, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	models := make([]*Model, 0)
	for _, model := range m.models {
		if model.LLMResourceID == llmResourceID {
			models = append(models, model)
		}
	}
	return models, nil
}

// Endpoint methods
func (m *MemoryStorage) CreateEndpoint(endpoint *Endpoint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	endpoint.ID = generateID()
	endpoint.CreatedAt = time.Now()
	endpoint.UpdatedAt = time.Now()
	m.endpoints[endpoint.ID] = endpoint
	return nil
}

func (m *MemoryStorage) GetEndpoint(id string) (*Endpoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if endpoint, exists := m.endpoints[id]; exists {
		return endpoint, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) ListEndpoints() ([]*Endpoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	endpoints := make([]*Endpoint, 0, len(m.endpoints))
	for _, endpoint := range m.endpoints {
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, nil
}

// Request methods
func (m *MemoryStorage) CreateRequest(request *Request) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	request.ID = generateID()
	request.CreatedAt = time.Now()
	m.requests[request.ID] = request
	return nil
}

func (m *MemoryStorage) GetRequest(id string) (*Request, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if request, exists := m.requests[id]; exists {
		return request, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) ListRequests(limit int) ([]*Request, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	requests := make([]*Request, 0, len(m.requests))
	for _, request := range m.requests {
		requests = append(requests, request)
	}

	// 简单的分页
	if limit > 0 && len(requests) > limit {
		return requests[:limit], nil
	}
	return requests, nil
}

// Response methods
func (m *MemoryStorage) CreateResponse(response *Response) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	response.ID = generateID()
	response.CreatedAt = time.Now()
	m.responses[response.ID] = response
	return nil
}

func (m *MemoryStorage) GetResponse(id string) (*Response, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if response, exists := m.responses[id]; exists {
		return response, nil
	}
	return nil, ErrNotFound
}

// Quota methods
func (m *MemoryStorage) CreateQuota(quota *Quota) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	quota.ID = generateID()
	quota.CreatedAt = time.Now()
	quota.UpdatedAt = time.Now()
	m.quotas[quota.ID] = quota
	return nil
}

func (m *MemoryStorage) GetQuota(id string) (*Quota, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if quota, exists := m.quotas[id]; exists {
		return quota, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) GetQuotaByUserID(userID string) (*Quota, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, quota := range m.quotas {
		if quota.UserID == userID {
			return quota, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) UpdateQuota(quota *Quota) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.quotas[quota.ID]; !exists {
		return ErrNotFound
	}

	quota.UpdatedAt = time.Now()
	m.quotas[quota.ID] = quota
	return nil
}

// Statistics methods
func (m *MemoryStorage) CreateStatistics(stats *Statistics) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats.ID = generateID()
	stats.CreatedAt = time.Now()
	stats.UpdatedAt = time.Now()
	m.statistics[stats.ID] = stats
	return nil
}

func (m *MemoryStorage) GetStatistics(id string) (*Statistics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if stats, exists := m.statistics[id]; exists {
		return stats, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) GetStatisticsByUserID(userID string) (*Statistics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, stats := range m.statistics {
		if stats.UserID == userID {
			return stats, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MemoryStorage) UpdateStatistics(stats *Statistics) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.statistics[stats.ID]; !exists {
		return ErrNotFound
	}

	stats.UpdatedAt = time.Now()
	m.statistics[stats.ID] = stats
	return nil
}

// 辅助函数
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
