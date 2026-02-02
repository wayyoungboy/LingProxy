package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/service"
	"github.com/lingproxy/lingproxy/pkg/logger"
)

// TokenHandler Token处理器
type TokenHandler struct {
	tokenService *service.TokenService
}

// NewTokenHandler 创建新的Token处理器
func NewTokenHandler(tokenService *service.TokenService) *TokenHandler {
	return &TokenHandler{
		tokenService: tokenService,
	}
}

// CreateTokenRequest 创建Token请求
type CreateTokenRequest struct {
	Name      string `json:"name" binding:"required"`
	ExpiresAt string `json:"expires_at,omitempty"` // ISO 8601格式
}

// UpdateTokenRequest 更新Token请求
type UpdateTokenRequest struct {
	Name   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"` // active/inactive
}

// CreateToken 创建Token
// @Summary Create a new token
// @Description Create a new API token for authentication
// @Tags tokens
// @Accept json
// @Produce json
// @Param token body CreateTokenRequest true "Token information"
// @Success 201 {object} map[string]interface{} "Created token"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/tokens [post]
func (h *TokenHandler) CreateToken(c *gin.Context) {
	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("创建Token失败", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			logger.Error("解析过期时间失败", "error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format, use ISO 8601"})
			return
		}
		expiresAt = &parsed
	}

	token, err := h.tokenService.CreateToken(req.Name, expiresAt)
	if err != nil {
		if err == service.ErrTokenNameExists {
			c.JSON(http.StatusConflict, gin.H{"error": "token name already exists"})
			return
		}
		logger.Error("创建Token失败", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("创建Token成功", "id", token.ID, "name", token.Name)
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":         token.ID,
			"name":       token.Name,
			"token":      token.Token, // 只在创建时返回完整Token
			"prefix":     token.Prefix,
			"status":     token.Status,
			"expires_at": token.ExpiresAt,
			"created_at": token.CreatedAt,
		},
	})
}

// ListTokens 获取Token列表
// @Summary List all tokens
// @Description Get a list of all API tokens
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of tokens"
// @Router /api/v1/tokens [get]
func (h *TokenHandler) ListTokens(c *gin.Context) {
	tokens, err := h.tokenService.ListTokens()
	if err != nil {
		logger.Error("获取Token列表失败", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 隐藏完整Token值，只显示前缀
	tokenList := make([]gin.H, 0, len(tokens))
	for _, token := range tokens {
		tokenList = append(tokenList, gin.H{
			"id":           token.ID,
			"name":         token.Name,
			"token":        token.Prefix, // 只显示前缀
			"status":       token.Status,
			"policy_id":    token.PolicyID,
			"last_used_at": token.LastUsedAt,
			"expires_at":   token.ExpiresAt,
			"created_at":   token.CreatedAt,
		})
	}

	logger.Info("获取Token列表成功", "count", len(tokens))
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"items": tokenList, "total": len(tokens)}})
}

// GetToken 获取单个Token
// @Summary Get token by ID
// @Description Get a specific token by ID
// @Tags tokens
// @Accept json
// @Produce json
// @Param id path string true "Token ID"
// @Success 200 {object} map[string]interface{} "Token details"
// @Failure 404 {object} map[string]string "Token not found"
// @Router /api/v1/tokens/{id} [get]
func (h *TokenHandler) GetToken(c *gin.Context) {
	id := c.Param("id")
	token, err := h.tokenService.GetToken(id)
	if err != nil {
		if err == service.ErrTokenNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		logger.Error("获取Token失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":           token.ID,
			"name":         token.Name,
			"token":        token.Prefix, // 只显示前缀
			"status":       token.Status,
			"policy_id":    token.PolicyID,
			"last_used_at": token.LastUsedAt,
			"expires_at":   token.ExpiresAt,
			"created_at":   token.CreatedAt,
		},
	})
}

// UpdateToken 更新Token
// @Summary Update token
// @Description Update token information
// @Tags tokens
// @Accept json
// @Produce json
// @Param id path string true "Token ID"
// @Param token body UpdateTokenRequest true "Token update information"
// @Success 200 {object} map[string]interface{} "Updated token"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Token not found"
// @Router /api/v1/tokens/{id} [put]
func (h *TokenHandler) UpdateToken(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("更新Token失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.tokenService.UpdateToken(id, req.Name, req.Status)
	if err != nil {
		if err == service.ErrTokenNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		if err == service.ErrTokenNameExists {
			c.JSON(http.StatusConflict, gin.H{"error": "token name already exists"})
			return
		}
		logger.Error("更新Token失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("更新Token成功", "id", id, "name", token.Name)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":           token.ID,
			"name":         token.Name,
			"token":        token.Prefix,
			"status":       token.Status,
			"last_used_at": token.LastUsedAt,
			"expires_at":   token.ExpiresAt,
			"updated_at":   token.UpdatedAt,
		},
	})
}

// DeleteToken 删除Token
// @Summary Delete token
// @Description Delete a token
// @Tags tokens
// @Accept json
// @Produce json
// @Param id path string true "Token ID"
// @Success 204 "Token deleted"
// @Failure 404 {object} map[string]string "Token not found"
// @Router /api/v1/tokens/{id} [delete]
func (h *TokenHandler) DeleteToken(c *gin.Context) {
	id := c.Param("id")
	if err := h.tokenService.DeleteToken(id); err != nil {
		if err == service.ErrTokenNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		logger.Error("删除Token失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("删除Token成功", "id", id)
	c.Status(http.StatusNoContent)
}

// ResetToken 重置Token
// @Summary Reset token
// @Description Reset token value (generate new token)
// @Tags tokens
// @Accept json
// @Produce json
// @Param id path string true "Token ID"
// @Success 200 {object} map[string]interface{} "Reset token with new value"
// @Failure 404 {object} map[string]string "Token not found"
// @Router /api/v1/tokens/{id}/reset [post]
func (h *TokenHandler) ResetToken(c *gin.Context) {
	id := c.Param("id")
	token, err := h.tokenService.ResetToken(id)
	if err != nil {
		if err == service.ErrTokenNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		logger.Error("重置Token失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("重置Token成功", "id", id, "name", token.Name)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":         token.ID,
			"name":       token.Name,
			"token":      token.Token, // 重置时返回完整Token
			"prefix":     token.Prefix,
			"status":     token.Status,
			"policy_id":  token.PolicyID,
			"updated_at": token.UpdatedAt,
		},
	})
}

// SetTokenPolicyRequest 设置Token策略请求
type SetTokenPolicyRequest struct {
	PolicyID string `json:"policy_id"`
}

// SetTokenPolicy 设置Token的策略
// @Summary Set token policy
// @Description Assign a policy to a token
// @Tags tokens
// @Accept json
// @Produce json
// @Param id path string true "Token ID"
// @Param request body SetTokenPolicyRequest true "Policy assignment"
// @Success 200 {object} map[string]interface{} "Token with policy"
// @Failure 404 {object} map[string]string "Token or policy not found"
// @Router /api/v1/tokens/{id}/policy [put]
func (h *TokenHandler) SetTokenPolicy(c *gin.Context) {
	id := c.Param("id")
	var req SetTokenPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("设置Token策略失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.tokenService.UpdateTokenPolicy(id, req.PolicyID)
	if err != nil {
		if err == service.ErrTokenNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		logger.Error("设置Token策略失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("设置Token策略成功", "id", id, "policy_id", req.PolicyID)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        token.ID,
			"name":      token.Name,
			"policy_id": token.PolicyID,
			"status":    token.Status,
		},
	})
}

// RemoveTokenPolicy 移除Token的策略
// @Summary Remove token policy
// @Description Remove policy assignment from a token
// @Tags tokens
// @Accept json
// @Produce json
// @Param id path string true "Token ID"
// @Success 200 {object} map[string]interface{} "Token without policy"
// @Failure 404 {object} map[string]string "Token not found"
// @Router /api/v1/tokens/{id}/policy [delete]
func (h *TokenHandler) RemoveTokenPolicy(c *gin.Context) {
	id := c.Param("id")
	token, err := h.tokenService.RemoveTokenPolicy(id)
	if err != nil {
		if err == service.ErrTokenNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		logger.Error("移除Token策略失败", "error", err.Error(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("移除Token策略成功", "id", id)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        token.ID,
			"name":      token.Name,
			"policy_id": token.PolicyID,
			"status":    token.Status,
		},
	})
}
