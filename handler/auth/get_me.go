package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

func (h *AuthHandler) GetMe(c *gin.Context) {
	user := middleware.UserFromContext(c.Request.Context())
	c.JSON(http.StatusOK, user)
}
