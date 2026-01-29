package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupWebRoutes 设置Web界面路由
func SetupWebRoutes(router *gin.Engine) {
	// 静态文件服务
	router.Static("/static", "./web/static")
	router.StaticFile("/favicon.ico", "./web/static/images/favicon.ico")

	// Web管理界面
	router.GET("/admin", func(c *gin.Context) {
		c.File("./web/templates/index.html")
	})

	// 登录页面
	router.GET("/login", func(c *gin.Context) {
		c.File("./web/templates/login.html")
	})

	// 默认重定向到管理界面
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/admin")
	})

	// API文档页面
	router.GET("/docs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API文档",
			"endpoints": []string{
				"GET /api/v1/health",
				"GET /api/v1/user/profile",
				"PUT /api/v1/user/profile",
				"GET /api/v1/user/quota",
				"GET /api/v1/admin/users",
				"POST /api/v1/admin/users",
				"GET /api/v1/admin/providers",
				"POST /api/v1/admin/providers",
			},
		})
	})
}