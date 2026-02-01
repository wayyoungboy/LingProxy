package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lingproxy/lingproxy/docs"
	"github.com/lingproxy/lingproxy/internal/handler"
	"github.com/lingproxy/lingproxy/internal/middleware"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
	"github.com/lingproxy/lingproxy/pkg/logger"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, storage *storage.StorageFacade, userService *service.UserService, policyService *service.PolicyService) {
	logger.Info("Starting route setup")

	// 创建处理器
	logger.Info("Initializing handlers")
	userHandler := handler.NewUserHandler(storage)
	llmResourceHandler := handler.NewLLMResourceHandler(storage)
	modelHandler := handler.NewModelHandler(storage)
	endpointHandler := handler.NewEndpointHandler(storage)
	requestHandler := handler.NewRequestHandler(storage)
	statsHandler := handler.NewStatsHandler(storage)
	openaiHandler := handler.NewOpenAIHandler(storage)
	authMiddleware := middleware.NewAuthMiddleware(storage)
	logger.Info("Handlers initialized successfully")

	// Swagger文档
	logger.Debug("Adding Swagger documentation route")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// OpenAI兼容API路由组
	logger.Info("Adding OpenAI compatible API routes")
	openai := r.Group("/v1")
	{
		logger.Debug("Adding OpenAI models routes")
		openai.GET("/models", openaiHandler.ListModels)
		openai.GET("/models/:model", openaiHandler.GetModel)
		logger.Debug("Adding OpenAI chat completions route")
		openai.POST("/chat/completions", openaiHandler.CreateChatCompletion)
		logger.Debug("Adding OpenAI completions route")
		openai.POST("/completions", openaiHandler.CreateCompletion)
	}

	// API路由组
	logger.Info("Adding API routes")
	api := r.Group("/api/v1")
	{
		// 公开路由
		logger.Debug("Adding health check route")
		api.GET("/health", handler.HealthHandler)

		// 用户路由
		logger.Debug("Adding user management routes")
		api.GET("/users", userHandler.ListUsers)
		api.GET("/users/:id", userHandler.GetUser)
		api.POST("/users", userHandler.CreateUser)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)

		// LLM资源路由
		logger.Debug("Adding LLM resource management routes")
		api.GET("/llm-resources", llmResourceHandler.ListLLMResources)
		api.GET("/llm-resources/:id", llmResourceHandler.GetLLMResource)
		api.POST("/llm-resources", llmResourceHandler.CreateLLMResource)
		api.PUT("/llm-resources/:id", llmResourceHandler.UpdateLLMResource)
		api.DELETE("/llm-resources/:id", llmResourceHandler.DeleteLLMResource)

		// 模型路由
		logger.Debug("Adding model management routes")
		api.GET("/models", modelHandler.ListModels)
		api.GET("/models/types", modelHandler.ListModelTypes)
		api.GET("/models/:id", modelHandler.GetModel)
		api.POST("/models", modelHandler.CreateModel)
		api.PUT("/models/:id", modelHandler.UpdateModel)
		api.DELETE("/models/:id", modelHandler.DeleteModel)
		api.GET("/models/:id/pricing", modelHandler.GetModelPricing)
		api.GET("/llm-resources/:id/models", modelHandler.ListModelsByLLMResource)

		// 端点路由
		logger.Debug("Adding endpoint management routes")
		api.GET("/endpoints", endpointHandler.ListEndpoints)
		api.GET("/endpoints/:id", endpointHandler.GetEndpoint)
		api.POST("/endpoints", endpointHandler.CreateEndpoint)
		api.PUT("/endpoints/:id", endpointHandler.UpdateEndpoint)
		api.DELETE("/endpoints/:id", endpointHandler.DeleteEndpoint)

		// 请求路由
		logger.Debug("Adding request management routes")
		api.GET("/requests", requestHandler.ListRequests)
		api.GET("/requests/:id", requestHandler.GetRequest)
		api.POST("/requests", requestHandler.CreateRequest)

		// 统计路由
		logger.Debug("Adding statistics routes")
		api.GET("/stats/system", statsHandler.GetSystemStats)
		api.GET("/stats/llm-resources/:id", statsHandler.GetLLMResourceStats)
		api.GET("/stats/users/:id", statsHandler.GetUserStats)

		// 需要认证的路由
		logger.Debug("Adding authenticated routes")
		auth := api.Group("")
		auth.Use(authMiddleware.RequireAuth())
		{
			// 代理路由
			logger.Debug("Adding proxy endpoint route")
			auth.Any("/proxy/*path", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "proxy endpoint", "path": c.Param("path")})
			})
		}
	}



	logger.Info("Route setup completed successfully")
}
