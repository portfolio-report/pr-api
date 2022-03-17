package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRoot returns static ok message
func (h *Handler) GetRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "ok"})
}
