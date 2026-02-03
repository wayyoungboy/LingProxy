package storage

import (
	"testing"
)

func TestStorageFacade(t *testing.T) {
	// 初始化内存存储
	memoryStorage := NewMemoryStorage()

	// 初始化存储门面
	facade := NewStorageFacade(memoryStorage)

	// 测试用户管理
	t.Run("User Management", func(t *testing.T) {
		// 创建用户
		user := &User{
			Username: "facade-test-user",
			APIKey:   "facade-test-api-key",
			Status:   "active",
		}

		if err := facade.CreateUser(user); err != nil {
			t.Errorf("CreateUser failed: %v", err)
		}

		// 获取用户
		retrievedUser, err := facade.GetUser(user.ID)
		if err != nil {
			t.Errorf("GetUser failed: %v", err)
		}
		if retrievedUser.Username != user.Username {
			t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
		}

		// 通过API Key获取用户
		userByAPIKey, err := facade.GetUserByAPIKey(user.APIKey)
		if err != nil {
			t.Errorf("GetUserByAPIKey failed: %v", err)
		}
		if userByAPIKey.Username != user.Username {
			t.Errorf("Expected username %s, got %s", user.Username, userByAPIKey.Username)
		}

		// 更新用户
		retrievedUser.Username = "facade-updated-user"
		retrievedUser.Status = "inactive"
		if err := facade.UpdateUser(retrievedUser); err != nil {
			t.Errorf("UpdateUser failed: %v", err)
		}

		updatedUser, err := facade.GetUser(user.ID)
		if err != nil {
			t.Errorf("GetUser after update failed: %v", err)
		}
		if updatedUser.Username != "facade-updated-user" {
			t.Errorf("Expected username facade-updated-user, got %s", updatedUser.Username)
		}
		if updatedUser.Status != "inactive" {
			t.Errorf("Expected status inactive, got %s", updatedUser.Status)
		}

		// 列出用户
		users, err := facade.ListUsers()
		if err != nil {
			t.Errorf("ListUsers failed: %v", err)
		}
		if len(users) == 0 {
			t.Error("Expected at least one user, got none")
		}

		// 删除用户
		if err := facade.DeleteUser(user.ID); err != nil {
			t.Errorf("DeleteUser failed: %v", err)
		}

		// 验证用户已删除
		_, err = facade.GetUser(user.ID)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// 测试LLM资源管理
	t.Run("LLM Resource Management", func(t *testing.T) {
		// 创建LLM资源
		resource := &LLMResource{
			Name:     "Facade Test Resource",
			Type:     "chat",
			Driver: "openai",
			Model:    "gpt-3.5-turbo",
			BaseURL:  "https://api.openai.com/v1",
			APIKey:   "facade-test-api-key",
			Status:   "active",
		}

		if err := facade.CreateLLMResource(resource); err != nil {
			t.Errorf("CreateLLMResource failed: %v", err)
		}

		// 获取LLM资源
		retrievedResource, err := facade.GetLLMResource(resource.ID)
		if err != nil {
			t.Errorf("GetLLMResource failed: %v", err)
		}
		if retrievedResource.Name != resource.Name {
			t.Errorf("Expected name %s, got %s", resource.Name, retrievedResource.Name)
		}

		// 更新LLM资源
		retrievedResource.Name = "Facade Updated Resource"
		if err := facade.UpdateLLMResource(retrievedResource); err != nil {
			t.Errorf("UpdateLLMResource failed: %v", err)
		}

		updatedResource, err := facade.GetLLMResource(resource.ID)
		if err != nil {
			t.Errorf("GetLLMResource after update failed: %v", err)
		}
		if updatedResource.Name != "Facade Updated Resource" {
			t.Errorf("Expected name Facade Updated Resource, got %s", updatedResource.Name)
		}

		// 列出LLM资源
		resources, err := facade.ListLLMResources()
		if err != nil {
			t.Errorf("ListLLMResources failed: %v", err)
		}
		if len(resources) == 0 {
			t.Error("Expected at least one resource, got none")
		}

		// 删除LLM资源
		if err := facade.DeleteLLMResource(resource.ID); err != nil {
			t.Errorf("DeleteLLMResource failed: %v", err)
		}

		// 验证资源已删除
		_, err = facade.GetLLMResource(resource.ID)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// 测试模型管理
	t.Run("Model Management", func(t *testing.T) {
		// 创建模型
		model := &Model{
			Name:          "Facade Test Model",
			LLMResourceID: "facade-test-resource",
			ModelID:       "facade-test-model",
			Type:          "chat",
			Category:      "gpt",
			Version:       "1.0",
			Description:   "Facade test model",
			Capabilities:  "[\"text-generation\"]",
			Pricing:       "{\"input_token_price\":0.001,\"output_token_price\":0.002,\"currency\":\"USD\"}",
			Limits:        "{\"max_tokens\":4096,\"context_window\":8192}",
			Parameters:    "{\"temperature\":0.7,\"max_tokens\":1000}",
			Features:      "{\"streaming\":true,\"function_calling\":false}",
			Status:        "active",
			Metadata:      "{\"key\":\"value\"}",
		}

		if err := facade.CreateModel(model); err != nil {
			t.Errorf("CreateModel failed: %v", err)
		}

		// 获取模型
		retrievedModel, err := facade.GetModel(model.ID)
		if err != nil {
			t.Errorf("GetModel failed: %v", err)
		}
		if retrievedModel.Name != model.Name {
			t.Errorf("Expected name %s, got %s", model.Name, retrievedModel.Name)
		}

		// 更新模型
		retrievedModel.Name = "Facade Updated Model"
		if err := facade.UpdateModel(retrievedModel); err != nil {
			t.Errorf("UpdateModel failed: %v", err)
		}

		updatedModel, err := facade.GetModel(model.ID)
		if err != nil {
			t.Errorf("GetModel after update failed: %v", err)
		}
		if updatedModel.Name != "Facade Updated Model" {
			t.Errorf("Expected name Facade Updated Model, got %s", updatedModel.Name)
		}

		// 列出模型
		models, err := facade.ListModels()
		if err != nil {
			t.Errorf("ListModels failed: %v", err)
		}
		if len(models) == 0 {
			t.Error("Expected at least one model, got none")
		}

		// 列出指定LLM资源的模型
		_, err = facade.ListModelsByLLMResource("facade-test-resource")
		if err != nil {
			t.Errorf("ListModelsByLLMResource failed: %v", err)
		}

		// 删除模型
		if err := facade.DeleteModel(model.ID); err != nil {
			t.Errorf("DeleteModel failed: %v", err)
		}

		// 验证模型已删除
		_, err = facade.GetModel(model.ID)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// 测试端点管理
	t.Run("Endpoint Management", func(t *testing.T) {
		// 创建端点
		endpoint := &Endpoint{
			LLMResourceID: "facade-test-resource",
			Path:          "/facade-test",
			Method:        "POST",
			Status:        "active",
		}

		if err := facade.CreateEndpoint(endpoint); err != nil {
			t.Errorf("CreateEndpoint failed: %v", err)
		}

		// 获取端点
		retrievedEndpoint, err := facade.GetEndpoint(endpoint.ID)
		if err != nil {
			t.Errorf("GetEndpoint failed: %v", err)
		}
		if retrievedEndpoint.Path != endpoint.Path {
			t.Errorf("Expected path %s, got %s", endpoint.Path, retrievedEndpoint.Path)
		}

		// 更新端点
		retrievedEndpoint.Path = "/facade-updated"
		if err := facade.UpdateEndpoint(retrievedEndpoint); err != nil {
			t.Errorf("UpdateEndpoint failed: %v", err)
		}

		updatedEndpoint, err := facade.GetEndpoint(endpoint.ID)
		if err != nil {
			t.Errorf("GetEndpoint after update failed: %v", err)
		}
		if updatedEndpoint.Path != "/facade-updated" {
			t.Errorf("Expected path /facade-updated, got %s", updatedEndpoint.Path)
		}

		// 列出端点
		endpoints, err := facade.ListEndpoints()
		if err != nil {
			t.Errorf("ListEndpoints failed: %v", err)
		}
		if len(endpoints) == 0 {
			t.Error("Expected at least one endpoint, got none")
		}

		// 删除端点
		if err := facade.DeleteEndpoint(endpoint.ID); err != nil {
			t.Errorf("DeleteEndpoint failed: %v", err)
		}

		// 验证端点已删除
		_, err = facade.GetEndpoint(endpoint.ID)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// 测试请求管理
	t.Run("Request Management", func(t *testing.T) {
		// 创建请求
		request := &Request{
			UserID:   "facade-test-user",
			Endpoint: "/facade-test",
			Method:   "POST",
			Status:   "success",
			Duration: 150,
			Tokens:   150,
		}

		if err := facade.CreateRequest(request); err != nil {
			t.Errorf("CreateRequest failed: %v", err)
		}

		// 获取请求
		retrievedRequest, err := facade.GetRequest(request.ID)
		if err != nil {
			t.Errorf("GetRequest failed: %v", err)
		}
		if retrievedRequest.Endpoint != request.Endpoint {
			t.Errorf("Expected endpoint %s, got %s", request.Endpoint, retrievedRequest.Endpoint)
		}

		// 列出请求
		requests, err := facade.ListRequests(5)
		if err != nil {
			t.Errorf("ListRequests failed: %v", err)
		}
		if len(requests) == 0 {
			t.Error("Expected at least one request, got none")
		}
	})
}
