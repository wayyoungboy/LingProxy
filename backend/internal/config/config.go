package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	App          AppConfig          `mapstructure:"app"`
	Storage      StorageConfig      `mapstructure:"storage"`
	Log          LogConfig          `mapstructure:"log"`
	Security     SecurityConfig     `mapstructure:"security"`
	LoadBalancer LoadBalancerConfig `mapstructure:"load_balancer"`
	Provider     ProviderConfig     `mapstructure:"provider"`
	Cache        CacheConfig        `mapstructure:"cache"`
	Admin        AdminConfig        `mapstructure:"admin"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Port        int    `mapstructure:"port"`
	Host        string `mapstructure:"host"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type string     `mapstructure:"type"`
	GORM GORMConfig `mapstructure:"gorm"`
}

// GORMConfig GORM存储配置
type GORMConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	TTL     time.Duration `mapstructure:"ttl"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	Auth      AuthConfig      `mapstructure:"auth"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Enabled bool `mapstructure:"enabled"` // 是否启用认证
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string        `mapstructure:"secret"`
	ExpireHours time.Duration `mapstructure:"expire_hours"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	Enabled      bool     `mapstructure:"enabled"`
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
	AllowHeaders []string `mapstructure:"allow_headers"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
}

// LoadBalancerConfig 负载均衡配置
type LoadBalancerConfig struct {
	DefaultStrategy string            `mapstructure:"default_strategy"`
	HealthCheck     HealthCheckConfig `mapstructure:"health_check"`
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	Interval    time.Duration `mapstructure:"interval"`
	Timeout     time.Duration `mapstructure:"timeout"`
	MaxFailures int           `mapstructure:"max_failures"`
}

// ProviderConfig Provider配置
type ProviderConfig struct {
	Timeout         time.Duration `mapstructure:"timeout"`
	MaxRetries      int           `mapstructure:"max_retries"`
	RetryDelay      time.Duration `mapstructure:"retry_delay"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxConnsPerHost int           `mapstructure:"max_conns_per_host"`
	IdleConnTimeout time.Duration `mapstructure:"idle_conn_timeout"`
}

// AdminConfig 管理员配置
type AdminConfig struct {
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	APIKey     string `mapstructure:"api_key"`
	AutoCreate bool   `mapstructure:"auto_create"`
}

var (
	// C 全局配置实例
	C *Config
)

// Init 初始化配置
func Init(configPath string) error {
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	// 设置配置文件路径
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置环境变量前缀
	viper.SetEnvPrefix("LINGPROXY")
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("配置文件 %s 未找到，使用默认配置", configPath)
		} else {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 解析配置
	C = &Config{}
	if err := viper.Unmarshal(C); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	// 创建必要的目录
	if err := createDirectories(); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	log.Printf("配置加载成功，环境: %s", C.App.Environment)
	return nil
}

// setDefaults 设置默认值
func setDefaults() {
	// 应用默认配置
	viper.SetDefault("app.name", "LingProxy")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.host", "0.0.0.0")

	// 存储默认配置
	viper.SetDefault("storage.type", "memory")
	viper.SetDefault("storage.gorm.driver", "sqlite")
	viper.SetDefault("storage.gorm.dsn", "lingproxy.db")

	// 缓存默认配置
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.ttl", "5m")

	// 日志默认配置（默认同时输出到stdout和文件，实现持久化）
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "both") // both: 同时输出到stdout和文件
	viper.SetDefault("log.file_path", "./logs/lingproxy.log")
	viper.SetDefault("log.max_size", 100)    // MB
	viper.SetDefault("log.max_backups", 3)   // 保留3个备份文件
	viper.SetDefault("log.max_age", 28)      // 保留28天
	viper.SetDefault("log.compress", true)  // 压缩旧日志

	// 安全默认配置
	viper.SetDefault("security.auth.enabled", true)
	viper.SetDefault("security.jwt.secret", "your-jwt-secret-key-change-this")
	viper.SetDefault("security.jwt.expire_hours", "24h")
	viper.SetDefault("security.cors.enabled", true)
	viper.SetDefault("security.cors.allow_origins", []string{"*"})
	viper.SetDefault("security.cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("security.cors.allow_headers", []string{"*"})
	viper.SetDefault("security.rate_limit.enabled", true)
	viper.SetDefault("security.rate_limit.requests_per_minute", 1000)

	// 负载均衡默认配置
	viper.SetDefault("load_balancer.default_strategy", "round_robin")
	viper.SetDefault("load_balancer.health_check.enabled", true)
	viper.SetDefault("load_balancer.health_check.interval", "30s")
	viper.SetDefault("load_balancer.health_check.timeout", "5s")
	viper.SetDefault("load_balancer.health_check.max_failures", 3)

	// 管理员默认配置
	viper.SetDefault("admin.username", "admin")
	// 注意：默认密码应该通过配置文件或环境变量设置，不在代码中硬编码
	viper.SetDefault("admin.password", "")
	viper.SetDefault("admin.api_key", "")
	viper.SetDefault("admin.auto_create", true)

	// Provider默认配置
	viper.SetDefault("provider.timeout", "30s")
	viper.SetDefault("provider.max_retries", 3)
	viper.SetDefault("provider.retry_delay", "1s")
	viper.SetDefault("provider.max_idle_conns", 100)
	viper.SetDefault("provider.max_conns_per_host", 100)
	viper.SetDefault("provider.idle_conn_timeout", "90s")
}

// validateConfig 验证配置
func validateConfig() error {
	if C.App.Port <= 0 || C.App.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", C.App.Port)
	}

	if C.Storage.Type != "memory" && C.Storage.Type != "gorm" {
		return fmt.Errorf("不支持的存储类型: %s", C.Storage.Type)
	}

	// 如果是gorm存储类型，验证driver配置
	if C.Storage.Type == "gorm" {
		if C.Storage.GORM.Driver != "sqlite" && C.Storage.GORM.Driver != "mysql" {
			return fmt.Errorf("不支持的GORM驱动类型: %s", C.Storage.GORM.Driver)
		}
		if C.Storage.GORM.DSN == "" {
			return fmt.Errorf("GORM DSN不能为空")
		}
	}

	return nil
}

// createDirectories 创建必要的目录
func createDirectories() error {
	dirs := []string{
		"./data",
		"./logs",
		"./temp",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}

	return nil
}

// IsDevelopment 判断是否为开发环境
func IsDevelopment() bool {
	return C.App.Environment == "development"
}

// IsProduction 判断是否为生产环境
func IsProduction() bool {
	return C.App.Environment == "production"
}

// Load 加载配置
func Load() (*Config, error) {
	if err := Init(""); err != nil {
		return nil, err
	}
	return C, nil
}
