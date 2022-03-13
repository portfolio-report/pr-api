package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/service"
)

type registerUserRequest struct {
	Username string `json:"username" binding:"required,min=6,max=100,ValidUsername"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var request registerUserRequest

	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	ctx := c.Request.Context()

	user, err := h.UserService.Create(request.Username)
	if err != nil {
		if err == service.UserExistsAlreadyError {
			libs.HandleBadRequestError(c, err.Error())
			return
		}

		panic(err)
	}

	if err := h.UserService.UpdatePassword(ctx, user, request.Password); err != nil {
		panic(err)
	}

	session, err := h.SessionService.CreateSession(user, c.GetHeader("User-Agent"))
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, session)
}
