package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHealth mounts a liveness/readiness endpoint.
func RegisterHealth(rg *gin.RouterGroup) {
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
