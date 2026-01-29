package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the LingProxy service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Service status"
// @Router /api/v1/health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "LingProxy is running",
		"version": "1.0.0",
	})
}