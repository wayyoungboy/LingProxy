package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// UserHandler 用户处理器
type UserHandler struct {
	storage     *storage.StorageFacade
	userService *service.UserService
}

// NewUserHandler 创建新的用户处理器
func NewUserHandler(storage *storage.StorageFacade, userService *service.UserService) *UserHandler {
	return &UserHandler{
		storage:     storage,
		userService: userService,
	}
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of users"
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.storage.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetUser 获取单个用户
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.storage.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &storage.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := h.userService.CreateUser(user); err != nil {
		if err == service.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 不返回密码哈希
	user.PasswordHash = ""
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"api_key":  user.APIKey,
			"role":     user.Role,
			"status":   user.Status,
		},
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 返回用户信息和token（使用API Key作为token）
	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"data": gin.H{
			"token": user.APIKey,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"role":     user.Role,
				"status":   user.Status,
			},
		},
	})
}

// CreateUser 创建用户（管理员功能）
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &storage.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := h.userService.CreateUser(user); err != nil {
		if err == service.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user storage.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	if err := h.storage.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// ResetAPIKey 重置API Key
func (h *UserHandler) ResetAPIKey(c *gin.Context) {
	id := c.Param("id")

	newAPIKey, err := h.userService.ResetAPIKey(id)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key reset successfully",
		"data": gin.H{
			"api_key": newAPIKey,
		},
	})
}

// UpdatePassword 更新密码
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.UpdatePassword(id, req.OldPassword, req.NewPassword); err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if err == service.ErrInvalidPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid old password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

// GetCurrentUser 获取当前用户信息
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	u := user.(*storage.User)
	u.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"data": u})
}
