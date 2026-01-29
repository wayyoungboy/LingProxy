package storage

// Storage 存储接口
type Storage interface {
	// User methods
	CreateUser(user *User) error
	GetUser(id string) (*User, error)
	GetUserByAPIKey(apiKey string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
	ListUsers() ([]*User, error)

	// LLMResource methods
	CreateLLMResource(resource *LLMResource) error
	GetLLMResource(id string) (*LLMResource, error)
	UpdateLLMResource(resource *LLMResource) error
	DeleteLLMResource(id string) error
	ListLLMResources() ([]*LLMResource, error)

	// Model methods
	CreateModel(model *Model) error
	GetModel(id string) (*Model, error)
	UpdateModel(model *Model) error
	DeleteModel(id string) error
	ListModels() ([]*Model, error)
	ListModelsByLLMResource(llmResourceID string) ([]*Model, error)

	// Endpoint methods
	CreateEndpoint(endpoint *Endpoint) error
	GetEndpoint(id string) (*Endpoint, error)
	ListEndpoints() ([]*Endpoint, error)

	// Request methods
	CreateRequest(request *Request) error
	GetRequest(id string) (*Request, error)
	ListRequests(limit int) ([]*Request, error)

	// Response methods
	CreateResponse(response *Response) error
	GetResponse(id string) (*Response, error)

	// Quota methods
	CreateQuota(quota *Quota) error
	GetQuota(id string) (*Quota, error)
	GetQuotaByUserID(userID string) (*Quota, error)
	UpdateQuota(quota *Quota) error

	// Statistics methods
	CreateStatistics(stats *Statistics) error
	GetStatistics(id string) (*Statistics, error)
	GetStatisticsByUserID(userID string) (*Statistics, error)
	UpdateStatistics(stats *Statistics) error
}
