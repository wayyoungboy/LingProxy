package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/config"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/service"
)

// SettingsHandler 设置处理器
type SettingsHandler struct {
	settingsService *service.SettingsService
}

// NewSettingsHandler 创建新的设置处理器
func NewSettingsHandler(settingsService *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
	}
}

// GetSettings 获取系统设置
// @Summary Get system settings
// @Description Get all system settings
// @Tags settings
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "System settings"
// @Router /api/v1/settings [get]
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	cfg := config.C

	settings := gin.H{
		"basic": gin.H{
			"system_name":    cfg.App.Name,
			"system_version": cfg.App.Version,
			"api_url":        "http://localhost:8080", // 可以从请求中获取
			"debug_mode":     cfg.App.Environment == "development",
			"environment":    cfg.App.Environment,
			"port":           cfg.App.Port,
			"host":           cfg.App.Host,
		},
		"cache": gin.H{
			"enabled":    cfg.Cache.Enabled,
			"type":       "memory", // 当前只支持内存缓存
			"ttl":        int(cfg.Cache.TTL.Seconds()),
			"size_limit": 100, // 默认值
		},
		"rate_limit": gin.H{
			"enabled":             cfg.Security.RateLimit.Enabled,
			"requests_per_minute": cfg.Security.RateLimit.RequestsPerMinute,
			"concurrency":         50, // 默认值
		},
		"security": gin.H{
			"auth_enabled":     cfg.Security.Auth.Enabled,
			"jwt_secret":       "******", // 脱敏显示
			"token_expiration": int(cfg.Security.JWT.ExpireHours.Hours()),
			"https_enabled":    false, // 当前未实现HTTPS
			"cors": gin.H{
				"enabled":       cfg.Security.CORS.Enabled,
				"allow_origins": cfg.Security.CORS.AllowOrigins,
				"allow_methods": cfg.Security.CORS.AllowMethods,
				"allow_headers": cfg.Security.CORS.AllowHeaders,
			},
		},
		"log": gin.H{
			"level":       cfg.Log.Level,
			"format":      cfg.Log.Format,
			"output":      cfg.Log.Output,
			"file_path":   cfg.Log.FilePath,
			"max_size":    cfg.Log.MaxSize,
			"max_backups": cfg.Log.MaxBackups,
			"max_age":     cfg.Log.MaxAge,
			"compress":    cfg.Log.Compress,
		},
		"load_balancer": gin.H{
			"default_strategy": cfg.LoadBalancer.DefaultStrategy,
			"health_check": gin.H{
				"enabled":      cfg.LoadBalancer.HealthCheck.Enabled,
				"interval":     int(cfg.LoadBalancer.HealthCheck.Interval.Seconds()),
				"timeout":      int(cfg.LoadBalancer.HealthCheck.Timeout.Seconds()),
				"max_failures": cfg.LoadBalancer.HealthCheck.MaxFailures,
			},
		},
		"provider": gin.H{
			"timeout":     int(cfg.Provider.Timeout.Seconds()),
			"max_retries": cfg.Provider.MaxRetries,
			"retry_delay": int(cfg.Provider.RetryDelay.Seconds()),
		},
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// UpdateSettings 更新系统设置
// @Summary Update system settings
// @Description Update system settings (some settings require restart)
// @Tags settings
// @Accept json
// @Produce json
// @Param settings body UpdateSettingsRequest true "Settings to update"
// @Success 200 {object} map[string]interface{} "Update result"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /api/v1/settings [put]
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var req service.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("更新设置失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restartRequiredFields := make([]string, 0)

	// 更新设置
	if err := h.settingsService.UpdateSettings(&req, &restartRequiredFields); err != nil {
		logger.Error("更新设置失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	requiresRestart := len(restartRequiredFields) > 0

	logger.Info("设置更新成功", logger.F("requires_restart", requiresRestart), logger.F("restart_fields", restartRequiredFields))
	c.JSON(http.StatusOK, gin.H{
		"message":                 "设置已更新",
		"requires_restart":        requiresRestart,
		"restart_required_fields": restartRequiredFields,
	})
}
