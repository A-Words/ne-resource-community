package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health returns simple status for readiness probes.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
