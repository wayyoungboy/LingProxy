package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/lingproxy/lingproxy/internal/config"
	"github.com/lingproxy/lingproxy/internal/handler"
	"github.com/lingproxy/lingproxy/internal/middleware"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
	_ "github.com/lingproxy/lingproxy/swagger"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, storage *storage.StorageFacade, userService *service.UserService, policyService *service.PolicyService, cfg *config.Config) {
	logger.Info("Starting route setup")

	// 添加全局中间件：RequestID、RequestLogger 和 CORS（必须在最前面）
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS())

	// 创建处理器
	logger.Info("Initializing handlers")
	apiKeyService := service.NewAPIKeyService(storage)
	templateService := service.NewPolicyTemplateService(storage)
	settingsService := service.NewSettingsService("configs/config.yaml")
	// policyService 已通过参数传入
	adminHandler := handler.NewAdminHandler(storage)
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeyService)
	tokenHandler := apiKeyHandler // 保持向后兼容
	tokenService := apiKeyService // 保持向后兼容
	policyHandler := handler.NewPolicyHandler(policyService, templateService)
	settingsHandler := handler.NewSettingsHandler(settingsService)
	systemHandler := handler.NewSystemHandler()
	logHandler := handler.NewLogHandler()
	llmResourceHandler := handler.NewLLMResourceHandler(storage)
	modelHandler := handler.NewModelHandler(storage)
	requestHandler := handler.NewRequestHandler(storage)
	statsHandler := handler.NewStatsHandler(storage)
	openaiHandler := handler.NewOpenAIHandler(storage, policyService, tokenService)
	authMiddleware := middleware.NewAuthMiddleware(storage, tokenService, cfg)
	logger.Info("Handlers initialized successfully")

	// Swagger文档
	logger.Debug("Adding Swagger documentation route")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// OpenAI兼容API路由组（根据认证开关决定是否需要认证）
	logger.Info("Adding OpenAI compatible API routes")
	openai := r.Group("/llm/v1")
	if cfg.Security.Auth.Enabled {
		openai.Use(authMiddleware.RequireAuth())
		logger.Info("OpenAI API routes require authentication")
	} else {
		logger.Warn("OpenAI API routes do NOT require authentication (auth disabled)")
	}
	{
		logger.Debug("Adding OpenAI models routes")
		openai.GET("/models", openaiHandler.ListModels)
		openai.GET("/models/:model", openaiHandler.GetModel)
		logger.Debug("Adding OpenAI chat completions route")
		openai.POST("/chat/completions", openaiHandler.CreateChatCompletion)
		logger.Debug("Adding OpenAI completions route")
		openai.POST("/completions", openaiHandler.CreateCompletion)
		logger.Debug("Adding OpenAI embeddings route")
		openai.POST("/embeddings", openaiHandler.CreateEmbedding)
	}

	// API路由组
	logger.Info("Adding API routes")
	api := r.Group("/api/v1")
	{
		// 公开路由
		logger.Debug("Adding health check route")
		api.GET("/health", handler.HealthHandler)

		// 认证路由（不需要认证）
		logger.Debug("Adding auth routes")
		userHandler := handler.NewUserHandler(storage, userService)
		api.POST("/auth/login", userHandler.Login)

		// 需要认证的路由
		logger.Debug("Adding authenticated routes")
		auth := api.Group("")
		auth.Use(authMiddleware.RequireAuth())
		{
			// 管理员路由
			logger.Debug("Adding admin routes")
			auth.GET("/admin/info", adminHandler.GetAdminInfo)
			auth.PUT("/admin/info", adminHandler.UpdateAdminInfo)    // 更新管理员信息（用户名和/或密码）
			auth.PUT("/admin/username", adminHandler.UpdateUsername) // 更新管理员用户名
			auth.PUT("/admin/password", adminHandler.UpdatePassword) // 更新管理员密码
			auth.PUT("/admin/api-key", adminHandler.ResetAPIKey)     // 重置API Key

			// API Key管理路由
			logger.Debug("Adding API key management routes")
			auth.GET("/api-keys", tokenHandler.ListAPIKeys)
			auth.GET("/api-keys/:id", tokenHandler.GetAPIKey)
			auth.POST("/api-keys", tokenHandler.CreateAPIKey)
			auth.PUT("/api-keys/:id", tokenHandler.UpdateAPIKey)
			auth.DELETE("/api-keys/:id", tokenHandler.DeleteAPIKey)
			auth.POST("/api-keys/:id/reset", tokenHandler.ResetAPIKey)
			auth.PUT("/api-keys/:id/policy", tokenHandler.SetAPIKeyPolicy)
			auth.DELETE("/api-keys/:id/policy", tokenHandler.RemoveAPIKeyPolicy)
			// 保持向后兼容的旧路径
			auth.GET("/tokens", tokenHandler.ListTokens)
			auth.GET("/tokens/:id", tokenHandler.GetToken)
			auth.POST("/tokens", tokenHandler.CreateToken)
			auth.PUT("/tokens/:id", tokenHandler.UpdateToken)
			auth.DELETE("/tokens/:id", tokenHandler.DeleteToken)
			auth.POST("/tokens/:id/reset", tokenHandler.ResetToken)
			auth.PUT("/tokens/:id/policy", tokenHandler.SetTokenPolicy)
			auth.DELETE("/tokens/:id/policy", tokenHandler.RemoveTokenPolicy)

			// 策略模板路由
			logger.Debug("Adding policy template routes")
			auth.GET("/policy-templates", policyHandler.ListPolicyTemplates)
			auth.GET("/policy-templates/:id", policyHandler.GetPolicyTemplate)

			// 策略管理路由
			logger.Debug("Adding policy management routes")
			auth.GET("/policies", policyHandler.ListPolicies)
			auth.GET("/policies/:id", policyHandler.GetPolicy)
			auth.POST("/policies", policyHandler.CreatePolicy)
			auth.PUT("/policies/:id", policyHandler.UpdatePolicy)
			auth.DELETE("/policies/:id", policyHandler.DeletePolicy)

			// 系统设置路由
			logger.Debug("Adding settings routes")
			auth.GET("/settings", settingsHandler.GetSettings)
			auth.PUT("/settings", settingsHandler.UpdateSettings)

			// 系统信息路由
			logger.Debug("Adding system info routes")
			auth.GET("/system/info", systemHandler.GetSystemInfo)

			// 日志管理路由
			logger.Debug("Adding log management routes")
			auth.GET("/logs/files", logHandler.ListLogFiles)
			auth.GET("/logs", logHandler.GetLogs)
			auth.GET("/logs/files/:file/download", logHandler.DownloadLogFile)
			auth.POST("/logs/clear", logHandler.ClearLogs)

			// 代理路由
			logger.Debug("Adding proxy endpoint route")
			auth.Any("/proxy/*path", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "proxy endpoint", "path": c.Param("path")})
			})
		}

		// LLM资源路由（根据认证开关决定是否需要认证）
		logger.Debug("Adding LLM resource management routes")
		llmResourceRoutes := api.Group("/llm-resources")
		if cfg.Security.Auth.Enabled {
			llmResourceRoutes.Use(authMiddleware.RequireAuth())
		}
		llmResourceRoutes.GET("", llmResourceHandler.ListLLMResources)
		llmResourceRoutes.GET("/:id", llmResourceHandler.GetLLMResource)
		llmResourceRoutes.POST("", llmResourceHandler.CreateLLMResource)
		llmResourceRoutes.PUT("/:id", llmResourceHandler.UpdateLLMResource)
		llmResourceRoutes.DELETE("/:id", llmResourceHandler.DeleteLLMResource)
		llmResourceRoutes.POST("/:id/test", llmResourceHandler.TestLLMResource)
		llmResourceRoutes.POST("/import", llmResourceHandler.ImportLLMResources)
		llmResourceRoutes.GET("/import/template", llmResourceHandler.DownloadImportTemplate)

		// 模型路由（根据认证开关决定是否需要认证）
		logger.Debug("Adding model management routes")
		modelRoutes := api.Group("/models")
		if cfg.Security.Auth.Enabled {
			modelRoutes.Use(authMiddleware.RequireAuth())
		}
		modelRoutes.GET("", modelHandler.ListModels)
		modelRoutes.GET("/types", modelHandler.ListModelTypes)
		modelRoutes.GET("/:id", modelHandler.GetModel)
		modelRoutes.POST("", modelHandler.CreateModel)
		modelRoutes.PUT("/:id", modelHandler.UpdateModel)
		modelRoutes.DELETE("/:id", modelHandler.DeleteModel)
		modelRoutes.GET("/:id/pricing", modelHandler.GetModelPricing)
		// 注意：这个路由在llm-resources组下，需要单独处理
		if cfg.Security.Auth.Enabled {
			api.GET("/llm-resources/:id/models", authMiddleware.RequireAuth(), modelHandler.ListModelsByLLMResource)
		} else {
			api.GET("/llm-resources/:id/models", modelHandler.ListModelsByLLMResource)
		}

		// 请求路由（根据认证开关决定是否需要认证）
		logger.Debug("Adding request management routes")
		requestRoutes := api.Group("/requests")
		if cfg.Security.Auth.Enabled {
			requestRoutes.Use(authMiddleware.RequireAuth())
		}
		requestRoutes.GET("", requestHandler.ListRequests)
		requestRoutes.GET("/:id", requestHandler.GetRequest)
		requestRoutes.POST("", requestHandler.CreateRequest)

		// 统计路由（根据认证开关决定是否需要认证）
		logger.Debug("Adding statistics routes")
		statsRoutes := api.Group("/stats")
		if cfg.Security.Auth.Enabled {
			statsRoutes.Use(authMiddleware.RequireAuth())
		}
		statsRoutes.GET("/system", statsHandler.GetSystemStats)
		statsRoutes.GET("/llm-resources/:id", statsHandler.GetLLMResourceStats)
		statsRoutes.GET("/llm-resources/usage", statsHandler.GetLLMResourceUsageStats)
		statsRoutes.GET("/users/:id", statsHandler.GetUserStats)
	}

	logger.Info("Route setup completed successfully")
}
