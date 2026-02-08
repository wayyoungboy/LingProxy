package storage

type Storage interface {
	CreateUser(user *User) error
	GetUser(id string) (*User, error)
	GetUserByAPIKey(apiKey string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
	ListUsers() ([]*User, error)

	CreateLLMResource(resource *LLMResource) error
	GetLLMResource(id string) (*LLMResource, error)
	UpdateLLMResource(resource *LLMResource) error
	DeleteLLMResource(id string) error
	ListLLMResources() ([]*LLMResource, error)

	CreateModel(model *Model) error
	GetModel(id string) (*Model, error)
	UpdateModel(model *Model) error
	DeleteModel(id string) error
	ListModels() ([]*Model, error)
	ListModelsByLLMResource(llmResourceID string) ([]*Model, error)

	CreateEndpoint(endpoint *Endpoint) error
	GetEndpoint(id string) (*Endpoint, error)
	UpdateEndpoint(endpoint *Endpoint) error
	DeleteEndpoint(id string) error
	ListEndpoints() ([]*Endpoint, error)

	CreateRequest(request *Request) error
	GetRequest(id string) (*Request, error)
	ListRequests(params *RequestQueryParams) ([]*Request, error)

	// API Key methods
	CreateAPIKey(apiKey *APIKey) error
	GetAPIKey(id string) (*APIKey, error)
	GetAPIKeyByValue(apiKeyValue string) (*APIKey, error)
	UpdateAPIKey(apiKey *APIKey) error
	DeleteAPIKey(id string) error
	ListAPIKeys() ([]*APIKey, error)

	// 保持向后兼容的方法别名
	// Deprecated: 使用 CreateAPIKey 代替
	CreateToken(token *APIKey) error
	// Deprecated: 使用 GetAPIKey 代替
	GetToken(id string) (*APIKey, error)
	// Deprecated: 使用 GetAPIKeyByValue 代替
	GetTokenByValue(tokenValue string) (*APIKey, error)
	// Deprecated: 使用 UpdateAPIKey 代替
	UpdateToken(token *APIKey) error
	// Deprecated: 使用 DeleteAPIKey 代替
	DeleteToken(id string) error
	// Deprecated: 使用 ListAPIKeys 代替
	ListTokens() ([]*APIKey, error)

	// PolicyTemplate methods
	CreatePolicyTemplate(template *PolicyTemplate) error
	GetPolicyTemplate(id string) (*PolicyTemplate, error)
	GetPolicyTemplateByType(policyType string) (*PolicyTemplate, error)
	UpdatePolicyTemplate(template *PolicyTemplate) error
	DeletePolicyTemplate(id string) error
	ListPolicyTemplates() ([]*PolicyTemplate, error)

	// Policy methods
	CreatePolicy(policy *Policy) error
	GetPolicy(id string) (*Policy, error)
	UpdatePolicy(policy *Policy) error
	DeletePolicy(id string) error
	ListPolicies() ([]*Policy, error)
}
