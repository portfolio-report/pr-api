package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

// DeleteMe removes current user
func (h *authHandler) DeleteMe(c *gin.Context) {
	user := middleware.UserFromContext(c.Request.Context())

	if err := h.UserService.Delete(user); err != nil {
		panic(err)
	}

	c.Status(http.StatusNoContent)
}
