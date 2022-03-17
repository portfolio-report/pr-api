package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LogoutUser logs out current user by deleting the session
func (h *authHandler) LogoutUser(c *gin.Context) {
	token := h.SessionService.GetSessionToken(c)
	if _, err := h.SessionService.DeleteSession(token); err != nil {
		panic(err)
	}
	c.Status(http.StatusNoContent)
}
