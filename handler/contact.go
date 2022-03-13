package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
)

type ContactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (h *Handler) Contact(c *gin.Context) {
	var request ContactRequest

	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	if h.MailerService == nil {
		panic(errors.New("could not send email, not configured"))
	}

	err := h.MailerService.SendContactMail(
		request.Email,
		request.Name,
		request.Subject+" (via Portfolio Report)",
		request.Message,
		c.ClientIP())
	if err != nil {
		panic(err)
	}

	fmt.Println("Mail sent from", request.Email, "with subject:", request.Subject)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
