package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
)

// PostSecurity creates new security
func (h *securitiesHandler) PostSecurity(c *gin.Context) {
	var request model.SecurityInput
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	security, err := h.SecurityService.CreateSecurity(&request)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, security)
}
