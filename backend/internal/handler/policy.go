package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/service"
)

// PolicyHandler 策略处理器
type PolicyHandler struct {
	policyService   *service.PolicyService
	templateService *service.PolicyTemplateService
}

// NewPolicyHandler 创建新的策略处理器
func NewPolicyHandler(policyService *service.PolicyService, templateService *service.PolicyTemplateService) *PolicyHandler {
	return &PolicyHandler{
		policyService:   policyService,
		templateService: templateService,
	}
}

// CreatePolicyRequest 创建策略请求
type CreatePolicyRequest struct {
	Name       string                 `json:"name" binding:"required"`
	TemplateID string                 `json:"template_id" binding:"required"`
	Parameters map[string]interface{} `json:"parameters" binding:"required"`
}

// UpdatePolicyRequest 更新策略请求
type UpdatePolicyRequest struct {
	Name       *string                `json:"name,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Enabled    *bool                  `json:"enabled,omitempty"`
}

// ListPolicyTemplates 获取策略模板列表
// @Summary List policy templates
// @Description Get all available policy templates
// @Tags policies
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of policy templates"
// @Router /api/v1/policy-templates [get]
func (h *PolicyHandler) ListPolicyTemplates(c *gin.Context) {
	templates, err := h.templateService.ListTemplates()
	if err != nil {
		logger.Error("获取策略模板列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("获取策略模板列表成功", logger.F("count", len(templates)))
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"items": templates, "total": len(templates)}})
}

// GetPolicyTemplate 获取策略模板
// @Summary Get policy template by ID
// @Description Get a specific policy template by ID
// @Tags policies
// @Accept json
// @Produce json
// @Param id path string true "Policy Template ID"
// @Success 200 {object} map[string]interface{} "Policy template details"
// @Failure 404 {object} map[string]string "Template not found"
// @Router /api/v1/policy-templates/{id} [get]
func (h *PolicyHandler) GetPolicyTemplate(c *gin.Context) {
	id := c.Param("id")
	template, err := h.templateService.GetTemplate(id)
	if err != nil {
		if err == service.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy template not found"})
			return
		}
		logger.Error("获取策略模板失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": template})
}

// ListPolicies 获取策略列表
// @Summary List all policies
// @Description Get a list of all policies
// @Tags policies
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of policies"
// @Router /api/v1/policies [get]
func (h *PolicyHandler) ListPolicies(c *gin.Context) {
	policies, err := h.policyService.ListPolicies()
	if err != nil {
		logger.Error("获取策略列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("获取策略列表成功", logger.F("count", len(policies)))
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"items": policies, "total": len(policies)}})
}

// GetPolicy 获取单个策略
// @Summary Get policy by ID
// @Description Get a specific policy by ID
// @Tags policies
// @Accept json
// @Produce json
// @Param id path string true "Policy ID"
// @Success 200 {object} map[string]interface{} "Policy details"
// @Failure 404 {object} map[string]string "Policy not found"
// @Router /api/v1/policies/{id} [get]
func (h *PolicyHandler) GetPolicy(c *gin.Context) {
	id := c.Param("id")
	policy, err := h.policyService.GetPolicy(id)
	if err != nil {
		if err == service.ErrPolicyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}
		logger.Error("获取策略失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 解析参数为JSON对象
	var params map[string]interface{}
	if policy.Parameters != "" {
		if err := json.Unmarshal([]byte(policy.Parameters), &params); err != nil {
			logger.Warn("解析策略参数失败", logger.F("error", err.Error()), logger.F("id", id))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          policy.ID,
			"name":        policy.Name,
			"template_id": policy.TemplateID,
			"type":        policy.Type,
			"parameters":  params,
			"enabled":     policy.Enabled,
			"created_at":  policy.CreatedAt,
			"updated_at":  policy.UpdatedAt,
		},
	})
}

// CreatePolicy 创建策略
// @Summary Create a new policy
// @Description Create a new routing policy
// @Tags policies
// @Accept json
// @Produce json
// @Param policy body CreatePolicyRequest true "Policy information"
// @Success 201 {object} map[string]interface{} "Created policy"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /api/v1/policies [post]
func (h *PolicyHandler) CreatePolicy(c *gin.Context) {
	var req CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("创建策略失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy, err := h.policyService.CreatePolicy(req.Name, req.TemplateID, req.Parameters)
	if err != nil {
		if err == service.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy template not found"})
			return
		}
		logger.Error("创建策略失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 解析参数
	var params map[string]interface{}
	if policy.Parameters != "" {
		json.Unmarshal([]byte(policy.Parameters), &params)
	}

	logger.Info("创建策略成功", logger.F("id", policy.ID), logger.F("name", policy.Name))
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":          policy.ID,
			"name":        policy.Name,
			"template_id": policy.TemplateID,
			"type":        policy.Type,
			"parameters":  params,
			"enabled":     policy.Enabled,
			"created_at":  policy.CreatedAt,
		},
	})
}

// UpdatePolicy 更新策略
// @Summary Update policy
// @Description Update policy information
// @Tags policies
// @Accept json
// @Produce json
// @Param id path string true "Policy ID"
// @Param policy body UpdatePolicyRequest true "Policy update information"
// @Success 200 {object} map[string]interface{} "Updated policy"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Policy not found"
// @Router /api/v1/policies/{id} [put]
func (h *PolicyHandler) UpdatePolicy(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("更新策略失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy, err := h.policyService.UpdatePolicy(id, req.Name, req.Parameters, req.Enabled)
	if err != nil {
		if err == service.ErrPolicyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}
		if err == service.ErrTemplateNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "policy template not found"})
			return
		}
		logger.Error("更新策略失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 解析参数
	var params map[string]interface{}
	if policy.Parameters != "" {
		json.Unmarshal([]byte(policy.Parameters), &params)
	}

	logger.Info("更新策略成功", logger.F("id", id), logger.F("name", policy.Name))
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":          policy.ID,
			"name":        policy.Name,
			"template_id": policy.TemplateID,
			"type":        policy.Type,
			"parameters":  params,
			"enabled":     policy.Enabled,
			"updated_at":  policy.UpdatedAt,
		},
	})
}

// DeletePolicy 删除策略
// @Summary Delete policy
// @Description Delete a policy
// @Tags policies
// @Accept json
// @Produce json
// @Param id path string true "Policy ID"
// @Success 204 "Policy deleted"
// @Failure 404 {object} map[string]string "Policy not found"
// @Router /api/v1/policies/{id} [delete]
func (h *PolicyHandler) DeletePolicy(c *gin.Context) {
	id := c.Param("id")
	if err := h.policyService.DeletePolicy(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
			return
		}
		if strings.Contains(err.Error(), "cannot delete builtin") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete builtin policy"})
			return
		}
		logger.Error("删除策略失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("删除策略成功", logger.F("id", id))
	c.Status(http.StatusNoContent)
}
