package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/pkg/password"
	"github.com/lingproxy/lingproxy/internal/storage"
	"github.com/lingproxy/lingproxy/pkg/logger"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	storage *storage.StorageFacade
}

// NewAdminHandler 创建新的管理员处理器
func NewAdminHandler(storage *storage.StorageFacade) *AdminHandler {
	return &AdminHandler{
		storage: storage,
	}
}

// GetAdminInfo 获取管理员信息
// @Summary Get admin information
// @Description Get admin user information
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Admin information"
// @Failure 404 {object} map[string]string "Admin not found"
// @Router /api/v1/admin/info [get]
func (h *AdminHandler) GetAdminInfo(c *gin.Context) {
	users, err := h.storage.ListUsers()
	if err != nil {
		logger.Error("获取管理员信息失败", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查找管理员用户
	var adminUser *storage.User
	for _, u := range users {
		if u.Role == "admin" {
			adminUser = u
			break
		}
	}

	if adminUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin user not found"})
		return
	}

	// 隐藏敏感信息
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        adminUser.ID,
			"username":  adminUser.Username,
			"api_key":   getPrefix(adminUser.APIKey), // 只显示前缀
			"role":      adminUser.Role,
			"status":    adminUser.Status,
			"created_at": adminUser.CreatedAt,
		},
	})
}

// ResetAPIKeyRequest 重置API Key请求
type ResetAPIKeyRequest struct {
	// 可以添加确认字段等
}

// ResetAPIKey 重置管理员API Key
// @Summary Reset admin API key
// @Description Reset admin user API key
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "New API key"
// @Failure 404 {object} map[string]string "Admin not found"
// @Router /api/v1/admin/api-key [put]
func (h *AdminHandler) ResetAPIKey(c *gin.Context) {
	users, err := h.storage.ListUsers()
	if err != nil {
		logger.Error("重置管理员API Key失败", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查找管理员用户
	var adminUser *storage.User
	for _, u := range users {
		if u.Role == "admin" {
			adminUser = u
			break
		}
	}

	if adminUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin user not found"})
		return
	}

	// 生成新的API Key
	newAPIKey := password.GenerateAPIKey()
	adminUser.APIKey = newAPIKey

	if err := h.storage.UpdateUser(adminUser); err != nil {
		logger.Error("重置管理员API Key失败", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("重置管理员API Key成功", "username", adminUser.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "API key reset successfully",
		"data": gin.H{
			"api_key": newAPIKey, // 重置时返回完整API Key
		},
	})
}

// getPrefix 获取Token前缀
func getPrefix(tokenValue string) string {
	if len(tokenValue) > 12 {
		return tokenValue[:12] + "..."
	}
	return tokenValue
}
