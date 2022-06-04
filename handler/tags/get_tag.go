package tags

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTag returns single tag with associated securities
func (h *tagsHandler) GetTag(c *gin.Context) {
	name := c.Param("name")

	securities := h.SecurityService.GetSecuritiesByTag(name)

	c.JSON(http.StatusOK, gin.H{
		"name":       name,
		"securities": securities,
	})
}
