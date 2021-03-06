package stats

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
)

// CountClientupdate counts each request as an update of Portfolio Performance
func (h *statsHandler) CountClientupdate(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")

	timestamp := time.Now().UTC()
	version := c.Param("version")
	country := h.GeoipService.GetCountryFromIp(c.ClientIP())
	useragent := c.GetHeader("User-Agent")

	clientupdate := db.Clientupdate{
		Timestamp: timestamp,
		Version:   version,
		Country:   &country,
		Useragent: &useragent,
	}

	if err := h.DB.Create(&clientupdate).Error; err != nil {
		panic(err)
	}
}
