package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/service"
)

// APIKeyHandler API Key处理器
type APIKeyHandler struct {
	apiKeyService *service.APIKeyService
}

// NewAPIKeyHandler 创建新的API Key处理器
func NewAPIKeyHandler(apiKeyService *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// TokenHandler 保持向后兼容的类型别名
// Deprecated: 使用 APIKeyHandler 代替
type TokenHandler = APIKeyHandler

// NewTokenHandler 保持向后兼容的函数别名
// Deprecated: 使用 NewAPIKeyHandler 代替
func NewTokenHandler(tokenService *service.APIKeyService) *APIKeyHandler {
	return NewAPIKeyHandler(tokenService)
}

// CreateTokenRequest 创建API Key请求
type CreateTokenRequest struct {
	Name          string   `json:"name" binding:"required"`
	ExpiresAt     string   `json:"expires_at,omitempty"`     // ISO 8601格式
	AllowedModels []string `json:"allowed_models,omitempty"` // 允许使用的模型ID列表（空列表表示允许所有模型）
	
	// 按类型配置的策略
	ChatPolicyID     string `json:"chat_policy_id,omitempty"`
	EmbeddingPolicyID string `json:"embedding_policy_id,omitempty"`
	RerankPolicyID    string `json:"rerank_policy_id,omitempty"`
	ImagePolicyID     string `json:"image_policy_id,omitempty"`
	AudioPolicyID     string `json:"audio_policy_id,omitempty"`
	VideoPolicyID     string `json:"video_policy_id,omitempty"`
}

// UpdateTokenRequest 更新API Key请求
type UpdateTokenRequest struct {
	Name          *string  `json:"name,omitempty"`
	Status        *string  `json:"status,omitempty"`         // active/inactive
	AllowedModels []string `json:"allowed_models,omitempty"` // 允许使用的模型ID列表
	
	// 按类型配置的策略
	ChatPolicyID     *string `json:"chat_policy_id,omitempty"`
	EmbeddingPolicyID *string `json:"embedding_policy_id,omitempty"`
	RerankPolicyID    *string `json:"rerank_policy_id,omitempty"`
	ImagePolicyID     *string `json:"image_policy_id,omitempty"`
	AudioPolicyID     *string `json:"audio_policy_id,omitempty"`
	VideoPolicyID     *string `json:"video_policy_id,omitempty"`
}

// CreateAPIKey 创建API Key
// @Summary Create a new API key
// @Description Create a new API key for authentication
// @Tags api-keys
// @Accept json
// @Produce json
// @Param apiKey body CreateTokenRequest true "API key information"
// @Success 201 {object} map[string]interface{} "Created API key"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/api-keys [post]
func (h *APIKeyHandler) CreateAPIKey(c *gin.Context) {
	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("创建API Key失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			logger.Error("解析过期时间失败", logger.F("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format, use ISO 8601"})
			return
		}
		expiresAt = &parsed
	}

	apiKey, err := h.apiKeyService.CreateAPIKey(req.Name, expiresAt)
	if err != nil {
		if err == service.ErrAPIKeyNameExists {
			c.JSON(http.StatusConflict, gin.H{"error": "API key name already exists"})
			return
		}
		logger.Error("创建API Key失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置模型许可
	if req.AllowedModels != nil {
		if err := apiKey.SetAllowedModels(req.AllowedModels); err != nil {
			logger.Error("设置模型许可失败", logger.F("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set allowed models"})
			return
		}
	}

	// 设置按类型配置的策略
	if req.ChatPolicyID != "" {
		apiKey.ChatPolicyID = req.ChatPolicyID
	}
	if req.EmbeddingPolicyID != "" {
		apiKey.EmbeddingPolicyID = req.EmbeddingPolicyID
	}
	if req.RerankPolicyID != "" {
		apiKey.RerankPolicyID = req.RerankPolicyID
	}
	if req.ImagePolicyID != "" {
		apiKey.ImagePolicyID = req.ImagePolicyID
	}
	if req.AudioPolicyID != "" {
		apiKey.AudioPolicyID = req.AudioPolicyID
	}
	if req.VideoPolicyID != "" {
		apiKey.VideoPolicyID = req.VideoPolicyID
	}

	// 如果有配置需要更新，使用完整更新方法
	if req.AllowedModels != nil || req.ChatPolicyID != "" || req.EmbeddingPolicyID != "" || req.RerankPolicyID != "" || req.ImagePolicyID != "" || req.AudioPolicyID != "" || req.VideoPolicyID != "" {
		var chatPolicyID, embeddingPolicyID, rerankPolicyID, imagePolicyID, audioPolicyID, videoPolicyID *string
		if req.ChatPolicyID != "" {
			chatPolicyID = &req.ChatPolicyID
		}
		if req.EmbeddingPolicyID != "" {
			embeddingPolicyID = &req.EmbeddingPolicyID
		}
		if req.RerankPolicyID != "" {
			rerankPolicyID = &req.RerankPolicyID
		}
		if req.ImagePolicyID != "" {
			imagePolicyID = &req.ImagePolicyID
		}
		if req.AudioPolicyID != "" {
			audioPolicyID = &req.AudioPolicyID
		}
		if req.VideoPolicyID != "" {
			videoPolicyID = &req.VideoPolicyID
		}
		
		apiKey, err = h.apiKeyService.UpdateAPIKeyFull(
			apiKey.ID,
			nil, // name 不变
			nil, // status 不变
			req.AllowedModels,
			chatPolicyID,
			embeddingPolicyID,
			rerankPolicyID,
			imagePolicyID,
			audioPolicyID,
			videoPolicyID,
		)
		if err != nil {
			logger.Error("更新API Key配置失败", logger.F("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update API key configuration"})
			return
		}
	}

	logger.Info("创建API Key成功", logger.F("id", apiKey.ID), logger.F("name", apiKey.Name))
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":                apiKey.ID,
			"name":              apiKey.Name,
			"api_key":           apiKey.APIKey, // 只在创建时返回完整API Key
			"prefix":            apiKey.Prefix,
			"status":            apiKey.Status,
			"allowed_models":   apiKey.GetAllowedModels(),
			"chat_policy_id":    apiKey.ChatPolicyID,
			"embedding_policy_id": apiKey.EmbeddingPolicyID,
			"rerank_policy_id":   apiKey.RerankPolicyID,
			"image_policy_id":    apiKey.ImagePolicyID,
			"audio_policy_id":    apiKey.AudioPolicyID,
			"video_policy_id":    apiKey.VideoPolicyID,
			"expires_at":        apiKey.ExpiresAt,
			"created_at":        apiKey.CreatedAt,
		},
	})
}

// CreateToken 保持向后兼容的函数别名
// Deprecated: 使用 CreateAPIKey 代替
func (h *APIKeyHandler) CreateToken(c *gin.Context) {
	h.CreateAPIKey(c)
}

// ListAPIKeys 获取API Key列表
// @Summary List all API keys
// @Description Get a list of all API keys
// @Tags api-keys
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of API keys"
// @Router /api/v1/api-keys [get]
func (h *APIKeyHandler) ListAPIKeys(c *gin.Context) {
	apiKeys, err := h.apiKeyService.ListAPIKeys()
	if err != nil {
		logger.Error("获取API Key列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 隐藏完整API Key值，只显示前缀
	apiKeyList := make([]gin.H, 0, len(apiKeys))
	for _, apiKey := range apiKeys {
		apiKeyList = append(apiKeyList, gin.H{
			"id":                apiKey.ID,
			"name":              apiKey.Name,
			"prefix":            apiKey.Prefix, // 显示前缀
			"status":            apiKey.Status,
			"policy_id":         apiKey.PolicyID, // 向后兼容
			"allowed_models":    apiKey.GetAllowedModels(),
			"chat_policy_id":    apiKey.ChatPolicyID,
			"embedding_policy_id": apiKey.EmbeddingPolicyID,
			"rerank_policy_id":   apiKey.RerankPolicyID,
			"image_policy_id":    apiKey.ImagePolicyID,
			"audio_policy_id":    apiKey.AudioPolicyID,
			"video_policy_id":    apiKey.VideoPolicyID,
			"last_used_at":      apiKey.LastUsedAt,
			"expires_at":        apiKey.ExpiresAt,
			"created_at":        apiKey.CreatedAt,
		})
	}

	logger.Info("获取API Key列表成功", logger.F("count", len(apiKeys)))
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"items": apiKeyList, "total": len(apiKeys)}})
}

// ListTokens 保持向后兼容的函数别名
// Deprecated: 使用 ListAPIKeys 代替
func (h *APIKeyHandler) ListTokens(c *gin.Context) {
	apiKeys, err := h.apiKeyService.ListTokens()
	if err != nil {
		logger.Error("获取API Key列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 隐藏完整API Key值，只显示前缀
	tokenList := make([]gin.H, 0, len(apiKeys))
	for _, apiKey := range apiKeys {
		tokenList = append(tokenList, gin.H{
			"id":           apiKey.ID,
			"name":         apiKey.Name,
			"prefix":      apiKey.Prefix, // 显示前缀
			"status":       apiKey.Status,
			"policy_id":    apiKey.PolicyID,
			"last_used_at": apiKey.LastUsedAt,
			"expires_at":   apiKey.ExpiresAt,
			"created_at":   apiKey.CreatedAt,
		})
	}

	logger.Info("获取API Key列表成功", logger.F("count", len(apiKeys)))
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"items": tokenList, "total": len(apiKeys)}})
}

// GetToken 获取单个API Key
// @Summary Get API key by ID
// @Description Get a specific API key by ID
// @Tags tokens
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 200 {object} map[string]interface{} "API key details"
// @Failure 404 {object} map[string]string "API key not found"
// @Router /api/v1/tokens/{id} [get]
// GetAPIKey 获取单个API Key
// @Summary Get API key by ID
// @Description Get a specific API key by ID
// @Tags api-keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 200 {object} map[string]interface{} "API key details"
// @Failure 404 {object} map[string]string "API key not found"
// @Router /api/v1/api-keys/{id} [get]
func (h *APIKeyHandler) GetAPIKey(c *gin.Context) {
	id := c.Param("id")
	apiKey, err := h.apiKeyService.GetAPIKey(id)
	if err != nil {
		if err == service.ErrAPIKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
			return
		}
		logger.Error("获取API Key失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":                apiKey.ID,
			"name":              apiKey.Name,
			"api_key":           apiKey.APIKey, // 返回完整API Key（管理员权限）
			"prefix":            apiKey.Prefix, // 同时返回前缀用于显示
			"status":            apiKey.Status,
			"policy_id":         apiKey.PolicyID, // 向后兼容
			"allowed_models":    apiKey.GetAllowedModels(),
			"chat_policy_id":    apiKey.ChatPolicyID,
			"embedding_policy_id": apiKey.EmbeddingPolicyID,
			"rerank_policy_id":   apiKey.RerankPolicyID,
			"image_policy_id":    apiKey.ImagePolicyID,
			"audio_policy_id":    apiKey.AudioPolicyID,
			"video_policy_id":    apiKey.VideoPolicyID,
			"last_used_at":      apiKey.LastUsedAt,
			"expires_at":        apiKey.ExpiresAt,
			"created_at":        apiKey.CreatedAt,
		},
	})
}

// GetToken 保持向后兼容的函数别名
// Deprecated: 使用 GetAPIKey 代替
func (h *APIKeyHandler) GetToken(c *gin.Context) {
	id := c.Param("id")
	apiKey, err := h.apiKeyService.GetToken(id)
	if err != nil {
		if err == service.ErrTokenNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
			return
		}
		logger.Error("获取API Key失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":           apiKey.ID,
			"name":         apiKey.Name,
			"token":        apiKey.APIKey, // 返回完整API Key（管理员权限）
			"prefix":       apiKey.Prefix, // 同时返回前缀用于显示
			"status":       apiKey.Status,
			"policy_id":    apiKey.PolicyID,
			"last_used_at": apiKey.LastUsedAt,
			"expires_at":   apiKey.ExpiresAt,
			"created_at":   apiKey.CreatedAt,
		},
	})
}

// UpdateAPIKey 更新API Key
// @Summary Update API key
// @Description Update API key information
// @Tags api-keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Param apiKey body UpdateTokenRequest true "API key update information"
// @Success 200 {object} map[string]interface{} "Updated API key"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "API key not found"
// @Router /api/v1/api-keys/{id} [put]
func (h *APIKeyHandler) UpdateAPIKey(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("更新API Key失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey, err := h.apiKeyService.UpdateAPIKeyFull(
		id,
		req.Name,
		req.Status,
		req.AllowedModels,
		req.ChatPolicyID,
		req.EmbeddingPolicyID,
		req.RerankPolicyID,
		req.ImagePolicyID,
		req.AudioPolicyID,
		req.VideoPolicyID,
	)
	if err != nil {
		if err == service.ErrAPIKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
			return
		}
		if err == service.ErrAPIKeyNameExists {
			c.JSON(http.StatusConflict, gin.H{"error": "API key name already exists"})
			return
		}
		logger.Error("更新API Key失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("更新API Key成功", logger.F("id", id), logger.F("name", apiKey.Name))
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":                apiKey.ID,
			"name":              apiKey.Name,
			"prefix":            apiKey.Prefix,
			"status":            apiKey.Status,
			"policy_id":         apiKey.PolicyID, // 向后兼容
			"allowed_models":    apiKey.GetAllowedModels(),
			"chat_policy_id":    apiKey.ChatPolicyID,
			"embedding_policy_id": apiKey.EmbeddingPolicyID,
			"rerank_policy_id":   apiKey.RerankPolicyID,
			"image_policy_id":    apiKey.ImagePolicyID,
			"audio_policy_id":    apiKey.AudioPolicyID,
			"video_policy_id":    apiKey.VideoPolicyID,
			"last_used_at":      apiKey.LastUsedAt,
			"expires_at":        apiKey.ExpiresAt,
			"updated_at":        apiKey.UpdatedAt,
		},
	})
}

// UpdateToken 保持向后兼容的函数别名
// Deprecated: 使用 UpdateAPIKey 代替
func (h *APIKeyHandler) UpdateToken(c *gin.Context) {
	h.UpdateAPIKey(c)
}

// DeleteAPIKey 删除API Key
// @Summary Delete API key
// @Description Delete an API key
// @Tags api-keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 204 "API key deleted"
// @Failure 404 {object} map[string]string "API key not found"
// @Router /api/v1/api-keys/{id} [delete]
func (h *APIKeyHandler) DeleteAPIKey(c *gin.Context) {
	id := c.Param("id")
	if err := h.apiKeyService.DeleteAPIKey(id); err != nil {
		if err == service.ErrAPIKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
			return
		}
		logger.Error("删除API Key失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("删除API Key成功", logger.F("id", id))
	c.Status(http.StatusNoContent)
}

// DeleteToken 保持向后兼容的函数别名
// Deprecated: 使用 DeleteAPIKey 代替
func (h *APIKeyHandler) DeleteToken(c *gin.Context) {
	h.DeleteAPIKey(c)
}

// ResetAPIKey 重置API Key
// @Summary Reset API key
// @Description Reset API key value (generate new API key)
// @Tags api-keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 200 {object} map[string]interface{} "Reset API key with new value"
// @Failure 404 {object} map[string]string "API key not found"
// @Router /api/v1/api-keys/{id}/reset [post]
func (h *APIKeyHandler) ResetAPIKey(c *gin.Context) {
	id := c.Param("id")
	apiKey, err := h.apiKeyService.ResetAPIKey(id)
	if err != nil {
		if err == service.ErrAPIKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
			return
		}
		logger.Error("重置API Key失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("重置API Key成功", logger.F("id", id), logger.F("name", apiKey.Name))
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":         apiKey.ID,
			"name":       apiKey.Name,
			"api_key":    apiKey.APIKey, // 重置时返回完整API Key
			"prefix":     apiKey.Prefix,
			"status":     apiKey.Status,
			"policy_id":  apiKey.PolicyID,
			"updated_at": apiKey.UpdatedAt,
		},
	})
}

// ResetToken 保持向后兼容的函数别名
// Deprecated: 使用 ResetAPIKey 代替
func (h *APIKeyHandler) ResetToken(c *gin.Context) {
	id := c.Param("id")
	apiKey, err := h.apiKeyService.ResetToken(id)
	if err != nil {
		if err == service.ErrTokenNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
			return
		}
		logger.Error("重置API Key失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("重置API Key成功", logger.F("id", id), logger.F("name", apiKey.Name))
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":         apiKey.ID,
			"name":       apiKey.Name,
			"token":      apiKey.APIKey, // 重置时返回完整API Key
			"prefix":     apiKey.Prefix,
			"status":     apiKey.Status,
			"policy_id":  apiKey.PolicyID,
			"updated_at": apiKey.UpdatedAt,
		},
	})
}

// SetTokenPolicyRequest 设置API Key策略请求
type SetTokenPolicyRequest struct {
	PolicyID string `json:"policy_id"`
}

// SetAPIKeyPolicy 设置API Key的策略
// @Summary Set API key policy
// @Description Assign a policy to an API key
// @Tags api-keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Param request body SetTokenPolicyRequest true "Policy assignment"
// @Success 200 {object} map[string]interface{} "API key with policy"
// @Failure 404 {object} map[string]string "API key or policy not found"
// @Router /api/v1/api-keys/{id}/policy [put]
func (h *APIKeyHandler) SetAPIKeyPolicy(c *gin.Context) {
	id := c.Param("id")
	var req SetTokenPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("设置API Key策略失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey, err := h.apiKeyService.UpdateAPIKeyPolicy(id, req.PolicyID)
	if err != nil {
		if err == service.ErrAPIKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
			return
		}
		logger.Error("设置API Key策略失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("设置API Key策略成功", logger.F("id", id), logger.F("policy_id", req.PolicyID))
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        apiKey.ID,
			"name":      apiKey.Name,
			"policy_id": apiKey.PolicyID,
			"status":    apiKey.Status,
		},
	})
}

// SetTokenPolicy 保持向后兼容的函数别名
// Deprecated: 使用 SetAPIKeyPolicy 代替
func (h *APIKeyHandler) SetTokenPolicy(c *gin.Context) {
	h.SetAPIKeyPolicy(c)
}

// RemoveAPIKeyPolicy 移除API Key的策略
// @Summary Remove API key policy
// @Description Remove policy assignment from an API key
// @Tags api-keys
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 200 {object} map[string]interface{} "API key without policy"
// @Failure 404 {object} map[string]string "API key not found"
// @Router /api/v1/api-keys/{id}/policy [delete]
func (h *APIKeyHandler) RemoveAPIKeyPolicy(c *gin.Context) {
	id := c.Param("id")
	apiKey, err := h.apiKeyService.RemoveAPIKeyPolicy(id)
	if err != nil {
		if err == service.ErrAPIKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
			return
		}
		logger.Error("移除API Key策略失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("移除API Key策略成功", logger.F("id", id))
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        apiKey.ID,
			"name":      apiKey.Name,
			"policy_id": apiKey.PolicyID,
			"status":    apiKey.Status,
		},
	})
}

// RemoveTokenPolicy 保持向后兼容的函数别名
// Deprecated: 使用 RemoveAPIKeyPolicy 代替
func (h *APIKeyHandler) RemoveTokenPolicy(c *gin.Context) {
	h.RemoveAPIKeyPolicy(c)
}
