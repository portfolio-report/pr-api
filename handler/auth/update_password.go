package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
)

type updatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=255"`
}

// UpdatePassword changes the password for the current user
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	user := middleware.UserFromContext(c.Request.Context())

	var request updatePasswordRequest
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	ctx := c.Request.Context()

	valid, err := h.UserService.VerifyPassword(ctx, user, request.OldPassword)
	if err != nil {
		panic(err)
	}
	if !valid {
		libs.HandleForbiddenError(c, "Password is wrong.")
		return
	}

	if err := h.UserService.UpdatePassword(ctx, user, request.NewPassword); err != nil {
		panic(err)
	}

	c.Status(http.StatusCreated)
}
