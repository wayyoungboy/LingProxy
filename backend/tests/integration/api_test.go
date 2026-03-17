package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/config"
	"github.com/lingproxy/lingproxy/internal/handler"
	"github.com/lingproxy/lingproxy/internal/middleware"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestServer represents a test server instance
type TestServer struct {
	Router   *gin.Engine
	Storage  *storage.StorageFacade
	Services *TestServices
}

// TestServices holds service instances for testing
type TestServices struct {
	TokenService  *service.APIKeyService
	UserService   *service.UserService
	PolicyService *service.PolicyService
}

// NewTestServer creates a new test server
func NewTestServer(t *testing.T) *TestServer {
	memStorage := storage.NewMemoryStorage()
	facade := storage.NewStorageFacade(memStorage)

	tokenService := service.NewTokenService(facade)
	userService := service.NewUserService(facade)
	policyService := service.NewPolicyService(facade)

	cfg := &config.Config{
		Security: config.SecurityConfig{
			Auth: config.AuthConfig{
				Enabled: true,
			},
			CORS: config.CORSConfig{
				Enabled: true,
			},
		},
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.CORS())

	// Setup routes
	api := r.Group("/api")
	{
		// Health check (no auth required)
		api.GET("/health", handler.HealthHandler)

		// Protected routes
		authMiddleware := middleware.NewAuthMiddleware(facade, tokenService, cfg)
		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/resources", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "protected resource"})
			})
		}
	}

	return &TestServer{
		Router:  r,
		Storage: facade,
		Services: &TestServices{
			TokenService:  tokenService,
			UserService:   userService,
			PolicyService: policyService,
		},
	}
}

func TestHealthCheck(t *testing.T) {
	server := NewTestServer(t)

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestProtectedRouteWithoutAuth(t *testing.T) {
	server := NewTestServer(t)

	req := httptest.NewRequest("GET", "/api/resources", nil)
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProtectedRouteWithValidAPIKey(t *testing.T) {
	server := NewTestServer(t)

	// Create a user with API key
	user := &storage.User{
		Username: "testuser",
		APIKey:   "test-api-key-12345",
		Status:   "active",
	}
	err := server.Storage.CreateUser(user)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/resources", nil)
	req.Header.Set("Authorization", "Bearer test-api-key-12345")
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProtectedRouteWithInvalidAPIKey(t *testing.T) {
	server := NewTestServer(t)

	req := httptest.NewRequest("GET", "/api/resources", nil)
	req.Header.Set("Authorization", "Bearer invalid-key")
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCORSPreflight(t *testing.T) {
	server := NewTestServer(t)

	req := httptest.NewRequest("OPTIONS", "/api/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestAPICreation(t *testing.T) {
	server := NewTestServer(t)

	// Create an API key using the service
	apiKey, err := server.Services.TokenService.CreateAPIKey("test-key", nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey.APIKey)
	assert.Contains(t, apiKey.APIKey, "ling-")
	assert.Equal(t, "active", apiKey.Status)
}

func TestUserManagement(t *testing.T) {
	server := NewTestServer(t)

	// Create user
	user := &storage.User{
		Username: "integration-test-user",
		APIKey:   "integration-test-api-key",
		Status:   "active",
	}

	err := server.Storage.CreateUser(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	// Get user
	retrievedUser, err := server.Storage.GetUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, retrievedUser.Username)

	// Update user
	retrievedUser.Username = "updated-user"
	err = server.Storage.UpdateUser(retrievedUser)
	assert.NoError(t, err)

	// Verify update
	updatedUser, err := server.Storage.GetUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updated-user", updatedUser.Username)

	// Delete user
	err = server.Storage.DeleteUser(user.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = server.Storage.GetUser(user.ID)
	assert.Error(t, err)
}

func TestLLMResourceManagement(t *testing.T) {
	server := NewTestServer(t)

	// Create LLM resource
	resource := &storage.LLMResource{
		Name:    "test-resource",
		Driver:  "openai",
		Model:   "gpt-4",
		BaseURL: "https://api.openai.com/v1",
		Status:  "active",
		Type:    "chat",
	}
	err := server.Storage.CreateLLMResource(resource)
	assert.NoError(t, err)
	assert.NotEmpty(t, resource.ID)

	// Get resource
	retrievedResource, err := server.Storage.GetLLMResource(resource.ID)
	assert.NoError(t, err)
	assert.Equal(t, resource.Name, retrievedResource.Name)

	// List resources
	resources, err := server.Storage.ListLLMResources()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(resources), 1)
}

func TestPolicyManagement(t *testing.T) {
	server := NewTestServer(t)

	// Create policy
	policy := &storage.Policy{
		Name:    "test-policy",
		Type:    "round_robin",
		Enabled: true,
		Builtin: false,
	}
	err := server.Storage.CreatePolicy(policy)
	assert.NoError(t, err)
	assert.NotEmpty(t, policy.ID)

	// Get policy
	retrievedPolicy, err := server.Storage.GetPolicy(policy.ID)
	assert.NoError(t, err)
	assert.Equal(t, policy.Name, retrievedPolicy.Name)
}

func TestAPIKeyValidation(t *testing.T) {
	server := NewTestServer(t)

	// Create an API key
	apiKey, err := server.Services.TokenService.CreateAPIKey("validation-test-key", nil)
	assert.NoError(t, err)

	// Validate the key
	validatedKey, err := server.Services.TokenService.ValidateAPIKey(apiKey.APIKey)
	assert.NoError(t, err)
	assert.Equal(t, apiKey.ID, validatedKey.ID)

	// Test invalid key
	_, err = server.Services.TokenService.ValidateAPIKey("invalid-key-xyz")
	assert.Error(t, err)
}

func TestInactiveAPIKey(t *testing.T) {
	server := NewTestServer(t)

	// Create an API key
	apiKey, err := server.Services.TokenService.CreateAPIKey("inactive-key", nil)
	assert.NoError(t, err)

	// Deactivate the key
	inactive := "inactive"
	_, err = server.Services.TokenService.UpdateAPIKey(apiKey.ID, nil, &inactive)
	assert.NoError(t, err)

	// Try to validate inactive key
	_, err = server.Services.TokenService.ValidateAPIKey(apiKey.APIKey)
	assert.Error(t, err)
	assert.Equal(t, service.ErrAPIKeyInactive, err)
}

// BenchmarkHealthCheck benchmarks the health check endpoint
func BenchmarkHealthCheck(b *testing.B) {
	gin.SetMode(gin.TestMode)
	server := NewTestServer(&testing.T{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/health", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)
	}
}