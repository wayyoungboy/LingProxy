package storage

import (
	"testing"
)

// TestMemoryStorage_CreateUser 测试创建用户
func TestMemoryStorage_CreateUser(t *testing.T) {
	storage := NewMemoryStorage()
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		APIKey:   "test-api-key",
		Status:   "active",
	}

	err := storage.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser failed: %v", err)
	}

	if user.ID == "" {
		t.Error("User ID should not be empty")
	}

	if user.CreatedAt.IsZero() {
		t.Error("User CreatedAt should not be zero")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("User UpdatedAt should not be zero")
	}
}

// TestMemoryStorage_GetUser 测试获取用户
func TestMemoryStorage_GetUser(t *testing.T) {
	storage := NewMemoryStorage()
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		APIKey:   "test-api-key",
		Status:   "active",
	}

	// 创建用户
	if err := storage.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// 获取用户
	retrievedUser, err := storage.GetUser(user.ID)
	if err != nil {
		t.Errorf("GetUser failed: %v", err)
	}

	if retrievedUser == nil {
		t.Error("Retrieved user should not be nil")
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("User ID mismatch: expected %s, got %s", user.ID, retrievedUser.ID)
	}

	if retrievedUser.Username != user.Username {
		t.Errorf("User Username mismatch: expected %s, got %s", user.Username, retrievedUser.Username)
	}
}

// TestMemoryStorage_UpdateUser 测试更新用户
func TestMemoryStorage_UpdateUser(t *testing.T) {
	storage := NewMemoryStorage()
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		APIKey:   "test-api-key",
		Status:   "active",
	}

	// 创建用户
	if err := storage.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// 更新用户
	user.Username = "updateduser"
	user.Email = "updated@example.com"

	if err := storage.UpdateUser(user); err != nil {
		t.Errorf("UpdateUser failed: %v", err)
	}

	// 验证更新
	retrievedUser, err := storage.GetUser(user.ID)
	if err != nil {
		t.Errorf("GetUser failed: %v", err)
	}

	if retrievedUser.Username != "updateduser" {
		t.Errorf("User Username not updated: expected updateduser, got %s", retrievedUser.Username)
	}

	if retrievedUser.Email != "updated@example.com" {
		t.Errorf("User Email not updated: expected updated@example.com, got %s", retrievedUser.Email)
	}

	if retrievedUser.UpdatedAt.Before(user.UpdatedAt) {
		t.Error("User UpdatedAt should be updated")
	}
}

// TestMemoryStorage_DeleteUser 测试删除用户
func TestMemoryStorage_DeleteUser(t *testing.T) {
	storage := NewMemoryStorage()
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		APIKey:   "test-api-key",
		Status:   "active",
	}

	// 创建用户
	if err := storage.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// 删除用户
	if err := storage.DeleteUser(user.ID); err != nil {
		t.Errorf("DeleteUser failed: %v", err)
	}

	// 验证删除
	retrievedUser, err := storage.GetUser(user.ID)
	if err == nil {
		t.Error("GetUser should return error for deleted user")
	}

	if retrievedUser != nil {
		t.Error("Retrieved user should be nil for deleted user")
	}
}

// TestMemoryStorage_ListUsers 测试列出所有用户
func TestMemoryStorage_ListUsers(t *testing.T) {
	storage := NewMemoryStorage()

	// 创建多个用户
	createdIDs := make([]string, 3)
	for i := 0; i < 3; i++ {
		user := &User{
			Username: "user" + string(rune('0'+i)),
			Email:    "user" + string(rune('0'+i)) + "@example.com",
			APIKey:   "api-key-" + string(rune('0'+i)),
			Status:   "active",
		}

		if err := storage.CreateUser(user); err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}
		createdIDs[i] = user.ID
		t.Logf("Created user %d with ID: %s", i, user.ID)
	}

	// 列出用户
	users, err := storage.ListUsers()
	if err != nil {
		t.Errorf("ListUsers failed: %v", err)
	}

	t.Logf("Retrieved %d users", len(users))
	for i, u := range users {
		t.Logf("User %d: ID=%s, Username=%s", i, u.ID, u.Username)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}

// TestMemoryStorage_CreateLLMResource 测试创建LLM资源
func TestMemoryStorage_CreateLLMResource(t *testing.T) {
	storage := NewMemoryStorage()
	resource := &LLMResource{
		Name:    "test-resource",
		Type:    "openai",
		Model:   "gpt-3.5-turbo",
		BaseURL: "https://api.openai.com/v1",
		APIKey:  "test-api-key",
		Status:  "active",
	}

	err := storage.CreateLLMResource(resource)
	if err != nil {
		t.Errorf("CreateLLMResource failed: %v", err)
	}

	if resource.ID == "" {
		t.Error("LLMResource ID should not be empty")
	}

	if resource.CreatedAt.IsZero() {
		t.Error("LLMResource CreatedAt should not be zero")
	}

	if resource.UpdatedAt.IsZero() {
		t.Error("LLMResource UpdatedAt should not be zero")
	}
}

// TestMemoryStorage_GetLLMResource 测试获取LLM资源
func TestMemoryStorage_GetLLMResource(t *testing.T) {
	storage := NewMemoryStorage()
	resource := &LLMResource{
		Name:    "test-resource",
		Type:    "openai",
		Model:   "gpt-3.5-turbo",
		BaseURL: "https://api.openai.com/v1",
		APIKey:  "test-api-key",
		Status:  "active",
	}

	// 创建资源
	if err := storage.CreateLLMResource(resource); err != nil {
		t.Fatalf("CreateLLMResource failed: %v", err)
	}

	// 获取资源
	retrievedResource, err := storage.GetLLMResource(resource.ID)
	if err != nil {
		t.Errorf("GetLLMResource failed: %v", err)
	}

	if retrievedResource == nil {
		t.Error("Retrieved resource should not be nil")
	}

	if retrievedResource.ID != resource.ID {
		t.Errorf("LLMResource ID mismatch: expected %s, got %s", resource.ID, retrievedResource.ID)
	}

	if retrievedResource.Name != resource.Name {
		t.Errorf("LLMResource Name mismatch: expected %s, got %s", resource.Name, retrievedResource.Name)
	}
}

// TestMemoryStorage_ListLLMResources 测试列出所有LLM资源
func TestMemoryStorage_ListLLMResources(t *testing.T) {
	storage := NewMemoryStorage()

	// 创建多个资源
	createdIDs := make([]string, 3)
	for i := 0; i < 3; i++ {
		resource := &LLMResource{
			Name:    "resource" + string(rune('0'+i)),
			Type:    "openai",
			Model:   "gpt-3.5-turbo",
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "api-key-" + string(rune('0'+i)),
			Status:  "active",
		}

		if err := storage.CreateLLMResource(resource); err != nil {
			t.Fatalf("CreateLLMResource failed: %v", err)
		}
		createdIDs[i] = resource.ID
		t.Logf("Created resource %d with ID: %s", i, resource.ID)
	}

	// 列出资源
	resources, err := storage.ListLLMResources()
	if err != nil {
		t.Errorf("ListLLMResources failed: %v", err)
	}

	t.Logf("Retrieved %d resources", len(resources))
	for i, r := range resources {
		t.Logf("Resource %d: ID=%s, Name=%s", i, r.ID, r.Name)
	}

	if len(resources) != 3 {
		t.Errorf("Expected 3 resources, got %d", len(resources))
	}
}
