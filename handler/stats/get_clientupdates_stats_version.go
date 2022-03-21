package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
)

// GetClientupdatesStatsVersion returns statistics for updates to a certain version
func (h *statsHandler) GetClientupdatesStatsVersion(c *gin.Context) {
	version := c.Param("version")

	type ByDate struct {
		Date  model.Date `json:"date"`
		Count int        `json:"count"`
	}

	var byDate []ByDate
	err := h.DB.Raw(`
		SELECT date(timestamp) AS date, count(*) AS count
		FROM clientupdates
		WHERE version = ?
		GROUP BY date
		ORDER BY date ASC
		`, version).Scan(&byDate).Error

	if err != nil {
		panic(err)
	}

	type ByCountry struct {
		Country string `json:"country"`
		Count   int    `json:"count"`
	}

	var byCountry []ByCountry
	err = h.DB.Raw(`
		SELECT COALESCE(country, '') AS country, count(*) AS count
		FROM clientupdates
		WHERE version = ?
		GROUP BY country
		`, version).Scan(&byCountry).Error

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"byDate": byDate, "byCountry": byCountry})
}
