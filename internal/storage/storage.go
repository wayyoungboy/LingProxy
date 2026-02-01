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
	ListRequests(limit int) ([]*Request, error)
}
