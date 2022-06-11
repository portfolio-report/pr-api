package tags

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/libs"
)

type securityTagRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

type putTagRequest struct {
	Securities []securityTagRequest `json:"securities"`
}

// PutTag creates or updates tag with associated securities
func (h *tagsHandler) PutTag(c *gin.Context) {
	name := c.Param("name")

	var request putTagRequest
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	uuids := make([]uuid.UUID, len(request.Securities))
	for i := range request.Securities {
		uuids[i] = request.Securities[i].UUID
	}

	securities, err := h.SecurityService.UpsertTag(name, uuids)
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":       name,
		"securities": securities,
	})
}
