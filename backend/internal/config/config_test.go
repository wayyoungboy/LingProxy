package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Reset viper for each test
	viper.Reset()
}

func TestSetDefaults(t *testing.T) {
	viper.Reset()
	setDefaults()

	// Verify application defaults
	assert.Equal(t, "LingProxy", viper.GetString("app.name"))
	assert.Equal(t, "1.0.0", viper.GetString("app.version"))
	assert.Equal(t, "development", viper.GetString("app.environment"))
	assert.Equal(t, 8080, viper.GetInt("app.port"))
	assert.Equal(t, "0.0.0.0", viper.GetString("app.host"))

	// Verify storage defaults
	assert.Equal(t, "memory", viper.GetString("storage.type"))
	assert.Equal(t, "sqlite", viper.GetString("storage.gorm.driver"))
	assert.Equal(t, "lingproxy.db", viper.GetString("storage.gorm.dsn"))

	// Verify cache defaults
	assert.Equal(t, true, viper.GetBool("cache.enabled"))

	// Verify security defaults
	assert.Equal(t, true, viper.GetBool("security.auth.enabled"))
	assert.Equal(t, true, viper.GetBool("security.cors.enabled"))
	assert.Equal(t, true, viper.GetBool("security.rate_limit.enabled"))

	// Verify load balancer defaults
	assert.Equal(t, "round_robin", viper.GetString("load_balancer.default_strategy"))
	assert.Equal(t, true, viper.GetBool("load_balancer.health_check.enabled"))
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid config",
			config: &Config{
				App: AppConfig{
					Port: 8080,
				},
				Storage: StorageConfig{
					Type: "memory",
				},
			},
			expectErr: false,
		},
		{
			name: "invalid port - negative",
			config: &Config{
				App: AppConfig{
					Port: -1,
				},
				Storage: StorageConfig{
					Type: "memory",
				},
			},
			expectErr: true,
			errMsg:    "无效的端口号",
		},
		{
			name: "invalid port - too high",
			config: &Config{
				App: AppConfig{
					Port: 70000,
				},
				Storage: StorageConfig{
					Type: "memory",
				},
			},
			expectErr: true,
			errMsg:    "无效的端口号",
		},
		{
			name: "invalid storage type",
			config: &Config{
				App: AppConfig{
					Port: 8080,
				},
				Storage: StorageConfig{
					Type: "invalid",
				},
			},
			expectErr: true,
			errMsg:    "不支持的存储类型",
		},
		{
			name: "gorm with valid sqlite",
			config: &Config{
				App: AppConfig{
					Port: 8080,
				},
				Storage: StorageConfig{
					Type: "gorm",
					GORM: GORMConfig{
						Driver: "sqlite",
						DSN:    "test.db",
					},
				},
			},
			expectErr: false,
		},
		{
			name: "gorm with valid mysql",
			config: &Config{
				App: AppConfig{
					Port: 8080,
				},
				Storage: StorageConfig{
					Type: "gorm",
					GORM: GORMConfig{
						Driver: "mysql",
						DSN:    "user:pass@tcp(localhost:3306)/db",
					},
				},
			},
			expectErr: false,
		},
		{
			name: "gorm with invalid driver",
			config: &Config{
				App: AppConfig{
					Port: 8080,
				},
				Storage: StorageConfig{
					Type: "gorm",
					GORM: GORMConfig{
						Driver: "postgres",
						DSN:    "test",
					},
				},
			},
			expectErr: true,
			errMsg:    "不支持的GORM驱动类型",
		},
		{
			name: "gorm with empty dsn",
			config: &Config{
				App: AppConfig{
					Port: 8080,
				},
				Storage: StorageConfig{
					Type: "gorm",
					GORM: GORMConfig{
						Driver: "sqlite",
						DSN:    "",
					},
				},
			},
			expectErr: true,
			errMsg:    "GORM DSN不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			C = tt.config
			err := validateConfig()

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateDirectories(t *testing.T) {
	// Create temp test directories
	err := createDirectories()
	assert.NoError(t, err)

	// Verify directories exist
	dirs := []string{"./data", "./logs", "./temp"}
	for _, dir := range dirs {
		_, err := os.Stat(dir)
		assert.NoError(t, err, "Directory %s should exist", dir)
	}
}

func TestIsDevelopment(t *testing.T) {
	C = &Config{
		App: AppConfig{
			Environment: "development",
		},
	}
	assert.True(t, IsDevelopment())
	assert.False(t, IsProduction())

	C.App.Environment = "production"
	assert.False(t, IsDevelopment())
	assert.True(t, IsProduction())
}

func TestIsProduction(t *testing.T) {
	C = &Config{
		App: AppConfig{
			Environment: "production",
		},
	}
	assert.True(t, IsProduction())
	assert.False(t, IsDevelopment())
}

func TestConfigStructures(t *testing.T) {
	// Test that all config structures can be properly initialized
	cfg := &Config{
		App: AppConfig{
			Name:        "TestApp",
			Version:     "2.0.0",
			Environment: "test",
			Port:        9090,
			Host:        "localhost",
		},
		Storage: StorageConfig{
			Type: "gorm",
			GORM: GORMConfig{
				Driver: "sqlite",
				DSN:    "test.db",
			},
		},
		Cache: CacheConfig{
			Enabled: true,
			TTL:     5 * time.Minute,
		},
		Log: LogConfig{
			Level:      "debug",
			Format:     "text",
			Output:     "stdout",
			FilePath:   "/var/log/test.log",
			MaxSize:    50,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		},
		Security: SecurityConfig{
			Auth: AuthConfig{
				Enabled: true,
			},
			JWT: JWTConfig{
				Secret:      "test-secret",
				ExpireHours: 24 * time.Hour,
			},
			CORS: CORSConfig{
				Enabled:      true,
				AllowOrigins: []string{"http://localhost:3000"},
				AllowMethods: []string{"GET", "POST"},
				AllowHeaders: []string{"Content-Type"},
			},
			RateLimit: RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 100,
			},
		},
		LoadBalancer: LoadBalancerConfig{
			DefaultStrategy: "random",
			HealthCheck: HealthCheckConfig{
				Enabled:     true,
				Interval:    10 * time.Second,
				Timeout:     3 * time.Second,
				MaxFailures: 5,
			},
		},
		Provider: ProviderConfig{
			Timeout:         60 * time.Second,
			MaxRetries:      5,
			RetryDelay:      2 * time.Second,
			MaxIdleConns:    50,
			MaxConnsPerHost: 20,
			IdleConnTimeout: 30 * time.Second,
		},
	}

	assert.Equal(t, "TestApp", cfg.App.Name)
	assert.Equal(t, "gorm", cfg.Storage.Type)
	assert.True(t, cfg.Cache.Enabled)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.True(t, cfg.Security.Auth.Enabled)
	assert.Equal(t, "random", cfg.LoadBalancer.DefaultStrategy)
	assert.Equal(t, 5, cfg.Provider.MaxRetries)
}