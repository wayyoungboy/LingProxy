package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/config"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthMiddleware_RequireAuth_AuthDisabled(t *testing.T) {
	// Setup
	memStorage := storage.NewMemoryStorage()
	facade := storage.NewStorageFacade(memStorage)
	tokenService := service.NewTokenService(facade)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			Auth: config.AuthConfig{
				Enabled: false,
			},
		},
	}

	middleware := NewAuthMiddleware(facade, tokenService, cfg)

	// Create test router
	router := gin.New()
	router.Use(middleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request without auth header
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should pass through without auth check
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_RequireAuth_MissingHeader(t *testing.T) {
	// Setup
	memStorage := storage.NewMemoryStorage()
	facade := storage.NewStorageFacade(memStorage)
	tokenService := service.NewTokenService(facade)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			Auth: config.AuthConfig{
				Enabled: true,
			},
		},
	}

	middleware := NewAuthMiddleware(facade, tokenService, cfg)

	// Create test router
	router := gin.New()
	router.Use(middleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request without auth header
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_RequireAuth_InvalidHeaderFormat(t *testing.T) {
	// Setup
	memStorage := storage.NewMemoryStorage()
	facade := storage.NewStorageFacade(memStorage)
	tokenService := service.NewTokenService(facade)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			Auth: config.AuthConfig{
				Enabled: true,
			},
		},
	}

	middleware := NewAuthMiddleware(facade, tokenService, cfg)

	// Create test router
	router := gin.New()
	router.Use(middleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		name   string
		header string
	}{
		{"no bearer prefix", "invalid-token"},
		{"wrong prefix", "Basic token123"},
		{"missing token", "Bearer "},
		{"empty bearer", "Bearer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestAuthMiddleware_RequireAuth_InvalidToken(t *testing.T) {
	// Setup
	memStorage := storage.NewMemoryStorage()
	facade := storage.NewStorageFacade(memStorage)
	tokenService := service.NewTokenService(facade)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			Auth: config.AuthConfig{
				Enabled: true,
			},
		},
	}

	middleware := NewAuthMiddleware(facade, tokenService, cfg)

	// Create test router
	router := gin.New()
	router.Use(middleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request with invalid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-12345")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_RequireAuth_ValidUserAPIKey(t *testing.T) {
	// Setup
	memStorage := storage.NewMemoryStorage()
	facade := storage.NewStorageFacade(memStorage)
	tokenService := service.NewTokenService(facade)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			Auth: config.AuthConfig{
				Enabled: true,
			},
		},
	}

	// Create active user with API key
	user := &storage.User{
		Username: "testuser",
		APIKey:   "test-api-key-123",
		Status:   "active",
	}
	err := facade.CreateUser(user)
	assert.NoError(t, err)

	middleware := NewAuthMiddleware(facade, tokenService, cfg)

	// Create test router
	router := gin.New()
	router.Use(middleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request with valid API key
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer test-api-key-123")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should pass through
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_RequireAuth_InactiveUser(t *testing.T) {
	// Setup
	memStorage := storage.NewMemoryStorage()
	facade := storage.NewStorageFacade(memStorage)
	tokenService := service.NewTokenService(facade)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			Auth: config.AuthConfig{
				Enabled: true,
			},
		},
	}

	// Create inactive user with API key
	user := &storage.User{
		Username: "inactive-user",
		APIKey:   "inactive-api-key",
		Status:   "inactive",
	}
	err := facade.CreateUser(user)
	assert.NoError(t, err)

	middleware := NewAuthMiddleware(facade, tokenService, cfg)

	// Create test router
	router := gin.New()
	router.Use(middleware.RequireAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request with inactive user's API key
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer inactive-api-key")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return unauthorized because user is inactive
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestNewAuthMiddleware(t *testing.T) {
	memStorage := storage.NewMemoryStorage()
	facade := storage.NewStorageFacade(memStorage)
	tokenService := service.NewTokenService(facade)
	cfg := &config.Config{}

	middleware := NewAuthMiddleware(facade, tokenService, cfg)

	assert.NotNil(t, middleware)
	assert.Equal(t, facade, middleware.storage)
	assert.Equal(t, tokenService, middleware.tokenService)
	assert.Equal(t, cfg, middleware.config)
}