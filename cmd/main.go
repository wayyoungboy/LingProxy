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

	// 初始化服务
	logger.Info("Initializing services")
	userService := service.NewUserService(storageFacade)
	policyService := service.NewPolicyService(storageFacade)
	logger.Info("Services initialized successfully")

	// 设置路由
	logger.Info("Setting up routes")
	engine := gin.Default()
	router.SetupRoutes(engine, storageFacade, userService, policyService)
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
