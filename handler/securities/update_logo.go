package securities

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/libs"
	"gorm.io/gorm"
)

// UpdateLogo stores new/changed logo of security
func (h *securitiesHandler) UpdateLogo(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		libs.HandleBadRequestError(c, "No file received")
		return
	}

	extension := filepath.Ext(file.Filename)
	if extension != ".jpg" && extension != ".png" {
		libs.HandleBadRequestError(c, "Unknown file extension")
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer openedFile.Close()

	logoUrl, err := h.SecurityService.UpdateLogo(uuid, openedFile, extension)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			libs.HandleNotFoundError(c)
			return
		}

		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"logoUrl": logoUrl})
}
