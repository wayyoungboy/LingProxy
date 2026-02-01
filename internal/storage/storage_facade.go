package storage

type StorageFacade struct {
	storage Storage
}

func NewStorageFacade(storage Storage) *StorageFacade {
	return &StorageFacade{
		storage: storage,
	}
}

func (f *StorageFacade) CreateUser(user *User) error {
	return f.storage.CreateUser(user)
}

func (f *StorageFacade) GetUser(id string) (*User, error) {
	return f.storage.GetUser(id)
}

func (f *StorageFacade) GetUserByAPIKey(apiKey string) (*User, error) {
	return f.storage.GetUserByAPIKey(apiKey)
}

func (f *StorageFacade) UpdateUser(user *User) error {
	return f.storage.UpdateUser(user)
}

func (f *StorageFacade) DeleteUser(id string) error {
	return f.storage.DeleteUser(id)
}

func (f *StorageFacade) ListUsers() ([]*User, error) {
	return f.storage.ListUsers()
}

func (f *StorageFacade) CreateLLMResource(resource *LLMResource) error {
	return f.storage.CreateLLMResource(resource)
}

func (f *StorageFacade) GetLLMResource(id string) (*LLMResource, error) {
	return f.storage.GetLLMResource(id)
}

func (f *StorageFacade) UpdateLLMResource(resource *LLMResource) error {
	return f.storage.UpdateLLMResource(resource)
}

func (f *StorageFacade) DeleteLLMResource(id string) error {
	return f.storage.DeleteLLMResource(id)
}

func (f *StorageFacade) ListLLMResources() ([]*LLMResource, error) {
	return f.storage.ListLLMResources()
}

func (f *StorageFacade) CreateEndpoint(endpoint *Endpoint) error {
	return f.storage.CreateEndpoint(endpoint)
}

func (f *StorageFacade) GetEndpoint(id string) (*Endpoint, error) {
	return f.storage.GetEndpoint(id)
}

func (f *StorageFacade) UpdateEndpoint(endpoint *Endpoint) error {
	return f.storage.UpdateEndpoint(endpoint)
}

func (f *StorageFacade) DeleteEndpoint(id string) error {
	return f.storage.DeleteEndpoint(id)
}

func (f *StorageFacade) ListEndpoints() ([]*Endpoint, error) {
	return f.storage.ListEndpoints()
}

func (f *StorageFacade) CreateRequest(request *Request) error {
	return f.storage.CreateRequest(request)
}

func (f *StorageFacade) GetRequest(id string) (*Request, error) {
	return f.storage.GetRequest(id)
}

func (f *StorageFacade) ListRequests(limit int) ([]*Request, error) {
	return f.storage.ListRequests(limit)
}

func (f *StorageFacade) CreateModel(model *Model) error {
	return f.storage.CreateModel(model)
}

func (f *StorageFacade) GetModel(id string) (*Model, error) {
	return f.storage.GetModel(id)
}

func (f *StorageFacade) UpdateModel(model *Model) error {
	return f.storage.UpdateModel(model)
}

func (f *StorageFacade) DeleteModel(id string) error {
	return f.storage.DeleteModel(id)
}

func (f *StorageFacade) ListModels() ([]*Model, error) {
	return f.storage.ListModels()
}

func (f *StorageFacade) ListModelsByLLMResource(llmResourceID string) ([]*Model, error) {
	return f.storage.ListModelsByLLMResource(llmResourceID)
}
