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

func (g *GormStorage) ListRequests(limit int) ([]*Request, error) {
	var requests []*Request
	query := g.db.Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
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
