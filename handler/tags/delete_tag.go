package tags

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteTag removes tag
func (h *tagsHandler) DeleteTag(c *gin.Context) {
	name := c.Param("name")

	h.SecurityService.DeleteTag(name)

	c.Status(http.StatusNoContent)
}
