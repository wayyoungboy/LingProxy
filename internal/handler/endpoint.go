package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/storage"
)

// EndpointHandler 端点处理器
type EndpointHandler struct {
	storage *storage.StorageFacade
}

// NewEndpointHandler 创建新的端点处理器
func NewEndpointHandler(storage *storage.StorageFacade) *EndpointHandler {
	return &EndpointHandler{
		storage: storage,
	}
}

// ListEndpoints godoc
// @Summary List all endpoints
// @Description Get a list of all API endpoints
// @Tags endpoints
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of endpoints"
// @Router /api/v1/endpoints [get]
func (h *EndpointHandler) ListEndpoints(c *gin.Context) {
	endpoints, err := h.storage.ListEndpoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": endpoints})
}

// GetEndpoint godoc
// @Summary Get endpoint by ID
// @Description Get a specific endpoint by ID
// @Tags endpoints
// @Accept json
// @Produce json
// @Param id path string true "Endpoint ID"
// @Success 200 {object} map[string]interface{} "Endpoint details"
// @Failure 404 {object} map[string]string "Endpoint not found"
// @Router /api/v1/endpoints/{id} [get]
func (h *EndpointHandler) GetEndpoint(c *gin.Context) {
	id := c.Param("id")
	endpoint, err := h.storage.GetEndpoint(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": endpoint})
}

// CreateEndpoint godoc
// @Summary Create a new endpoint
// @Description Create a new API endpoint configuration
// @Tags endpoints
// @Accept json
// @Produce json
// @Param endpoint body storage.Endpoint true "Endpoint configuration"
// @Success 201 {object} map[string]interface{} "Created endpoint"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/endpoints [post]
func (h *EndpointHandler) CreateEndpoint(c *gin.Context) {
	var endpoint storage.Endpoint
	if err := c.ShouldBindJSON(&endpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.CreateEndpoint(&endpoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": endpoint})
}