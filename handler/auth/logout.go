package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) LogoutUser(c *gin.Context) {
	token := h.SessionService.GetSessionToken(c)
	if _, err := h.SessionService.DeleteSession(token); err != nil {
		panic(err)
	}
	c.Status(http.StatusNoContent)
}
