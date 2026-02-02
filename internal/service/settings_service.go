package service

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/lingproxy/lingproxy/internal/config"
)

// UpdateSettingsRequest 更新设置请求
type UpdateSettingsRequest struct {
	Basic          *BasicSettings          `json:"basic,omitempty"`
	Cache          *CacheSettings          `json:"cache,omitempty"`
	RateLimit      *RateLimitSettings      `json:"rate_limit,omitempty"`
	Security       *SecuritySettings       `json:"security,omitempty"`
	Log            *LogSettings            `json:"log,omitempty"`
	LoadBalancer   *LoadBalancerSettings   `json:"load_balancer,omitempty"`
	CircuitBreaker *CircuitBreakerSettings `json:"circuit_breaker,omitempty"`
}

// BasicSettings 基本设置
type BasicSettings struct {
	SystemName string `json:"system_name,omitempty"`
	APIPort    *int   `json:"port,omitempty"`
	Host       string `json:"host,omitempty"`
	DebugMode  *bool  `json:"debug_mode,omitempty"`
}

// CacheSettings 缓存设置
type CacheSettings struct {
	Enabled   *bool `json:"enabled,omitempty"`
	TTL       *int  `json:"ttl,omitempty"`
	SizeLimit *int  `json:"size_limit,omitempty"`
}

// RateLimitSettings 限流设置
type RateLimitSettings struct {
	Enabled           *bool `json:"enabled,omitempty"`
	RequestsPerMinute *int  `json:"requests_per_minute,omitempty"`
	Concurrency       *int  `json:"concurrency,omitempty"`
}

// SecuritySettings 安全设置
type SecuritySettings struct {
	AuthEnabled     *bool         `json:"auth_enabled,omitempty"`
	JWTSecret       string        `json:"jwt_secret,omitempty"`
	TokenExpiration *int          `json:"token_expiration,omitempty"`
	HTTPSEnabled    *bool         `json:"https_enabled,omitempty"`
	CORS            *CORSSettings `json:"cors,omitempty"`
}

// CORSSettings CORS设置
type CORSSettings struct {
	Enabled      *bool    `json:"enabled,omitempty"`
	AllowOrigins []string `json:"allow_origins,omitempty"`
	AllowMethods []string `json:"allow_methods,omitempty"`
	AllowHeaders []string `json:"allow_headers,omitempty"`
}

// LogSettings 日志设置
type LogSettings struct {
	Level      string `json:"level,omitempty"`
	Format     string `json:"format,omitempty"`
	Output     string `json:"output,omitempty"`
	FilePath   string `json:"file_path,omitempty"`
	MaxSize    *int   `json:"max_size,omitempty"`
	MaxBackups *int   `json:"max_backups,omitempty"`
	MaxAge     *int   `json:"max_age,omitempty"`
	Compress   *bool  `json:"compress,omitempty"`
}

// LoadBalancerSettings 负载均衡设置
type LoadBalancerSettings struct {
	DefaultStrategy string               `json:"default_strategy,omitempty"`
	HealthCheck     *HealthCheckSettings `json:"health_check,omitempty"`
}

// HealthCheckSettings 健康检查设置
type HealthCheckSettings struct {
	Enabled     *bool `json:"enabled,omitempty"`
	Interval    *int  `json:"interval,omitempty"`
	Timeout     *int  `json:"timeout,omitempty"`
	MaxFailures *int  `json:"max_failures,omitempty"`
}

// CircuitBreakerSettings 熔断器设置
type CircuitBreakerSettings struct {
	Enabled          *bool `json:"enabled,omitempty"`
	FailureThreshold *int  `json:"failure_threshold,omitempty"`
	Timeout          *int  `json:"timeout,omitempty"`
	MaxRequests      *int  `json:"max_requests,omitempty"`
	Interval         *int  `json:"interval,omitempty"`
}

// SettingsService 设置服务
type SettingsService struct {
	configPath string
}

// NewSettingsService 创建新的设置服务
func NewSettingsService(configPath string) *SettingsService {
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	return &SettingsService{
		configPath: configPath,
	}
}

// UpdateSettings 更新设置
func (s *SettingsService) UpdateSettings(req *UpdateSettingsRequest, restartRequiredFields *[]string) error {
	cfg := config.C

	// 更新基本设置
	if req.Basic != nil {
		if req.Basic.SystemName != "" {
			cfg.App.Name = req.Basic.SystemName
			viper.Set("app.name", req.Basic.SystemName)
		}
		if req.Basic.APIPort != nil {
			cfg.App.Port = *req.Basic.APIPort
			viper.Set("app.port", *req.Basic.APIPort)
			*restartRequiredFields = append(*restartRequiredFields, "port")
		}
		if req.Basic.Host != "" {
			cfg.App.Host = req.Basic.Host
			viper.Set("app.host", req.Basic.Host)
			*restartRequiredFields = append(*restartRequiredFields, "host")
		}
		if req.Basic.DebugMode != nil {
			if *req.Basic.DebugMode {
				cfg.App.Environment = "development"
			} else {
				cfg.App.Environment = "production"
			}
			viper.Set("app.environment", cfg.App.Environment)
		}
	}

	// 更新缓存设置
	if req.Cache != nil {
		if req.Cache.Enabled != nil {
			cfg.Cache.Enabled = *req.Cache.Enabled
			viper.Set("cache.enabled", *req.Cache.Enabled)
		}
		if req.Cache.TTL != nil {
			cfg.Cache.TTL = time.Duration(*req.Cache.TTL) * time.Second
			viper.Set("cache.ttl", fmt.Sprintf("%ds", *req.Cache.TTL))
		}
	}

	// 更新限流设置
	if req.RateLimit != nil {
		if req.RateLimit.Enabled != nil {
			cfg.Security.RateLimit.Enabled = *req.RateLimit.Enabled
			viper.Set("security.rate_limit.enabled", *req.RateLimit.Enabled)
		}
		if req.RateLimit.RequestsPerMinute != nil {
			cfg.Security.RateLimit.RequestsPerMinute = *req.RateLimit.RequestsPerMinute
			viper.Set("security.rate_limit.requests_per_minute", *req.RateLimit.RequestsPerMinute)
		}
	}

	// 更新安全设置
	if req.Security != nil {
		if req.Security.AuthEnabled != nil {
			cfg.Security.Auth.Enabled = *req.Security.AuthEnabled
			viper.Set("security.auth.enabled", *req.Security.AuthEnabled)
			*restartRequiredFields = append(*restartRequiredFields, "auth_enabled")
		}
		if req.Security.JWTSecret != "" && req.Security.JWTSecret != "******" {
			cfg.Security.JWT.Secret = req.Security.JWTSecret
			viper.Set("security.jwt.secret", req.Security.JWTSecret)
			*restartRequiredFields = append(*restartRequiredFields, "jwt_secret")
		}
		if req.Security.TokenExpiration != nil {
			cfg.Security.JWT.ExpireHours = time.Duration(*req.Security.TokenExpiration) * time.Hour
			viper.Set("security.jwt.expire_hours", fmt.Sprintf("%dh", *req.Security.TokenExpiration))
		}
		if req.Security.CORS != nil {
			if req.Security.CORS.Enabled != nil {
				cfg.Security.CORS.Enabled = *req.Security.CORS.Enabled
				viper.Set("security.cors.enabled", *req.Security.CORS.Enabled)
			}
			if req.Security.CORS.AllowOrigins != nil {
				cfg.Security.CORS.AllowOrigins = req.Security.CORS.AllowOrigins
				viper.Set("security.cors.allow_origins", req.Security.CORS.AllowOrigins)
			}
			if req.Security.CORS.AllowMethods != nil {
				cfg.Security.CORS.AllowMethods = req.Security.CORS.AllowMethods
				viper.Set("security.cors.allow_methods", req.Security.CORS.AllowMethods)
			}
			if req.Security.CORS.AllowHeaders != nil {
				cfg.Security.CORS.AllowHeaders = req.Security.CORS.AllowHeaders
				viper.Set("security.cors.allow_headers", req.Security.CORS.AllowHeaders)
			}
		}
	}

	// 更新日志设置
	if req.Log != nil {
		if req.Log.Level != "" {
			cfg.Log.Level = req.Log.Level
			viper.Set("log.level", req.Log.Level)
		}
		if req.Log.Format != "" {
			cfg.Log.Format = req.Log.Format
			viper.Set("log.format", req.Log.Format)
		}
		if req.Log.Output != "" {
			cfg.Log.Output = req.Log.Output
			viper.Set("log.output", req.Log.Output)
		}
		if req.Log.FilePath != "" {
			cfg.Log.FilePath = req.Log.FilePath
			viper.Set("log.file_path", req.Log.FilePath)
		}
		if req.Log.MaxSize != nil {
			cfg.Log.MaxSize = *req.Log.MaxSize
			viper.Set("log.max_size", *req.Log.MaxSize)
		}
		if req.Log.MaxBackups != nil {
			cfg.Log.MaxBackups = *req.Log.MaxBackups
			viper.Set("log.max_backups", *req.Log.MaxBackups)
		}
		if req.Log.MaxAge != nil {
			cfg.Log.MaxAge = *req.Log.MaxAge
			viper.Set("log.max_age", *req.Log.MaxAge)
		}
		if req.Log.Compress != nil {
			cfg.Log.Compress = *req.Log.Compress
			viper.Set("log.compress", *req.Log.Compress)
		}
	}

	// 更新负载均衡设置
	if req.LoadBalancer != nil {
		if req.LoadBalancer.DefaultStrategy != "" {
			cfg.LoadBalancer.DefaultStrategy = req.LoadBalancer.DefaultStrategy
			viper.Set("load_balancer.default_strategy", req.LoadBalancer.DefaultStrategy)
		}
		if req.LoadBalancer.HealthCheck != nil {
			if req.LoadBalancer.HealthCheck.Enabled != nil {
				cfg.LoadBalancer.HealthCheck.Enabled = *req.LoadBalancer.HealthCheck.Enabled
				viper.Set("load_balancer.health_check.enabled", *req.LoadBalancer.HealthCheck.Enabled)
			}
			if req.LoadBalancer.HealthCheck.Interval != nil {
				cfg.LoadBalancer.HealthCheck.Interval = time.Duration(*req.LoadBalancer.HealthCheck.Interval) * time.Second
				viper.Set("load_balancer.health_check.interval", fmt.Sprintf("%ds", *req.LoadBalancer.HealthCheck.Interval))
			}
			if req.LoadBalancer.HealthCheck.Timeout != nil {
				cfg.LoadBalancer.HealthCheck.Timeout = time.Duration(*req.LoadBalancer.HealthCheck.Timeout) * time.Second
				viper.Set("load_balancer.health_check.timeout", fmt.Sprintf("%ds", *req.LoadBalancer.HealthCheck.Timeout))
			}
			if req.LoadBalancer.HealthCheck.MaxFailures != nil {
				cfg.LoadBalancer.HealthCheck.MaxFailures = *req.LoadBalancer.HealthCheck.MaxFailures
				viper.Set("load_balancer.health_check.max_failures", *req.LoadBalancer.HealthCheck.MaxFailures)
			}
		}
	}

	// 更新熔断器设置
	if req.CircuitBreaker != nil {
		if req.CircuitBreaker.Enabled != nil {
			cfg.CircuitBreaker.Enabled = *req.CircuitBreaker.Enabled
			viper.Set("circuit_breaker.enabled", *req.CircuitBreaker.Enabled)
		}
		if req.CircuitBreaker.FailureThreshold != nil {
			cfg.CircuitBreaker.FailureThreshold = *req.CircuitBreaker.FailureThreshold
			viper.Set("circuit_breaker.failure_threshold", *req.CircuitBreaker.FailureThreshold)
		}
		if req.CircuitBreaker.Timeout != nil {
			cfg.CircuitBreaker.Timeout = time.Duration(*req.CircuitBreaker.Timeout) * time.Second
			viper.Set("circuit_breaker.timeout", fmt.Sprintf("%ds", *req.CircuitBreaker.Timeout))
		}
		if req.CircuitBreaker.MaxRequests != nil {
			cfg.CircuitBreaker.MaxRequests = *req.CircuitBreaker.MaxRequests
			viper.Set("circuit_breaker.max_requests", *req.CircuitBreaker.MaxRequests)
		}
		if req.CircuitBreaker.Interval != nil {
			cfg.CircuitBreaker.Interval = time.Duration(*req.CircuitBreaker.Interval) * time.Second
			viper.Set("circuit_breaker.interval", fmt.Sprintf("%ds", *req.CircuitBreaker.Interval))
		}
	}

	// 保存到配置文件
	return s.saveConfigToFile()
}

// saveConfigToFile 保存配置到文件
func (s *SettingsService) saveConfigToFile() error {
	// 写入文件
	if err := viper.WriteConfigAs(s.configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
