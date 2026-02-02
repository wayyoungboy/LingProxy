package storage

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGormStorage(t *testing.T) {
	// 创建内存SQLite数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化GORM存储
	storage := NewGormStorage(db)

	// 测试用户管理
	t.Run("User Management", func(t *testing.T) {
		// 创建用户
		user := &User{
			Username: "testuser",
			APIKey:   "test-api-key",
			Status:   "active",
		}

		if err := storage.CreateUser(user); err != nil {
			t.Errorf("CreateUser failed: %v", err)
		}

		// 获取用户
		retrievedUser, err := storage.GetUser(user.ID)
		if err != nil {
			t.Errorf("GetUser failed: %v", err)
		}
		if retrievedUser.Username != user.Username {
			t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
		}

		// 通过API Key获取用户
		userByAPIKey, err := storage.GetUserByAPIKey(user.APIKey)
		if err != nil {
			t.Errorf("GetUserByAPIKey failed: %v", err)
		}
		if userByAPIKey.Username != user.Username {
			t.Errorf("Expected username %s, got %s", user.Username, userByAPIKey.Username)
		}

		// 更新用户
		retrievedUser.Username = "updateduser"
		if err := storage.UpdateUser(retrievedUser); err != nil {
			t.Errorf("UpdateUser failed: %v", err)
		}

		updatedUser, err := storage.GetUser(user.ID)
		if err != nil {
			t.Errorf("GetUser after update failed: %v", err)
		}
		if updatedUser.Username != "updateduser" {
			t.Errorf("Expected username updateduser, got %s", updatedUser.Username)
		}

		// 列出用户
		users, err := storage.ListUsers()
		if err != nil {
			t.Errorf("ListUsers failed: %v", err)
		}
		if len(users) == 0 {
			t.Error("Expected at least one user, got none")
		}

		// 删除用户
		if err := storage.DeleteUser(user.ID); err != nil {
			t.Errorf("DeleteUser failed: %v", err)
		}

		// 验证用户已删除
		_, err = storage.GetUser(user.ID)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// 测试LLM资源管理
	t.Run("LLM Resource Management", func(t *testing.T) {
		// 创建LLM资源
		resource := &LLMResource{
			Name:     "Test Resource",
			Type:     "chat",
			Provider: "openai",
			Model:    "gpt-3.5-turbo",
			BaseURL:  "https://api.openai.com/v1",
			APIKey:   "test-api-key",
			Status:   "active",
		}

		if err := storage.CreateLLMResource(resource); err != nil {
			t.Errorf("CreateLLMResource failed: %v", err)
		}

		// 获取LLM资源
		retrievedResource, err := storage.GetLLMResource(resource.ID)
		if err != nil {
			t.Errorf("GetLLMResource failed: %v", err)
		}
		if retrievedResource.Name != resource.Name {
			t.Errorf("Expected name %s, got %s", resource.Name, retrievedResource.Name)
		}

		// 更新LLM资源
		retrievedResource.Name = "Updated Resource"
		if err := storage.UpdateLLMResource(retrievedResource); err != nil {
			t.Errorf("UpdateLLMResource failed: %v", err)
		}

		updatedResource, err := storage.GetLLMResource(resource.ID)
		if err != nil {
			t.Errorf("GetLLMResource after update failed: %v", err)
		}
		if updatedResource.Name != "Updated Resource" {
			t.Errorf("Expected name Updated Resource, got %s", updatedResource.Name)
		}

		// 列出LLM资源
		resources, err := storage.ListLLMResources()
		if err != nil {
			t.Errorf("ListLLMResources failed: %v", err)
		}
		if len(resources) == 0 {
			t.Error("Expected at least one resource, got none")
		}

		// 删除LLM资源
		if err := storage.DeleteLLMResource(resource.ID); err != nil {
			t.Errorf("DeleteLLMResource failed: %v", err)
		}

		// 验证资源已删除
		_, err = storage.GetLLMResource(resource.ID)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// 测试模型管理
	t.Run("Model Management", func(t *testing.T) {
		// 创建模型
		model := &Model{
			Name:          "Test Model",
			LLMResourceID: "test-resource",
			ModelID:       "test-model",
			Type:          "chat",
			Category:      "gpt",
			Version:       "1.0",
			Description:   "Test model",
			Capabilities:  "[\"text-generation\"]",
			Pricing:       "{\"input_token_price\":0.001,\"output_token_price\":0.002,\"currency\":\"USD\"}",
			Limits:        "{\"max_tokens\":4096,\"context_window\":8192}",
			Parameters:    "{\"temperature\":0.7,\"max_tokens\":1000}",
			Features:      "{\"streaming\":true,\"function_calling\":false}",
			Status:        "active",
			Metadata:      "{\"key\":\"value\"}",
		}

		if err := storage.CreateModel(model); err != nil {
			t.Errorf("CreateModel failed: %v", err)
		}

		// 获取模型
		retrievedModel, err := storage.GetModel(model.ID)
		if err != nil {
			t.Errorf("GetModel failed: %v", err)
		}
		if retrievedModel.Name != model.Name {
			t.Errorf("Expected name %s, got %s", model.Name, retrievedModel.Name)
		}

		// 更新模型
		retrievedModel.Name = "Updated Model"
		if err := storage.UpdateModel(retrievedModel); err != nil {
			t.Errorf("UpdateModel failed: %v", err)
		}

		updatedModel, err := storage.GetModel(model.ID)
		if err != nil {
			t.Errorf("GetModel after update failed: %v", err)
		}
		if updatedModel.Name != "Updated Model" {
			t.Errorf("Expected name Updated Model, got %s", updatedModel.Name)
		}

		// 列出模型
		models, err := storage.ListModels()
		if err != nil {
			t.Errorf("ListModels failed: %v", err)
		}
		if len(models) == 0 {
			t.Error("Expected at least one model, got none")
		}

		// 列出指定LLM资源的模型
		_, err = storage.ListModelsByLLMResource("test-resource")
		if err != nil {
			t.Errorf("ListModelsByLLMResource failed: %v", err)
		}

		// 删除模型
		if err := storage.DeleteModel(model.ID); err != nil {
			t.Errorf("DeleteModel failed: %v", err)
		}

		// 验证模型已删除
		_, err = storage.GetModel(model.ID)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// 测试端点管理
	t.Run("Endpoint Management", func(t *testing.T) {
		// 创建端点
		endpoint := &Endpoint{
			LLMResourceID: "test-resource",
			Path:          "/test",
			Method:        "POST",
			Status:        "active",
		}

		if err := storage.CreateEndpoint(endpoint); err != nil {
			t.Errorf("CreateEndpoint failed: %v", err)
		}

		// 获取端点
		retrievedEndpoint, err := storage.GetEndpoint(endpoint.ID)
		if err != nil {
			t.Errorf("GetEndpoint failed: %v", err)
		}
		if retrievedEndpoint.Path != endpoint.Path {
			t.Errorf("Expected path %s, got %s", endpoint.Path, retrievedEndpoint.Path)
		}

		// 更新端点
		retrievedEndpoint.Path = "/updated"
		if err := storage.UpdateEndpoint(retrievedEndpoint); err != nil {
			t.Errorf("UpdateEndpoint failed: %v", err)
		}

		updatedEndpoint, err := storage.GetEndpoint(endpoint.ID)
		if err != nil {
			t.Errorf("GetEndpoint after update failed: %v", err)
		}
		if updatedEndpoint.Path != "/updated" {
			t.Errorf("Expected path /updated, got %s", updatedEndpoint.Path)
		}

		// 列出端点
		endpoints, err := storage.ListEndpoints()
		if err != nil {
			t.Errorf("ListEndpoints failed: %v", err)
		}
		if len(endpoints) == 0 {
			t.Error("Expected at least one endpoint, got none")
		}

		// 删除端点
		if err := storage.DeleteEndpoint(endpoint.ID); err != nil {
			t.Errorf("DeleteEndpoint failed: %v", err)
		}

		// 验证端点已删除
		_, err = storage.GetEndpoint(endpoint.ID)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// 测试请求管理
	t.Run("Request Management", func(t *testing.T) {
		// 创建请求
		request := &Request{
			UserID:   "test-user",
			Endpoint: "/test",
			Method:   "POST",
			Status:   "success",
			Duration: 100,
			Tokens:   100,
		}

		if err := storage.CreateRequest(request); err != nil {
			t.Errorf("CreateRequest failed: %v", err)
		}

		// 获取请求
		retrievedRequest, err := storage.GetRequest(request.ID)
		if err != nil {
			t.Errorf("GetRequest failed: %v", err)
		}
		if retrievedRequest.Endpoint != request.Endpoint {
			t.Errorf("Expected endpoint %s, got %s", request.Endpoint, retrievedRequest.Endpoint)
		}

		// 列出请求
		requests, err := storage.ListRequests(10)
		if err != nil {
			t.Errorf("ListRequests failed: %v", err)
		}
		if len(requests) == 0 {
			t.Error("Expected at least one request, got none")
		}
	})
}
