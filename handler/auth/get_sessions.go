package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

// GetSessions returns all sessions of current user
func (h *authHandler) GetSessions(c *gin.Context) {
	user := middleware.UserFromContext(c.Request.Context())
	sessions, err := h.SessionService.GetAllOfUser(user)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, sessions)
}
