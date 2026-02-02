package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/config"
	"github.com/lingproxy/lingproxy/internal/pkg/password"
	"github.com/lingproxy/lingproxy/internal/router"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
	"github.com/lingproxy/lingproxy/pkg/logger"
	mysql "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	logger.Info("Starting application initialization")
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", "error", err)
	}
	logger.Info("Config loaded successfully", "environment", cfg.App.Environment, "port", cfg.App.Port)

	// 检查认证开关状态
	if !cfg.Security.Auth.Enabled {
		logger.Warn("Authentication is DISABLED. All APIs (except login) are accessible without authentication.")
		logger.Warn("This is NOT recommended for production environments.")
	} else {
		logger.Info("Authentication is ENABLED. APIs require authentication.")
	}

	// 初始化日志
	logger.Init(cfg.Log)
	logger.Info("Logger initialized", "level", cfg.Log.Level, "output", cfg.Log.Output)

	// 设置Gin模式
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化存储
	var storageImpl storage.Storage

	switch cfg.Storage.Type {
	case "memory":
		// 使用内存存储
		storageImpl = storage.NewMemoryStorage()
		logger.Info("Using memory storage")

	case "gorm":
		// 使用GORM存储
		var db *gorm.DB
		var err error

		switch cfg.Storage.GORM.Driver {
		case "sqlite":
			db, err = gorm.Open(sqlite.Open(cfg.Storage.GORM.DSN), &gorm.Config{})
		case "mysql":
			db, err = gorm.Open(mysql.Open(cfg.Storage.GORM.DSN), &gorm.Config{})
		default:
			logger.Fatal("Unsupported GORM driver", "driver", cfg.Storage.GORM.Driver)
		}

		if err != nil {
			logger.Fatal("Failed to connect to database", "error", err)
		}

		storageImpl = storage.NewGormStorage(db)
		logger.Info("Using GORM storage", "driver", cfg.Storage.GORM.Driver, "dsn", cfg.Storage.GORM.DSN)

	default:
		logger.Fatal("Unsupported storage type", "type", cfg.Storage.Type)
	}

	// 初始化存储门面
	storageFacade := storage.NewStorageFacade(storageImpl)

	// 初始化管理员用户
	if cfg.Admin.AutoCreate {
		if err := initAdminUser(storageFacade, cfg); err != nil {
			logger.Fatal("Failed to initialize admin user", "error", err)
		}
	}

	// 初始化服务
	logger.Info("Initializing services")
	userService := service.NewUserService(storageFacade)
	policyService := service.NewPolicyService(storageFacade)

	// 初始化内置策略模板
	policyTemplateService := service.NewPolicyTemplateService(storageFacade)
	if err := policyTemplateService.InitBuiltinTemplates(); err != nil {
		logger.Warn("Failed to initialize builtin policy templates", "error", err)
	} else {
		logger.Info("Builtin policy templates initialized successfully")
	}

	logger.Info("Services initialized successfully")

	// 设置路由
	logger.Info("Setting up routes")
	engine := gin.Default()
	router.SetupRoutes(engine, storageFacade, userService, policyService, cfg)
	logger.Info("Routes setup completed")

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: engine,
	}

	// 启动服务器
	go func() {
		logger.Info("Starting server", "port", cfg.App.Port, "mode", cfg.App.Environment, "host", cfg.App.Host)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err, "port", cfg.App.Port)
		}
	}()

	// 等待中断信号优雅关闭
	logger.Info("Server started successfully, waiting for requests")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Received shutdown signal, initiating graceful shutdown")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	} else {
		logger.Info("Server shutdown completed gracefully")
	}

	logger.Info("Application exited")
}

// initAdminUser 初始化管理员用户
func initAdminUser(storageFacade *storage.StorageFacade, cfg *config.Config) error {
	logger.Info("Initializing admin user", "username", cfg.Admin.Username)

	// 检查管理员是否存在
	users, err := storageFacade.ListUsers()
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	// 查找管理员用户
	var adminUser *storage.User
	for _, u := range users {
		if u.Username == cfg.Admin.Username && u.Role == "admin" {
			adminUser = u
			break
		}
	}

	// 如果不存在，创建管理员
	if adminUser == nil {
		adminUser = &storage.User{
			Username: cfg.Admin.Username,
			Role:     "admin",
			Status:   "active",
		}

		// 设置密码
		if cfg.Admin.Password != "" {
			hash, err := password.HashPassword(cfg.Admin.Password)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}
			adminUser.PasswordHash = hash
		}

		// 如果配置了API Key，使用配置的；否则自动生成
		if cfg.Admin.APIKey != "" {
			adminUser.APIKey = cfg.Admin.APIKey
		} else {
			adminUser.APIKey = password.GenerateAPIKey()
		}

		if err := storageFacade.CreateUser(adminUser); err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		logger.Info("Admin user created successfully",
			"username", adminUser.Username,
			"api_key", adminUser.APIKey,
			"id", adminUser.ID)
		logger.Warn("IMPORTANT: Save this API key securely. It will not be shown again.")
		if cfg.Admin.Password == "" {
			logger.Warn("⚠️  WARNING: Admin password is not set in config. Please set it in config.yaml or use the password reset feature.")
		} else {
			logger.Info("Password is set from config. Please change it after first login for security.")
		}
	} else {
		// 如果管理员已存在但没有密码，设置默认密码
		if adminUser.PasswordHash == "" && cfg.Admin.Password != "" {
			hash, err := password.HashPassword(cfg.Admin.Password)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}
			adminUser.PasswordHash = hash
			if err := storageFacade.UpdateUser(adminUser); err != nil {
				return fmt.Errorf("failed to update admin password: %w", err)
			}
			logger.Info("Admin password set successfully", "username", adminUser.Username)
		}
		logger.Info("Admin user already exists", "username", adminUser.Username, "id", adminUser.ID)
	}

	return nil
}
