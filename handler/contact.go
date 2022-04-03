package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
)

type contactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// Contact sends an email to the configured contact address
func (h *rootHandler) Contact(c *gin.Context) {
	var request contactRequest

	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	if h.MailerService == nil {
		panic(errors.New("could not send email, not configured"))
	}

	// Sanitize input
	request.Subject = strings.ReplaceAll(request.Subject, "\n", "")
	request.Subject = strings.ReplaceAll(request.Subject, "\r", "")
	request.Name = strings.ReplaceAll(request.Name, "\n", "")
	request.Name = strings.ReplaceAll(request.Name, "\r", "")
	request.Email = strings.ReplaceAll(request.Email, "\n", "")
	request.Email = strings.ReplaceAll(request.Email, "\r", "")

	err := h.MailerService.SendContactMail(
		request.Email,
		request.Name,
		request.Subject,
		request.Message,
		c.ClientIP())
	if err != nil {
		panic(err)
	}

	fmt.Println("Mail sent from", request.Email, "with subject:", request.Subject)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
