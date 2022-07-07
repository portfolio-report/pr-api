package securities

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/libs"
	"gorm.io/gorm"
)

// DeleteLogo deletes logo of security
func (h *securitiesHandler) DeleteLogo(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	h.SecurityService.DeleteLogo(uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			libs.HandleNotFoundError(c)
			return
		}

		panic(err)
	}

	c.Status(http.StatusNoContent)
}
