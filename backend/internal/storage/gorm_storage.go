package storage

import (
	"time"

	"gorm.io/gorm"
)

// GormStorage GORM存储实现
type GormStorage struct {
	db *gorm.DB
}

// NewGormStorage 创建新的GORM存储实例
func NewGormStorage(db *gorm.DB) *GormStorage {
	// 自动迁移数据库表结构
	db.AutoMigrate(
		&User{},
		&APIKey{},
		&PolicyTemplate{},
		&Policy{},
		&LLMResource{},
		&Model{},
		&Endpoint{},
		&Request{},
		&Response{},
		&Quota{},
		&Statistics{},
	)

	return &GormStorage{
		db: db,
	}
}

// User methods
func (g *GormStorage) CreateUser(user *User) error {
	user.ID = generateID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return g.db.Create(user).Error
}

func (g *GormStorage) GetUser(id string) (*User, error) {
	var user User
	if err := g.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (g *GormStorage) GetUserByAPIKey(apiKey string) (*User, error) {
	var user User
	if err := g.db.Where("api_key = ?", apiKey).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (g *GormStorage) UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()
	if err := g.db.Save(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (g *GormStorage) DeleteUser(id string) error {
	if err := g.db.Where("id = ?", id).Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormStorage) ListUsers() ([]*User, error) {
	var users []*User
	if err := g.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// LLMResource methods
func (g *GormStorage) CreateLLMResource(resource *LLMResource) error {
	resource.ID = generateID()
	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()
	return g.db.Create(resource).Error
}

func (g *GormStorage) GetLLMResource(id string) (*LLMResource, error) {
	var resource LLMResource
	if err := g.db.Where("id = ?", id).First(&resource).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &resource, nil
}

func (g *GormStorage) UpdateLLMResource(resource *LLMResource) error {
	resource.UpdatedAt = time.Now()
	if err := g.db.Save(resource).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (g *GormStorage) DeleteLLMResource(id string) error {
	if err := g.db.Where("id = ?", id).Delete(&LLMResource{}).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormStorage) ListLLMResources() ([]*LLMResource, error) {
	var resources []*LLMResource
	if err := g.db.Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

// Model methods
func (g *GormStorage) CreateModel(model *Model) error {
	model.ID = generateID()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return g.db.Create(model).Error
}

func (g *GormStorage) GetModel(id string) (*Model, error) {
	var model Model
	if err := g.db.Where("id = ?", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &model, nil
}

func (g *GormStorage) UpdateModel(model *Model) error {
	model.UpdatedAt = time.Now()
	if err := g.db.Save(model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (g *GormStorage) DeleteModel(id string) error {
	if err := g.db.Where("id = ?", id).Delete(&Model{}).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormStorage) ListModels() ([]*Model, error) {
	var models []*Model
	if err := g.db.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (g *GormStorage) ListModelsByLLMResource(llmResourceID string) ([]*Model, error) {
	var models []*Model
	if err := g.db.Where("llm_resource_id = ?", llmResourceID).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

// Endpoint methods
func (g *GormStorage) CreateEndpoint(endpoint *Endpoint) error {
	endpoint.ID = generateID()
	endpoint.CreatedAt = time.Now()
	endpoint.UpdatedAt = time.Now()
	return g.db.Create(endpoint).Error
}

func (g *GormStorage) GetEndpoint(id string) (*Endpoint, error) {
	var endpoint Endpoint
	if err := g.db.Where("id = ?", id).First(&endpoint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &endpoint, nil
}

func (g *GormStorage) UpdateEndpoint(endpoint *Endpoint) error {
	endpoint.UpdatedAt = time.Now()
	if err := g.db.Save(endpoint).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (g *GormStorage) DeleteEndpoint(id string) error {
	if err := g.db.Where("id = ?", id).Delete(&Endpoint{}).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormStorage) ListEndpoints() ([]*Endpoint, error) {
	var endpoints []*Endpoint
	if err := g.db.Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}

// Request methods
func (g *GormStorage) CreateRequest(request *Request) error {
	request.ID = generateID()
	request.CreatedAt = time.Now()
	return g.db.Create(request).Error
}

func (g *GormStorage) GetRequest(id string) (*Request, error) {
	var request Request
	if err := g.db.Where("id = ?", id).First(&request).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &request, nil
}

func (g *GormStorage) ListRequests(params *RequestQueryParams) ([]*Request, error) {
	var requests []*Request
	query := g.db.Model(&Request{})
	
	// 时间范围过滤
	if params.StartTime != nil {
		query = query.Where("created_at >= ?", *params.StartTime)
	}
	if params.EndTime != nil {
		query = query.Where("created_at <= ?", *params.EndTime)
	}
	
	// 请求路径过滤（支持模糊匹配）
	if params.Endpoint != "" {
		query = query.Where("endpoint LIKE ?", "%"+params.Endpoint+"%")
	}
	
	// 状态过滤
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	
	// 排序和分页
	query = query.Order("created_at desc")
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	
	if err := query.Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// Response methods
func (g *GormStorage) CreateResponse(response *Response) error {
	response.ID = generateID()
	response.CreatedAt = time.Now()
	return g.db.Create(response).Error
}

func (g *GormStorage) GetResponse(id string) (*Response, error) {
	var response Response
	if err := g.db.Where("id = ?", id).First(&response).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &response, nil
}

// Quota methods
func (g *GormStorage) CreateQuota(quota *Quota) error {
	quota.ID = generateID()
	quota.CreatedAt = time.Now()
	quota.UpdatedAt = time.Now()
	return g.db.Create(quota).Error
}

func (g *GormStorage) GetQuota(id string) (*Quota, error) {
	var quota Quota
	if err := g.db.Where("id = ?", id).First(&quota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &quota, nil
}

func (g *GormStorage) GetQuotaByUserID(userID string) (*Quota, error) {
	var quota Quota
	if err := g.db.Where("user_id = ?", userID).First(&quota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &quota, nil
}

func (g *GormStorage) UpdateQuota(quota *Quota) error {
	quota.UpdatedAt = time.Now()
	if err := g.db.Save(quota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// Statistics methods
func (g *GormStorage) CreateStatistics(stats *Statistics) error {
	stats.ID = generateID()
	stats.CreatedAt = time.Now()
	stats.UpdatedAt = time.Now()
	return g.db.Create(stats).Error
}

func (g *GormStorage) GetStatistics(id string) (*Statistics, error) {
	var stats Statistics
	if err := g.db.Where("id = ?", id).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &stats, nil
}

func (g *GormStorage) GetStatisticsByUserID(userID string) (*Statistics, error) {
	var stats Statistics
	if err := g.db.Where("user_id = ?", userID).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &stats, nil
}

func (g *GormStorage) UpdateStatistics(stats *Statistics) error {
	stats.UpdatedAt = time.Now()
	if err := g.db.Save(stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// API Key methods
func (g *GormStorage) CreateAPIKey(apiKey *APIKey) error {
	apiKey.ID = generateID()
	apiKey.CreatedAt = time.Now()
	apiKey.UpdatedAt = time.Now()
	return g.db.Create(apiKey).Error
}

func (g *GormStorage) GetAPIKey(id string) (*APIKey, error) {
	var apiKey APIKey
	if err := g.db.Where("id = ?", id).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &apiKey, nil
}

func (g *GormStorage) GetAPIKeyByValue(apiKeyValue string) (*APIKey, error) {
	var apiKey APIKey
	if err := g.db.Where("api_key = ?", apiKeyValue).First(&apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &apiKey, nil
}

func (g *GormStorage) UpdateAPIKey(apiKey *APIKey) error {
	apiKey.UpdatedAt = time.Now()
	if err := g.db.Save(apiKey).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (g *GormStorage) DeleteAPIKey(id string) error {
	if err := g.db.Where("id = ?", id).Delete(&APIKey{}).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormStorage) ListAPIKeys() ([]*APIKey, error) {
	var apiKeys []*APIKey
	if err := g.db.Find(&apiKeys).Error; err != nil {
		return nil, err
	}
	return apiKeys, nil
}

// 保持向后兼容的方法别名
// Deprecated: 使用 CreateAPIKey 代替
func (g *GormStorage) CreateToken(token *APIKey) error {
	return g.CreateAPIKey(token)
}

// Deprecated: 使用 GetAPIKey 代替
func (g *GormStorage) GetToken(id string) (*APIKey, error) {
	return g.GetAPIKey(id)
}

// Deprecated: 使用 GetAPIKeyByValue 代替
func (g *GormStorage) GetTokenByValue(tokenValue string) (*APIKey, error) {
	return g.GetAPIKeyByValue(tokenValue)
}

// Deprecated: 使用 UpdateAPIKey 代替
func (g *GormStorage) UpdateToken(token *APIKey) error {
	return g.UpdateAPIKey(token)
}

// Deprecated: 使用 DeleteAPIKey 代替
func (g *GormStorage) DeleteToken(id string) error {
	return g.DeleteAPIKey(id)
}

// Deprecated: 使用 ListAPIKeys 代替
func (g *GormStorage) ListTokens() ([]*APIKey, error) {
	return g.ListAPIKeys()
}

// PolicyTemplate methods
func (g *GormStorage) CreatePolicyTemplate(template *PolicyTemplate) error {
	template.ID = generateID()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	return g.db.Create(template).Error
}

func (g *GormStorage) GetPolicyTemplate(id string) (*PolicyTemplate, error) {
	var template PolicyTemplate
	if err := g.db.Where("id = ?", id).First(&template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &template, nil
}

func (g *GormStorage) GetPolicyTemplateByType(policyType string) (*PolicyTemplate, error) {
	var template PolicyTemplate
	if err := g.db.Where("type = ?", policyType).First(&template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &template, nil
}

func (g *GormStorage) UpdatePolicyTemplate(template *PolicyTemplate) error {
	template.UpdatedAt = time.Now()
	if err := g.db.Save(template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (g *GormStorage) DeletePolicyTemplate(id string) error {
	if err := g.db.Where("id = ?", id).Delete(&PolicyTemplate{}).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormStorage) ListPolicyTemplates() ([]*PolicyTemplate, error) {
	var templates []*PolicyTemplate
	if err := g.db.Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// Policy methods
func (g *GormStorage) CreatePolicy(policy *Policy) error {
	policy.ID = generateID()
	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()
	return g.db.Create(policy).Error
}

func (g *GormStorage) GetPolicy(id string) (*Policy, error) {
	var policy Policy
	if err := g.db.Where("id = ?", id).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &policy, nil
}

func (g *GormStorage) UpdatePolicy(policy *Policy) error {
	policy.UpdatedAt = time.Now()
	if err := g.db.Save(policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (g *GormStorage) DeletePolicy(id string) error {
	if err := g.db.Where("id = ?", id).Delete(&Policy{}).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormStorage) ListPolicies() ([]*Policy, error) {
	var policies []*Policy
	if err := g.db.Find(&policies).Error; err != nil {
		return nil, err
	}
	return policies, nil
}
