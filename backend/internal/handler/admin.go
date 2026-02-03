package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/pkg/password"
	"github.com/lingproxy/lingproxy/internal/storage"
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
		logger.Error("获取管理员信息失败", logger.F("error", err.Error()))
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
			"id":         adminUser.ID,
			"username":   adminUser.Username,
			"api_key":    getPrefix(adminUser.APIKey), // 只显示前缀
			"role":       adminUser.Role,
			"status":     adminUser.Status,
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
		logger.Error("重置管理员API Key失败", logger.F("error", err.Error()))
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
		logger.Error("重置管理员API Key失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("重置管理员API Key成功", logger.F("username", adminUser.Username))
	c.JSON(http.StatusOK, gin.H{
		"message": "API key reset successfully",
		"data": gin.H{
			"api_key": newAPIKey, // 重置时返回完整API Key
		},
	})
}

// UpdateAdminPasswordRequest 更新管理员密码请求
type UpdateAdminPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdatePassword 更新管理员密码
// @Summary Update admin password
// @Description Update admin user password (password is stored as hash in database)
// @Tags admin
// @Accept json
// @Produce json
// @Param request body UpdateAdminPasswordRequest true "Password update request"
// @Success 200 {object} map[string]string "Password updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Admin not found"
// @Router /api/v1/admin/password [put]
func (h *AdminHandler) UpdatePassword(c *gin.Context) {
	var req UpdateAdminPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新管理员密码失败：请求参数无效", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.storage.ListUsers()
	if err != nil {
		logger.Error("更新管理员密码失败", logger.F("error", err.Error()))
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
		logger.Warn("更新管理员密码失败：管理员不存在")
		c.JSON(http.StatusNotFound, gin.H{"error": "admin user not found"})
		return
	}

	// 验证旧密码
	if adminUser.PasswordHash != "" {
		valid, err := password.VerifyPassword(req.OldPassword, adminUser.PasswordHash)
		if err != nil || !valid {
			logger.Warn("更新管理员密码失败：旧密码不正确", logger.F("username", adminUser.Username))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid old password"})
			return
		}
	}

	// 加密新密码
	hash, err := password.HashPassword(req.NewPassword)
	if err != nil {
		logger.Error("更新管理员密码失败：密码加密失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	adminUser.PasswordHash = hash
	if err := h.storage.UpdateUser(adminUser); err != nil {
		logger.Error("更新管理员密码失败", logger.F("error", err.Error()), logger.F("username", adminUser.Username))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("更新管理员密码成功", logger.F("username", adminUser.Username))
	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

// UpdateAdminUsernameRequest 更新管理员用户名请求
type UpdateAdminUsernameRequest struct {
	Password string `json:"password" binding:"required"` // 需要密码验证
	Username string `json:"username" binding:"required,min=3,max=50"`
}

// UpdateUsername 更新管理员用户名
// @Summary Update admin username
// @Description Update admin user username (requires password verification)
// @Tags admin
// @Accept json
// @Produce json
// @Param request body UpdateAdminUsernameRequest true "Username update request"
// @Success 200 {object} map[string]string "Username updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Admin not found"
// @Router /api/v1/admin/username [put]
func (h *AdminHandler) UpdateUsername(c *gin.Context) {
	var req UpdateAdminUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新管理员用户名失败：请求参数无效", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证用户名格式
	if len(req.Username) < 3 || len(req.Username) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be between 3 and 50 characters"})
		return
	}

	users, err := h.storage.ListUsers()
	if err != nil {
		logger.Error("更新管理员用户名失败", logger.F("error", err.Error()))
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
		logger.Warn("更新管理员用户名失败：管理员不存在")
		c.JSON(http.StatusNotFound, gin.H{"error": "admin user not found"})
		return
	}

	// 验证密码
	if adminUser.PasswordHash != "" {
		valid, err := password.VerifyPassword(req.Password, adminUser.PasswordHash)
		if err != nil || !valid {
			logger.Warn("更新管理员用户名失败：密码不正确", logger.F("username", adminUser.Username))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
			return
		}
	} else {
		// 如果没有设置密码，不允许更新用户名
		logger.Warn("更新管理员用户名失败：管理员未设置密码")
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be set before updating username"})
		return
	}

	// 检查新用户名是否已被使用
	for _, u := range users {
		if u.Username == req.Username && u.ID != adminUser.ID {
			logger.Warn("更新管理员用户名失败：用户名已被使用", logger.F("username", req.Username))
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}
	}

	oldUsername := adminUser.Username
	adminUser.Username = req.Username

	if err := h.storage.UpdateUser(adminUser); err != nil {
		logger.Error("更新管理员用户名失败", logger.F("error", err.Error()), logger.F("old_username", oldUsername), logger.F("new_username", req.Username))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("更新管理员用户名成功", logger.F("old_username", oldUsername), logger.F("new_username", req.Username))
	c.JSON(http.StatusOK, gin.H{
		"message": "username updated successfully",
		"data": gin.H{
			"username": adminUser.Username,
		},
	})
}

// UpdateAdminInfoRequest 更新管理员信息请求（可同时更新用户名和密码）
type UpdateAdminInfoRequest struct {
	Password    string `json:"password" binding:"required"`                      // 当前密码（用于验证）
	NewUsername string `json:"new_username,omitempty"`                           // 新用户名（可选）
	NewPassword string `json:"new_password,omitempty" binding:"omitempty,min=6"` // 新密码（可选，至少6位）
}

// UpdateAdminInfo 更新管理员信息（可同时更新用户名和密码）
// @Summary Update admin information
// @Description Update admin user information (username and/or password, requires current password verification)
// @Tags admin
// @Accept json
// @Produce json
// @Param request body UpdateAdminInfoRequest true "Admin info update request"
// @Success 200 {object} map[string]interface{} "Admin info updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Admin not found"
// @Router /api/v1/admin/info [put]
func (h *AdminHandler) UpdateAdminInfo(c *gin.Context) {
	var req UpdateAdminInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新管理员信息失败：请求参数无效", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 至少需要更新一项
	if req.NewUsername == "" && req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one of new_username or new_password must be provided"})
		return
	}

	users, err := h.storage.ListUsers()
	if err != nil {
		logger.Error("更新管理员信息失败", logger.F("error", err.Error()))
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
		logger.Warn("更新管理员信息失败：管理员不存在")
		c.JSON(http.StatusNotFound, gin.H{"error": "admin user not found"})
		return
	}

	// 验证当前密码
	if adminUser.PasswordHash != "" {
		valid, err := password.VerifyPassword(req.Password, adminUser.PasswordHash)
		if err != nil || !valid {
			logger.Warn("更新管理员信息失败：密码不正确", logger.F("username", adminUser.Username))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
			return
		}
	} else {
		// 如果没有设置密码，不允许更新
		logger.Warn("更新管理员信息失败：管理员未设置密码")
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be set before updating admin info"})
		return
	}

	oldUsername := adminUser.Username
	updated := false

	// 更新用户名
	if req.NewUsername != "" {
		// 验证用户名格式
		if len(req.NewUsername) < 3 || len(req.NewUsername) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username must be between 3 and 50 characters"})
			return
		}

		// 检查新用户名是否已被使用
		for _, u := range users {
			if u.Username == req.NewUsername && u.ID != adminUser.ID {
				logger.Warn("更新管理员信息失败：用户名已被使用", logger.F("username", req.NewUsername))
				c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
				return
			}
		}

		adminUser.Username = req.NewUsername
		updated = true
	}

	// 更新密码
	if req.NewPassword != "" {
		hash, err := password.HashPassword(req.NewPassword)
		if err != nil {
			logger.Error("更新管理员信息失败：密码加密失败", logger.F("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		adminUser.PasswordHash = hash
		updated = true
	}

	if !updated {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no changes to update"})
		return
	}

	if err := h.storage.UpdateUser(adminUser); err != nil {
		logger.Error("更新管理员信息失败", logger.F("error", err.Error()), logger.F("old_username", oldUsername))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("更新管理员信息成功",
		logger.F("old_username", oldUsername),
		logger.F("new_username", adminUser.Username),
		logger.F("password_updated", req.NewPassword != ""))

	c.JSON(http.StatusOK, gin.H{
		"message": "admin info updated successfully",
		"data": gin.H{
			"username":         adminUser.Username,
			"password_updated": req.NewPassword != "",
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
