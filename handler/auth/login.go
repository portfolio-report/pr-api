package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
	"gorm.io/gorm"
)

type loginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var request loginUserRequest

	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	ctx := c.Request.Context()

	user, err := h.UserService.GetUserByUsername(ctx, request.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			libs.HandleUnauthorizedError(c)
			return
		}

		panic(err)
	}

	valid, err := h.UserService.VerifyPassword(ctx, user, request.Password)
	if err != nil {
		panic(err)
	}
	if !valid {
		libs.HandleUnauthorizedError(c)
		return
	}

	session, err := h.SessionService.CreateSession(user, c.GetHeader("User-Agent"))
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, session)

	go func() {
		err := h.SessionService.CleanupExpiredSessions()
		if err != nil {
			fmt.Println("Error in background processing:", err.Error())
		}
	}()
}
