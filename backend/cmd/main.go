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
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/pkg/password"
	"github.com/lingproxy/lingproxy/internal/router"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
	mysql "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志（在配置加载后立即初始化，确保后续日志都能持久化）
	logger.Init(cfg.Log)
	logger.Info("Starting application initialization")
	logger.Info("Config loaded successfully", logger.F("environment", cfg.App.Environment), logger.F("port", cfg.App.Port))
	logger.Info("Logger initialized", logger.F("level", cfg.Log.Level), logger.F("output", cfg.Log.Output), logger.F("file_path", cfg.Log.FilePath))

	// 检查认证开关状态
	if !cfg.Security.Auth.Enabled {
		logger.Warn("Authentication is DISABLED. All APIs (except login) are accessible without authentication.")
		logger.Warn("This is NOT recommended for production environments.")
	} else {
		logger.Info("Authentication is ENABLED. APIs require authentication.")
	}

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
			logger.Fatal("Unsupported GORM driver", logger.F("driver", cfg.Storage.GORM.Driver))
		}

		if err != nil {
			logger.Fatal("Failed to connect to database", logger.F("error", err))
		}

		storageImpl = storage.NewGormStorage(db)
		logger.Info("Using GORM storage", logger.F("driver", cfg.Storage.GORM.Driver), logger.F("dsn", cfg.Storage.GORM.DSN))

	default:
		logger.Fatal("Unsupported storage type", logger.F("type", cfg.Storage.Type))
	}

	// 初始化存储门面
	storageFacade := storage.NewStorageFacade(storageImpl)

	// 初始化管理员用户（默认用户名: admin, 密码: admin123）
	if err := initAdminUser(storageFacade); err != nil {
		logger.Fatal("Failed to initialize admin user", logger.F("error", err))
	}

	// 初始化服务
	logger.Info("Initializing services")
	userService := service.NewUserService(storageFacade)
	policyService := service.NewPolicyService(storageFacade)

	// 初始化内置策略模板
	policyTemplateService := service.NewPolicyTemplateService(storageFacade)
	if err := policyTemplateService.InitBuiltinTemplates(); err != nil {
		logger.Warn("Failed to initialize builtin policy templates", logger.F("error", err))
	} else {
		logger.Info("Builtin policy templates initialized successfully")
	}

	// 初始化内置策略
	if err := policyService.InitBuiltinPolicies(); err != nil {
		logger.Warn("Failed to initialize builtin policies", logger.F("error", err))
	} else {
		logger.Info("Builtin policies initialized successfully")
	}

	logger.Info("Services initialized successfully")

	// 设置路由
	logger.Info("Setting up routes")
	engine := gin.Default()
	router.SetupRoutes(engine, storageFacade, userService, policyService, cfg)
	logger.Info("Routes setup completed")

	// 创建HTTP服务器
	// 使用配置中的host和port，如果host为空则监听所有接口
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	if cfg.App.Host == "" {
		addr = fmt.Sprintf(":%d", cfg.App.Port)
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// 启动服务器
	go func() {
		logger.Info("Starting server", logger.F("port", cfg.App.Port), logger.F("mode", cfg.App.Environment), logger.F("host", cfg.App.Host))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.F("error", err), logger.F("port", cfg.App.Port))
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
		logger.Error("Server forced to shutdown", logger.F("error", err))
	} else {
		logger.Info("Server shutdown completed gracefully")
	}

	logger.Info("Application exited")
}

// initAdminUser 初始化管理员用户
// 默认用户名: admin, 密码: admin123
func initAdminUser(storageFacade *storage.StorageFacade) error {
	const (
		adminUsername = "admin"
		adminPassword = "admin123"
	)

	logger.Info("Initializing admin user", logger.F("username", adminUsername))

	// 检查管理员是否存在
	users, err := storageFacade.ListUsers()
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	// 查找管理员用户
	var adminUser *storage.User
	for _, u := range users {
		if u.Username == adminUsername && u.Role == "admin" {
			adminUser = u
			break
		}
	}

	// 如果不存在，创建管理员
	if adminUser == nil {
		adminUser = &storage.User{
			Username: adminUsername,
			Role:     "admin",
			Status:   "active",
		}

		// 设置密码
		hash, err := password.HashPassword(adminPassword)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		adminUser.PasswordHash = hash

		// 自动生成API Key
		adminUser.APIKey = password.GenerateAPIKey()

		if err := storageFacade.CreateUser(adminUser); err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		logger.Info("Admin user created successfully",
			logger.F("username", adminUser.Username),
			logger.F("api_key", adminUser.APIKey),
			logger.F("id", adminUser.ID))
		logger.Warn("IMPORTANT: Save this API key securely. It will not be shown again.")
		logger.Info("Default admin credentials: username=admin, password=admin123")
		logger.Warn("⚠️  WARNING: Please change the default password after first login for security.")
	} else {
		// 如果管理员已存在但没有密码，设置默认密码
		if adminUser.PasswordHash == "" {
			hash, err := password.HashPassword(adminPassword)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}
			adminUser.PasswordHash = hash
			if err := storageFacade.UpdateUser(adminUser); err != nil {
				return fmt.Errorf("failed to update admin password: %w", err)
			}
			logger.Info("Admin password set successfully", logger.F("username", adminUser.Username))
		}
		logger.Info("Admin user already exists", logger.F("username", adminUser.Username), logger.F("id", adminUser.ID))
	}

	return nil
}
